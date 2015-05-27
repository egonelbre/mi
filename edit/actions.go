package edit

func Move(b *Buffer, dx, dy int) {
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

	for i := range b.Regions {
		r := &b.Regions[i]
		MoveCursor(&r.Start)
		r.End = r.Start
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
