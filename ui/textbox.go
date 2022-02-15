package ui

import (
	"fmt"
	"log"
)

// Everything here is a little experimental and prototypish

const (
	initialLineBufferSize = 50
	initialTokenCap       = 20
	textCursorWidth       = 2
	blinkTime             = 45
	rulerWidth            = 40
	rulerAlpha            = 155
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
		Cap           int
		charBuf       []rune
		charCount     int
		lines         []line
		lineCount     int
		caret         int
		lineIndex     int
		currentLine   *line
		currentIndent int

		activeRect  Rectangle
		Margin      float64
		LinePadding float64
		TabSize     int
		AutoIndent  bool
		Font        Font
		TextSize    float64
		TextClr     Color

		HasRuler  bool
		rulerRect Rectangle

		showCursor      bool
		cursor          Rectangle
		blinkTimer      int
		ShowCurrentLine bool

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
		start     int
		end       int
		indentEnd int

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
	t.lines[0].tokens = make([]token, initialTokenCap)
	if t.HasSyntaxHighlight {
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
	inBoxBounds := t.activeRect.pointInBounds(mPos)
	if inBoxBounds {
		setCursorShape(CursorShapeText)
	} else {
		setCursorShape(CursorShapeDefault)
	}
	if isMouseJustPressed() {
		if inBoxBounds {
			if !t.focused {
				t.focused = true
			}
			t.moveCursorToMouse(mPos)
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
				if t.caret == t.currentLine.indentEnd {
					t.deleteIndent()
				} else {
					t.deleteChar()
				}
			}
		}
		if isKeyRepeated(keyEnter) {
			t.insertLine()
		}
		if isKeyRepeated(keyTab) {
			t.insertIndent()
		}
		switch {
		case isKeyRepeated(keyUp):
			t.moveCursorUp()
		case isKeyRepeated(keyDown):
			t.moveCursorDown()
		case isKeyRepeated(keyLeft):
			if isKeyPressed(keyCtlr) {
				t.moveCursorToPreviousWord()
			} else {
				t.moveCursorLeft()
			}
		case isKeyRepeated(keyRight):
			if isKeyPressed(keyCtlr) {
				t.moveCursorToNextWord()
			} else {
				t.moveCursorRight()
			}
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

	if t.ShowCurrentLine {
		buf.addEntry(RenderEntry{
			Kind: RenderRectangle,
			Rect: Rectangle{
				X:      t.currentLine.origin[0],
				Y:      t.currentLine.origin[1],
				Width:  t.activeRect.Width,
				Height: t.TextSize,
			},
			Clr: Color{t.TextClr[0], t.TextClr[1], t.TextClr[2], rulerAlpha},
		})
	}

	for i := 0; i < t.lineCount; i += 1 {
		line := &t.lines[i]
		var xptr float64 = 0
		for j := 0; j < line.count; j += 1 {
			var clr Color
			token := line.tokens[j]
			text := string(t.charBuf[line.start+token.start : line.start+token.end])
			switch t.HasSyntaxHighlight {
			case true:
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
			case false:
				clr = t.TextClr
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
		// case false:
		// 	end := line.end
		// 	text := string(t.charBuf[line.start:end])
		// 	textEntry := RenderEntry{
		// 		Kind: RenderText,
		// 		Rect: Rectangle{
		// 			X:      line.origin[0],
		// 			Y:      line.origin[1],
		// 			Height: t.TextSize,
		// 		},
		// 		Clr:  t.TextClr,
		// 		Font: t.Font,
		// 		Text: text,
		// 	}
		// 	buf.addEntry(textEntry)
		// }
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

func (t *TextBox) LoadBufferData(data []rune) error {
	if len(data) > t.Cap {
		log.SetPrefix("[UI Debug]: ")
		log.Println("Given file is too big compared to the TextBox capacity")
		return fmt.Errorf("given file is too big")
	}
	t.charCount = len(data)
	t.lineCount = 0
	t.lineIndex = 0
	t.caret = 0

	t.lines[0] = line{
		id:    0,
		text:  fmt.Sprint(1),
		start: 0,
		end:   0,
		origin: Point{
			t.activeRect.X,
			t.activeRect.Y,
		},
		tokens: make([]token, initialTokenCap),
	}
	t.lineCount += 1

	var current int = 0
	var c rune
	for {
		if current >= len(data) {
			break
		}
		c = data[current]
		t.charBuf[current] = c
		current += 1
		if c == '\r' && data[current] == '\n' {
			t.charBuf[current] = '\n'
			current += 1
			i := current

			t.lineIndex += 1
			t.lines[t.lineCount] = line{
				id:        t.lineCount,
				text:      fmt.Sprint(t.lineCount + 1),
				start:     i,
				end:       i,
				indentEnd: i,
				origin: Point{
					t.lines[t.lineCount-1].origin[0],
					t.lines[t.lineCount-1].origin[1] + t.TextSize + t.LinePadding,
				},
				tokens: make([]token, initialTokenCap),
			}
			t.lineCount += 1
			continue
		}
		t.lines[t.lineIndex].end += 1

	}
	t.lineIndex = 0
	t.currentLine = &t.lines[t.lineIndex]
	for i := 0; i < t.lineCount; i += 1 {
		t.lexLine(&t.lines[i])
	}
	return nil
}

func (t *TextBox) GetCharBuffer() []rune {
	return t.charBuf[:t.charCount]
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

	t.lexLine(t.currentLine)
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

	t.lexLine(t.currentLine)
}

func (t *TextBox) insertNewline() {
	copy(t.charBuf[t.caret+2:], t.charBuf[t.caret:t.charCount])
	t.charBuf[t.caret] = '\r'
	t.charBuf[t.caret+1] = '\n'
	t.charCount += 2
	for i := t.currentLine.id + 1; i < t.lineCount; i += 1 {
		t.lines[i].start += 2
		t.lines[i].end += 2
	}
}

func (t *TextBox) insertIndent() {
	// copy(t.charBuf[t.caret+t.TabSize:], t.charBuf[t.caret:t.charCount])
	if t.caret == t.currentLine.indentEnd {
		t.currentIndent += 1
		t.currentLine.indentEnd += 1
		// t.currentLine.indent += 1
	}
	t.insertChar('\t')
}

func (t *TextBox) deleteIndent() {
	if t.currentLine.start == t.currentLine.indentEnd {
		t.deleteChar()
	} else {
		t.currentIndent -= 1
		t.currentLine.indentEnd -= 1
		t.deleteChar()
	}
}

func (t *TextBox) insertLine() {
	t.insertNewline()
	newlineStart := t.caret + 2
	newlineEnd := t.currentLine.end + 2
	t.currentLine.end = t.caret
	t.lexLine(t.currentLine)

	// Move all the line by one to make room for the new line
	t.lineCount += 1
	for i := t.lineCount - 1; i >= t.lineIndex+2; i -= 1 {
		t.lines[i] = t.lines[i-1]
		t.lines[i].id += 1
		t.lines[i].text = fmt.Sprint(i + 1)
		t.lines[i].origin[1] += t.TextSize + t.LinePadding
	}
	t.lineIndex += 1
	t.currentLine = &t.lines[t.lineIndex]
	t.lines[t.lineIndex] = line{
		id:        t.lineIndex,
		text:      fmt.Sprint(t.lineIndex + 1),
		start:     newlineStart,
		end:       newlineEnd,
		indentEnd: newlineStart,
		origin: Point{
			t.lines[t.lineIndex-1].origin[0],
			t.lines[t.lineIndex-1].origin[1] + t.TextSize + t.LinePadding,
		},
	}
	t.currentLine.tokens = make([]token, initialTokenCap)
	t.moveCursorLineStart()
	if t.AutoIndent {
		for i := 0; i < t.currentIndent; i += 1 {
			t.insertChar('\t')
			t.currentLine.indentEnd += 1
		}
	}
	t.lexLine(t.currentLine)
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

func (t *TextBox) moveCursorUp() {
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
}

func (t *TextBox) moveCursorDown() {
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

func (t *TextBox) moveCursorRight() {
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
}

func (t *TextBox) moveCursorLeft() {
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

func (t *TextBox) moveCursorToNextWord() {
	if t.caret+1 <= t.charCount {
		c := t.charBuf[t.caret]
		if isTerminalSymbol(c) {
			t.cursor.X += t.Font.GlyphAdvance(c, t.TextSize)
			t.caret += 1
		}
	}
	for t.caret+1 <= t.charCount {
		c := t.charBuf[t.caret]
		if isTerminalSymbol(c) {
			break
		}
		t.cursor.X += t.Font.GlyphAdvance(c, t.TextSize)
		t.caret += 1
	}
}

func (t *TextBox) moveCursorToPreviousWord() {
	if t.caret-1 >= 0 {
		c := t.charBuf[t.caret-1]
		if isTerminalSymbol(c) {
			t.cursor.X -= t.Font.GlyphAdvance(c, t.TextSize)
			t.caret -= 1
		}
	}
	for t.caret-1 >= 0 {
		c := t.charBuf[t.caret-1]
		if isTerminalSymbol(c) {
			break
		}
		t.cursor.X -= t.Font.GlyphAdvance(c, t.TextSize)
		t.caret -= 1
	}
}

func (t *TextBox) moveCursorToMouse(mPos Point) {
	relPos := mPos[1] - t.activeRect.Y
	t.lineIndex = int(relPos / (t.TextSize + t.LinePadding))
	if t.lineIndex >= 0 && t.lineIndex < t.lineCount {
		t.currentLine = &t.lines[t.lineIndex]

		// Search for the correct rune to position the cursor to
		selectedLine := &t.lines[t.lineIndex]
		currentXStartPos := selectedLine.origin[0]
		currentXEndPos := currentXStartPos
		t.moveCursorLineStart()
		for j := selectedLine.start; j < selectedLine.end; j += 1 {
			advance := t.Font.GlyphAdvance(t.charBuf[j], t.TextSize)
			currentXEndPos += advance
			if mPos[0] >= currentXStartPos && mPos[0] <= currentXEndPos {
				break
			}
			t.caret = j + 1
			t.cursor.X += advance
			currentXStartPos = currentXEndPos
		}
	} else {
		t.lineIndex = t.lineCount - 1
		t.currentLine = &t.lines[t.lineIndex]
		t.moveCursorLineEnd()
	}
}

func (t *TextBox) moveCursorLineStart() {
	t.caret = t.currentLine.start
	t.cursor.X = t.currentLine.origin[0]
	t.cursor.Y = t.currentLine.origin[1]
}

func (t *TextBox) moveCursorLineEnd() {
	t.caret = t.currentLine.end
	var lineAdvance float64
	for i := t.currentLine.start; i < t.currentLine.end; i += 1 {
		lineAdvance += t.Font.GlyphAdvance(t.charBuf[i], t.TextSize)
	}
	// lineSize := t.Font.MeasureText(
	// 	string(t.charBuf[t.currentLine.start:t.currentLine.end]),
	// 	t.TextSize,
	// )
	t.cursor.X = t.currentLine.origin[0] + lineAdvance
	t.cursor.Y = t.currentLine.origin[1]
}

func (t *TextBox) CurrentLine() int {
	return t.lineIndex + 1
}

func (t *TextBox) CurrentColumn() int {
	return t.caret - t.currentLine.start
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
		case '\n', '\r':
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
		newbuf := make([]token, len(l.tokens)*2)
		copy(newbuf[:], l.tokens[:len(l.tokens)])
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

func isTerminalSymbol(r rune) bool {
	return r == ' ' || r == '.' || r == '/' || r == '{' || r == '[' || r == '('
}
