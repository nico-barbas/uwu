package ui

// Everything here is a little experimental and prototypish

const (
	initialLineBufferSize = 50
	textCursorWidth       = 2
	blinkTime             = 45
)

const (
	cursorUp = iota
	cursorDown
	cursorLeft
	cursorRight
)

type (
	TextBox struct {
		widgetRoot

		Background Background

		// This should be a dynamic array
		// in the off-chance that the edited text
		// is incredibly long and its size can't be
		// accounted beforehand
		Cap         int
		charBuf     []rune
		charCount   int
		lines       []line
		lineCount   int
		caret       int
		lineIndex   int
		currentLine *line

		Margin      float64
		LinePadding float64
		Font        Font
		TextSize    float64
		TextClr     Color

		showCursor  bool
		cursor      Rectangle
		blinkTimer  int
		newlineSize float64
	}

	// NOTE: This level of granularity should be sufficient for now
	// If not (i.e. needing text highlighting), can break up into words
	line struct {
		id int
		// A sub slice of the backing buffer
		// with a record of the start and end in
		// the editor's buffer
		start int
		end   int

		// For graphical display
		origin Point
	}

	cursorDir uint8
)

func (t *TextBox) init() {
	t.charBuf = make([]rune, t.Cap)
	t.lines = make([]line, initialLineBufferSize)
	t.lines[0] = line{
		id:    0,
		start: 0,
		end:   0,
		origin: Point{
			t.rect.X,
			t.rect.Y,
		},
	}
	t.lineIndex = 0
	t.currentLine = &t.lines[0]
	t.lineCount += 1

	t.cursor = Rectangle{
		X: t.rect.X, Y: t.rect.Y,
		Width: textCursorWidth, Height: t.TextSize,
	}
	// nl := t.Font.MeasureText("\n", t.TextSize)
	// t.newlineSize = nl[1] - t.TextSize
}

func (t *TextBox) update() {
	keys := pressedChars()
	if len(keys) > 0 {
		for _, k := range keys {
			t.insertChar(k)
		}
		// newTextSize := t.Font.MeasureText(string(keys), t.TextSize)
		// t.cursor.X += newTextSize[0]
	}
	if isKeyRepeated(keyDelete) {
		t.deleteChar()
	}
	if isKeyRepeated(keyEnter) {
		t.insertChar('\n')
		t.insertLine()
	}

	switch {
	case isKeyRepeated(keyUp):
	case isKeyRepeated(keyDown):
	case isKeyRepeated(keyLeft):
		t.moveCursorH(cursorLeft)
	case isKeyRepeated(keyRight):
		t.moveCursorH(cursorRight)
	}
	t.blinkTimer += 1
	if t.blinkTimer == blinkTime {
		t.blinkTimer = 0
		t.showCursor = !t.showCursor
	}
}

func (t *TextBox) draw(buf *renderBuffer) {
	bgEntry := t.Background.entry(t.rect)
	buf.addEntry(bgEntry)

	for i := 0; i < t.lineCount; i += 1 {
		line := &t.lines[i]
		end := line.end
		if t.charBuf[line.end] == '\n' {
			end -= 1
		}
		text := string(t.charBuf[line.start:end])
		textEntry := RenderEntry{
			Kind: RenderText,
			Rect: Rectangle{
				X:      line.origin[0],
				Y:      line.origin[1],
				Height: t.TextSize,
			},
			Clr:  t.TextClr,
			Font: t.Font,
			Text: text,
		}
		buf.addEntry(textEntry)
	}
	if t.showCursor {
		buf.addEntry(RenderEntry{
			Kind: RenderRectangle,
			Rect: t.cursor,
			Clr:  t.TextClr,
		})
	}

}

func (t *TextBox) insertChar(r rune) {
	copy(t.charBuf[t.caret+1:], t.charBuf[t.caret:t.charCount])
	t.charBuf[t.caret] = r
	t.charCount += 1
	t.currentLine.end += 1
	for i := t.currentLine.id + 1; i < t.lineCount; i += 1 {
		t.lines[i+1].start += 1
		t.lines[i+1].end += 1
	}
	t.cursor.X += t.Font.MeasureText(string(r), t.TextSize)[0]
	t.caret += 1
}

