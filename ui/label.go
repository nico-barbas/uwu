package ui

type Label struct {
	widgetRoot
	Background Background
	Font       Font
	Text       string
	Clr        Color
	Size       float64
}

func (l *Label) draw(buf *renderBuffer) {
	bgEntry := l.Background.entry(l.rect)
	buf.addEntry(bgEntry)

	textSize := l.Font.MeasureText(l.Text, l.Size)
	textEntry := RenderEntry{
		Kind: RenderText,
		Rect: Rectangle{
			X:      l.rect.X + (l.rect.Width/2 - textSize[0]/2),
			Y:      l.rect.Y + (l.rect.Height/2 - textSize[1]/2),
			Height: l.Size,
		},
		Clr:  l.Clr,
		Font: l.Font,
		Text: l.Text,
	}
	buf.addEntry(textEntry)
}

func (l *Label) SetText(text string) {
	l.Text = text
}
