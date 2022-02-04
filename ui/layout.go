package ui

type Layout struct {
	WidgetRoot
	Background Background
	Style      Style
	widgets    WidgetList
}

func (l *Layout) init() {
	l.widgets.initList(l.Style)
}

func (l *Layout) draw(buf *RenderBuffer) {
	bgEntry := l.Background.entry(l.rect)
	buf.addEntry(bgEntry)

	l.widgets.drawWidgets(buf)
}
