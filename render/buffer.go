package render

import (
	"unicode/utf8"

	"github.com/loov/mi/edit"
	"github.com/nsf/termbox-go"
)

var (
	NewLine  = ' ' // '↓'
	LineFeed = '←'
)

func Buffer(b *edit.Buffer) {
	w, h := termbox.Size()
	last := b.TopLine + h
	if last > len(b.Lines) {
		last = len(b.Lines)
	}

	tw := b.TabWidth
	for y, line := range b.Lines[b.TopLine:last] {
		x, p := 0, 0
		for p < len(line) && x < w {
			c, cw := utf8.DecodeRune([]byte(line[p:]))
			p += cw
			switch c {
			case '\t':
				if tw <= 0 {
					continue
				}
				for {
					termbox.SetCell(x, y, ' ', termbox.ColorDefault, termbox.ColorDefault)
					x++
					if x%tw == 0 {
						break
					}
				}
			case '\n':
				termbox.SetCell(x, y, NewLine, termbox.ColorDefault, termbox.ColorDefault)
				x++
			case '\r':
				termbox.SetCell(x, y, LineFeed, termbox.ColorDefault, termbox.ColorDefault)
				x++
			default:
				termbox.SetCell(x, y, c, termbox.ColorDefault, termbox.ColorDefault)
				x++
			}
		}
		for x < w {
			termbox.SetCell(x, y, ' ', termbox.ColorDefault, termbox.ColorDefault)
			x++
		}
	}
}
