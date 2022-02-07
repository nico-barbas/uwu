package ui

type TextBox struct {
	widgetRoot

	Background Background

	Cap       int
	charBuf   []rune
	charCount int

	Margin   float64
	Font     Font
	TextSize float64
	TextClr  Color
}

func (t *TextBox) init() {
	t.charBuf = make([]rune, t.Cap)
}

func (t *TextBox) update() {
	keys := pressedChars()
	if len(keys) > 0 {
		for _, k := range keys {
			t.charBuf[t.charCount] = k
			t.charCount += 1
		}
	}
	if isKeyRepeated(keyDelete) {
		if t.charCount > 0 {
			t.charCount -= 1
		}
	}
	if isKeyRepeated(keyEnter) {
		t.charBuf[t.charCount] = '\n'
		t.charCount += 1
	}
}

func (t *TextBox) draw(buf *renderBuffer) {
	bgEntry := t.Background.entry(t.rect)
	buf.addEntry(bgEntry)

	text := string(t.charBuf[:t.charCount])
	textEntry := RenderEntry{
		Kind: RenderText,
		Rect: Rectangle{
			X:      t.rect.X + t.Margin,
			Y:      t.rect.Y + t.Margin,
			Height: t.TextSize,
		},
		Clr:  t.TextClr,
		Font: t.Font,
		Text: text,
	}
	buf.addEntry(textEntry)
}