func (t *TextBox) deleteChar() {
	if t.charCount > 0 && t.caret > 0 {
		r := t.charBuf[t.caret-1]
		if t.caret < t.charCount {
			copy(t.charBuf[t.caret-1:], t.charBuf[t.caret:t.charCount])
		}
		t.charCount -= 1
		t.currentLine.end -= 1
		for i := t.currentLine.id + 1; i < t.lineCount; i += 1 {
			t.lines[i+1].start -= 1
			t.lines[i+1].end -= 1
		}
		t.cursor.X -= t.Font.MeasureText(string(r), t.TextSize)[0]
		t.caret -= 1
	}
}

func (t *TextBox) insertLine() {
	t.lineIndex += 1
	cur := t.lineIndex
	for i := cur; i < t.lineCount; i += 1 {
		t.lines[i+1] = t.lines[i]
		t.lines[i+1].id += 1
		t.lines[i+1].start += 1
		t.lines[i+1].end += 1
		t.lines[i+1].origin[1] += t.TextSize + t.LinePadding
	}
	var o Point
	if cur == t.lineCount {
		o = t.lines[cur-1].origin
		o[1] += t.TextSize + t.LinePadding
	} else {
		o = t.lines[cur].origin
	}
	t.lines[cur] = line{
		id:     cur,
		start:  t.charCount,
		end:    t.charCount,
		origin: o,
	}
	t.currentLine = &t.lines[cur]
	t.lineCount += 1
	t.cursor.X = t.currentLine.origin[0]
	t.cursor.Y = t.currentLine.origin[1]

	// need to split if caret is in the middle of the line
}

func (t *TextBox) moveCursorV(dir cursorDir) {
	switch dir {
	case cursorUp:
		if t.lineIndex > 0 {
			col := t.caret - t.currentLine.start
			t.lineIndex -= 1
			t.currentLine = &t.lines[t.lineIndex-1]
			if t.currentLine.start+col < t.currentLine.end {
				t.caret = t.currentLine.start + col
			} else {
				t.moveCursorLineEnd()
			}
			t.cursor.Y = t.currentLine.origin[1]
		}
	case cursorDown:
		if t.lineIndex <= t.lineCount {
			col := t.caret - t.currentLine.start
			t.lineIndex += 1
			t.currentLine = &t.lines[t.lineIndex+1]
			if t.currentLine.start+col < t.currentLine.end {
				t.caret = t.currentLine.start + col
			} else {
				t.moveCursorLineEnd()
			}
			t.cursor.Y = t.currentLine.origin[1]
		}
	}
}

func (t *TextBox) moveCursorH(dir cursorDir) {
	switch dir {
	case cursorRight:
		if t.caret+1 <= t.charCount {
			if t.caret+1 > t.currentLine.end {

			} else {
				c := t.charBuf[t.caret]
				t.cursor.X += t.Font.MeasureText(string(c), t.TextSize)[0]
				t.caret += 1
			}
		}
	case cursorLeft:
		if t.caret-1 > 0 {
			if t.caret-1 < t.currentLine.start {

			} else if t.caret > 0 {
				c := t.charBuf[t.caret-1]
				t.cursor.X -= t.Font.MeasureText(string(c), t.TextSize)[0]
				t.caret -= 1
			}
		}
	}
}

func (t *TextBox) moveCursorLineStart() {
	t.caret = t.currentLine.start
	t.cursor.X = t.currentLine.origin[0]
}

func (t *TextBox) moveCursorLineEnd() {
	t.caret = t.currentLine.end
	lineSize := t.Font.MeasureText(
		string(t.charBuf[t.currentLine.start:t.currentLine.end]),
		t.TextSize,
	)
	t.cursor.X = t.currentLine.origin[0] + lineSize[0]
}
