package ui

type Widget interface {
	Node
	setHandle(h Handle)
	setParent(n Node)
	setRect(r Rectangle)
	init()
	update()
	draw(buf *RenderBuffer)
}

//
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

func (w *widgetRoot) init()                  {}
func (w *widgetRoot) update()                {}
func (w *widgetRoot) draw(buf *RenderBuffer) {}

type DebugWidget struct {
	widgetRoot
}

func (d *DebugWidget) draw(buf *RenderBuffer) {
	buf.addEntry(RenderEntry{
		Kind: RenderRectangle,
		Rect: d.rect,
		Clr:  Color{255, 0, 0, 255},
	})
}
