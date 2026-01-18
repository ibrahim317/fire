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
	rl.SetExitKey(0) // Disable the default ESC exit key
	defer rl.CloseWindow()

	// Load assets and set up cleanup
	game.LoadAssets()
	defer game.UnloadAssets()

	// Initialize UI components (after window creation)
	game.InitUI()

	// Load map
	game.Map = core.InitMap(&game)


	// Game loop
	for !rl.WindowShouldClose() {
		switch game.Mode {
		case core.ModeMainMenu:
			rl.BeginDrawing()
			game.MainMenu.Draw(game.Bg)
			rl.EndDrawing()

			newMode := game.MainMenu.Update()
			if newMode != core.ModeMainMenu {
				// Transitioning to a new mode
				if newMode == core.ModeGame {
					// Always reload map when entering game mode to see design changes
					game.ReloadMap()
					game.ResetHeroPosition()
				}
				game.Mode = newMode
			}

		case core.ModeGame:
			// Handle escape to return to menu
			if rl.IsKeyPressed(rl.KeyEscape) {
				game.Mode = core.ModeMainMenu
				continue
			}

			logic.Update(&game)
			render.Draw(&game)

		case core.ModeDesigner:
			rl.BeginDrawing()
			game.Designer.Draw(&game)
			rl.EndDrawing()

			newMode := game.Designer.Update(&game)
			if newMode != core.ModeDesigner {
				game.Mode = newMode
			}

		case core.ModeSettings:
			rl.BeginDrawing()
			game.Settings.Draw(game.Bg)
			rl.EndDrawing()

			newMode := game.Settings.Update()
			if newMode != core.ModeSettings {
				game.Mode = newMode
			}
		}
	}
}
