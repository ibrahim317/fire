package systems

import (
	"image/color"
	"unsafe"

	"fire/internal/components"
	"fire/internal/ecs"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// AnimationSystem updates sprite animations based on entity state.
type AnimationSystem struct{}

// NewAnimationSystem creates a new AnimationSystem.
func NewAnimationSystem() *AnimationSystem {
	return &AnimationSystem{}
}

// Update advances animation frames for all entities with SpriteComponent.
func (s *AnimationSystem) Update(world *ecs.World, dt float32) {
	spriteStore, ok := ecs.GetStore[*components.SpriteComponent](world.Components)
	if !ok {
		return
	}

	transformStore, _ := ecs.GetStore[*components.TransformComponent](world.Components)
	physicsStore, _ := ecs.GetStore[*components.PhysicsComponent](world.Components)
	inputStore, _ := ecs.GetStore[*components.InputComponent](world.Components)

	for _, id := range spriteStore.All() {
		entity := world.GetEntity(id)
		if entity == nil || !entity.Active {
			continue
		}

		sprite, _ := spriteStore.Get(id)

		// Determine animation state based on physics and input
		newAnim := components.AnimIdle

		hasPhysics := physicsStore != nil && physicsStore.Has(id)
		hasInput := inputStore != nil && inputStore.Has(id)

		if hasPhysics && hasInput {
			physics, _ := physicsStore.Get(id)
			input, _ := inputStore.Get(id)

			if !physics.IsOnGround {
				if transformStore != nil {
					transform, ok := transformStore.Get(id)
					if ok && transform.Velocity.Y > 0 {
						newAnim = components.AnimFalling
					} else {
						newAnim = components.AnimJumping
					}
				} else {
					newAnim = components.AnimFalling
				}
			} else if input.MoveX != 0 {
				newAnim = components.AnimRunning
			}
		}

		// Update animation state
		sprite.CurrentAnim = newAnim

		// Advance animation frame
		animData := sprite.GetCurrentAnimation()
		if animData != nil {
			updateAnimationFrame(animData)
		}
	}
}

// updateAnimationFrame advances the animation frame counter and updates textures.
func updateAnimationFrame(animData *components.AnimationData) {
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

		if !animData.IsSpriteSheet {
			// Calculate memory offset for current frame (for GIF animations)
			nextFrameDataOffset := uint32(animData.Image.Width * animData.Image.Height * 4 * animData.CurrentFrame)

			// Update GPU texture with current frame data
			rl.UpdateTexture(animData.Texture,
				unsafe.Slice((*color.RGBA)(unsafe.Pointer(uintptr(animData.Image.Data)+uintptr(nextFrameDataOffset))), animData.FrameSize))
		}

		animData.FrameCounter = 0
	}
}
