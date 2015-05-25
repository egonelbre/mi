package main

import (
	"io/ioutil"
	"strings"
	"sort"

	"github.com/nsf/termbox-go"
)

type Buffer struct {
	Lines   []Line
	TopLine int
	Regions []Region
}

func (b *Buffer) RegionsChanged() {
	sort.Sort(byPosition(b.Regions))
	//TODO: merge regions
}

func BufferFromFile(filename string) (*Buffer, error) {
	buf := &Buffer{}

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

func (r byPosition) Len() int           { return len(r) }
func (r byPosition) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r byPosition) Less(i, j int) bool {
	a, b := &r[i], &r[j]
	if a.Start == b.Start {
		return a.End.Less(b.End)
	}
	return a.Start.Less(b.Start)
}


func Render(b *Buffer) {
	_, h := termbox.Size()
	last := b.TopLine + h
	if last > len(b.Lines) {
		last = len(b.Lines)
	}
	for y, line := range b.Lines[b.TopLine:last] {
		for x, c := range line {
			termbox.SetCell(x, y, c, termbox.ColorDefault, termbox.ColorDefault)
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
		case termbox.EventError:
			panic(ev.Err)
		}

		Render(buffer)
		termbox.Flush()
	}
}
