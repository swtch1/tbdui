package component

import (
	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/swtch1/tbdui/conf"
	"github.com/swtch1/tbdui/logger"
)

// List records and displays input from a user.
type List struct {
	ls *widgets.List

	selected   bool
	dimensions Dimensions

	borderColor   termui.Color
	selectedColor termui.Color

	logger *logger.UILogger
}

// NewList initializes an input box.
func NewList(title string, c conf.Config, d Dimensions) *List {
	l := widgets.NewList()
	l.Title = title
	l.SetRect(d.X1, d.Y1, d.X2, d.Y2)
	// just hard coding these for now
	l.SelectedRowStyle.Bg = termui.ColorWhite
	l.SelectedRowStyle.Fg = termui.ColorBlack
	return &List{
		ls:            l,
		dimensions:    d,
		borderColor:   c.DefaultPrimaryColor,
		selectedColor: c.DefaultSecondaryColor,
		logger:        logger.NewUILogger(), // start with an empty logger so it can be enables selectively
	}
}

// AddRow to the list.
func (l *List) AddRow(text string) {
	l.ls.Rows = append(l.ls.Rows, text)
}

// SetLogger for the component.
func (l *List) SetLogger(log *logger.UILogger) {
	l.logger = log
}

// Widget returns the underlying termui widget.
func (l *List) Widget() *widgets.List {
	l.preRender()
	return l.ls
}

// Dimensions returns the current dimensions of the component.
func (l *List) Dimensions() Dimensions {
	// TODO: this could be done with the MIN MAX points on the pg, but this is simpler.
	return l.dimensions
}

// Render registers the object's state with the UI.
func (l *List) Render() {
	l.preRender()
	termui.Render(l.ls)
}

// preRender does all the work of translating the local component into the termui component before rendering.
func (l *List) preRender() {
	// selection color
	if l.selected {
		l.ls.BorderStyle.Fg = l.selectedColor
	} else {
		l.ls.BorderStyle.Fg = l.borderColor
	}
}

// Select marks the component as actively selected.
func (l *List) Select() {
	l.selected = true
}

// Deselect marks the component as unselected.
func (l *List) Deselect() {
	l.selected = false
}

// Next selects the next row in the list and returns the index of the selected row.
func (l *List) Next() int {
	if len(l.ls.Rows) == 0 {
		return -1
	}
	if l.ls.SelectedRow == len(l.ls.Rows)-1 {
		l.ls.SelectedRow = 0
	} else {
		l.ls.SelectedRow++
	}
	return l.ls.SelectedRow
}

// Previous selects the previous row in the list and returns the intext of the selected row.
func (l *List) Previous() int {
	if len(l.ls.Rows) == 0 {
		return -1
	}
	if l.ls.SelectedRow == 0 {
		l.ls.SelectedRow = len(l.ls.Rows) - 1
	} else {
		l.ls.SelectedRow--
	}
	return l.ls.SelectedRow
}

// Write to a list exists to satisfy interfaces with other types, but here it is a no op.
func (l *List) Write(character string) {}

// Flush all text in the component.
func (l *List) Flush() {
	l.ls.Rows = []string{}
}
