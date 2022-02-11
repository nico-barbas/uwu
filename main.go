package main

import (
	"fmt"
	"image"
	_ "image/png"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/nico-ec/uwu/ui"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var (
	uwuBackgroundClr = ui.Color{247, 231, 230, 255}
	uwuTextClr       = ui.Color{255, 95, 131, 255}
	uwuKeywordClr    = ui.Color{200, 106, 255, 255}
	uwuDigitClr      = ui.Color{213, 133, 128, 255}
)

func main() {
	ebiten.SetWindowSize(1600, 900)
	ebiten.SetWindowDecorated(false)
	ebiten.SetRunnableOnUnfocused(false)

	uwu := new(uwu)
	uwu.ui = ui.NewContext()
	ui.MakeContextCurrent(uwu.ui)
	uwu.font = NewFont("assets/CozetteVector.ttf", 72, []int{14})

	i, _, err := ebitenutil.NewImageFromFile("assets/uiHeader.png")
	if err != nil {
		panic(err)
	}
	uwu.header = img{
		data: i,
	}

	i, _, err = ebitenutil.NewImageFromFile("assets/uiLayout.png")
	if err != nil {
		panic(err)
	}
	uwu.layout = img{
		data: i,
	}
	hdl := ui.AddWindow(
		ui.Window{
			Active: true,
			Rect:   ui.Rectangle{0, 0, 1600, 900},
			Style: ui.Style{
				Ordering: ui.StyleOrderRow,
				Padding:  0,
				Margin:   ui.Point{0, 0},
			},
			Background: ui.Background{
				Visible: true,
				Kind:    ui.BackgroundSolidColor,
				Clr:     uwuBackgroundClr,
			},
			HasHeader:    true,
			HeaderHeight: 25,
			HeaderBackground: ui.Background{
				Visible: true,
				Kind:    ui.BackgroundImageSlice,
				Clr:     ui.Color{232, 152, 168, 255},
				Img:     &uwu.header,
				Constr:  ui.Constraint{2, 2, 2, 2},
			},
			HasHeaderTitle: true,
			HeaderTitle:    "UwU",
			HeaderFont:     &uwu.font,
			HeaderFontSize: 14,
			HeaderFontClr:  uwuTextClr,
			// HasCloseBtn:    true,
			// CloseBtn: ui.Background{
			// 	Visible: true,
			// 	Kind:    ui.BackgroundImageSlice,
			// 	Clr:     ui.Color{255, 255, 255, 255},
			// 	Img:     &uiBtn,
			// 	Constr:  ui.Constraint{2, 2, 2, 2},
			// },
		},
	)

	rem := ui.ContainerRemainingLength(hdl)
	lyt := ui.AddWidget(hdl, &ui.Layout{
		Background: ui.Background{
			Visible: false,
		},
		Style: ui.Style{
			Ordering: ui.StyleOrderColumn,
			Padding:  0,
			Margin:   ui.Point{0, 0},
		},
	}, rem-20)
	tree := &ui.List{
		Background: ui.Background{
			Visible: true,
			Kind:    ui.BackgroundImageSlice,
			Clr:     ui.Color{232, 152, 168, 255},
			Img:     &uwu.layout,
			Constr:  ui.Constraint{2, 2, 2, 2},
		},
		Style: ui.Style{
			Padding: 3,
			Margin:  ui.Point{5, 0},
		},

		Name:       "Root",
		Font:       &uwu.font,
		TextSize:   14,
		TextClr:    uwuTextClr,
		IndentSize: 10,
	}
	ui.AddWidget(lyt, tree, 140)
	subFolder := ui.NewSubList("subFolder")
	tree.AddItem(&subFolder)
	subFolder.AddItem(&ui.ListItem{Name: "file1"})
	subFolder.AddItem(&ui.ListItem{Name: "file2"})
	subFolder.AddItem(&ui.ListItem{Name: "file3"})
	tree.AddItem(&ui.ListItem{Name: "file1"})
	tree.AddItem(&ui.ListItem{Name: "file2"})
	tree.AddItem(&ui.ListItem{Name: "file3"})
	tree.AddItem(&ui.ListItem{Name: "file4"})
	tree.AddItem(&ui.ListItem{Name: "file5"})

	// Text editor
	editor := &ui.TextBox{
		Background: ui.Background{
			Visible: false,
		},
		Cap:                500,
		Margin:             10,
		Font:               &uwu.font,
		TextSize:           14,
		HasRuler:           true,
		HasSyntaxHighlight: true,
	}
	editor.SetLexKeywords([]string{
		"func",
	})
	editor.SetSyntaxColors(ui.ColorStyle{
		Normal:  uwuTextClr,
		Keyword: uwuKeywordClr,
		Digit:   uwuDigitClr,
	})
	ui.AddWidget(lyt, editor, ui.FitContainer)

	// Status bar
	ui.AddWidget(hdl, &ui.Layout{
		Background: ui.Background{
			Kind: ui.BackgroundSolidColor,
			Clr:  uwuTextClr,
		},
	}, ui.FitContainer)

	if err := ebiten.RunGame(uwu); err != nil {
		panic(err)
	}
}

