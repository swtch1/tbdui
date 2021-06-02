package component

import (
	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

// Info is a read-only text space with no border.
type Info struct {
	pg   *widgets.Paragraph
	text string

	dimensions Dimensions

	textColor termui.Color
}

// NewInfo initializes an info bar.
func NewInfo(text string, d Dimensions) *Info {
	p := widgets.NewParagraph()
	p.Border = false
	p.SetRect(d.X1, d.Y1, d.X2, d.Y2)
	return &Info{
		pg:         p,
		text:       text,
		dimensions: d,
	}
}

// Widget returns the underlying termui widget.
func (i *Info) Widget() *widgets.Paragraph {
	i.preRender()
	return i.pg
}

// Dimensions returns the current dimensions of the component.
func (i *Info) Dimensions() Dimensions {
	// TODO: this could be done with the MIN MAX points on the pg, but this is simpler.
	return i.dimensions
}

// Render registers the object's state with the UI.
func (i *Info) Render() {
	i.preRender()
	termui.Render(i.pg)
}

// preRender does all the work of translating the local component into the termui component before rendering.
func (i *Info) preRender() {
	i.pg.Text = i.text
}

// Write a single input character to the component.
func (i *Info) Write(character string) {
	i.text += character
}

// Overwrite any existing text in the component.
func (i *Info) Overwrite(text string) {
	i.text = text
}

// Flush all text in the component.
func (i *Info) Flush() {
	i.text = ""
}

// Contents returns the contents of the component.  Empty box default text will never be returned.
func (i *Info) Contents() string {
	return i.text
}
