package main

import (
	"image/color"
	"log"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/nico-ec/feelsgood/ui"
)

var (
	uwuBackgroundClr = ui.Color{247, 231, 230, 255}
	uwuTextClr       = ui.Color{255, 95, 131, 255}
)

type Image struct {
	rl.Texture2D
}

func (i *Image) GetWidth() float64 {
	return float64(i.Width)
}

func (i *Image) GetHeight() float64 {
	return float64(i.Height)
}

type Font struct {
	rl.Font
}

func (f *Font) MeasureText(t string, size float64) ui.Point {
	tSize := rl.MeasureTextEx(f.Font, t, float32(size), 0)
	return ui.Point{float64(tSize.X), float64(tSize.Y)}
}

func main() {
	log.SetFlags(0)
	log.SetFlags(log.Lshortfile)
	rl.SetConfigFlags(rl.FlagWindowUndecorated)
	rl.InitWindow(1600, 900, "UwU")
	rl.SetTargetFPS(60)

	uiHeader := Image{
		Texture2D: rl.LoadTexture("assets/uiHeader.png"),
	}

	uiLayout := Image{
		Texture2D: rl.LoadTexture("assets/uiLayout.png"),
	}
	// uiBtn := Image{
	// 	Texture2D: rl.LoadTexture("assets/uiBtn.png"),
	// }
	uiFont := Font{
		Font: rl.LoadFontEx("assets/monogram.ttf", 32, nil, 250),
	}
	rl.SetTextureFilter(uiFont.Font.Texture, rl.FilterPoint)

	ctx := ui.NewContext()
	ui.MakeContextCurrent(ctx)

	hdl := ui.AddWindow(
		ui.Window{
			Active: true,
			Rect:   ui.Rectangle{0, 0, 1600, 900},
			Style: ui.Style{
				Ordering: ui.StyleOrderColumn,
				Padding:  0,
				Margin:   ui.Point{0, 0},
			},
			Background: ui.Background{
				Visible: true,
				Kind:    ui.BackgroundSolidColor,
				Clr:     uwuBackgroundClr,
			},
			// HasBorders:   true,
			// BorderWidth:  2,
			// BorderColor:  ui.Color{232, 152, 168, 255},
			HasHeader:    true,
			HeaderHeight: 25,
			HeaderBackground: ui.Background{
				Visible: true,
				Kind:    ui.BackgroundImageSlice,
				Clr:     ui.Color{232, 152, 168, 255},
				Img:     &uiHeader,
				Constr:  ui.Constraint{2, 2, 2, 2},
			},
			HasHeaderTitle: true,
			HeaderTitle:    "UwU",
			HeaderFont:     &uiFont,
			HeaderFontSize: 16,
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
	// lyt := ui.AddWidget(hdl, &ui.Layout{
	// 	Background: ui.Background{
	// 		Visible: true,
	// 		Kind:    ui.BackgroundImageSlice,
	// 		Clr:     ui.Color{232, 152, 168, 255},
	// 		Img:     &uiLayout,
	// 		Constr:  ui.Constraint{2, 2, 2, 2},
	// 	},
	// 	Style: ui.Style{
	// 		Ordering: ui.StyleOrderRow,
	// 		Padding:  3,
	// 		Margin:   ui.Point{5, 5},
	// 	},
	// }, 80)
	// ui.AddWidget(lyt, &ui.Label{
	// 	Background: ui.Background{
	// 		Visible: true,
	// 		Kind:    ui.BackgroundImageSlice,
	// 		Clr:     ui.Color{255, 255, 255, 255},
	// 		Img:     &uiBtn,
	// 		Constr:  ui.Constraint{7, 7, 7, 7},
	// 	},
	// 	Font: &uiFont,
	// 	Text: "Hello",
	// 	Clr:  uwuTextClr,
	// 	Size: 16,
	// }, 20)

	// ui.AddWidget(lyt, &ui.Label{
	// 	Background: ui.Background{
	// 		Visible: true,
	// 		Kind:    ui.BackgroundImageSlice,
	// 		Clr:     ui.Color{255, 255, 255, 255},
	// 		Img:     &uiBtn,
	// 		Constr:  ui.Constraint{7, 7, 7, 7},
	// 	},
	// 	Font: &uiFont,
	// 	Text: "World",
	// 	Clr:  uwuTextClr,
	// 	Size: 16,
	// }, 20)

	// ui.AddWidget(lyt, &ui.Button{
	// 	Background: ui.Background{
	// 		Visible: true,
	// 		Kind:    ui.BackgroundImageSlice,
	// 		Img:     &uiBtn,
	// 		Constr:  ui.Constraint{7, 7, 7, 7},
	// 	},
	// 	Clr:          ui.Color{255, 255, 255, 255},
	// 	HighlightClr: ui.Color{255, 0, 255, 255},
	// 	PressedClr:   ui.Color{255, 255, 0, 255},
	// 	HasText:      true,
	// 	Font:         &uiFont,
	// 	Text:         "!",
	// 	TextClr:      uwuTextClr,
	// 	TextSize:     16,
	// }, 20)

	tree := &ui.List{
		Background: ui.Background{
			Visible: true,
			Kind:    ui.BackgroundImageSlice,
			Clr:     ui.Color{232, 152, 168, 255},
			Img:     &uiLayout,
			Constr:  ui.Constraint{2, 2, 2, 2},
		},
		// 	Style: ui.Style{
		// 		Ordering: ui.StyleOrderRow,
		// 		Padding:  3,
		// 		Margin:   ui.Point{5, 5},
		// 	},
		// }

		Name:       "Root",
		Font:       &uiFont,
		TextSize:   14,
		TextClr:    uwuTextClr,
		IndentSize: 10,
	}
	ui.AddWidget(hdl, tree, 140)
	tree.AddItem(&ui.ListItem{Name: "file1"})
	tree.AddItem(&ui.ListItem{Name: "file2"})
	tree.AddItem(&ui.ListItem{Name: "file3"})
	tree.AddItem(&ui.ListItem{Name: "file4"})
	tree.AddItem(&ui.ListItem{Name: "file5"})

	ui.AddWidget(hdl, &ui.TextBox{
		Background: ui.Background{
			Visible: false,
			// Kind:    ui.BackgroundImageSlice,
			// Clr:     ui.Color{232, 152, 168, 255},
			// Img:     &uiHeader,
			// Constr:  ui.Constraint{2, 2, 2, 2},
		},
		Cap:      500,
		Margin:   10,
		Font:     &uiFont,
		TextClr:  uwuTextClr,
		TextSize: 16,
	}, ui.FitContainer)

	for !rl.WindowShouldClose() {
		key := rl.GetCharPressed()
		for key > 0 {
			ctx.AppendCharPressed(key)
			key = rl.GetCharPressed()
		}
		mpos := rl.GetMousePosition()
		mleft := rl.IsMouseButtonDown(rl.MouseLeftButton)
		ctx.UpdateUI(ui.Input{
			MPos:  ui.Point{float64(mpos.X), float64(mpos.Y)},
			MLeft: mleft,
			Enter: rl.IsKeyPressed(rl.KeyEnter) || rl.IsKeyPressed(rl.KeyKpEnter),
			Del:   rl.IsKeyDown(rl.KeyBackspace),
			Left:  rl.IsKeyDown(rl.KeyLeft),
			Right: rl.IsKeyDown(rl.KeyRight),
		})

		rl.BeginDrawing()

		rl.ClearBackground(rl.Black)

		uiBuf := ctx.DrawUI()
		for _, e := range uiBuf {
			switch e.Kind {
			case ui.RenderText:
				font := e.Font.(*Font)
				rl.DrawTextEx(
					font.Font,
					e.Text,
					rl.Vector2{float32(e.Rect.X), float32(e.Rect.Y)},
					float32(e.Rect.Height),
					0,
					color.RGBA{e.Clr[0], e.Clr[1], e.Clr[2], e.Clr[3]},
				)
			case ui.RenderRectangle:
				rl.DrawRectangleRec(
					rl.Rectangle{
						float32(e.Rect.X), float32(e.Rect.Y),
						float32(e.Rect.Width), float32(e.Rect.Height),
					},
					color.RGBA{e.Clr[0], e.Clr[1], e.Clr[2], e.Clr[3]},
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

				img := e.Img.(*Image)
				for i := 0; i < 9; i += 1 {
					rl.DrawTexturePro(
						img.Texture2D,
						rl.Rectangle{float32(srcRects[i].X), float32(srcRects[i].Y), float32(srcRects[i].Width), float32(srcRects[i].Height)},
						rl.Rectangle{float32(dstRects[i].X), float32(dstRects[i].Y), float32(dstRects[i].Width), float32(dstRects[i].Height)},
						rl.Vector2{}, 0,
						color.RGBA{e.Clr[0], e.Clr[1], e.Clr[2], e.Clr[3]},
					)
				}
			default:
			}
		}

		rl.EndDrawing()
	}

	rl.CloseWindow()
}
