package ui

type (
	TabViewer struct {
		widgetRoot

		Background Background

		tabs           []tab
		HeaderHeight   int
		HeaderFont     Font
		HeaderTextSize float64
	}

	tab struct {
		name   string
		widget Widget
	}
)

func (t *TabViewer) draw(buf *renderBuffer) {
	bgEntry := t.Background.entry(t.rect)
	buf.addEntry(bgEntry)

	// var xptr float64
}
