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

const widgetListCap = 10

type WidgetList struct {
	style   Style
	widgets [widgetListCap]Widget
	gens    [widgetListCap]uint
	count   int
	ptr     float64
}

func (w *WidgetList) initList(style Style) {
	w.style = style
	switch w.style.Ordering {
	case StyleOrderingRow:
		w.ptr = style.Margin[0]
	case StyleOrderingColumn:
		w.ptr = style.Margin[1]
	}

}

func (w *WidgetList) addWidget(parent Node, pRect Rectangle, wgt Widget, len int) Handle {
	handle := Handle{node: wgt, id: w.count, gen: w.gens[w.count]}
	wgt.setHandle(handle)
	wgt.setParent(parent)

	rect := Rectangle{}
	switch w.style.Ordering {
	case StyleOrderingRow:
		rect = Rectangle{
			X: pRect.X + w.style.Margin[0], Y: pRect.Y + w.ptr,
			Width: pRect.Width - (w.style.Margin[0] * 2), Height: float64(len),
		}
	case StyleOrderingColumn:
		rect = Rectangle{
			X: pRect.X + w.ptr, Y: pRect.Y + w.style.Margin[1],
			Width: float64(len), Height: pRect.Height - (w.style.Margin[1] * 2),
		}
	}
	wgt.setRect(rect)
	w.widgets[w.count] = wgt
	w.gens[w.count] += 1
	w.count += 1
	w.ptr += float64(len) + w.style.Padding
	w.widgets[w.count-1].init()
	return handle
}

func (w *WidgetList) drawWidgets(buf *RenderBuffer) {
	for i := 0; i < w.count; i += 1 {
		w.widgets[i].draw(buf)
	}
}

type WidgetRoot struct {
	handle  Handle
	wParent Node
	rect    Rectangle
}

func (w *WidgetRoot) parent() Node {
	return w.wParent
}

func (w *WidgetRoot) setHandle(h Handle) {
	w.handle = h
}

func (w *WidgetRoot) setParent(n Node) {
	w.wParent = n
}

func (w *WidgetRoot) setRect(r Rectangle) {
	w.rect = r
}

func (w *WidgetRoot) init()                  {}
func (w *WidgetRoot) update()                {}
func (w *WidgetRoot) draw(buf *RenderBuffer) {}

type DebugWidget struct {
	WidgetRoot
}

func (d *DebugWidget) draw(buf *RenderBuffer) {
	buf.addEntry(RenderEntry{
		Kind: RenderRectangle,
		Rect: d.rect,
		Clr:  Color{255, 0, 0, 255},
	})
}
