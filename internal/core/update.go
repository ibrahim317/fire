package core

import (
	"image/color"
	"unsafe"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Updates the animation frame for a given character state
func (g *Game) UpdateCharacterAnimation(state CharacterState) {
	animData := g.Hero.States[state]
	if animData.FrameCount <= 1 {
		return // No animation needed for single frame
	}

	animData.FrameCounter++

	if animData.FrameCounter >= animData.FrameDelay {
		// Move to next frame
		animData.CurrentFrame++
		if animData.CurrentFrame >= animData.FrameCount {
			animData.CurrentFrame = 0
		}

		// Calculate memory offset for current frame
		nextFrameDataOffset := uint32(animData.Image.Width * animData.Image.Height * 4 * animData.CurrentFrame)

		// Update GPU texture with current frame data
		rl.UpdateTexture(animData.Texture,
			unsafe.Slice((*color.RGBA)(unsafe.Pointer(uintptr(animData.Image.Data)+uintptr(nextFrameDataOffset))), animData.FrameSize))

		animData.FrameCounter = 0
	}

	// Update the animation data in the map
	g.Hero.States[state] = animData
}

func (c *Character) UpdateMovementDirection(direction MovementDirection) {
	c.MovementDirection = direction
}

func (c *Character) UpdatePosition() {
	switch c.MovementDirection {
	case Left:
		c.Position.X = c.Position.X - 2.5
	case Right:
		c.Position.X = c.Position.X + 2.5
	case Up:
		c.Position.Y = c.Position.Y - 2.5
	case UpRight:
		c.Position.X = c.Position.X + 2.5
		c.Position.Y = c.Position.Y - 2.5
	case UpLeft:
		c.Position.X = c.Position.X - 2.5
		c.Position.Y = c.Position.Y - 2.5
	case Down:
		c.Position.Y = c.Position.Y + 2.5
	case DownRight:
		c.Position.X = c.Position.X + 2.5
		c.Position.Y = c.Position.Y + 2.5
	case DownLeft:
		c.Position.X = c.Position.X - 2.5
		c.Position.Y = c.Position.Y + 2.5
	}
}
