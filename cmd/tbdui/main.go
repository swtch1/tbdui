package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/gizak/termui/v3"
	"github.com/sirupsen/logrus"
	"github.com/swtch1/tbdui/char"
	"github.com/swtch1/tbdui/component"
	"github.com/swtch1/tbdui/conf"
	"github.com/swtch1/tbdui/dynamodb"
	"github.com/swtch1/tbdui/logger"
)

var (
	// outputLog will enable dumping the application log to the output box instead of the standard output text
	outputLog = false
	debugLog  = true
)

const (
	defaultBorderColor   = termui.ColorWhite
	defaultSelectedColor = termui.ColorGreen
)

func main() {
	log := logger.NewUILogger()

	dynDB, err := newDynamoDB(log)
	if err != nil {
		fmt.Println("failed to setup connection to DynamoDB:", err)
		os.Exit(1)
	}

	if err := termui.Init(); err != nil {
		logrus.WithError(err).Fatal("failed to initialize termui")
	}
	defer termui.Close()

	input := make(chan string)
	tui := newTUI(dynDB, input, log)

	errCh := make(chan error)
	go func() {
		// if the tui errors signal for app termination
		if err := tui.Run(conf.NewDefault()); err != nil {
			errCh <- err
		}
	}()

	uiEvents := termui.PollEvents()
	for {
		select {
		case err := <-errCh:
			logrus.WithError(err).Fatal("main UI exited")
		case e := <-uiEvents:
			switch e.ID {
			case char.CTRL_C:
				return
			default:
				input <- e.ID
			}
		}
	}
}

func newDynamoDB(l *logger.UILogger) (*dynamodb.DB, error) {
	accessKeyID, ok := os.LookupEnv("AWS_ACCESS_KEY_ID")
	if !ok {
		return nil, fmt.Errorf("AWS_ACCESS_KEY_ID env var is required")
	}
	secretAccessKey, ok := os.LookupEnv("AWS_SECRET_ACCESS_KEY")
	if !ok {
		return nil, fmt.Errorf("AWS_SECRET_ACCESS_KEY env var is required")
	}
	region, ok := os.LookupEnv("AWS_REGION")
	if !ok {
		return nil, fmt.Errorf("AWS_REGION env var is required")
	}
	environment, ok := os.LookupEnv("ENVIRONMENT")
	if !ok {
		return nil, fmt.Errorf("ENVIRONMENT env var is required")
	}

	db, err := dynamodb.NewDB(credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""), region, environment)
	if err != nil {
		return nil, err
	}
	db.SetLogger(l)
	return db, nil
}

// TUI is a terminal user interface.
type TUI struct {
	db      *dynamodb.DB
	inputCh chan string
	logger  *logger.UILogger
}

