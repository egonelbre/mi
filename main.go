package main

import (
	"io/ioutil"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/nsf/termbox-go"
)

const (
	RuneNewLine  = ' ' // '↓'
	RuneLineFeed = '←'
)

type Buffer struct {
	Lines   []Line
	TopLine int
	Regions []Region

	TabWidth int
}

func NewBuffer() *Buffer {
	return &Buffer{
		TabWidth: 4,
	}
}

func (b *Buffer) RegionsChanged() {
	sort.Sort(byPosition(b.Regions))
	//TODO: merge regions
}

func BufferFromFile(filename string) (*Buffer, error) {
	buf := NewBuffer()

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return buf, nil
	}

	lines := strings.SplitAfter(string(data), "\n")
	for _, line := range lines {
		buf.Lines = append(buf.Lines, NewLine(line))
	}

	return buf, nil
}

type Line string

func NewLine(text string) Line { return Line(text) }

type Region struct{ Start, End Cursor }

type Cursor struct{ Line, Column int }

func (a Cursor) Less(b Cursor) bool {
	if a.Line == b.Line {
		return a.Column < b.Column
	}
	return a.Line < b.Line
}

type byPosition []Region

func (r byPosition) Len() int      { return len(r) }
func (r byPosition) Swap(i, j int) { r[i], r[j] = r[j], r[i] }
func (r byPosition) Less(i, j int) bool {
	a, b := &r[i], &r[j]
	if a.Start == b.Start {
		return a.End.Less(b.End)
	}
	return a.Start.Less(b.Start)
}

func Render(b *Buffer) {
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
				termbox.SetCell(x, y, RuneNewLine, termbox.ColorDefault, termbox.ColorDefault)
				x++
			case '\r':
				termbox.SetCell(x, y, RuneLineFeed, termbox.ColorDefault, termbox.ColorDefault)
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

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	buffer, err := BufferFromFile("main.go")
	if err != nil {
		panic(err)
	}

	Render(buffer)
	termbox.Flush()

inputloop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				break inputloop
			}
			switch ev.Ch {
			case 'q', 'Q':
				break inputloop
			}
		case termbox.EventError:
			panic(ev.Err)
		}

		Render(buffer)
		termbox.Flush()
	}
}
