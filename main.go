package main

import (
	"fmt"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/nico-ec/uwu/editor"
	"github.com/pkg/profile"
)

func main() {
	defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()
	ebiten.SetWindowSize(1600, 900)
	ebiten.SetWindowDecorated(false)
	ebiten.SetRunnableOnUnfocused(false)
	ebiten.SetMaxTPS(30)

	ed := editor.NewEditor()

	if err := ebiten.RunGame(ed); err != nil {
		fmt.Println(err)
	}
	// fmt.Println("Serving pprof on port 8080")
	// if err := http.ListenAndServe(":8080", nil); err != nil {
	// 	panic(err)
	// }
}
