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
	// fix regions if needed
	for i := range b.Regions {
		r := &b.Regions[i]
		if r.End.Before(r.Start) {
			r.End = r.Start
		}
	}

	// ensure that regions are sorted by their position
	sort.Sort(byPosition(b.Regions))

	// create a region if there is none
	if len(b.Regions) == 0 {
		b.Regions = append(b.Regions, Region{})
	}

	regions := []Region{}
	last := b.Regions[0]
	for _, r := range b.Regions[1:] {
		if last.Overlaps(&r) {
			last.Merge(&r)
		} else {
			regions = append(regions, last)
			last = r
		}
	}
	regions = append(regions, last)
	b.Regions = regions
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

func (a *Region) Overlaps(b *Region) bool {
	return a.Contains(b.Start) || a.Contains(b.End) ||
		b.Contains(a.Start) || b.Contains(a.Start)
}

func (a *Region) Merge(b *Region) {
	if b.Start.Before(a.Start) {
		a.Start = b.Start
	}
	if b.End.After(a.End) {
		a.End = b.End
	}
}

type Cursor struct{ Line, Column int }

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

func (a *Cursor) Offset(dx, dy int) {
	a.Line += dy
	a.Column += dx
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
