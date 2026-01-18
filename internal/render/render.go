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

	// Draw character with animation
	game.DrawCharacter()
	game.DrawMap()
	game.DrawHealth()
	game.DrawMob()
	rl.EndDrawing()
}