// Run starts a continuous loop that will draw the screen
func (ui TUI) Run(c conf.Config) error {
	termWidth, termHeight := termui.TerminalDimensions()

	borderWidth := 1
	startHeight := 0
	leftBorder := borderWidth
	rightBorder := termWidth - borderWidth
	inputBoxHeight := 3
	inputBoxWidth := termWidth / 3

	// text across the top
	topText := component.NewInfo("<Ctrl + c> to quit", component.Dimensions{
		X1: leftBorder,
		Y1: 0,
		X2: termWidth - borderWidth*2,
		Y2: startHeight + inputBoxHeight,
	})

	// integrations search box on the left
	start := 3
	searchBox := component.NewInputBox("Search Integrations", ":type to search", c, component.Dimensions{
		X1: leftBorder,
		Y1: start,
		X2: inputBoxWidth,
		Y2: start + inputBoxHeight,
	})

	// company filter box on the left
	start = 6
	companyFilterBox := component.NewInputBox("Filter Company", ":type company ID to filter results", c, component.Dimensions{
		X1: leftBorder,
		Y1: start,
		X2: inputBoxWidth,
		Y2: start + inputBoxHeight,
	})

	// table filter box on the left
	start = 9
	tableFilterBox := component.NewInputBox("Filter Table", ":type partial table name to filter", c, component.Dimensions{
		X1: leftBorder,
		Y1: start,
		X2: inputBoxWidth,
		Y2: start + inputBoxHeight,
	})

	start = 12
	tableList := component.NewList("Select Table", c, component.Dimensions{
		X1: leftBorder,
		Y1: start,
		X2: inputBoxWidth,
		Y2: termHeight - borderWidth,
	})

	// output box, on the right
	outputBox := component.NewInputBox("", "try searching...", c, component.Dimensions{
		X1: leftBorder + inputBoxWidth,
		Y1: searchBox.Dimensions().Y1,
		X2: rightBorder,
		Y2: termHeight - borderWidth,
	})
	outputBox.HideUnselectedText = false

	// set all types to be rendered here so we can switch things on and off
	mr := NewMassRenderer([]Renderable{
		topText,
		searchBox,
		companyFilterBox,
		tableFilterBox,
		tableList,
		outputBox,
	})

	// do you want tabs? because this is how you get tabs!
	tabOrder := []Selectable{searchBox, companyFilterBox, tableFilterBox, tableList, outputBox}
	sh := NewSelectionHandler(tabOrder, ui.logger)

	var selected Writer = sh.Next()
	for {
		mr.Render()
		select {
		case c := <-ui.inputCh:
			// cheap debug logging
			if debugLog {
				ui.Log("received input: %v", c)
			}

			switch c {

			// switch between elements
			case char.TAB:
				selected = sh.Next()

			// display output text, or app log depending, to the output box
			case char.ENTER:
				// log debug info when necessary
				if outputLog {
					outputBox.Overwrite(ui.logger.Dump())
					continue
				}

				// get all the integrations that match
				// names, err := ui.db.MatchingIntegrationIDs(searchBox.Contents())
				// if err != nil {
				// 	outputBox.Overwrite(err.Error())
				// 	continue
				// }
				integrations, err := ui.db.AllIntegrations()
				if err != nil {
					outputBox.Overwrite(err.Error())
					continue
				}

				b, err := json.MarshalIndent(integrations[0], "", "  ")
				if err != nil {
					outputBox.Overwrite(err.Error())
					continue
				}
				outputBox.Overwrite(string(b))

			// flush the app log
			case char.CTRL_F:
				ui.logger.Flush()

			// toggle output log
			case char.CTRL_L:
				if outputLog {
					outputLog = false
					continue
				}
				outputLog = true

			// write text to the selected box
			default:
				selected.Write((c))
			}
		}
	}
}

// Renderable types can be rendered.
type Renderable interface {
	Render()
}

// MassRenderer holds types and renders them all at once.
type MassRenderer struct {
	all []Renderable
}

// NewMassRenderer instantiates a new MassRenderer.
func NewMassRenderer(r []Renderable) *MassRenderer {
	return &MassRenderer{
		all: r,
	}
}

// Render all types contained.
func (r *MassRenderer) Render() {
	for i := range r.all {
		r.all[i].Render()
	}
}

func renderAll(r []interface{ Render() }) {
	for _, x := range r {
		x.Render()
	}
}

// Writer can write text.
type Writer interface {
	Write(string)
}

// SelectionHandler controls element selection.
type SelectionHandler struct {
	components  []Selectable
	selectedIdx int
	logger      *logger.UILogger
}

// NewSelectionHandler instantiates a new TabHandler.
func NewSelectionHandler(s []Selectable, l *logger.UILogger) *SelectionHandler {
	return &SelectionHandler{
		components: s,
		logger:     l,
	}
}

// Next selects the next component in the list and deselects all others.
func (h *SelectionHandler) Next() Writer {
	for _, c := range h.components {
		c.Deselect()
	}
	defer func() {
		if h.selectedIdx == len(h.components)-1 {
			h.selectedIdx = 0
			return
		}
		h.selectedIdx++
	}()
	h.logger.Write("selection handler", "selecting component at index %d", h.selectedIdx)
	h.components[h.selectedIdx].Select()
	return h.components[h.selectedIdx]
}

// Previous selects the previous component in the list and deselects all others.
func (h *SelectionHandler) Previous() Writer {
	for _, c := range h.components {
		c.Deselect()
	}
	defer func() {
		if h.selectedIdx == 0 {
			h.selectedIdx = len(h.components) - 1
			return
		}
		h.selectedIdx--
	}()
	h.logger.Write("selection handler", "selection component at index %d", h.selectedIdx)
	return h.components[h.selectedIdx]
}

func newTUI(db *dynamodb.DB, input chan string, l *logger.UILogger) *TUI {
	return &TUI{
		db:      db,
		inputCh: input,
		logger:  l,
	}
}

// Log writes a log message to the UI log.
func (ui TUI) Log(msg string, args ...interface{}) {
	ui.logger.Write("tui", msg, args...)
}

// Selectable describes UI items that can be selected.  When
// you tab to, or move to, a component it is selected.
type Selectable interface {
	Select()
	Deselect()
	Writer
}
