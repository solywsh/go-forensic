package printer

import (
	"bytes"
	"strings"
)

// WrapText
// Split the input text into lines according to the specified width.
func WrapText(text string, width int) string {
	var buffer bytes.Buffer
	words := strings.Fields(text) // Separate words by spaces
	line := ""

	for _, word := range words {
		// If the length of the current line plus the length of
		// the new word exceeds the specified width, a line break is made.
		if len(line)+len(word)+1 > width {
			buffer.WriteString(line + "\n")
			line = word
		} else {
			if line != "" {
				line += " "
			}
			line += word
		}
	}

	// Finally, add the last collected line to the buffer.
	if line != "" {
		buffer.WriteString(line)
	}

	return buffer.String()
}
