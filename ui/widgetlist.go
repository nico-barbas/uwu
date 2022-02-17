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
		w.ptr += style.Margin[1]
	case StyleOrderColumn:
		w.ptr += style.Margin[0]
	}
}

func (w *WidgetList) addWidget(wgt Widget, pRect Rectangle, l int) {
	len := float64(l)

	rect := Rectangle{}
	switch w.style.Ordering {
	case StyleOrderRow:
		if l == FitContainer {
			len = pRect.Height - w.ptr - w.style.Margin[1]
		}
		rect = Rectangle{
			X: pRect.X + w.style.Margin[0], Y: pRect.Y + w.ptr,
			Width: pRect.Width - (w.style.Margin[0] * 2), Height: len,
		}
	case StyleOrderColumn:
		if l == FitContainer {
			len = pRect.Width - w.ptr - w.style.Margin[0]
		}
		rect = Rectangle{
			X: pRect.X + w.ptr, Y: pRect.Y + w.style.Margin[1],
			Width: len, Height: pRect.Height - (w.style.Margin[1] * 2),
		}
	}
	wgt.setRect(rect)
	w.widgets[w.count] = wgt
	w.gens[w.count] += 1
	w.count += 1
	w.ptr += len + w.style.Padding
	w.widgets[w.count-1].init()
}

func (w *WidgetList) updateWidgets() {
	for i := 0; i < w.count; i += 1 {
		w.widgets[i].update()
	}
}

func (w *WidgetList) drawWidgets(buf *renderBuffer) {
	for i := 0; i < w.count; i += 1 {
		w.widgets[i].draw(buf)
	}
}

func (w *WidgetList) getRemainingLen(pRect Rectangle) int {
	switch w.style.Ordering {
	case StyleOrderRow:
		return int(pRect.Height - w.ptr)
	case StyleOrderColumn:
		return int(pRect.Width - w.ptr)
	}
	return -1
}
