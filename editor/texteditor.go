package editor

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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
				Clr:     theme.backgroundClr2,
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

func (t *textEditor) updateTextEditor() {
	textBox, ok := t.tabViewer.ActiveTab().(*ui.TextBox)
	if !ok {
		return
	}
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

	if ebiten.IsKeyPressed(ebiten.KeyControl) && inpututil.IsKeyJustPressed(ebiten.KeyS) {
		t.saveNode()
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
		t.tabViewer.AddTab(name, textBox)
		textBox.LoadBufferData(d)
	} else {
		t.tabViewer.SetActiveTab(name)
	}
}
