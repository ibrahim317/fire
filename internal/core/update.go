package core

import (
	"image/color"
	"unsafe"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Updates the animation frame for a given character state
func (g *Game) UpdateCharacterAnimation() {
	animData := g.Hero.States[g.Hero.CurrentState]
	UpdateAnimation(&animData)
	g.Hero.States[g.Hero.CurrentState] = animData
}

func (g *Game) UpdateMobAnimation() {
	animData := g.Mob.AnimationData
	UpdateAnimation(&animData)
	g.Mob.AnimationData = animData
}

func UpdateAnimation(animData *AnimationData) {

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

}

func (c *Character) UpdateMovementDirection(direction MovementDirection) {
	c.MovementDirection = direction
}

func (c *Character) UpdatePosition() {
	c.Position.Y += c.Velocity.Y
	c.Position.X += c.Velocity.X
}

func (c *Character) UpdateVelocity(velocity rl.Vector2) {
	c.Velocity.X = velocity.X*2 + c.Acceleration.X
	c.Velocity.Y = velocity.Y*2 + c.Acceleration.Y
}

func (g *Game) UpdateHeroAcceleration(impactForce rl.Vector2) {
	if g.Hero.IsOnGround {
		g.Hero.Acceleration.Y = 0
	} else {
		g.Hero.Acceleration.Y += impactForce.Y + g.Gravity*0.2
	}
	g.Hero.Acceleration.X += impactForce.X
}

func CheckCollisionDirection(r1, r2 rl.Rectangle) rl.Vector2 {
	r1Left, r1Top, r1Right, r1Bottom := getEdges(r1)
	r2Left, r2Top, r2Right, r2Bottom := getEdges(r2)

	// No collision, return zero vector
	if r1Left > r2Right || r1Right < r2Left || r1Bottom < r2Top || r1Top > r2Bottom {
		return rl.Vector2{X: 0, Y: 0}
	}

	// Calculate overlap on each side
	overlapLeft := r1Right - r2Left
	overlapRight := r2Right - r1Left
	overlapTop := r1Bottom - r2Top
	overlapBottom := r2Bottom - r1Top

	// Find the minimum overlap direction and return corresponding unit vector
	minOverlap := overlapLeft
	unitVec := rl.Vector2{X: -1, Y: 0} // left

	if overlapRight < minOverlap {
		minOverlap = overlapRight
		unitVec = rl.Vector2{X: 1, Y: 0} // right
	}
	if overlapTop < minOverlap {
		minOverlap = overlapTop
		unitVec = rl.Vector2{X: 0, Y: -1} // top
	}
	if overlapBottom < minOverlap {
		unitVec = rl.Vector2{X: 0, Y: 1} // bottom (ground)
	}

	return unitVec
}

// return the edges of Rectangle clock-wise order
func getEdges(r rl.Rectangle) (float32, float32, float32, float32) {
	return r.X, r.Y, r.X + r.Width, r.Y + r.Height
}

func (g *Game) CheckCollisionWithMap() {

	g.Hero.IsOnGround = false
	for _, tile := range g.Map.Tiles {
		heroRect := rl.Rectangle{X: g.Hero.Position.X,
			Y:      g.Hero.Position.Y,
			Width:  float32(g.Hero.States[g.Hero.CurrentState].Texture.Width) * g.HeroScaling,
			Height: float32(g.Hero.States[g.Hero.CurrentState].Texture.Height) * g.HeroScaling,
		}
		tileRect := rl.Rectangle{X: tile.X, Y: tile.Y, Width: float32(g.GrassTile.Width), Height: float32(g.GrassTile.Height)}
		collisionRec := rl.GetCollisionRec(heroRect, tileRect)
		collisionRecVector := rl.Vector2{X: collisionRec.Width, Y: collisionRec.Height}
		collisionDirection := CheckCollisionDirection(heroRect, tileRect)

		if collisionDirection.Y != 0 {

			correctionY := collisionDirection.Y * collisionRecVector.Y
			g.Hero.Position.Y += correctionY
			if collisionDirection.Y*g.Hero.Velocity.Y < 0 {
				g.Hero.UpdateVelocity(rl.Vector2{X: g.Hero.Velocity.X, Y: 0})
			}

			if collisionDirection.Y == -1 {
				g.Hero.IsOnGround = true
			}
			continue
		}
		if collisionDirection.X != 0 {
			correctionX := collisionDirection.X * collisionRecVector.X
			g.Hero.Position.X += correctionX
			if collisionDirection.X*g.Hero.Velocity.X < 0 {
				g.Hero.UpdateVelocity(rl.Vector2{X: 0, Y: g.Hero.Velocity.Y})
			}
			continue
		}
	}
}

// After collision, resultVelocity is (2, 0)

func XAxisCollision(r1, r2 rl.Rectangle) bool {
	r1LeftEdge, _, r1RightEdge, _ := getEdges(r1)
	r2LeftEdge, _, r2RightEdge, _ := getEdges(r2)

	if r1LeftEdge > r2RightEdge || r1RightEdge < r2LeftEdge {
		return true
	}

	return false
}

func (g *Game) HandleMovement(velocity rl.Vector2) {
	if velocity.X < 0 {
		if velocity.Y < 0 {
			g.Hero.UpdateMovementDirection(UpLeft)
		} else if velocity.Y == 0 {
			g.Hero.UpdateMovementDirection(Left)
		}
	} else if velocity.X > 0 {
		if velocity.Y < 0 {
			g.Hero.UpdateMovementDirection(UpRight)
		} else if velocity.Y == 0 {
			g.Hero.UpdateMovementDirection(Right)
		}
	} else {
		if velocity.Y < 0 {
			g.Hero.UpdateMovementDirection(Up)
		}
	}
	g.Hero.UpdateVelocity(velocity)
}
