package ui

// Everything here is a little experimental and prototypish

const (
	initialLineBufferSize = 50
	textCursorWidth       = 2
)

type TextBox struct {
	widgetRoot

	Background Background

	// This should be a dynamic array
	// in the off-chance that the edited text
	// is incredibly long and its size can't be
	// accounted beforehand
	Cap          int
	charBuf      []rune
	charCount    int
	lines        []line
	lineCount    int
	current      *line
	currentIndex int

	Margin      float64
	LinePadding float64
	Font        Font
	TextSize    float64
	TextClr     Color

	cursor      Rectangle
	newlineSize float64
}

// NOTE: This level of breaking up should be sufficient for now
// If not (i.e. needing text highlighting), can break up into words
type line struct {
	// A sub slice of the backing buffer
	data     []rune
	origin   Point
	end      Point
	current  rune
	previous rune
}

func (t *TextBox) init() {
	t.charBuf = make([]rune, t.Cap)
	t.lines = make([]line, 0, initialLineBufferSize)

	t.cursor = Rectangle{
		X: t.rect.X, Y: t.rect.Y,
		Width: textCursorWidth, Height: t.TextSize,
	}
	nl := t.Font.MeasureText("\n", t.TextSize)
	t.newlineSize = nl[1] - t.TextSize
}

func (t *TextBox) update() {
	keys := pressedChars()
	if len(keys) > 0 {
		for _, k := range keys {
			t.charBuf[t.charCount] = k
			t.charCount += 1
		}
		newTextSize := t.Font.MeasureText(string(keys), t.TextSize)
		t.cursor.X += newTextSize[0]
	}
	if isKeyRepeated(keyDelete) {
		if t.charCount > 0 {
			t.charCount -= 1
			delChar := t.charBuf[t.charCount]
			switch delChar {
			case '\n':
				t.cursor.Y -= t.newlineSize
			default:
				delCharSize := t.Font.MeasureText(string(delChar), t.TextSize)
				t.cursor.X -= delCharSize[0]
			}
		}
	}
	if isKeyRepeated(keyEnter) {
		t.charBuf[t.charCount] = '\n'
		t.charCount += 1

		t.cursor.X = t.rect.X
		t.cursor.Y += t.newlineSize
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
	cursorEntry := RenderEntry{
		Kind: RenderRectangle,
		Rect: t.cursor,
		Clr:  t.TextClr,
	}
	buf.addEntry(cursorEntry)
}

func (t *TextBox) addLine() {
	t.lines[t.lineCount] = line{
		origin: Point{
			t.rect.X,
			t.rect.Y + float64(t.lineCount)*(t.TextSize+t.LinePadding),
		},
	}
	t.lineCount += 1
}
