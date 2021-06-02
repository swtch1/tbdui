package component

import (
	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/swtch1/tbdui/char"
	"github.com/swtch1/tbdui/conf"
	"github.com/swtch1/tbdui/logger"
)

// InputBox records and displays input from a user.
type InputBox struct {
	pg        *widgets.Paragraph
	blankText string
	text      string

	selected bool
	// HideUnselectedText ensures no text is shown when the component is unselected. True by default.
	HideUnselectedText bool
	// AllowWrite from users.  True by default.
	AllowWrite bool
	dimensions Dimensions

	borderColor   termui.Color
	selectedColor termui.Color

	logger *logger.UILogger
}

// NewInputBox initializes an input box.
func NewInputBox(title, defaultText string, c conf.Config, d Dimensions) *InputBox {
	p := widgets.NewParagraph()
	p.Title = title
	p.SetRect(d.X1, d.Y1, d.X2, d.Y2)
	return &InputBox{
		pg:                 p,
		blankText:          defaultText,
		text:               defaultText,
		HideUnselectedText: true,
		AllowWrite:         true,
		dimensions:         d,
		borderColor:        c.DefaultPrimaryColor,
		selectedColor:      c.DefaultSecondaryColor,
		logger:             logger.NewUILogger(), // start with an empty logger so it can be enables selectively
	}
}

// SetLogger for the component.
func (b *InputBox) SetLogger(l *logger.UILogger) {
	b.logger = l
}

// Widget returns the underlying termui widget.
func (b *InputBox) Widget() *widgets.Paragraph {
	b.preRender()
	return b.pg
}

// Dimensions returns the current dimensions of the component.
func (b *InputBox) Dimensions() Dimensions {
	// TODO: this could be done with the MIN MAX points on the pg, but this is simpler.
	return b.dimensions
}

// Render registers the object's state with the UI.
func (b *InputBox) Render() {
	b.preRender()
	termui.Render(b.pg)
}

// preRender does all the work of translating the local component into the termui component before rendering.
func (b *InputBox) preRender() {
	// selection color
	if b.selected {
		b.pg.BorderStyle.Fg = b.selectedColor
	} else {
		b.pg.BorderStyle.Fg = b.borderColor
	}

	// optionally reset the text
	if len(b.text) == 0 {
		b.Flush()
	}

	// TODO: I know this is hacky.. find a better way.  But I am le tired.
	if b.HideUnselectedText && !b.selected && b.text == b.blankText {
		b.pg.Text = ""
		return
	}
	b.pg.Text = b.text
}

// Select marks the component as actively selected.
func (b *InputBox) Select() {
	b.selected = true
}

// Deselect marks the component as unselected.
func (b *InputBox) Deselect() {
	b.selected = false
}

// Write a single input character to the component.
func (b *InputBox) Write(character string) {
	if !b.AllowWrite {
		return
	}

	if b.text == b.blankText {
		b.text = ""
	}

	switch character {
	case char.SPACE:
		b.text += " "
	case char.BACKSPACE:
		if len(b.text) <= 1 {
			b.Flush()
			return
		}
		b.text = b.text[0 : len(b.text)-1]
	default:
		b.text += character

	// no ops below
	case char.TAB, char.ENTER:
	case char.HOME, char.END, char.ESCAPE:
	case char.UP, char.DOWN, char.LEFT, char.RIGHT:
	case char.NEXT, char.PREVIOUS:
	}
}

// Overwrite any existing text in the component.
func (b *InputBox) Overwrite(text string) {
	b.text = text
}

// Flush all text in the component.
func (b *InputBox) Flush() {
	b.text = b.blankText
}

// Contents returns the contents of the component.  Empty box default text will never be returned.
func (b *InputBox) Contents() string {
	if b.text == b.blankText {
		return ""
	}
	return b.text
}
