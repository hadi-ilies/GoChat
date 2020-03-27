package client

type dimension struct {
	x0 int
	y0 int
	x1 int
	y1 int
}

func newDimension(x0, y0, x1, y1 int) dimension {
	return dimension{x0: x0, y0: y0, x1: x1, y1: y1}
}
