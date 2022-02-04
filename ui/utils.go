package ui

type Handle struct {
	node Node
	id   int
	gen  uint
}

type Node interface {
	parent() Node
}

// Those are types used to keep the library
// backend independant. They are used to place the
// different elements, but mainly to provide a wrapper
// over the user native types for fonts and images
type (
	Font interface {
		MeasureText(text string) Point
	}

	Image interface {
		GetWidth() float64
		GetHeight() float64
	}
)

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

	StyleOrderingKind int
	Style             struct {
		Ordering StyleOrderingKind
		Padding  float64
		Margin   Point
	}

	Constraint struct {
		Left, Right float64
		Up, Down    float64
	}

	BackgroundKind int
	Background     struct {
		Visible bool
		Kind    BackgroundKind
		Clr     Color
		Img     Image
		Constr  Constraint
	}
)

const (
	StyleOrderingRow StyleOrderingKind = iota
	StyleOrderingColumn
)

const (
	BackgroundSolidColor BackgroundKind = iota
	BackgroundImageSlice
)

func (b Background) entry(rect Rectangle) RenderEntry {
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
		result.Constr = b.Constr
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
		Kind   RenderCommand
		Rect   Rectangle
		Clr    Color
		Img    Image
		Constr Constraint
		Text   string
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

func (r *RenderBuffer) flushBuffer() []RenderEntry {
	result := r.data[:r.count]
	r.count = 0
	return result
}
