package logic

import (
	"fire/internal/core"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func Update(game *core.Game) {
	// Update character animations
	game.UpdateCharacterAnimation(game.Hero.CurrentState)
	game.CheckCollisionWithMap()
	if rl.IsKeyDown(rl.KeyUp) {
		if rl.IsKeyDown(rl.KeyLeft) {
			game.Hero.UpdateMovementDirection(core.UpLeft)
			game.UpdateHeroAcceleration(rl.Vector2{X: -1, Y: -1.6})
		} else if rl.IsKeyDown(rl.KeyRight) {
			game.Hero.UpdateMovementDirection(core.UpRight)
			game.UpdateHeroAcceleration(rl.Vector2{X: 1, Y: -1})
		} else {
			game.Hero.UpdateMovementDirection(core.Up)
			game.UpdateHeroAcceleration(rl.Vector2{X: 0, Y: -1})
		}
		game.Hero.CurrentState = core.Jumping
	} else if rl.IsKeyDown(rl.KeyLeft) || rl.IsKeyDown(rl.KeyRight) {
		if rl.IsKeyDown(rl.KeyLeft) {
			game.Hero.UpdateVelocity(rl.Vector2{X: -2.5, Y: 0})
			game.Hero.UpdateMovementDirection(core.Left)
		} else {
			game.Hero.UpdateMovementDirection(core.Right)
			game.Hero.UpdateVelocity(rl.Vector2{X: 2.5, Y: 0})
		}
		game.Hero.CurrentState = core.Running
	} else {
		game.Hero.CurrentState = core.Idle
		game.Hero.UpdateVelocity(rl.Vector2{X: 0, Y: 0})
	}

	if !game.Hero.IsOnGround {
		game.Hero.CurrentState = core.Falling
	}
	game.UpdateHeroAcceleration(rl.Vector2{X: 0, Y: 0})
	game.Hero.UpdatePosition()
}
