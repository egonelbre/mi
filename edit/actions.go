package edit

func MoveRegion(b *Buffer, r *Region, dx, dy int) {
	MoveCursor := func(cursor *Cursor) {
		cursor.Line += dy
		cursor.Column += dx
		if cursor.Column < 0 {
			cursor.Line--
			cursor.Column = 0
		}
		if cursor.Line < 0 {
			cursor.Line = 0
		}
	}
	
	MoveCursor(&r.Start)
	r.End = r.Start
}

func Move(b *Buffer, dx, dy int) {
	for i := range b.Regions {
		MoveRegion(b, &b.Regions[i], dx, dy)
	}
	b.RegionsChanged()
}

func Type(b *Buffer, text string) {
	TypeInRegion := func(r *Region, ri int) {
		line := b.Lines[r.End.Line]
		line = line[:r.Start.Column] + Line(text) + line[r.End.Column:]
		b.Lines[r.End.Line] = line

		dy, dx := 0, len(text)
		for i := ri; i < len(b.Regions); i++ {
			x := &b.Regions[i]
			x.Start.Offset(dx, dy)
			x.End.Offset(dx, dy)
		}
	}

	for ri := range b.Regions {
		TypeInRegion(&b.Regions[ri], ri)
	}

	b.RegionsChanged()
}

func AddMoveRegion(b *Buffer, dy int) {
	regions := []Region{}
	for _, r := range b.Regions {
		MoveRegion(b, &r, 0, dy)
		regions = append(regions, r)	
	}
	
	b.Regions = append(b.Regions, regions...)
	b.RegionsChanged()	
}
