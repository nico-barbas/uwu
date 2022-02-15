package main

import (
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/nico-ec/uwu/editor"
)

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
