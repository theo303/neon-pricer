package svg

type Bounds struct {
	minX, maxX float64
	minY, maxY float64
}

func (b Bounds) Width() float64 {
	return b.maxX - b.minX
}

func (b Bounds) Height() float64 {
	return b.maxY - b.minY
}

func (b Bounds) Expand(nb Bounds) Bounds {
	return Bounds{
		minX: min(b.minX, nb.minX),
		maxX: max(b.maxX, nb.maxX),
		minY: min(b.minY, nb.minY),
		maxY: max(b.maxY, nb.maxY),
	}
}

func (b Bounds) expandPoint(p point) Bounds {
	return Bounds{
		minX: min(b.minX, p.x),
		maxX: max(b.maxX, p.x),
		minY: min(b.minY, p.y),
		maxY: max(b.maxY, p.y),
	}
}
