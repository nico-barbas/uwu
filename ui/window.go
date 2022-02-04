package ui

import "log"

type Window struct {
	handle     Handle
	Active     bool
	Rect       Rectangle
	activeRect Rectangle
	Style      Style
	Background Background
	widgets    WidgetList

	HasHeader        bool
	HeaderHeight     float64
	HeaderBackground Background
	headerRect       Rectangle

	HasCloseBtn  bool
	CloseBtn     Background
	closeBtnRect Rectangle
}

func (win *Window) parent() Node {
	return nil
}

func (win *Window) initWindow() {
	if win.HasCloseBtn && !win.HasHeader {
		// What is the best behavior here? Should the UI force a header on the window?
		// Or should it disable the Close button?
		log.SetPrefix("[UI Error]: ")
		log.Fatalln("Can not add a Close button on a headerless Window")
		win.HasCloseBtn = false
	}
	if win.HasHeader {
		if win.HasCloseBtn {
			width := win.HeaderHeight - (win.Style.Margin[1] * 2)
			win.closeBtnRect = Rectangle{
				X:      win.Rect.X + win.Rect.Width - width - win.Style.Margin[0],
				Y:      win.Rect.Y + win.Style.Margin[1],
				Width:  width,
				Height: width,
			}
		}
		win.headerRect = Rectangle{
			X: win.Rect.X, Y: win.Rect.Y,
			Width: win.Rect.Width, Height: win.HeaderHeight,
		}
		win.activeRect = Rectangle{
			X: win.Rect.X, Y: win.Rect.Y + win.HeaderHeight,
			Width: win.Rect.Width, Height: win.Rect.Height - win.HeaderHeight,
		}
	} else {
		win.activeRect = win.Rect
	}
	win.widgets.initList(win.Style)
}

func (win *Window) draw(buf *renderBuffer) {
	bgEntry := win.Background.entry(win.Rect)
	buf.addEntry(bgEntry)
	if win.HasHeader {
		hdrEntry := win.HeaderBackground.entry(win.headerRect)
		buf.addEntry(hdrEntry)
		if win.HasCloseBtn {
			// btnEntry := RenderEntry{
			// 	Kind: RenderImage,
			// 	Rect: win.closeBtnRect,
			// 	Clr:  Color{255, 255, 255, 255},
			// 	Img:  win.CloseBtnImg,
			// }
			btnEntry := win.CloseBtn.entry(win.closeBtnRect)
			buf.addEntry(btnEntry)
		}
	}

	win.widgets.drawWidgets(buf)
}
