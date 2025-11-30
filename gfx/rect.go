package gfx

type Edges struct{ Top, Right, Bottom, Left float64 }
type Rect struct{ Left, Top, Width, Height float64 }

func (r Rect) Right() float64  { return r.Left + r.Width - 1 }
func (r Rect) Bottom() float64 { return r.Top + r.Height - 1 }
func (r Rect) AddPadding(edges Edges) Rect {
	r.Top += edges.Top
	r.Left += edges.Left
	r.Width -= edges.Left + edges.Right
	r.Height -= edges.Top + edges.Bottom
	return r
}
