package editor

import (
	"fmt"

	"github.com/nico-ec/uwu/ui"
)

type statusBar struct {
	statusLayout *ui.Layout

	lineLabel *ui.Label
	colLabel  *ui.Label
}

func newStatusBar(parent ui.Container, font *Font) statusBar {
	s := statusBar{
		statusLayout: &ui.Layout{
			Background: ui.Background{
				Kind: ui.BackgroundSolidColor,
				Clr:  uwuTextClr,
			},
			Style: ui.Style{
				Ordering: ui.StyleOrderColumn,
				Padding:  10,
				Margin:   ui.Point{5, 0},
			},
		},
		lineLabel: &ui.Label{
			Background: ui.Background{
				Visible: false,
			},
			Font: font,
			Text: "",
			Clr:  uwuBackgroundClr,
			Size: 12,
		},
		colLabel: &ui.Label{
			Background: ui.Background{
				Visible: false,
			},
			Font: font,
			Text: "",
			Clr:  uwuBackgroundClr,
			Size: 12,
		},
	}
	parent.AddWidget(s.statusLayout, ui.FitContainer)
	s.statusLayout.AddWidget(s.lineLabel, int(font.MeasureText("line: 0000", 12)[0]))
	s.statusLayout.AddWidget(s.colLabel, int(font.MeasureText("column: 0000", 12)[0]))

	return s
}

func (s *statusBar) initStatusBar() {
	AddSignalListener(EditorLineChanged, s)
	AddSignalListener(EditorColumnChanged, s)
}

func (s *statusBar) OnSignal(signal Signal) {
	switch signal.Kind {
	case EditorLineChanged:
		s.lineLabel.SetText(
			fmt.Sprintf("line: %d", signal.Value),
		)
	case EditorColumnChanged:
		s.colLabel.SetText(
			fmt.Sprintf("column: %d", signal.Value),
		)
	}
}
