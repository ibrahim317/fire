package systems

import (
	"fire/internal/components"
	"fire/internal/ecs"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// PhysicsSystem applies gravity and handles movement.
type PhysicsSystem struct{}

// NewPhysicsSystem creates a new PhysicsSystem.
func NewPhysicsSystem() *PhysicsSystem {
	return &PhysicsSystem{}
}

// Update applies physics to all entities with Transform and Physics components.
func (s *PhysicsSystem) Update(world *ecs.World, dt float32) {
	transformStore, ok1 := ecs.GetStore[*components.TransformComponent](world.Components)
	physicsStore, ok2 := ecs.GetStore[*components.PhysicsComponent](world.Components)
	inputStore, _ := ecs.GetStore[*components.InputComponent](world.Components)

	if !ok1 || !ok2 {
		return
	}

	for _, id := range transformStore.All() {
		if !physicsStore.Has(id) {
			continue
		}

		entity := world.GetEntity(id)
		if entity == nil || !entity.Active {
			continue
		}

		transform, _ := transformStore.Get(id)
		physics, _ := physicsStore.Get(id)

		// Apply gravity when not on ground
		if !physics.IsOnGround {
			transform.Acceleration.Y += physics.Gravity * 0.2
		} else {
			transform.Acceleration.Y = 0
		}

		// Handle input if entity has InputComponent
		if inputStore != nil && inputStore.Has(id) {
			input, _ := inputStore.Get(id)

			// Horizontal movement
			velocity := rl.Vector2{X: input.MoveX, Y: 0}

			// Jump
			if input.JumpPressed && physics.IsOnGround {
				velocity.Y = -physics.JumpForce
				physics.IsOnGround = false
			}

			// Update facing direction
			if input.MoveX < 0 {
				transform.FacingRight = false
			} else if input.MoveX > 0 {
				transform.FacingRight = true
			}

			// Apply velocity with speed multiplier
			transform.Velocity.X = velocity.X*physics.MoveSpeed + transform.Acceleration.X
			transform.Velocity.Y = velocity.Y*physics.MoveSpeed + transform.Acceleration.Y
		} else {
			// Non-player entities just apply acceleration
			transform.Velocity.Y += transform.Acceleration.Y
		}

		// Update position
		transform.Position.X += transform.Velocity.X
		transform.Position.Y += transform.Velocity.Y
	}
}
