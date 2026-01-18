package systems

import (
	"fire/internal/components"
	"fire/internal/ecs"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// InputSystem reads keyboard input and updates InputComponent.
type InputSystem struct{}

// NewInputSystem creates a new InputSystem.
func NewInputSystem() *InputSystem {
	return &InputSystem{}
}

// Update reads input and updates all entities with InputComponent.
func (s *InputSystem) Update(world *ecs.World, dt float32) {
	inputStore, ok := ecs.GetStore[*components.InputComponent](world.Components)
	if !ok {
		return
	}

	for _, id := range inputStore.All() {
		entity := world.GetEntity(id)
		if entity == nil || !entity.Active {
			continue
		}

		input, _ := inputStore.Get(id)

		// Reset per-frame state
		input.JumpPressed = false

		// Horizontal movement
		input.MoveX = 0
		if rl.IsKeyDown(rl.KeyRight) {
			input.MoveX = 1
		} else if rl.IsKeyDown(rl.KeyLeft) {
			input.MoveX = -1
		}

		// Jump input
		if rl.IsKeyPressed(rl.KeyUp) || rl.IsKeyPressed(rl.KeySpace) {
			input.JumpPressed = true
		}
		input.JumpHeld = rl.IsKeyDown(rl.KeyUp) || rl.IsKeyDown(rl.KeySpace)
	}
}
