package edit

import (
	"io/ioutil"
	"sort"
	"strings"
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
