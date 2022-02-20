package ui

import (
	"log"
	"sync/atomic"
)

type Container interface {
	AddWidget(w Widget, length int)
	// RemoveWidget(w Widget)
	RemainingLength() int
}

// Those are types used to keep the library
// backend independant. They are used to place the
// different elements, but mainly to provide a wrapper
// over the user native types for fonts and images
type (
	Font interface {
		GlyphAdvance(r rune, size float64) float64
		MeasureText(text string, size float64) Point
	}

	Image interface {
		GetWidth() float64
		GetHeight() float64
	}

	Canvas interface {
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
		// Shared fields
		Visible bool
		Kind    BackgroundKind
		Clr     Color

		// ImageSlice fields
		Img    Image
		Constr Constraint
	}
)

func (c Color) RGBA() (r, g, b, a uint32) {
	r = uint32(c[0])
	r |= r << 8
	g = uint32(c[1])
	g |= g << 8
	b = uint32(c[2])
	b |= b << 8
	a = uint32(c[3])
	a |= a << 8
	return
}

func (r Rectangle) pointInBounds(p Point) bool {
	return (p[0] >= r.X && p[0] <= r.X+r.Width) && (p[1] >= r.Y && p[1] <= r.Y+r.Height)
}

const FitContainer = -1

const (
	StyleOrderRow StyleOrderingKind = iota
	StyleOrderColumn
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
// Rendering types and utility
//

const (
	RenderRectangle RenderCommand = iota
	RenderImage
	RenderImageSlice
	RenderText
)

type (
	RenderCommand int

	// A type describing a draw command.
	//
	// It is a "fat" struct since Go doesn't have
	// a sum type and the size of the struct is still
	// relatively small.
	//
	// <Kind> is the discriminator field
	RenderEntry struct {
		// The discriminator field
		Kind RenderCommand
		// Again no unions, so Rect is used for multiple
		// purposes depending on the Kind:
		// - For images, it is the destination rectangle (where to render the image)
		// - For texts, it is both the position, and the font size (pos at .X, .Y and font size at .Height)
		Rect   Rectangle
		Clr    Color
		Img    Image
		Constr Constraint
		Font   Font
		Text   string
	}

	// A buffer of RenderEntries
	renderBuffer struct {
		data  []RenderEntry
		cap   int
		count int
	}
)

func newRenderBuffer(cap int) renderBuffer {
	return renderBuffer{
		data: make([]RenderEntry, cap),
		cap:  cap,
	}
}

func (r *renderBuffer) addEntry(e RenderEntry) {
	r.data[r.count] = e
	r.count += 1
}

func (r *renderBuffer) flushBuffer() []RenderEntry {
	result := r.data[:r.count]
	r.count = 0
	return result
}

//
// Input types and utility
//

const charPressedCap = 50

type key int

const (
	keyEsc key = iota
	keyEnter
	keyDelete
	keyCtlr
	keyShift
	keySpace
	keyTab
	keyUp
	keyDown
	keyLeft
	keyRight
	keyMax
)

type (
	inputData struct {
		mPos              Point
		mLeft             bool
		previousmPos      Point
		previousmLeft     bool
		pressedChars      [charPressedCap]rune
		pressedCharsCount int32

		previousKeys [keyMax]bool
		keys         [keyMax]bool
		keyCounts    [keyMax]int
	}

	Input struct {
		MPos  Point
		MLeft bool

		// all the mods and keys the UI cares about
		Esc   bool
		Enter bool
		Del   bool
		Ctrl  bool
		Shift bool
		Space bool
		Tab   bool

		Up    bool
		Down  bool
		Left  bool
		Right bool
	}

	CursorShape int
)

const (
	CursorShapeDefault CursorShape = iota
	CursorShapeText
)

func mousePosition() Point {
	return ctx.input.mPos
}

func isMousePressed() bool {
	return ctx.input.mLeft
}

func isMouseJustPressed() bool {
	return ctx.input.mLeft && (ctx.input.mLeft != ctx.input.previousmLeft)
}

func isMouseJustReleased() bool {
	return !ctx.input.mLeft && (ctx.input.mLeft != ctx.input.previousmLeft)
}

func pressedChars() []rune {
	return ctx.input.pressedChars[:ctx.input.pressedCharsCount]
}

func isKeyPressed(k key) bool {
	return ctx.input.keys[k]
}

func isAnyKeyPressed(keys []key) bool {
	for _, k := range keys {
		if isKeyPressed(k) {
			return true
		}
	}
	return false
}

func isKeyRepeated(k key) bool {
	const (
		delay    = 15
		interval = 1
	)
	d := ctx.input.keyCounts[k]
	if d == 1 {
		return true
	}
	if d >= delay && (d-delay)%interval == 0 {
		return true
	}
	return false
}

func setCursorShape(s CursorShape) {
	if ctx.cursorShapeCallback != nil {
		ctx.cursorShapeCallback(s)
	} else {
		log.SetPrefix("[UI Error]: ")
		log.Fatalln("No Cursor shape callback was provided")
	}
}

// FIXME: Make this thread-safe
func (c *Context) AppendCharPressed(r rune) {
	index := atomic.SwapInt32(&c.input.pressedCharsCount, c.input.pressedCharsCount+1)
	c.input.pressedChars[index] = r
}
