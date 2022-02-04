package ui

type Window struct {
	handle     Handle
	Active     bool
	Rect       Rectangle
	Style      Style
	Background Background
	widgets    WidgetList

	HasHeader        bool
	HeaderHeight     float64
	HeaderBackground Background
}

func (win *Window) parent() Node {
	return nil
}

func (win *Window) initWindow() {
	win.widgets.initList(win.Style)
}

func (win *Window) draw(buf *renderBuffer) {
	bgEntry := win.Background.entry(win.Rect)
	buf.addEntry(bgEntry)

	win.widgets.drawWidgets(buf)
}
