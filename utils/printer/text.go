package printer

import (
	"github.com/charmbracelet/lipgloss"
	"regexp"
	"sort"
)

type TextX struct {
	content  string
	keywords []string
	// https://github.com/charmbracelet/lipgloss
	style lipgloss.Style
}

func NewText(content string, keywords ...string) *TextX {
	return &TextX{
		content:  content,
		keywords: keywords,
		style:    lipgloss.NewStyle(),
	}
}

func (t *TextX) SetStyle(style lipgloss.Style) *TextX {
	t.style = style
	return t
}

func (t *TextX) String() string {
	if len(t.keywords) == 0 {
		return t.style.Render(t.content)
	}
	type Match struct {
		Start   int
		End     int
		Keyword string
	}

	var matches []Match
	for _, keyword := range t.keywords {
		re := regexp.MustCompile(regexp.QuoteMeta(keyword))
		locs := re.FindAllStringIndex(t.content, -1)
		for _, loc := range locs {
			matches = append(matches, Match{
				Start:   loc[0],
				End:     loc[1],
				Keyword: keyword,
			})
		}
	}

	// sort by matching position from front to back.
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].Start < matches[j].Start
	})

	// replace from the back to the front
	// use rune to ensure that replacement does not break string indices.
	result := []rune(t.content)
	for i := len(matches) - 1; i >= 0; i-- {
		m := matches[i]
		formatted := t.style.Render(m.Keyword)
		// replace the matching part
		result = append(result[:m.Start], append([]rune(formatted), result[m.End:]...)...)
	}
	return string(result)
}
