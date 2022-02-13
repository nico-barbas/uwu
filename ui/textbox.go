package ui

import "fmt"

// Everything here is a little experimental and prototypish

const (
	initialLineBufferSize = 50
	initialTokenCap       = 20
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

type (
	TextBox struct {
		widgetRoot

		Background Background
		focused    bool
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

		HasSyntaxHighlight bool
		lexer              lexer
		clrStyle           ColorStyle
	}

	ColorStyle struct {
		Normal  Color
		Keyword Color
		Digit   Color
	}

	line struct {
		id   int
		text string
		// A sub slice of the backing buffer
		// with a record of the start and end in
		// the editor's buffer
		start int
		end   int

		tokens []token
		count  int
		// For graphical display
		origin Point
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
	if t.HasSyntaxHighlight {
		t.lines[0].tokens = make([]token, initialTokenCap)
		t.TextClr = t.clrStyle.Normal
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
	mPos := mousePosition()
	if isMouseJustPressed() {
		if t.rect.pointInBounds(mPos) {
			if !t.focused {
				t.focused = true
			}
			// Set the cursor position to the closest character
		} else {
			t.focused = false
		}
	}
	if t.focused {
		keys := pressedChars()
		if len(keys) > 0 {
			for _, k := range keys {
				t.insertChar(k)
			}
		}
		if isKeyRepeated(keyDelete) {
			if isKeyRepeated(keyCtlr) {
				// delete word
			} else {
				t.deleteChar()
			}
		}
		if isKeyRepeated(keyEnter) {
			// t.insertChar('\n')
			// t.insertNewline()
			t.insertLine()
		}

		switch {
		case isKeyRepeated(keyUp):
			t.moveCursorV(cursorUp)
		case isKeyRepeated(keyDown):
			t.moveCursorV(cursorDown)
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
}

func (t *TextBox) draw(buf *renderBuffer) {
	bgEntry := t.Background.entry(t.rect)
	buf.addEntry(bgEntry)

	for i := 0; i < t.lineCount; i += 1 {
		line := &t.lines[i]
		switch t.HasSyntaxHighlight {
		case true:
			var xptr float64 = 0
			for j := 0; j < line.count; j += 1 {
				var clr Color
				token := line.tokens[j]
				text := string(t.charBuf[line.start+token.start : line.start+token.end])
				switch token.kind {
				case tokenIdentifier:
					clr = t.clrStyle.Normal
				case tokenKeyword:
					clr = t.clrStyle.Keyword
				case tokenNumber:
					clr = t.clrStyle.Digit
				default:
					clr = t.clrStyle.Normal
				}
				buf.addEntry(RenderEntry{
					Kind: RenderText,
					Rect: Rectangle{
						X:      line.origin[0] + xptr,
						Y:      line.origin[1],
						Height: t.TextSize,
					},
					Clr:  clr,
					Font: t.Font,
					Text: text,
				})
				xptr += token.width
			}
		case false:
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
		}
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
	if t.showCursor && t.focused {
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
	t.cursor.X += t.Font.GlyphAdvance(r, t.TextSize)
	t.caret += 1

	if t.HasSyntaxHighlight {
		t.lexLine(t.currentLine)
	}
}

func (t *TextBox) deleteChar() {
	if t.charCount > 0 && t.caret > 0 {
		r := t.charBuf[t.caret-1]
		if t.caret < t.charCount {
			copy(t.charBuf[t.caret-1:], t.charBuf[t.caret:t.charCount])
		}
		for i := t.currentLine.id + 1; i < t.lineCount; i += 1 {
			t.lines[i].start -= 1
			t.lines[i].end -= 1
		}
		t.currentLine.end -= 1
		t.caret -= 1
		t.charCount -= 1
		if t.currentLine.end < t.currentLine.start {
			t.deleteLine()
		} else {
			t.cursor.X -= t.Font.GlyphAdvance(r, t.TextSize)
		}
	}

	if t.HasSyntaxHighlight {
		t.lexLine(t.currentLine)
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
	t.insertNewline()
	fmt.Println(t.charBuf[t.currentLine.start:t.currentLine.end])
	newlineStart := t.caret + 1
	newlineEnd := t.currentLine.end + 1
	t.currentLine.end = t.caret
	if t.HasSyntaxHighlight {
		t.lexLine(t.currentLine)
	}

	t.lineCount += 1
	for i := t.lineIndex + 2; i < t.lineCount; i += 1 {
		t.lines[i] = t.lines[i-1]
		t.lines[i].id += 1
		t.lines[i].text = fmt.Sprint(i + 1)
		t.lines[i].origin[1] += t.TextSize + t.LinePadding
	}
	t.lineIndex += 1
	t.currentLine = &t.lines[t.lineIndex]
	t.lines[t.lineIndex] = line{
		id:    t.lineIndex,
		text:  fmt.Sprint(t.lineIndex + 1),
		start: newlineStart,
		end:   newlineEnd,
		origin: Point{
			t.lines[t.lineIndex-1].origin[0],
			t.lines[t.lineIndex-1].origin[1] + t.TextSize + t.LinePadding,
		},
	}
	if t.HasSyntaxHighlight {
		t.currentLine.tokens = make([]token, initialTokenCap)
		t.lexLine(t.currentLine)
	}
	t.moveCursorLineStart()
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
			t.currentLine = &t.lines[t.lineIndex]
			if t.currentLine.start+col < t.currentLine.end {
				t.caret = t.currentLine.start + col
			} else {
				t.moveCursorLineEnd()
			}
			t.cursor.Y = t.currentLine.origin[1]
		}
	case cursorDown:
		if t.lineIndex < t.lineCount-1 {
			col := t.caret - t.currentLine.start
			t.lineIndex += 1
			t.currentLine = &t.lines[t.lineIndex]
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
				t.lineIndex += 1
				t.currentLine = &t.lines[t.lineIndex]
				t.moveCursorLineStart()
			} else {
				c := t.charBuf[t.caret]
				t.cursor.X += t.Font.GlyphAdvance(c, t.TextSize)
				t.caret += 1
			}
		}
	case cursorLeft:
		if t.caret-1 >= 0 {
			if t.caret-1 < t.currentLine.start {
				t.lineIndex -= 1
				t.currentLine = &t.lines[t.lineIndex]
				t.moveCursorLineEnd()
			} else if t.caret > 0 {
				c := t.charBuf[t.caret-1]
				t.cursor.X -= t.Font.GlyphAdvance(c, t.TextSize)
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

//
// Lexing
//

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
	lexer struct {
		input    []rune
		start    int
		current  int
		keywords []string
	}

	tokenKind uint8

	token struct {
		start int
		end   int
		width float64
		kind  tokenKind
	}
)

func (t *TextBox) SetLexKeywords(kw []string) {
	t.lexer.keywords = kw
}

func (t *TextBox) SetSyntaxColors(style ColorStyle) {
	t.clrStyle = style
}

func (t *TextBox) lexInit(start int, substr []rune) {
	t.lexer.input = substr
	t.lexer.start = start
	t.lexer.current = 0
}

func (t *TextBox) lexLine(l *line) {
	l.emptyTokens()
	t.lexInit(
		l.start,
		t.charBuf[l.start:l.end],
	)

lex:
	for {
		if t.lexer.eof() {
			break lex
		}
		tok := token{
			start: t.lexer.current,
		}
		start := t.lexer.current
		c := t.lexer.advance()
		switch c {
		case '\n':
			tok.kind = tokenNewline

		case ' ':
			wCount := 1
			for {
				if t.lexer.eof() {
					break
				}
				next := t.lexer.peek()
				if next != ' ' {
					break
				}
				t.lexer.advance()
				wCount += 1
			}
			tok.kind = tokenWhitespace

		default:
			switch {
			case isDigit(c):
				for {
					if t.lexer.eof() {
						break
					}
					next := t.lexer.peek()
					hasDecimal := false
					if !isDigit(next) {
						if next == '.' && !hasDecimal {
							hasDecimal = true
						} else {
							break
						}
					}
					t.lexer.advance()
				}
				tok.kind = tokenNumber

			case isLetter(c):
				for {
					if t.lexer.eof() {
						break
					}
					next := t.lexer.peek()
					if !isLetter(next) {
						break
					}
					t.lexer.advance()
				}
				word := t.lexer.input[start:t.lexer.current]
				if t.lexer.isKeyword(string(word)) {
					tok.kind = tokenKeyword
				} else {
					tok.kind = tokenIdentifier
				}

			default:
				// This is the default branch for all the symbols for now
				// may need more granularity
			}
		}

		tok.end = t.lexer.current
		for i := tok.start; i < tok.end; i += 1 {
			r := t.charBuf[l.start+i]
			tok.width += t.Font.GlyphAdvance(r, t.TextSize)
		}
		l.addToken(tok)
	}
}

func (l *line) addToken(t token) {
	if l.count >= len(l.tokens) {
		newbuf := make([]token, 0, len(l.tokens)*2)
		copy(newbuf[:], l.tokens[:])
		l.tokens = newbuf
	}
	l.tokens[l.count] = t
	l.count += 1
}

func (l *line) emptyTokens() {
	l.count = 0
}

func (l *lexer) advance() rune {
	l.current += 1
	return l.input[l.current-1]
}

func (l *lexer) peek() rune {
	return l.input[l.current]
}

func (l *lexer) eof() bool {
	return l.current >= len(l.input)
}

func (l *lexer) isKeyword(word string) bool {
	for _, keyword := range l.keywords {
		if string(word) == keyword {
			return true
		}
	}
	return false
}

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func isLetter(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}
