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

func Flush() {
	termbox.Flush()
}

func Line(b *edit.Buffer, index int) {
	w, _ := termbox.Size()
	y := index - b.TopLine

	line := b.Lines[index]

	tw := b.TabWidth
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

func Buffer(b *edit.Buffer) {
	_, h := termbox.Size()
	last := b.TopLine + h
	if last > len(b.Lines) {
		last = len(b.Lines)
	}

	for line := b.TopLine; line < last; line++ {
		Line(b, line)
	}
}
