package editor

import (
	"fmt"

	"github.com/nico-ec/uwu/ui"
)

type statusBar struct {
	statusHandle ui.Handle

	textHandle  ui.Handle
	lineLabel   *ui.Label
	colLabel    *ui.Label
	currentLine int
	currentCol  int
}

func newStatusBar(edWin, edText ui.Handle, font *Font) statusBar {
	s := statusBar{
		textHandle:  edText,
		currentLine: -1,
		currentCol:  -1,
	}
	s.statusHandle = ui.AddWidget(edWin, &ui.Layout{
		Background: ui.Background{
			Kind: ui.BackgroundSolidColor,
			Clr:  uwuTextClr,
		},
		Style: ui.Style{
			Ordering: ui.StyleOrderColumn,
			Padding:  10,
			Margin:   ui.Point{5, 0},
		},
	}, ui.FitContainer)
	s.lineLabel = &ui.Label{
		Background: ui.Background{
			Visible: false,
		},
		Font: font,
		Text: "",
		Clr:  uwuBackgroundClr,
		Size: 12,
	}
	ui.AddWidget(s.statusHandle, s.lineLabel, int(font.MeasureText("line: 0000", 12)[0]))
	s.colLabel = &ui.Label{
		Background: ui.Background{
			Visible: false,
		},
		Font: font,
		Text: "",
		Clr:  uwuBackgroundClr,
		Size: 12,
	}
	ui.AddWidget(s.statusHandle, s.colLabel, int(font.MeasureText("column: 0000", 12)[0]))

	return s
}

func (s *statusBar) updateStatus() {
	t := ui.GetWidget(s.textHandle).(*ui.TextBox)
	lineI := t.CurrentLine()
	colI := t.CurrentColumn()
	if s.currentLine != lineI {
		s.currentLine = lineI
		s.lineLabel.SetText(
			fmt.Sprintf("line: %d", s.currentLine),
		)
	}
	if s.currentCol != colI {
		s.currentCol = colI
		s.colLabel.SetText(
			fmt.Sprintf("column: %d", s.currentCol),
		)
	}
}