type img struct {
	data *ebiten.Image
}

func (i *img) GetWidth() float64 {
	return float64(i.data.Bounds().Dx())
}

func (i *img) GetHeight() float64 {
	return float64(i.data.Bounds().Dy())
}

type fnt struct {
	faces map[int]font.Face
}

func NewFont(path string, dpi float64, sizes []int) fnt {
	f := fnt{
		faces: make(map[int]font.Face, len(sizes)),
	}

	fontData, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	tt, err := opentype.Parse(fontData)
	if err != nil {
		panic(err)
	}

	for _, v := range sizes {
		face, err := opentype.NewFace(tt, &opentype.FaceOptions{
			Size:    float64(v),
			DPI:     dpi,
			Hinting: font.HintingNone,
		})
		if err != nil {
			panic(err)
		}
		f.faces[v] = face
	}
	return f
}

func (f *fnt) MeasureText(t string, size float64) ui.Point {
	measure := ui.Point{}

	if v, exist := f.faces[int(size)]; !exist {
		panic("No face of size in given Font")
	} else {
		r := text.BoundString(v, t)
		measure[0] = float64(r.Dx())
		measure[1] = float64(r.Dy())
	}

	return measure
}

type uwu struct {
	ui     *ui.Context
	font   fnt
	header img
	layout img
}

func (uwu *uwu) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return fmt.Errorf("closing editor")
	}
	var runes []rune
	runes = ebiten.AppendInputChars(runes[:0])
	for _, r := range runes {
		uwu.ui.AppendCharPressed(r)
		// key = rl.GetCharPressed()
	}
	mx, my := ebiten.CursorPosition()
	mleft := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	uwu.ui.UpdateUI(ui.Input{
		MPos:  ui.Point{float64(mx), float64(my)},
		MLeft: mleft,
		Enter: ebiten.IsKeyPressed(ebiten.KeyEnter) || ebiten.IsKeyPressed(ebiten.KeyKPEnter),
		Del:   ebiten.IsKeyPressed(ebiten.KeyBackspace),
		Left:  ebiten.IsKeyPressed(ebiten.KeyLeft),
		Right: ebiten.IsKeyPressed(ebiten.KeyRight),
		Up:    ebiten.IsKeyPressed(ebiten.KeyUp),
		Down:  ebiten.IsKeyPressed(ebiten.KeyDown),
	})
	return nil
}

