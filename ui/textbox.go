package ui

import "fmt"

// Everything here is a little experimental and prototypish

const (
	initialLineBufferSize = 50
	textCursorWidth       = 2
	blinkTime             = 45
	rulerWidth            = 40
	rulerAlpha            = 155
)

const (
	cursorUp cursorDir = iota
	cursorDown
	cursorLeft
	cursorRight
)

const (
	tokenNewline tokenKind = iota
	tokenWhitespace
	tokenKeyword
	tokenIdentifier
	tokenNumber
	// Maybe need more granularity?
	tokenSymbol // =+-/*%()
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

		activeRect  Rectangle
		Margin      float64
		LinePadding float64
		Font        Font
		TextSize    float64
		TextClr     Color

		HasRuler  bool
		rulerRect Rectangle

		showCursor  bool
		cursor      Rectangle
		blinkTimer  int
		newlineSize float64
	}

	line struct {
		id   int
		text string
		// A sub slice of the backing buffer
		// with a record of the start and end in
		// the editor's buffer
		start int
		end   int

		current *token
		tokens  []token
		count   int
		// For graphical display
		origin Point
	}

	lexer struct{}

	tokenKind uint8

	token struct {
		start int
		end   int
		width float64
		kind  tokenKind
	}

	cursorDir uint8
)

func (t *TextBox) init() {
	t.charBuf = make([]rune, t.Cap)
	t.lines = make([]line, initialLineBufferSize)
	t.activeRect = Rectangle{
		X:      t.rect.X + t.Margin,
		Y:      t.rect.Y + t.Margin,
		Width:  t.rect.Width - t.Margin*2,
		Height: t.rect.Height - t.Margin*2,
	}
	if t.HasRuler {
		t.activeRect.X += t.Margin + rulerWidth
		t.rulerRect = Rectangle{
			X:      t.rect.X + t.Margin,
			Y:      t.rect.Y + t.Margin,
			Width:  rulerWidth,
			Height: t.rect.Height - (t.Margin * 2),
		}
	}
	t.lines[0] = line{
		id:    0,
		text:  fmt.Sprint(1),
		start: 0,
		end:   0,
		origin: Point{
			t.activeRect.X,
			t.activeRect.Y,
		},
	}
	t.lineIndex = 0
	t.currentLine = &t.lines[0]
	t.lineCount += 1

	t.cursor = Rectangle{
		X: t.activeRect.X, Y: t.activeRect.Y,
		Width: textCursorWidth, Height: t.TextSize,
	}
}

func (t *TextBox) update() {
	keys := pressedChars()
	if len(keys) > 0 {
		for _, k := range keys {
			t.insertChar(k)
		}
	}
	if isKeyRepeated(keyDelete) {
		t.deleteChar()
	}
	if isKeyRepeated(keyEnter) {
		// t.insertChar('\n')
		t.insertNewline()
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
		if t.HasRuler {

			lnWidth := t.Font.MeasureText(line.text, t.TextSize)
			buf.addEntry(RenderEntry{
				Kind: RenderText,
				Rect: Rectangle{
					X:      t.rulerRect.X + t.rulerRect.Width - lnWidth[0] - t.Margin,
					Y:      t.rulerRect.Y + (t.TextSize+t.LinePadding)*float64(i),
					Height: t.TextSize,
				},
				Clr:  Color{t.TextClr[0], t.TextClr[1], t.TextClr[2], rulerAlpha},
				Font: t.Font,
				Text: line.text,
			})
		}
	}
	if t.HasRuler {
		buf.addEntry(RenderEntry{
			Kind: RenderRectangle,
			Rect: Rectangle{
				X:      t.rulerRect.X + rulerWidth - 1,
				Y:      t.rulerRect.Y,
				Width:  1,
				Height: t.activeRect.Height,
			},
			Clr: Color{t.TextClr[0], t.TextClr[1], t.TextClr[2], rulerAlpha},
		})
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
		t.lines[i].start += 1
		t.lines[i].end += 1
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
		for i := t.currentLine.id + 1; i < t.lineCount; i += 1 {
			t.lines[i+1].start -= 1
			t.lines[i+1].end -= 1
		}
		t.currentLine.end -= 1
		t.caret -= 1
		t.charCount -= 1
		if t.currentLine.end < t.currentLine.start {
			t.deleteLine()
		} else {
			t.cursor.X -= t.Font.MeasureText(string(r), t.TextSize)[0]
		}
	}
}

func (t *TextBox) insertNewline() {
	copy(t.charBuf[t.caret+1:], t.charBuf[t.caret:t.charCount])
	t.charBuf[t.caret] = '\n'
	t.charCount += 1
	for i := t.currentLine.id + 1; i < t.lineCount; i += 1 {
		t.lines[i].start += 1
		t.lines[i].end += 1
	}
}

func (t *TextBox) insertLine() {
	t.lineIndex += 1
	cur := t.lineIndex
	// FIXME: This is not correct. Can end up out of bounds
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
		text:   fmt.Sprint(cur + 1),
		start:  t.charCount,
		end:    t.charCount,
		origin: o,
	}
	t.currentLine = &t.lines[cur]
	t.lineCount += 1
	t.moveCursorLineStart()

	// need to split if caret is in the middle of the line
}

// Do we assume that the carret is on the deleted line?
func (t *TextBox) deleteLine() {
	// FIXME: This is not correct. Can end up out of bounds
	for i := t.lineIndex; i < t.lineCount; i += 1 {
		t.lines[i] = t.lines[i+1]
		t.lines[i].id -= 1
		t.lines[i].origin[1] -= t.TextSize + t.LinePadding
	}
	t.lineIndex -= 1
	t.lineCount -= 1
	t.currentLine = &t.lines[t.lineIndex]
	// t.currentLine.end -= 1
	t.moveCursorLineEnd()
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
				t.lineIndex -= 1
				t.currentLine = &t.lines[t.lineIndex]
				t.moveCursorLineEnd()
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
	t.cursor.Y = t.currentLine.origin[1]
}

func (t *TextBox) moveCursorLineEnd() {
	t.caret = t.currentLine.end
	lineSize := t.Font.MeasureText(
		string(t.charBuf[t.currentLine.start:t.currentLine.end]),
		t.TextSize,
	)
	t.cursor.X = t.currentLine.origin[0] + lineSize[0]
	t.cursor.Y = t.currentLine.origin[1]
}

func (t *TextBox) lexLine(l *line) {
	for i := l.start; i < l.end; i += 1 {
		tok := token{
			start: i,
			end:   i,
			width: 0,
		}
		r := t.charBuf[i]
		switch r {
		case '\n':
			tok.kind = tokenNewline
		case ' ':
			whitespaceW := t.Font.MeasureText(" ", t.TextSize)
			peek := 1
			for ; ; peek += 1 {
				next := t.charBuf[i+peek]
				if next != ' ' {
					break
				}
			}
			tok.kind = tokenWhitespace
			tok.width = whitespaceW[0] * float64(peek)
		default:
			switch {
			case isDigit(r):
				for peek := 1; ; peek += 1 {
				}
			case isLetter(r):

			default:
				// still tag it as identifier for now
			}
		}
		l.addToken(tok)
	}
}

func (l *line) addToken(t token) {
	l.tokens[l.count] = t
	l.count += 1
}

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func isLetter(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}
