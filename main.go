package main

import (
	"fire/internal/core"
	"fire/internal/logic"
	"fire/internal/render"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	var game core.Game

	// Init game state
	game.Init()

	rl.InitWindow(game.ScreenWidth, game.ScreenHeight, "Fire")
	rl.SetTargetFPS(core.FPS)
	defer rl.CloseWindow() // This will close the window when main exits

	// Load assets and set up cleanup
	game.LoadAssets()
	defer game.UnloadAssets()

	game.Map = core.InitMap(&game)
	// Game loop
	for !rl.WindowShouldClose() {
		logic.Update(&game)
		render.Draw(&game)
	}
}
