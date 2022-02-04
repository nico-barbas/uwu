package ui

type Window struct {
	handle     Handle
	Active     bool
	Rect       Rectangle
	Background Background
	widgets    WidgetList
}

func (win *Window) parent() Node {
	return nil
}

func (win *Window) initWindow(style Style) {
	win.widgets.initList(style)
}

func (win *Window) draw(buf *RenderBuffer) {
	bgEntry := win.Background.entry(win.Rect)
	buf.addEntry(bgEntry)

	win.widgets.drawWidgets(buf)
}
