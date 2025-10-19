package logic

import (
	"fire/internal/core"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func Update(game *core.Game) {
	// Update character animations
	game.UpdateCharacterAnimation(game.Hero.CurrentState)
	if rl.IsKeyDown(rl.KeyUp) {
		if rl.IsKeyDown(rl.KeyLeft) {
			game.Hero.UpdateMovementDirection(core.UpLeft)
		} else if rl.IsKeyDown(rl.KeyRight) {
			game.Hero.UpdateMovementDirection(core.UpRight)
		} else {
			game.Hero.UpdateMovementDirection(core.Up)
		}
		game.Hero.CurrentState = core.Jumping
		game.Hero.UpdatePosition()
	} else if rl.IsKeyDown(rl.KeyLeft) || rl.IsKeyDown(rl.KeyRight) {
		if rl.IsKeyDown(rl.KeyLeft) {
			game.Hero.UpdateMovementDirection(core.Left)
		} else {
			game.Hero.UpdateMovementDirection(core.Right)
		}
		game.Hero.CurrentState = core.Running
		game.Hero.UpdatePosition()
	} else {
		game.Hero.CurrentState = core.Idle
	}
}
