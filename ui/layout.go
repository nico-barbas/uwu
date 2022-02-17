package ui

type Layout struct {
	widgetRoot
	Background Background
	Style      Style
	widgets    WidgetList
}

func (l *Layout) init() {
	l.widgets.initList(l.Style)
}

func (l *Layout) update() {
	l.widgets.updateWidgets()
}

func (l *Layout) draw(buf *renderBuffer) {
	bgEntry := l.Background.entry(l.rect)
	buf.addEntry(bgEntry)

	l.widgets.drawWidgets(buf)
}

func (l *Layout) AddWidget(wgt Widget, length int) {
	l.widgets.addWidget(wgt, l.rect, length)
}

func (l *Layout) RemainingLength() int {
	return l.widgets.getRemainingLen(l.rect)
}
