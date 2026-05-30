package internal

import (
	"strings"
	"unicode/utf8"
)

// Wrap wraps s into lines of at most width runes.
func Wrap(s string, width int) string {
	if width <= 0 {
		return s
	}
	var b strings.Builder
	b.Grow(len(s))

	for i, para := range strings.Split(s, "\n") {
		if i > 0 {
			b.WriteByte('\n')
		}
		wrapPara(&b, para, width)
	}
	return b.String()
}

func wrapPara(b *strings.Builder, text string, width int) {
	col := 0
	for i, word := range strings.Fields(text) {
		wlen := utf8.RuneCountInString(word)

		if i > 0 {
			if col+1+wlen <= width {
				b.WriteByte(' ')
				col++
			} else {
				b.WriteByte('\n')
				col = 0
			}
		}

		if col+wlen <= width {
			b.WriteString(word)
			col += wlen
			continue
		}

		col = hyphenate(b, []rune(word), col, width)
	}
}

func hyphenate(b *strings.Builder, runes []rune, col, width int) int {
	for len(runes) > 0 {
		avail := width - col

		if len(runes) <= avail {
			b.WriteString(string(runes))
			return col + len(runes)
		}

		if avail < 2 {
			if col > 0 {
				b.WriteByte('\n')
				col = 0
			} else {
				b.WriteRune(runes[0])
				runes = runes[1:]
				if len(runes) > 0 {
					b.WriteByte('\n')
				}
			}
			continue
		}

		for _, r := range runes[:avail-1] {
			b.WriteRune(r)
		}
		b.WriteString("-\n")
		runes = runes[avail-1:]
		col = 0
	}
	return col
}