func (uwu *uwu) Draw(screen *ebiten.Image) {
	uiBuf := uwu.ui.DrawUI()
	for _, e := range uiBuf {
		switch e.Kind {
		case ui.RenderText:
			font := e.Font.(*fnt)
			t := font.MeasureText(e.Text, e.Rect.Height)
			text.Draw(
				screen,
				e.Text,
				font.faces[int(e.Rect.Height)],
				int(e.Rect.X),
				int(e.Rect.Y+t[1]),
				e.Clr,
			)
		case ui.RenderRectangle:
			ebitenutil.DrawRect(
				screen,
				e.Rect.X, e.Rect.Y,
				e.Rect.Width, e.Rect.Height,
				e.Clr,
			)

		case ui.RenderImageSlice:
			dstRects := [9]ui.Rectangle{}
			srcRects := [9]ui.Rectangle{}

			l := e.Constr.Left
			r := e.Constr.Right
			u := e.Constr.Up
			d := e.Constr.Down

			imgW := e.Img.GetWidth()
			imgH := e.Img.GetHeight()

			srcX0 := float64(0)
			srcX1 := l
			srcX2 := imgW - r

			srcY0 := float64(0)
			srcY1 := u
			srcY2 := imgH - d

			dstL := l
			dstR := r
			dstU := u
			dstD := d

			// if scale > 0 {
			// 	dstL *= scale
			// 	dstR *= scale
			// 	dstU *= scale
			// 	dstD *= scale
			// }

			dstX0 := e.Rect.X
			dstX1 := e.Rect.X + dstL
			dstX2 := e.Rect.X + e.Rect.Width - dstR

			dstY0 := e.Rect.Y
			dstY1 := e.Rect.Y + dstU
			dstY2 := e.Rect.Y + e.Rect.Height - dstD

			// TOP
			dstRects[0] = ui.Rectangle{X: dstX0, Y: dstY0, Width: dstL, Height: dstU}
			srcRects[0] = ui.Rectangle{X: srcX0, Y: srcY0, Width: l, Height: u}
			//
			dstRects[1] = ui.Rectangle{X: dstX1, Y: dstY0, Width: e.Rect.Width - (dstL + dstR), Height: dstU}
			srcRects[1] = ui.Rectangle{X: srcX1, Y: srcY0, Width: imgW - (l + r), Height: u}
			//
			dstRects[2] = ui.Rectangle{X: dstX2, Y: dstY0, Width: dstR, Height: dstU}
			srcRects[2] = ui.Rectangle{X: srcX2, Y: srcY0, Width: r, Height: u}
			//
			// MIDDLE
			dstRects[3] = ui.Rectangle{X: dstX0, Y: dstY1, Width: dstL, Height: e.Rect.Height - (dstU + dstD)}
			srcRects[3] = ui.Rectangle{X: srcX0, Y: srcY1, Width: l, Height: imgH - (u + d)}
			//
			dstRects[4] = ui.Rectangle{X: dstX1, Y: dstY1, Width: e.Rect.Width - (dstL + dstR), Height: e.Rect.Height - (dstU + dstD)}
			srcRects[4] = ui.Rectangle{X: srcX1, Y: srcY1, Width: imgW - (l + r), Height: imgH - (u + d)}
			//
			dstRects[5] = ui.Rectangle{X: dstX2, Y: dstY1, Width: dstR, Height: e.Rect.Height - (dstU + dstD)}
			srcRects[5] = ui.Rectangle{X: srcX2, Y: srcY1, Width: r, Height: imgH - (u + d)}
			//
			// BOTTOM
			dstRects[6] = ui.Rectangle{X: dstX0, Y: dstY2, Width: dstL, Height: dstD}
			srcRects[6] = ui.Rectangle{X: srcX0, Y: srcY2, Width: l, Height: d}
			//
			dstRects[7] = ui.Rectangle{X: dstX1, Y: dstY2, Width: e.Rect.Width - (dstL + dstR), Height: dstD}
			srcRects[7] = ui.Rectangle{X: srcX1, Y: srcY2, Width: imgW - (l + r), Height: d}
			//
			dstRects[8] = ui.Rectangle{X: dstX2, Y: dstY2, Width: dstR, Height: dstD}
			srcRects[8] = ui.Rectangle{X: srcX2, Y: srcY2, Width: r, Height: d}

			img := e.Img.(*img)
			for i := 0; i < 9; i += 1 {
				opt := &ebiten.DrawImageOptions{}
				opt.GeoM.Scale(
					dstRects[i].Width/srcRects[i].Width,
					dstRects[i].Height/srcRects[i].Height,
				)
				opt.GeoM.Translate(dstRects[i].X, dstRects[i].Y)
				if e.Clr[3] != 0 {
					r, g, b, a := e.Clr.RGBA()
					opt.ColorM.Scale(
						float64(r)/float64(a),
						float64(g)/float64(a),
						float64(b)/float64(a),
						float64(a)/0xffff,
					)
				}
				r := image.Rect(
					int(srcRects[i].X), int(srcRects[i].Y),
					int(srcRects[i].X+srcRects[i].Width), int(srcRects[i].Y+srcRects[i].Height),
				)
				screen.DrawImage(img.data.SubImage(r).(*ebiten.Image), opt)
			}
		}

	}
}

