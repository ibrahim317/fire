---
trigger: always_on
glob:
description:
---
# Fire Platformer - Developer Guide

This project uses a **Composition-based Entity System (ECS-Lite)** architecture.

## System Overview

The game loop logic is decoupled from game objects.
- **Entities**: Just IDs with tags (e.g., "player", "enemy").
- **Components**: Data containers (e.g., `TransformComponent`, `SpriteComponent`).
- **Systems**: logic that iterates over entities with specific components.
- **World**: Container for all entities, components, and systems.

## Directory Structure

```
internal/
├── ecs/            # Core architecture (Entity, Component, World, EventBus)
├── components/     # All game data structs (Transform, Physics, AI, etc.)
├── systems/        # Logic implementations (Input, Physics, Collision, Render)
└── core/           # Game setup, assets, maps, and UI (MainMenu, Designer)
```

## How To...

### Add a New Feature (e.g., "Shield")

1.  **Define Data**: Create `ShieldComponent` in `internal/components/components.go`.
    ```go
    type ShieldComponent struct {
        Strength int
        Active   bool
    }
    ```

2.  **Implement Logic**: Create `internal/systems/shield.go`.
    ```go
    type ShieldSystem struct{}
    func (s *ShieldSystem) Update(world *ecs.World, dt float32) {
        // Query entities with ShieldComponent
        // Reduce strength over time, handle active state, etc.
    }
    ```

3.  **Register System**: Add `world.AddSystem(systems.NewShieldSystem())` in `main.go` (inside `initGameWorld`).

4.  **Add to Entity**: In `core/spawn.go`, add the component to `SpawnPlayer` or relevant entity.

### Add a New Enemy

1.  **Assets**: Load new texture/animation in `core/assets.go`.
2.  **Spawn Logic**: Create `SpawnNewEnemy` in `core/spawn.go`.
3.  **Components**: Attach `Transform`, `Sprite`, `Collider`, `Physics`, and `AI` components.
4.  **Behavior**: If it needs unique behavior, add a new `AIBehavior` type in `components.go` and handle it in `AISystem` (or create a specific system).

### Create a New Level

1.  Run the game and select **Design Mode**.
2.  Left-click to place tiles, Right-click to remove.
3.  Click **Save Map**.
4.  The map is loaded automatically in Game Mode.

## Debugging

- **Render Issues**: Check `RenderSystem.Update` and ensure entities have both `Transform` and `Sprite` components.
- **Physics/Collision**: Check `ColliderComponent` bounds and `PhysicsComponent` settings (gravity, IsOnGround).
- **Missing Entities**: Ensure `Spawn...` function is called and keys match in `Asset` loading.


