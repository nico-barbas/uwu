package editor

import (
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/nico-ec/uwu/ui"
)

type CmdPanel struct {
	window  ui.WinHandle
	textBox *ui.TextBox
}

func (c *CmdPanel) initCmdPanel() {
	theme := getTheme()
	//
	// Search window
	//
	c.window = ui.AddWindow(ui.Window{
		Active: false,
		Rect:   ui.Rectangle{550, 380, 500, 44},
		Style: ui.Style{
			Ordering: ui.StyleOrderRow,
			Padding:  0,
			Margin:   ui.Point{0, 0},
		},
		Background: ui.Background{
			Visible: true,
			Kind:    ui.BackgroundSolidColor,
			Clr:     theme.backgroundClr1,
		},
		HasHeader:    true,
		HeaderHeight: 20,
		HeaderBackground: ui.Background{
			Visible: true,
			Kind:    ui.BackgroundImageSlice,
			Clr:     theme.dividerClr,
			Img:     &ed.header,
			Constr:  ui.Constraint{2, 2, 2, 2},
		},
		HasHeaderTitle: true,
		HeaderTitle:    "Command",
		HeaderFont:     &ed.font,
		HeaderFontSize: 12,
		HeaderFontClr:  theme.normalTextClr,

		HasBorders:  true,
		BorderWidth: 1,
		BorderColor: theme.dividerClr,
	})
	c.window.SetCloseBtn(ui.Button{
		Background: ui.Background{
			Visible: true,
			Kind:    ui.BackgroundSolidColor,
		},
		UserID:       editorCloseBtn,
		Clr:          theme.backgroundClr3,
		HighlightClr: theme.backgroundClr3,
		PressedClr:   theme.backgroundClr3,
		HasIcon:      true,
		Icon:         &ed.cross,
		IconClr:      theme.backgroundClr1,
		Receiver:     c,
	})

	c.textBox = &ui.TextBox{
		Background: ui.Background{
			Visible: false,
		},
		Cap:       500,
		Margin:    3,
		Font:      &ed.font,
		TextSize:  12,
		TextClr:   theme.normalTextClr,
		Multiline: false,
	}
	c.window.AddWidget(c.textBox, ui.FitContainer)
	c.window.UnfocusWindow()
}

func (c *CmdPanel) updateCmdPanel() {
	if ebiten.IsKeyPressed(ebiten.KeyControl) && inpututil.IsKeyJustPressed(ebiten.KeyP) {
		c.window.SetActive(!c.window.IsActive())
		c.textBox.SetFocus(true)
	}
	if c.window.IsActive() {
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			cmd := c.textBox.GetCharBuffer()
			c.parseCommand(string(cmd))

			c.textBox.EmptyCharBuffer()
			c.window.SetActive(false)
		} else if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			c.window.SetActive(false)
		}
	}
}

func (c *CmdPanel) parseCommand(input string) {
	if len(input) == 0 {
		err := SignalError{
			Kind: editorWarning,
			Msg:  "Unknown command",
		}
		FireSignal(EditorErrorRaised, err)
		return
	}
	tokens := strings.Split(input, " ")

	switch tokens[0] {
	case ":openproject":
		if len(tokens) != 2 {
			err := SignalError{
				Kind: editorError,
				Msg:  "Invalid arguments for command ':openproject'",
			}
			FireSignal(EditorErrorRaised, err)
			return
		}
		FireSignal(EditorProjectOpened, SignalString(tokens[1]))
	case ":openprojectfile":
		if len(tokens) != 2 {
			err := SignalError{
				Kind: editorError,
				Msg:  "Invalid arguments for command ':openprojectfile'",
			}
			FireSignal(EditorErrorRaised, err)
			return
		}
	default:
		err := SignalError{
			Kind: editorWarning,
			Msg:  "Unknown command",
		}
		FireSignal(EditorErrorRaised, err)
	}
}

func (c *CmdPanel) OnButtonPressed(w ui.Widget, id ui.ButtonID) {
	c.window.SetActive(false)
}
