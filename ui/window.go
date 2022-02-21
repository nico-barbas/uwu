package ui

type Window struct {
	handle     WinHandle
	zIndex     int
	Active     bool
	Rect       Rectangle
	activeRect Rectangle
	Style      Style
	Background Background
	widgets    WidgetList

	HasBorders  bool
	BorderWidth float64
	BorderColor Color

	HasHeader        bool
	HeaderHeight     float64
	HeaderBackground Background
	headerRect       Rectangle
	HasHeaderTitle   bool
	HeaderTitle      string
	HeaderFont       Font
	HeaderFontSize   float64
	HeaderFontClr    Color
	headerTitlePos   Point

	MinimizeBtn Button
	CloseBtn    Button
}

func (win *Window) initWindow() {
	// if (win.HasHeaderBtn || win.HasHeaderTitle) && !win.HasHeader {
	// 	// What is the best behavior here? Should the UI force a header on the window?
	// 	// Or should it disable the Close button?
	// 	log.SetPrefix("[UI Error]: ")
	// 	log.Fatalln("Can not add a Close button on a headerless Window")
	// 	win.HasHeaderBtn = false
	// }
	if win.HasHeader {
		win.headerRect = Rectangle{
			X: win.Rect.X, Y: win.Rect.Y,
			Width: win.Rect.Width, Height: win.HeaderHeight,
		}
		win.activeRect = Rectangle{
			X: win.Rect.X, Y: win.Rect.Y + win.HeaderHeight,
			Width: win.Rect.Width, Height: win.Rect.Height - win.HeaderHeight,
		}
		if win.HasHeaderTitle {
			titleSize := win.HeaderFont.MeasureText(win.HeaderTitle, win.HeaderFontSize)
			win.headerTitlePos = Point{
				win.headerRect.X + (win.headerRect.Width/2 - titleSize[0]/2),
				win.headerRect.Y + (win.headerRect.Height/2 - titleSize[1]/2),
			}
		}
	} else {
		win.activeRect = win.Rect
	}
	win.widgets.initList(win.Style)
}

func (win *Window) update() {
	if !win.Active {
		return
	}

	focused := win.zIndex == 0
	win.widgets.updateWidgets(focused)
	win.MinimizeBtn.update(focused)
	win.CloseBtn.update(focused)
}

func (win *Window) draw(buf *renderBuffer) {
	if !win.Active {
		return
	}

	bgEntry := win.Background.entry(win.Rect)
	buf.addEntry(bgEntry)

	if win.HasHeader {
		hdrEntry := win.HeaderBackground.entry(win.headerRect)
		buf.addEntry(hdrEntry)
		if win.HasHeaderTitle {
			buf.addEntry(RenderEntry{
				Kind: RenderText,
				Rect: Rectangle{
					X:      win.headerTitlePos[0],
					Y:      win.headerTitlePos[1],
					Height: win.HeaderFontSize,
				},
				Clr:  win.HeaderFontClr,
				Font: win.HeaderFont,
				Text: win.HeaderTitle,
			})
		}
		win.MinimizeBtn.draw(buf)
		win.CloseBtn.draw(buf)
	}

	win.widgets.drawWidgets(buf)

	if win.HasBorders {
		buf.addEntry(RenderEntry{
			Kind: RenderRectangle,
			Rect: Rectangle{
				X: win.Rect.X, Y: win.Rect.Y,
				Width: win.BorderWidth, Height: win.Rect.Height,
			},
			Clr: win.BorderColor,
		})
		buf.addEntry(RenderEntry{
			Kind: RenderRectangle,
			Rect: Rectangle{
				X: win.Rect.X, Y: win.Rect.Y,
				Width: win.Rect.Width, Height: win.BorderWidth,
			},
			Clr: win.BorderColor,
		})
		buf.addEntry(RenderEntry{
			Kind: RenderRectangle,
			Rect: Rectangle{
				X: win.Rect.X + win.Rect.Width - win.BorderWidth, Y: win.Rect.Y,
				Width: win.BorderWidth, Height: win.Rect.Height,
			},
			Clr: win.BorderColor,
		})
		buf.addEntry(RenderEntry{
			Kind: RenderRectangle,
			Rect: Rectangle{
				X: win.Rect.X, Y: win.Rect.Y + win.Rect.Height - win.BorderWidth,
				Width: win.Rect.Width, Height: win.BorderWidth,
			},
			Clr: win.BorderColor,
		})
	}
}

func (w *Window) AddWidget(wgt Widget, length int) {
	w.widgets.addWidget(wgt, w.activeRect, length)
}

func (w *Window) RemainingLength() int {
	return w.widgets.getRemainingLen(w.activeRect)
}

func (w *Window) setCloseBtn(btn Button) {
	width := w.HeaderHeight - (w.Style.Margin[1] * 2)
	w.CloseBtn = btn
	w.CloseBtn.setRect(Rectangle{
		X:      w.Rect.X + w.Rect.Width - width - w.Style.Margin[0],
		Y:      w.Rect.Y + w.Style.Margin[1],
		Width:  width,
		Height: width - w.Style.Margin[0]*2,
	})
}

func (w *Window) setMinimizeBtn(btn Button) {
	width := w.HeaderHeight - (w.Style.Margin[1] * 2)
	w.MinimizeBtn = btn
	w.MinimizeBtn.setRect(Rectangle{
		X:      w.Rect.X + w.Rect.Width - width*2 - w.Style.Margin[0],
		Y:      w.Rect.Y + w.Style.Margin[1],
		Width:  width,
		Height: width - w.Style.Margin[0]*2,
	})
}

type WinHandle struct {
	id  int
	gen uint
}

func (h WinHandle) SetActive(active bool) {
	getWindow(h).Active = active
	switch active {
	case true:
		h.FocusWindow()
	case false:
		h.UnfocusWindow()
	}
}

func (h WinHandle) IsActive() bool {
	return getWindow(h).Active
}

func (h WinHandle) SetCloseBtn(btn Button) {
	getWindow(h).setCloseBtn(btn)
}

func (h WinHandle) SetMinimizeBtn(btn Button) {
	getWindow(h).setMinimizeBtn(btn)
}

func (h WinHandle) AddWidget(wgt Widget, length int) {
	getWindow(h).AddWidget(wgt, length)
}

func (h WinHandle) RemainingLength() int {
	return getWindow(h).RemainingLength()
}

func (h WinHandle) FocusWindow() {
	if getWindow(h).zIndex == 0 {
		return
	}
	for i := 0; i < ctx.count; i += 1 {
		win := ctx.actives[i]
		if win.handle.id == h.id && win.handle.gen == h.gen {
			win.zIndex = 0
		} else {
			win.zIndex += 1
		}
	}
}

func (h WinHandle) UnfocusWindow() {
	if getWindow(h).zIndex != 0 {
		return
	}
	for i := 0; i < ctx.count; i += 1 {
		win := ctx.actives[i]
		if win.handle.id == h.id && win.handle.gen == h.gen {
			win.zIndex = ctx.count - 1
		} else {
			win.zIndex -= 1
		}
	}
}
