package core

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// DrawCharacter draws the character with the current animation frame
func (g *Game) DrawCharacter() {
	animData := g.Hero.States[g.Hero.CurrentState]
	frameWidth := float32(animData.Texture.Width)
	frameHeight := float32(animData.Texture.Height)
	pos := g.Hero.Position

	var origin rl.Vector2 = rl.Vector2{X: 0, Y: 0}
	var sourceRec rl.Rectangle = rl.Rectangle{X: 0, Y: 0, Width: frameWidth, Height: frameHeight}
	var destRec rl.Rectangle = rl.Rectangle{
		X:      pos.X,
		Y:      pos.Y,
		Width:  g.HeroScaling * frameWidth,
		Height: g.HeroScaling * frameHeight,
	}

	if g.Hero.MovementDirection == Left ||
		g.Hero.MovementDirection == UpLeft ||
		g.Hero.MovementDirection == DownLeft {
		sourceRec = rl.Rectangle{X: 0, Y: 0, Width: -frameWidth, Height: frameHeight}
	}

	rl.DrawTexturePro(animData.Texture, sourceRec, destRec, origin, 0, rl.White)
	rl.DrawRectangleLines(int32(pos.X), int32(pos.Y), int32(destRec.Width), int32(destRec.Height), rl.Red)
}

func (g *Game) DrawMap() {
	for _, tile := range g.Map.Tiles {
		texture := g.textureForTile(tile.TileType)
		tileWidth := float32(texture.Width)
		tileHeight := float32(texture.Height)

		sourceRec := rl.Rectangle{X: 0, Y: 0, Width: tileWidth, Height: tileHeight}
		destRec := rl.Rectangle{X: tile.X, Y: tile.Y, Width: tileWidth, Height: tileHeight}
		rl.DrawTexturePro(texture, sourceRec, destRec, rl.Vector2{X: 0, Y: 0}, 0, rl.White)
	}
}

func (g *Game) textureForTile(tileType TileType) rl.Texture2D {
	switch tileType {
	case Stone, Water, Tree, Rock:
		// Placeholder until dedicated textures are added.
		return g.GrassTile
	default:
		return g.GrassTile
	}
}
