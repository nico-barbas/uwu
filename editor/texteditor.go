package editor

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"

	"github.com/nico-ec/uwu/ui"
)

type textEditor struct {
	currentEdit projectNode
	handle      ui.Handle
	textBox     *ui.TextBox
}

func newTextEditor(parent ui.Handle) textEditor {
	textEd := textEditor{}
	tabView := &ui.TabViewer{
		HeaderBackground: ui.Background{
			Visible: true,
			Kind:    ui.BackgroundImageSlice,
			Clr:     ui.Color{232, 152, 168, 255},
			Img:     &ed.header,
			Constr:  ui.Constraint{2, 2, 2, 2},
		},
		HeaderHeight: 25,
		TabFont:      &ed.font,
		TabTextSize:  12,
		TabClr:       uwuTextClr,
	}
	ui.AddWidget(parent, tabView, ui.FitContainer)

	textEd.textBox = &ui.TextBox{
		Background: ui.Background{
			Visible: false,
		},
		Cap:                500,
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
	textEd.textBox.SetLexKeywords([]string{
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
	textEd.textBox.SetSyntaxColors(ui.ColorStyle{
		Normal:  uwuTextClr,
		Keyword: uwuKeywordClr,
		Digit:   uwuDigitClr,
	})
	textEd.handle = tabView.AddTab("test.go", textEd.textBox)
	// textEd.handle = ui.AddWidget(parent, textEd.textBox, ui.FitContainer)

	return textEd
}

func (t *textEditor) saveCurrentNode() {
	if t.currentEdit == nil {
		return
	}
	path := t.currentEdit.path()
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, fs.ModeExclusive)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	buf := t.textBox.GetCharBuffer()
	fmt.Println(buf)
	_, err = file.WriteString(string(t.textBox.GetCharBuffer()))
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
	t.textBox.LoadBufferData(d)
}
