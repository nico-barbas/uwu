package ui

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
	case StyleOrderRow:
		w.ptr = style.Margin[0]
	case StyleOrderColumn:
		w.ptr = style.Margin[1]
	}

}

func (w *WidgetList) addWidget(parent Node, pRect Rectangle, wgt Widget, len int) Handle {
	handle := Handle{node: wgt, id: w.count, gen: w.gens[w.count]}
	wgt.setHandle(handle)
	wgt.setParent(parent)

	rect := Rectangle{}
	switch w.style.Ordering {
	case StyleOrderRow:
		rect = Rectangle{
			X: pRect.X + w.style.Margin[0], Y: pRect.Y + w.ptr,
			Width: pRect.Width - (w.style.Margin[0] * 2), Height: float64(len),
		}
	case StyleOrderColumn:
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

func (w *WidgetList) drawWidgets(buf *renderBuffer) {
	for i := 0; i < w.count; i += 1 {
		w.widgets[i].draw(buf)
	}
}