func (uwu *uwu) Layout(w, h int) (int, int) {
	return 1600, 900
}

// type Image struct {
// 	rl.Texture2D
// }

// func (i *Image) GetWidth() float64 {
// 	return float64(i.Width)
// }

// func (i *Image) GetHeight() float64 {
// 	return float64(i.Height)
// }

// type Font struct {
// 	rl.Font
// }

// func (f *Font) MeasureText(t string, size float64) ui.Point {
// 	tSize := rl.MeasureTextEx(f.Font, t, float32(size), 0)
// 	return ui.Point{float64(tSize.X), float64(tSize.Y)}
// }

// func main() {
// 	log.SetFlags(0)
// 	log.SetFlags(log.Lshortfile)
// 	rl.SetConfigFlags(rl.FlagWindowUndecorated)
// 	rl.InitWindow(1600, 900, "UwU")
// 	rl.SetTargetFPS(60)

// 	uiHeader := Image{
// 		Texture2D: rl.LoadTexture("assets/uiHeader.png"),
// 	}

// 	uiLayout := Image{
// 		Texture2D: rl.LoadTexture("assets/uiLayout.png"),
// 	}

// 	uiFont := Font{
// 		Font: rl.LoadFont("assets/monogram.fnt"),
// 	}
// 	// rl.SetTextureFilter(uiFont.Font.Texture, rl.FilterPoint)

// 	ctx := ui.NewContext()
// 	ui.MakeContextCurrent(ctx)

// 	hdl := ui.AddWindow(
// 		ui.Window{
// 			Active: true,
// 			Rect:   ui.Rectangle{0, 0, 1600, 900},
// 			Style: ui.Style{
// 				Ordering: ui.StyleOrderRow,
// 				Padding:  0,
// 				Margin:   ui.Point{0, 0},
// 			},
// 			Background: ui.Background{
// 				Visible: true,
// 				Kind:    ui.BackgroundSolidColor,
// 				Clr:     uwuBackgroundClr,
// 			},
// 			HasHeader:    true,
// 			HeaderHeight: 25,
// 			HeaderBackground: ui.Background{
// 				Visible: true,
// 				Kind:    ui.BackgroundImageSlice,
// 				Clr:     ui.Color{232, 152, 168, 255},
// 				Img:     &uiHeader,
// 				Constr:  ui.Constraint{2, 2, 2, 2},
// 			},
// 			HasHeaderTitle: true,
// 			HeaderTitle:    "UwU",
// 			HeaderFont:     &uiFont,
// 			HeaderFontSize: 12,
// 			HeaderFontClr:  uwuTextClr,
// 			// HasCloseBtn:    true,
// 			// CloseBtn: ui.Background{
// 			// 	Visible: true,
// 			// 	Kind:    ui.BackgroundImageSlice,
// 			// 	Clr:     ui.Color{255, 255, 255, 255},
// 			// 	Img:     &uiBtn,
// 			// 	Constr:  ui.Constraint{2, 2, 2, 2},
// 			// },
// 		},
// 	)

// 	rem := ui.ContainerRemainingLength(hdl)
// 	lyt := ui.AddWidget(hdl, &ui.Layout{
// 		Background: ui.Background{
// 			Visible: false,
// 		},
// 		Style: ui.Style{
// 			Ordering: ui.StyleOrderColumn,
// 			Padding:  0,
// 			Margin:   ui.Point{0, 0},
// 		},
// 	}, rem-20)
// 	tree := &ui.List{
// 		Background: ui.Background{
// 			Visible: true,
// 			Kind:    ui.BackgroundImageSlice,
// 			Clr:     ui.Color{232, 152, 168, 255},
// 			Img:     &uiLayout,
// 			Constr:  ui.Constraint{2, 2, 2, 2},
// 		},
// 		Style: ui.Style{
// 			Padding: 3,
// 			Margin:  ui.Point{5, 0},
// 		},

