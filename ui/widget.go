package ui

type Widget interface {
	Node
	setHandle(h Handle)
	setParent(n Node)
	setRect(r Rectangle)
	moveBy(offset Point)
	init()
	update()
	draw(buf *renderBuffer)
}

// A type to embed into every other widget types implementation
// to avoid having to redeclare all the boilerplate fields and
// methods
type widgetRoot struct {
	handle  Handle
	wParent Node
	rect    Rectangle
}

func (w *widgetRoot) parent() Node {
	return w.wParent
}

func (w *widgetRoot) setHandle(h Handle) {
	w.handle = h
}

func (w *widgetRoot) setParent(n Node) {
	w.wParent = n
}

func (w *widgetRoot) setRect(r Rectangle) {
	w.rect = r
}

func (w *widgetRoot) moveBy(offset Point) {
	w.rect.X += offset[0]
	w.rect.Y += offset[1]
}

// "Virtual" methods to avoid having to redeclare it
// for every widget implementation
func (w *widgetRoot) init()                  {}
func (w *widgetRoot) update()                {}
func (w *widgetRoot) draw(buf *renderBuffer) {}

//
// Simple Widget used for debugging purposes
type DebugWidget struct {
	widgetRoot
}

func (d *DebugWidget) draw(buf *renderBuffer) {
	buf.addEntry(RenderEntry{
		Kind: RenderRectangle,
		Rect: d.rect,
		Clr:  Color{255, 0, 0, 255},
	})
}
