package render

import (
	"fire/internal/core"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func Draw(game *core.Game) {
	rl.BeginDrawing()
	rl.DrawTextureEx(game.Bg, rl.Vector2{X: 0, Y: 0}, 0.0,
		2.7,
		rl.White)
	rl.DrawTexture(game.GrassBlock, 0, 558, rl.White)
	rl.DrawTexture(game.GrassBlock, 70, 558, rl.White)
	rl.DrawTexture(game.GrassBlock, 140, 558, rl.White)

	rl.DrawTexture(game.GrassBlock, 320, 500, rl.White)
	rl.DrawTexture(game.GrassBlock, 260, 500, rl.White)
	rl.DrawTexture(game.GrassBlock, 280, 500, rl.White)
	// Draw character with animation
	game.DrawCharacter(game.Hero.CurrentState)

	rl.EndDrawing()
}
