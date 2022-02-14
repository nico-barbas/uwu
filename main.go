package main

import (
	_ "image/png"

	// rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/nico-ec/uwu/editor"
)

// var (
// 	uwuBackgroundClr = ui.Color{247, 231, 230, 255}
// 	uwuTextClr       = ui.Color{255, 95, 131, 255}
// 	uwuKeywordClr    = ui.Color{200, 106, 255, 255}
// 	uwuDigitClr      = ui.Color{213, 133, 128, 255}
// )

func main() {
	ebiten.SetWindowSize(1600, 900)
	ebiten.SetWindowDecorated(false)
	ebiten.SetRunnableOnUnfocused(false)
	ebiten.SetMaxTPS(30)

	ed := editor.NewEditor()

	if err := ebiten.RunGame(ed); err != nil {
		panic(err)
	}
}
