package editor

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/nico-ec/uwu/clipboard"
	"github.com/nico-ec/uwu/ui"
)

const initialAddedBufferCap = 200

type textEditor struct {
	currentEdit projectNode
	tabViewer   *ui.TabViewer
	// textBox        *ui.TextBox
	previousLine   int
	previousColumn int
}

func newTextEditor(parent ui.Container) textEditor {
	theme := getTheme()
	textEd := textEditor{
		tabViewer: &ui.TabViewer{
			HeaderBackground: ui.Background{
				Visible: true,
				Kind:    ui.BackgroundImageSlice,
				Clr:     theme.dividerClr,
				Img:     &ed.header,
				Constr:  ui.Constraint{2, 2, 2, 2},
			},
			HeaderHeight:    25,
			TabFont:         &ed.font,
			TabTextSize:     12,
			TabBckgroundClr: theme.backgroundClr3,
			TabFontClr:      theme.normalTextClr2,
		},
	}
	parent.AddWidget(textEd.tabViewer, ui.FitContainer)

	return textEd
}

// Extend the features of ui.TextBox and handle more input kind
func (t *textEditor) updateTextEditor() {
	textBox, ok := t.tabViewer.ActiveTab().(*ui.TextBox)
	if !ok {
		return
	}

	// Check if line or column changed and fire signal
	ln, col := textBox.CurrentLine(), textBox.CurrentColumn()
	switch {
	case ln != t.previousLine:
		FireSignal(EditorLineChanged, SignalInt(ln))
		t.previousLine = ln
		fallthrough
	case col != t.previousColumn:
		FireSignal(EditorColumnChanged, SignalInt(col))
		t.previousColumn = col
	}

	// Advanced input handling that textbox doesn't handle
	if ebiten.IsKeyPressed(ebiten.KeyControl) {
		switch {
		case inpututil.IsKeyJustPressed(ebiten.KeyS):
			t.saveNode()
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyEnd) {
		textBox.MoveCursorLineEnd()
	} else if inpututil.IsKeyJustPressed(ebiten.KeyHome) {
		textBox.MoveCursorLineStart()
	}
}

func (t *textEditor) saveNode() {
	if t.currentEdit == nil {
		return
	}
	textBox, ok := t.tabViewer.ActiveTab().(*ui.TextBox)
	if !ok {
		return
	}

	path := t.currentEdit.path()
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, fs.ModeExclusive)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	buf := textBox.GetCharBuffer()
	fmt.Println(buf)
	_, err = file.WriteString(string(textBox.GetCharBuffer()))
	if err != nil {
		panic(err)
	}

}

func (t *textEditor) loadNode(node projectNode) {
	data, err := os.ReadFile(node.path())
	if err != nil {
		panic(err)
	}
	d := bytes.Runes(data)
	t.currentEdit = node
	name := node.name()

	if !t.tabViewer.ContainsTab(name) {
		theme := getTheme()
		textBox := &ui.TextBox{
			Background: ui.Background{
				Visible: false,
			},
			Cap:                len(d) + initialAddedBufferCap,
			Margin:             10,
			Font:               &ed.font,
			TextSize:           12,
			TabSize:            2,
			AutoIndent:         true,
			Multiline:          true,
			HasRuler:           true,
			HasSyntaxHighlight: true,
			ShowCurrentLine:    true,
		}
		// Temporary. Those are go keywords
		// Allow for user to set their prefered
		// language from a given .toml file
		textBox.SetLexKeywords([]string{
			"type",
			"struct",
			"interface",
			"func",
			"go",
			"return",
			"bool",
			"uint",
			"uint8",
			"uint16",
			"uint32",
			"uint64",
			"int",
			"int8",
			"int16",
			"int32",
			"int64",
			"float64",
			"float32",
		})
		textBox.SetSyntaxColors(ui.ColorStyle{
			Normal:  theme.syntaxNormalClr,
			Keyword: theme.syntaxKeywordClr,
			Digit:   theme.syntaxNumberClr,
		})
		textBox.SetClipboardCallback(t)
		t.tabViewer.AddTab(name, textBox)
		textBox.LoadBufferData(d)
	} else {
		t.tabViewer.SetActiveTab(name)
	}
}

func (t *textEditor) ReadClipboard() string {
	data, err := clipboard.ReadClipboard()
	if err != nil {
		panic(err)
	}
	return data
}

func (t *textEditor) WriteClipboard(s string) {
	panic("Writing to clipboard not implemented yet")
}
