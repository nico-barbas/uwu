package ui

type Handle struct {
	parent Node
	id     int
	gen    int
}

type Node interface {
	parent() Node
	child() Node
}

// All the utility types used for simple
// calculation and data transformations
type (
	// Represents a point in the world
	Point [2]float64

	Color [4]uint8

	Rectangle struct {
		X, Y          float64
		Width, Height float64
	}

	Font interface {
		MeasureText(text string) Point
	}

	Image interface {
		Width() float64
		Height() float64
	}

	Constraint struct {
		Left, Right float64
		Up, Bottom  float64
	}

	BackgroundKind int

	Background struct {
		Visible bool
		Kind    BackgroundKind
		Img     Image
		Clr     Color
	}
)

const (
	BackgroundSolidColor BackgroundKind = iota
	BackgroundImageSlice
)

func (b Background) backgroundEntry(rect Rectangle) RenderEntry {
	result := RenderEntry{
		Rect: rect,
		Clr:  b.Clr,
	}
	switch b.Kind {
	case BackgroundSolidColor:
		result.Kind = RenderRectangle
	case BackgroundImageSlice:
		result.Kind = RenderImageSlice
		result.Img = b.Img
	}
	return result
}

//

const (
	RenderRectangle RenderCommand = iota
	RenderImage
	RenderImageSlice
	RenderText
)

type (
	RenderCommand int

	RenderEntry struct {
		Kind  RenderCommand
		Rect  Rectangle
		Clr   Color
		Img   Image
		Const Constraint
		Text  string
	}

	RenderBuffer struct {
		data  []RenderEntry
		cap   int
		count int
	}
)

func newRenderBuffer(cap int) RenderBuffer {
	return RenderBuffer{
		data: make([]RenderEntry, cap),
		cap:  cap,
	}
}

func (r *RenderBuffer) addEntry(e RenderEntry) {
	r.data[r.count] = e
	r.count += 1
}
