package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/nico-ec/feelsgood/ui"
)

func main() {
	ctx := ui.NewContext()
	ui.MakeContextCurrent(ctx)

	ui.AddWindow(ui.Window{
		Active: true,
		Rect:   ui.Rectangle{100, 100, 100, 100},
		Background: ui.Background{
			Visible: true,
			Kind:    ui.BackgroundSolidColor,
			Clr:     ui.Color{255, 255, 0, 255},
		},
	})

	rl.InitWindow(800, 450, "raylib [core] example - basic window")

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		rl.DrawText("Congrats! You created your first window!", 190, 200, 20, rl.LightGray)

		rl.EndDrawing()
	}

	rl.CloseWindow()
}
