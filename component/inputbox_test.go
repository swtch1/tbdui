package component

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/swtch1/tbdui/char"
	"github.com/swtch1/tbdui/conf"
)

func TestWritingToInputBox(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		initialText string
		blankText   string
		toWrite     string
		expected    string
	}{
		{
			name:        "blank box with default text",
			initialText: "start here",
			blankText:   "start here",
			toWrite:     "a",
			expected:    "a",
		},
		{
			name:        "empty box",
			initialText: "",
			blankText:   "",
			toWrite:     "a",
			expected:    "a",
		},
		{
			name:        "backspace with existing text",
			initialText: "foo",
			blankText:   "",
			toWrite:     char.BACKSPACE,
			expected:    "fo",
		},
		{
			name:        "backspace with no text",
			initialText: "",
			blankText:   "",
			toWrite:     char.BACKSPACE,
			expected:    "",
		},
		{
			name:        "backspace with no text resets to blank",
			initialText: "",
			blankText:   "blank",
			toWrite:     char.BACKSPACE,
			expected:    "blank",
		},
		{
			name:        "backspace with one char resets to blank",
			initialText: "a",
			blankText:   "blank",
			toWrite:     char.BACKSPACE,
			expected:    "blank",
		},
		{
			name:        "first character blanks the box",
			initialText: ":start here",
			blankText:   ":start here",
			toWrite:     "a",
			expected:    "a",
		},
		{
			name:        "add a space",
			initialText: "hello",
			blankText:   ":start here",
			toWrite:     char.SPACE,
			expected:    "hello ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewInputBox("title", "", conf.Config{}, Dimensions{})
			b.text = tt.initialText
			b.blankText = tt.blankText
			b.Write(tt.toWrite)
			require.Equal(t, tt.expected, b.text)
		})
	}
}
