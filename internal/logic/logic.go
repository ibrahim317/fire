package logic

import (
	"fire/internal/core"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func Update(game *core.Game) {
	// Update character animations
	game.UpdateCharacterAnimation(game.Hero.CurrentState)

	if rl.IsKeyDown(rl.KeyLeft) || rl.IsKeyDown(rl.KeyRight) {
		game.Hero.UpdatePosition()
		if rl.IsKeyDown(rl.KeyLeft) {
			game.Hero.UpdateFacingDirection(core.Left)
		} else {
			game.Hero.UpdateFacingDirection(core.Right)
		}
		game.Hero.CurrentState = core.Running
	} else {
		game.Hero.CurrentState = core.Idle
	}
}
