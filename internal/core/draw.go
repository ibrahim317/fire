package core

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// DrawCharacter draws the character with the current animation frame
func (g *Game) DrawCharacter(state CharacterState) {
	animData := g.Hero.States[state]
	frameWidth := float32(animData.Texture.Width)
	frameHeight := float32(animData.Texture.Height)
	pos := g.Hero.Position

	var sourceRec rl.Rectangle
	var destRec rl.Rectangle
	var origin rl.Vector2

	// Source rectangle for the texture (entire frame)
	sourceRec = rl.Rectangle{X: 0, Y: 0, Width: frameWidth, Height: frameHeight}
	// Destination rectangle at the character's position, scaling width to -1.7 * width for facing left (mirrored), or 1.7 * width for right
	if g.Hero.FacingDirection == Left {
		sourceRec = rl.Rectangle{X: 0, Y: 0, Width: -frameWidth, Height: frameHeight}
		destRec = rl.Rectangle{
			X:      pos.X,
			Y:      pos.Y,
			Width:  1.7 * frameWidth,
			Height: 1.7 * frameHeight,
		}
	} else {
		destRec = rl.Rectangle{
			X:      pos.X,
			Y:      pos.Y,
			Width:  1.7 * frameWidth,
			Height: 1.7 * frameHeight,
		}
	}

	origin = rl.Vector2{X: 0, Y: 0}

	rl.DrawTexturePro(animData.Texture, sourceRec, destRec, origin, 0, rl.White)
}
