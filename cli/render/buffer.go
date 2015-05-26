package render

import (
	"fmt"
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

type Regions struct {
	Index   int
	Regions []edit.Region
}

func (r *Regions) Contains(line, column int) bool {	
	//TODO: optimize
	c := edit.Cursor{line, column}
	for i := range r.Regions {
		if r.Regions[i].Contains(c) {
			return true
		}
	}
	return false
}

func normal() (fg, bg termbox.Attribute) { return termbox.ColorDefault, termbox.ColorDefault }
func highlight() (fg, bg termbox.Attribute) { return termbox.ColorBlack, termbox.ColorWhite }

func Line(b *edit.Buffer, index int, r *Regions) {
	w, _ := termbox.Size()
	y := index - b.TopLine

	line := b.Lines[index]
	
	tw := b.TabWidth
	x, p := 0, 0
	for p < len(line) && x < w {
		fg, bg := normal()
		if r.Contains(index, x) {
			fg, bg = highlight()
		}
			
		c, cw := utf8.DecodeRune([]byte(line[p:]))
		p += cw
		switch c {
		case '\t':
			if tw <= 0 {
				continue
			}
			for {
				fg, bg := normal()
				if r.Contains(index, x) {
					fg, bg = highlight()
				}
				termbox.SetCell(x, y, ' ', fg, bg)
				x++
				if x%tw == 0 {
					break
				}
			}
		case '\n':
			termbox.SetCell(x, y, NewLine, fg, bg)
			x++
		case '\r':
			termbox.SetCell(x, y, LineFeed, fg, bg)
			x++
		default:
			termbox.SetCell(x, y, c, fg, bg)
			x++
		}
	}
	
	for x < w {
		fg, bg := normal()
		if r.Contains(index, x){
			fg, bg = highlight()
		}
		termbox.SetCell(x, y, ' ', fg, bg)
		x++
	}
}

func Debug(y int, msg string){
	fg, bg := highlight()
	
	w, _ := termbox.Size()
	x := 0
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
	for x < w {
		termbox.SetCell(x, y, ' ', fg, bg)
		x++
	}
}

func Buffer(b *edit.Buffer) {
	_, h := termbox.Size()
	last := b.TopLine + h - 1
	if last > len(b.Lines) {
		last = len(b.Lines)
	}
	
	regions := &Regions{0, b.Regions}
	for line := b.TopLine; line < last; line++ {
		Line(b, line, regions)
	}
	
	Debug(h-1, fmt.Sprintf("%+v", b.Regions))
}