// 		Name:       "Root",
// 		Font:       &uiFont,
// 		TextSize:   12,
// 		TextClr:    uwuTextClr,
// 		IndentSize: 10,
// 	}
// 	ui.AddWidget(lyt, tree, 140)
// 	subFolder := ui.NewSubList("subFolder")
// 	tree.AddItem(&subFolder)
// 	subFolder.AddItem(&ui.ListItem{Name: "file1"})
// 	subFolder.AddItem(&ui.ListItem{Name: "file2"})
// 	subFolder.AddItem(&ui.ListItem{Name: "file3"})
// 	tree.AddItem(&ui.ListItem{Name: "file1"})
// 	tree.AddItem(&ui.ListItem{Name: "file2"})
// 	tree.AddItem(&ui.ListItem{Name: "file3"})
// 	tree.AddItem(&ui.ListItem{Name: "file4"})
// 	tree.AddItem(&ui.ListItem{Name: "file5"})

// 	editor := &ui.TextBox{
// 		Background: ui.Background{
// 			Visible: false,
// 		},
// 		Cap:                500,
// 		Margin:             10,
// 		Font:               &uiFont,
// 		TextSize:           12,
// 		HasRuler:           true,
// 		HasSyntaxHighlight: true,
// 	}
// 	editor.SetLexKeywords([]string{
// 		"func",
// 	})
// 	editor.SetSyntaxColors(ui.ColorStyle{
// 		Normal:  uwuTextClr,
// 		Keyword: uwuKeywordClr,
// 		Digit:   uwuDigitClr,
// 	})
// 	ui.AddWidget(lyt, editor, ui.FitContainer)

// 	// Status bar
// 	ui.AddWidget(hdl, &ui.Layout{
// 		Background: ui.Background{
// 			Kind: ui.BackgroundSolidColor,
// 			Clr:  uwuTextClr,
// 		},
// 	}, ui.FitContainer)

// 	for !rl.WindowShouldClose() {
// 		key := rl.GetCharPressed()
// 		for key > 0 {
// 			ctx.AppendCharPressed(key)
// 			key = rl.GetCharPressed()
// 		}
// 		mpos := rl.GetMousePosition()
// 		mleft := rl.IsMouseButtonDown(rl.MouseLeftButton)
// 		ctx.UpdateUI(ui.Input{
// 			MPos:  ui.Point{float64(mpos.X), float64(mpos.Y)},
// 			MLeft: mleft,
// 			Enter: rl.IsKeyPressed(rl.KeyEnter) || rl.IsKeyPressed(rl.KeyKpEnter),
// 			Del:   rl.IsKeyDown(rl.KeyBackspace),
// 			Left:  rl.IsKeyDown(rl.KeyLeft),
// 			Right: rl.IsKeyDown(rl.KeyRight),
// 			Up:    rl.IsKeyDown(rl.KeyUp),
// 			Down:  rl.IsKeyDown(rl.KeyDown),
// 		})

// 		rl.BeginDrawing()

// 		rl.ClearBackground(rl.Black)

// 		uiBuf := ctx.DrawUI()
// 		for _, e := range uiBuf {
// 			switch e.Kind {
// 			case ui.RenderText:
// 				font := e.Font.(*Font)
// 				rl.DrawTextEx(
// 					font.Font,
// 					e.Text,
// 					rl.Vector2{float32(e.Rect.X), float32(e.Rect.Y)},
// 					float32(e.Rect.Height),
// 					0,
// 					color.RGBA{e.Clr[0], e.Clr[1], e.Clr[2], e.Clr[3]},
// 				)
// 			case ui.RenderRectangle:
// 				rl.DrawRectangleRec(
// 					rl.Rectangle{
// 						float32(e.Rect.X), float32(e.Rect.Y),
// 						float32(e.Rect.Width), float32(e.Rect.Height),
// 					},
// 					color.RGBA{e.Clr[0], e.Clr[1], e.Clr[2], e.Clr[3]},
// 				)
// 			case ui.RenderImageSlice:
// 				dstRects := [9]ui.Rectangle{}
// 				srcRects := [9]ui.Rectangle{}

// 				l := e.Constr.Left
// 				r := e.Constr.Right
// 				u := e.Constr.Up
// 				d := e.Constr.Down

// 				imgW := e.Img.GetWidth()
// 				imgH := e.Img.GetHeight()

// 				srcX0 := float64(0)
// 				srcX1 := l
// 				srcX2 := imgW - r

