package main

import (
	"image/color"
	"log"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/nico-ec/feelsgood/ui"
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
	rl.InitWindow(1600, 800, "Persephone")
	rl.SetTargetFPS(60)

	uiPatch := Image{
		Texture2D: rl.LoadTexture("assets/uiPatch.png"),
	}
	uiFont := Font{
		Font: rl.LoadFontEx("assets/monogram.ttf", 32, nil, 250),
	}
	rl.SetTextureFilter(uiFont.Font.Texture, rl.FilterPoint)

	ctx := ui.NewContext()
	ui.MakeContextCurrent(ctx)

	hdl := ui.AddWindow(
		ui.Window{
			Active: true,
			Rect:   ui.Rectangle{300, 100, 350, 200},
			Style: ui.Style{
				Ordering: ui.StyleOrderColumn,
				Padding:  3,
				Margin:   ui.Point{3, 3},
			},
			Background: ui.Background{
				Visible: true,
				Kind:    ui.BackgroundImageSlice,
				Clr:     ui.Color{255, 255, 255, 255},
				Img:     &uiPatch,
				Constr:  ui.Constraint{2, 2, 2, 2},
			},

			HasHeader:    true,
			HeaderHeight: 18,
			HeaderBackground: ui.Background{
				Visible: true,
				Kind:    ui.BackgroundImageSlice,
				Clr:     ui.Color{255, 142, 176, 255},
				Img:     &uiPatch,
				Constr:  ui.Constraint{2, 2, 2, 2},
			},
			HasCloseBtn: true,
			CloseBtn: ui.Background{
				Visible: true,
				Kind:    ui.BackgroundImageSlice,
				Clr:     ui.Color{255, 255, 255, 255},
				Img:     &uiPatch,
				Constr:  ui.Constraint{2, 2, 2, 2},
			},
		},
	)
	lyt := ui.AddWidget(hdl, &ui.Layout{
		Background: ui.Background{
			Visible: true,
			Kind:    ui.BackgroundImageSlice,
			Clr:     ui.Color{198, 56, 34, 255},
			Img:     &uiPatch,
			Constr:  ui.Constraint{2, 2, 2, 2},
		},
		Style: ui.Style{
			Ordering: ui.StyleOrderRow,
			Padding:  3,
			Margin:   ui.Point{5, 5},
		},
	}, 40)
	ui.AddWidget(lyt, &ui.Label{
		Background: ui.Background{
			Visible: true,
			Kind:    ui.BackgroundImageSlice,
			Clr:     ui.Color{255, 255, 255, 255},
			Img:     &uiPatch,
			Constr:  ui.Constraint{2, 2, 2, 2},
		},
		Font: &uiFont,
		Text: "Hello",
		Clr:  ui.Color{0, 0, 0, 255},
		Size: 12,
	}, 20)

	ui.AddWidget(lyt, &ui.Label{
		Background: ui.Background{
			Visible: true,
			Kind:    ui.BackgroundImageSlice,
			Clr:     ui.Color{255, 255, 255, 255},
			Img:     &uiPatch,
			Constr:  ui.Constraint{2, 2, 2, 2},
		},
		Font: &uiFont,
		Text: "World",
		Clr:  ui.Color{0, 0, 0, 255},
		Size: 12,
	}, 20)

	ui.AddWidget(lyt, &ui.Button{
		Background: ui.Background{
			Visible: true,
			Kind:    ui.BackgroundImageSlice,
			Img:     &uiPatch,
			Constr:  ui.Constraint{2, 2, 2, 2},
		},
		Clr:          ui.Color{255, 255, 255, 255},
		HighlightClr: ui.Color{255, 0, 255, 255},
		PressedClr:   ui.Color{255, 255, 0, 255},
		HasText:      true,
		Font:         &uiFont,
		Text:         "!",
		TextClr:      ui.Color{0, 0, 0, 255},
		TextSize:     12,
	}, 20)

	ui.AddWidget(hdl, &ui.Label{
		Background: ui.Background{
			Visible: true,
			Kind:    ui.BackgroundImageSlice,
			Clr:     ui.Color{198, 56, 34, 255},
			Img:     &uiPatch,
			Constr:  ui.Constraint{2, 2, 2, 2},
		},
		Font: &uiFont,
		Text: "CANVAS GOES HERE!!",
		Clr:  ui.Color{0, 0, 0, 255},
		Size: 24,
	}, ui.FitContainer)

	for !rl.WindowShouldClose() {
		mpos := rl.GetMousePosition()
		mleft := rl.IsMouseButtonDown(rl.MouseLeftButton)
		ctx.UpdateUI(ui.Point{float64(mpos.X), float64(mpos.Y)}, mleft)

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
