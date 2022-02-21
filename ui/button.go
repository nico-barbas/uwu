package ui

import (
	"log"
)

type (
	Button struct {
		widgetRoot
		UserID       ButtonID
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

		HasIcon bool
		Icon    Image
		IconClr Color
		// iconRect Rectangle
	}

	ButtonID int
)

type ButtonReceiver interface {
	OnButtonPressed(w Widget, id ButtonID)
}

func (btn *Button) init() {
	btn.Background.Clr = btn.Clr
}

func (btn *Button) update(parentFocused bool) {
	mPos := mousePosition()
	mLeft := isMousePressed()
	released := isMouseJustReleased()
	if btn.rect.pointInBounds(mPos) {
		btn.Background.Clr = btn.HighlightClr
		if released {
			if btn.Receiver != nil {
				btn.Receiver.OnButtonPressed(btn, btn.UserID)
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

	if btn.HasText {
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
	if btn.HasIcon {
		iconRect := Rectangle{
			X: btn.rect.X + (btn.rect.Width/2 - btn.Icon.GetWidth()/2),
			Y: btn.rect.Y + (btn.rect.Height/2 - btn.Icon.GetHeight()/2),
		}
		buf.addEntry(RenderEntry{
			Kind: RenderImage,
			Rect: iconRect,
			Img:  btn.Icon,
			Clr:  btn.IconClr,
		})
	}
}
