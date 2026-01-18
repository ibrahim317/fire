package main

import (
	"fire/internal/core"
	"fire/internal/systems"

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
					// Initialize ECS world for gameplay
					initGameWorld(&game)
				}
				game.Mode = newMode
			}

		case core.ModeGame:
			// Handle escape to return to menu
			if rl.IsKeyPressed(rl.KeyEscape) {
				game.Mode = core.ModeMainMenu
				continue
			}

			// Run ECS systems
			dt := rl.GetFrameTime()
			game.World.Update(dt)

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

// initGameWorld creates and initializes the ECS world for gameplay.
func initGameWorld(game *core.Game) {
	// Create new world
	game.InitWorld()

	// Spawn player entity
	core.SpawnPlayer(game.World, game)

	// Spawn mob entity
	core.SpawnMob(game.World, game, 500, 450)

	// Load and spawn map tiles
	core.LoadAndSpawnMap(game.World, game.GrassTile)

	// Register systems in execution order
	game.World.AddSystem(systems.NewInputSystem())
	game.World.AddSystem(systems.NewPhysicsSystem())
	game.World.AddSystem(systems.NewCollisionSystem())
	game.World.AddSystem(systems.NewAnimationSystem())
	game.World.AddSystem(systems.NewRenderSystem(systems.RenderConfig{
		Background:       game.Bg,
		GrassTile:        game.GrassTile,
		HealthHeart:      game.HealthHeart,
		HighlightBorders: &game.HighlightBorders,
	}))
}
