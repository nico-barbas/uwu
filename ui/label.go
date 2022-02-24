package ui

type Label struct {
	widgetRoot
	Background Background
	Font       Font
	Text       string
	Align      TextAlign
	Clr        Color
	Size       float64
}

func (l *Label) draw(buf *renderBuffer) {
	bgEntry := l.Background.entry(l.rect)
	buf.addEntry(bgEntry)

	textSize := l.Font.MeasureText(l.Text, l.Size)
	var textRect Rectangle

	switch l.Align {
	case TextAlignCenter:
		textRect = Rectangle{
			X:      l.rect.X + (l.rect.Width/2 - textSize[0]/2),
			Y:      l.rect.Y + (l.rect.Height/2 - textSize[1]/2),
			Height: l.Size,
		}
	case TextAlignCenterLeft:
		textRect = Rectangle{
			X:      l.rect.X,
			Y:      l.rect.Y + (l.rect.Height/2 - textSize[1]/2),
			Height: l.Size,
		}
	case TextAlignCenterRight:
		textRect = Rectangle{
			X:      l.rect.X + (l.rect.Width - textSize[0]),
			Y:      l.rect.Y + (l.rect.Height/2 - textSize[1]/2),
			Height: l.Size,
		}
	}

	textEntry := RenderEntry{
		Kind: RenderText,
		Rect: textRect,
		Clr:  l.Clr,
		Font: l.Font,
		Text: l.Text,
	}
	buf.addEntry(textEntry)
}

func (l *Label) SetText(text string) {
	l.Text = text
}
