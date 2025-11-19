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
	if !c.IsOnGround {
		c.Position.Y += c.Velocity.Y
	}
	c.Position.X += c.Velocity.X
}

func (c *Character) UpdateVelocity(velocity rl.Vector2) {
	c.Velocity.X = velocity.X + c.Acceleration.X
	c.Velocity.Y = velocity.Y + c.Acceleration.Y
}

func (g *Game) UpdateHeroAcceleration(impactForce rl.Vector2) {
	if g.Hero.IsOnGround {
		g.Hero.Acceleration.Y = 0
	} else {
		g.Hero.Acceleration.Y += impactForce.Y + g.Gravity*0.2
	}
	g.Hero.Acceleration.X += impactForce.X
}

func CheckCoiliotion(r1, r2 rl.Rectangle) bool {
	r1LeftEdge, r1TopEdge, r1RightEdge, r1BottomEdge := getEdges(r1)
	r2LeftEdge, r2TopEdge, r2RightEdge, r2BottomEdge := getEdges(r2)

	// if r1 is on the right of r2 OR if r1 is on the left of r2
	if r1LeftEdge > r2RightEdge || r1RightEdge < r2LeftEdge {
		return false
	}

	// if r1 is above r2 OR r1 is under r2
	if r1BottomEdge < r2TopEdge || r1TopEdge > r2BottomEdge {
		return false
	}

	return true
}

// return the edges of Rectangle clock-wise order
func getEdges(r rl.Rectangle) (float32, float32, float32, float32) {
	return r.X, r.Y, r.X + r.Width, r.Y + r.Height
}

func (g *Game) CheckCollisionWithMap() {
	for _, tile := range g.Map.Tiles {
		heroRect := rl.Rectangle{X: g.Hero.Position.X, Y: g.Hero.Position.Y, Width: float32(g.GrassTile.Width), Height: float32(g.GrassTile.Height)}
		heroRect.Width = g.HeroScaling * float32(g.Hero.States[g.Hero.CurrentState].Texture.Width)
		heroRect.Height = g.HeroScaling * float32(g.Hero.States[g.Hero.CurrentState].Texture.Height)
		tileRect := rl.Rectangle{X: tile.X, Y: tile.Y, Width: float32(g.GrassTile.Width), Height: float32(g.GrassTile.Height)}
		if CheckCoiliotion(heroRect, tileRect) {
			g.Hero.IsOnGround = true
			return
		} else {
			g.Hero.IsOnGround = false
		}
	}
}
