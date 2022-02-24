package editor

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/nico-ec/uwu/ui"
)

type statusBar struct {
	statusLayout *ui.Layout

	lineLabel *ui.Label
	colLabel  *ui.Label
	errIcon   *ui.Icon
	errLabel  *ui.Label

	errorRaisedRecently bool
	errorTimer          int
	errorDuration       int
}

func newStatusBar(parent ui.Container, font *Font) statusBar {
	theme := getTheme()
	s := statusBar{
		statusLayout: &ui.Layout{
			Background: ui.Background{
				Kind: ui.BackgroundSolidColor,
				Clr:  theme.backgroundClr3,
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
			Clr:  theme.normalTextClr2,
			Size: 12,
		},
		colLabel: &ui.Label{
			Background: ui.Background{
				Visible: false,
			},
			Font: font,
			Text: "",
			Clr:  theme.normalTextClr2,
			Size: 12,
		},
		errIcon: &ui.Icon{},
		errLabel: &ui.Label{
			Background: ui.Background{
				Visible: false,
			},
			Font:  font,
			Text:  "",
			Align: ui.TextAlignCenterLeft,
			Clr:   theme.normalTextClr2,
			Size:  12,
		},
	}
	parent.AddWidget(s.statusLayout, ui.FitContainer) // 20 units I think
	s.statusLayout.AddWidget(s.lineLabel, int(font.MeasureText("line: 0000", 12)[0]))
	s.statusLayout.AddWidget(s.colLabel, int(font.MeasureText("column: 0000", 12)[0]))
	s.statusLayout.AddWidget(s.errIcon, 20)
	s.statusLayout.AddWidget(s.errLabel, ui.FitContainer)

	return s
}

func (s *statusBar) initStatusBar() {
	AddSignalListener(EditorLineChanged, s)
	AddSignalListener(EditorColumnChanged, s)
	AddSignalListener(EditorErrorRaised, s)
}

func (s *statusBar) updateStatusBar() {
	if s.errorRaisedRecently {
		s.errorTimer += 1
		if s.errorTimer == s.errorDuration {
			s.errorRaisedRecently = false
			s.errorTimer = 0
			// TODO: Have a cleaner way to do that
			s.errIcon.Img = nil
			s.errLabel.Text = ""
		}
	}
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

	case EditorErrorRaised:
		err := signal.Value.(SignalError)
		var iconImg *Image
		switch err.Kind {
		case editorWarning:
			iconImg = &ed.warning
		case editorError:
			iconImg = &ed.err
		}
		s.errIcon.Img = iconImg
		s.errLabel.SetText(err.Msg)
		s.errorRaisedRecently = true
		s.errorDuration = ebiten.MaxTPS() * 5
	}
}
