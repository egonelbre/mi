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

func (b *Buffer) UpdateRegions() {
	sort.Sort(byPosition(b.Regions))
	if len(b.Regions) == 0 {
		b.Regions = append(b.Regions, Region{})
	}
}

func (b *Buffer) Move(dx, dy int) {
	for i := range b.Regions {
		r := &b.Regions[i]
		r.Start.Move(dx, dy)
		r.End = r.Start
	}
	b.UpdateRegions()
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

func (r *Region) Contains(c Cursor) bool {
	if c.Before(r.Start) || c.After(r.End) {
		return false
	}
	return true
}

type Cursor struct{ Line, Column int }

func (a *Cursor) Move(dx, dy int) {
	a.Line += dy
	a.Column += dx
	if a.Column < 0 {
		a.Line--
		a.Column = 0
	}
	if a.Line < 0 {
		a.Line = 0
	}
}

func (a Cursor) Before(b Cursor) bool {
	if a.Line == b.Line {
		return a.Column < b.Column
	}
	return a.Line < b.Line
}

func (a Cursor) After(b Cursor) bool {
	if a.Line == b.Line {
		return a.Column > b.Column
	}
	return a.Line > b.Line
}

type byPosition []Region

func (r byPosition) Len() int      { return len(r) }
func (r byPosition) Swap(i, j int) { r[i], r[j] = r[j], r[i] }
func (r byPosition) Less(i, j int) bool {
	a, b := &r[i], &r[j]
	if a.Start == b.Start {
		return a.End.Before(b.End)
	}
	return a.Start.Before(b.Start)
}
