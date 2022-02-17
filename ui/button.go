package ui

import (
	"log"
)

type Button struct {
	widgetRoot
	Background   Background
	Clr          Color
	HighlightClr Color
	PressedClr   Color

	Receiver ButtonReceiver
	pressed  bool

	HasText  bool
	Font     Font
	Text     string
	TextClr  Color
	TextSize float64
}

type ButtonReceiver interface {
	OnButtonPressed(w Widget)
}

func (btn *Button) init() {
	btn.Background.Clr = btn.Clr
}

func (btn *Button) update() {
	mPos := mousePosition()
	mLeft := isMousePressed()
	released := isMouseJustReleased()
	if btn.rect.pointInBounds(mPos) {
		btn.Background.Clr = btn.HighlightClr
		if released {
			if btn.Receiver != nil {
				btn.Receiver.OnButtonPressed(btn)
			} else {
				log.SetPrefix("[UI Debug]: ")
				log.Println("No Receiver attached to this button")
			}
			btn.pressed = false
		} else if mLeft {
			btn.pressed = true
		}
	} else {
		btn.Background.Clr = btn.Clr
		if released {
			btn.pressed = false
		}
	}

	if btn.pressed {
		btn.Background.Clr = btn.PressedClr
	} else if released {
		btn.Background.Clr = btn.Clr
	}
}

func (btn *Button) draw(buf *renderBuffer) {
	bgEntry := btn.Background.entry(btn.rect)
	buf.addEntry(bgEntry)

	textSize := btn.Font.MeasureText(btn.Text, btn.TextSize)
	textEntry := RenderEntry{
		Kind: RenderText,
		Rect: Rectangle{
			X:      btn.rect.X + (btn.rect.Width/2 - textSize[0]/2),
			Y:      btn.rect.Y + (btn.rect.Height/2 - textSize[1]/2),
			Height: btn.TextSize,
		},
		Clr:  btn.TextClr,
		Font: btn.Font,
		Text: btn.Text,
	}
	buf.addEntry(textEntry)
}
