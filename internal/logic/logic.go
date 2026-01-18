package logic

import (
	"fire/internal/core"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// common factor: ac+bc = c(a+b)

func Update(game *core.Game) {
	// Update character animations
	game.UpdateCharacterAnimation()
	game.UpdateMobAnimation()

	right := rl.IsKeyDown(rl.KeyRight)
	left := rl.IsKeyDown(rl.KeyLeft)
	up := rl.IsKeyDown(rl.KeyUp)

	// Handle input and update velocity
	if up {
		if left {
			game.HandleMovement(rl.Vector2{X: -1, Y: -1.6})
		} else if right {
			game.HandleMovement(rl.Vector2{X: 1, Y: -1.6})
		} else {
			game.HandleMovement(rl.Vector2{X: 0, Y: -1.6})
		}
		game.Hero.CurrentState = core.Jumping
	} else if left {
		game.HandleMovement(rl.Vector2{X: -1, Y: 0})
		game.Hero.CurrentState = core.Running
	} else if right {
		game.HandleMovement(rl.Vector2{X: 1, Y: 0})
		game.Hero.CurrentState = core.Running
	} else {
		game.Hero.CurrentState = core.Idle
		game.HandleMovement(rl.Vector2{X: 0, Y: 0})
	}

	// Move the hero first, then resolve collisions
	game.Hero.UpdatePosition()
	game.CheckCollisionWithMap()

	// Update acceleration and state based on collision results
	game.UpdateHeroAcceleration(rl.Vector2{X: 0, Y: 0})
	if !game.Hero.IsOnGround {
		game.Hero.CurrentState = core.Falling
	}
}