// 				srcY0 := float64(0)
// 				srcY1 := u
// 				srcY2 := imgH - d

// 				dstL := l
// 				dstR := r
// 				dstU := u
// 				dstD := d

// 				// if scale > 0 {
// 				// 	dstL *= scale
// 				// 	dstR *= scale
// 				// 	dstU *= scale
// 				// 	dstD *= scale
// 				// }

// 				dstX0 := e.Rect.X
// 				dstX1 := e.Rect.X + dstL
// 				dstX2 := e.Rect.X + e.Rect.Width - dstR

// 				dstY0 := e.Rect.Y
// 				dstY1 := e.Rect.Y + dstU
// 				dstY2 := e.Rect.Y + e.Rect.Height - dstD

// 				// TOP
// 				dstRects[0] = ui.Rectangle{X: dstX0, Y: dstY0, Width: dstL, Height: dstU}
// 				srcRects[0] = ui.Rectangle{X: srcX0, Y: srcY0, Width: l, Height: u}
// 				//
// 				dstRects[1] = ui.Rectangle{X: dstX1, Y: dstY0, Width: e.Rect.Width - (dstL + dstR), Height: dstU}
// 				srcRects[1] = ui.Rectangle{X: srcX1, Y: srcY0, Width: imgW - (l + r), Height: u}
// 				//
// 				dstRects[2] = ui.Rectangle{X: dstX2, Y: dstY0, Width: dstR, Height: dstU}
// 				srcRects[2] = ui.Rectangle{X: srcX2, Y: srcY0, Width: r, Height: u}
// 				//
// 				// MIDDLE
// 				dstRects[3] = ui.Rectangle{X: dstX0, Y: dstY1, Width: dstL, Height: e.Rect.Height - (dstU + dstD)}
// 				srcRects[3] = ui.Rectangle{X: srcX0, Y: srcY1, Width: l, Height: imgH - (u + d)}
// 				//
// 				dstRects[4] = ui.Rectangle{X: dstX1, Y: dstY1, Width: e.Rect.Width - (dstL + dstR), Height: e.Rect.Height - (dstU + dstD)}
// 				srcRects[4] = ui.Rectangle{X: srcX1, Y: srcY1, Width: imgW - (l + r), Height: imgH - (u + d)}
// 				//
// 				dstRects[5] = ui.Rectangle{X: dstX2, Y: dstY1, Width: dstR, Height: e.Rect.Height - (dstU + dstD)}
// 				srcRects[5] = ui.Rectangle{X: srcX2, Y: srcY1, Width: r, Height: imgH - (u + d)}
// 				//
// 				// BOTTOM
// 				dstRects[6] = ui.Rectangle{X: dstX0, Y: dstY2, Width: dstL, Height: dstD}
// 				srcRects[6] = ui.Rectangle{X: srcX0, Y: srcY2, Width: l, Height: d}
// 				//
// 				dstRects[7] = ui.Rectangle{X: dstX1, Y: dstY2, Width: e.Rect.Width - (dstL + dstR), Height: dstD}
// 				srcRects[7] = ui.Rectangle{X: srcX1, Y: srcY2, Width: imgW - (l + r), Height: d}
// 				//
// 				dstRects[8] = ui.Rectangle{X: dstX2, Y: dstY2, Width: dstR, Height: dstD}
// 				srcRects[8] = ui.Rectangle{X: srcX2, Y: srcY2, Width: r, Height: d}

// 				img := e.Img.(*Image)
// 				for i := 0; i < 9; i += 1 {
// 					rl.DrawTexturePro(
// 						img.Texture2D,
// 						rl.Rectangle{float32(srcRects[i].X), float32(srcRects[i].Y), float32(srcRects[i].Width), float32(srcRects[i].Height)},
// 						rl.Rectangle{float32(dstRects[i].X), float32(dstRects[i].Y), float32(dstRects[i].Width), float32(dstRects[i].Height)},
// 						rl.Vector2{}, 0,
// 						color.RGBA{e.Clr[0], e.Clr[1], e.Clr[2], e.Clr[3]},
// 					)
// 				}
// 			default:
// 			}
// 		}

// 		rl.EndDrawing()
// 	}

// 	rl.CloseWindow()
// }
