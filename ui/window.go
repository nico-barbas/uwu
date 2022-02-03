package ui

type Window struct {
	handle     Handle
	Active     bool
	Rect       Rectangle
	Background Background
}

func (w *Window) parent() Node {
	return nil
}

func (w *Window) child() Node {
	return nil
}

func (w *Window) draw(buf *RenderBuffer) {
	bgEntry := w.Background.backgroundEntry(w.Rect)
	buf.addEntry(bgEntry)
}
