# Internal Package Overview

This document describes each file under `internal/`, which holds the game’s core logic, ECS framework, and UI.

---

## `components/`

### `components.go`

**Purpose:** Defines all **ECS component types** used by gameplay entities.

- **TransformComponent** — Position, velocity, acceleration, and facing direction.
- **AnimationState** / **AnimationData** — Animation state enum and per-animation data (texture, frame count, delay, sprite sheet vs GIF).
- **SpriteComponent** — Map of animations by state, current animation, and scale.
- **ColliderComponent** — Axis-aligned bounds (relative to entity), trigger flag, and collision layer (`"player"`, `"enemy"`, `"ground"`, etc.). `GetWorldBounds()` converts to world coordinates.
- **InputComponent** — Player input: `MoveX`, `JumpPressed`, `JumpHeld`.
- **PhysicsComponent** — Gravity, jump force, move speed, and `IsOnGround`.
- **HealthComponent** — Current/max health with `TakeDamage()`, `Heal()`, `IsDead()`.
- **TileComponent** — Marks an entity as a static tile; holds `TileType` (Grass, Stone, Water, Tree, Rock).
- **AIComponent** — AI behavior (Patrol, Chase, Idle), patrol path, path index, and optional target entity ID.

Used by systems and by `core/spawn.go` and `core/maps.go` when creating entities.

---

## `core/`

Game configuration, assets, menus, map data, and entity spawning. Not part of the ECS loop; it sets up the world and runs non-ECS screens.

### `game.go`

**Purpose:** Central **game state** and initialization.

- **Game** struct: screen size, loaded assets (fonts, textures, hero/mob animation data), current **GameMode**, references to MainMenu, Designer, Settings, global settings (e.g. `HighlightBorders`, `Gravity`), and the ECS **World** (used only in game mode).
- **AnimationDataLegacy** / **CharacterStateLegacy** — Legacy animation types used during asset loading; converted to components when spawning.
- `Init()` — Default resolution, gravity, hero scale, mode = main menu.
- `InitUI()` — Creates MainMenu, Designer, Settings (call after window exists).
- `InitWorld()` — Creates a new ECS world for gameplay.
- `GetHeroAnimationData(state)` — Returns legacy hero animation data for a state.

### `mode.go`

**Purpose:** **Game mode** enum.

- **GameMode**: `ModeMainMenu`, `ModeGame`, `ModeDesigner`, `ModeSettings`.
- Used by `main.go` to branch the main loop (menu vs gameplay vs designer vs settings).

### `assets.go`

**Purpose:** **Load and unload** all game assets.

- `LoadAssets()` — Loads fonts, background, grass tile, heart texture; loads hero GIFs (idle, run, jump, fall) and mob sprite sheet into `Game`’s legacy animation fields.
- `UnloadAssets()` — Unloads those resources (used with `defer` in `main.go`).
- Helpers: `loadAnimatedGifData()`, `loadSpriteSheetData()` for animation data.
- **ResourcePath(rel)** — Resolves paths under project root (finds directory containing `resources/`).
- **resolveProjectRoot()** / **findProjectRoot()** — Locate project root by walking up from cwd.

### `ui.go`

**Purpose:** **Main menu** UI (non-ECS).

- **Button** — Rectangle, label, colors, hover state; `Update()`, `IsClicked()`, `Draw()`.
- **MainMenu** — Title, Play / Design / Settings buttons; `Update()` returns next `GameMode` on click; `Draw(bg)` draws background overlay, title, subtitle, buttons, and instructions.

### `settings.go`

**Purpose:** **Settings screen** UI (non-ECS).

- **Toggle** — Label, bound `*bool`, colors; `Update()` flips value on click; `Draw()` renders label and toggle.
- **SettingsMenu** — Title, Back button, “Highlight Object Borders” toggle; `Update()` returns `ModeMainMenu` on Back or Escape; `Draw(bg)` draws background, title, toggle, back button, and instructions.

### `designer.go`

**Purpose:** **Map designer** mode (non-ECS): edit tiles, save/load map.

- **Designer** — Map path, tile size, status message/color/expiry, Save/Back button rects, **LevelMap**, grass texture.
- `Update(game)` — Escape → main menu; left-click Save → save map; left-click Back → main menu; left-click canvas → add tile; right-click → remove tile. Returns next `GameMode`.
- `Draw(game)` — Background, grid, tiles, buttons, instructions, status.
- Helpers: `AddTileAt()`, `RemoveTileAt()`, `SnapToGrid()`, `DrawGrid()`, `DrawMap()`, `DrawButtons()`, `DrawInstructions()`, `DrawStatus()`, `SetStatus()`, `SaveMap()`, `ReloadMap()`.

### `maps.go`

**Purpose:** **Level map** data (JSON), load/save, and **spawning tile entities** in the ECS world.

- **TileJSON** / **LevelMapJSON** / **LevelMap** — JSON-friendly tile (x, y, tileType) and list of tiles.
- `LoadLevelMap(path)` — Load map from JSON.
- `LevelMap.Save(path)` — Write map to JSON (creates dir if needed).
- `AddTile()` / `RemoveTileAt()` — Modify in-memory map (used by designer).
- **SpawnTiles(world, levelMap, tileTexture)** — For each tile in map, creates an entity with Transform, Collider, and TileComponent.
- **LoadAndSpawnMap(world, tileTexture)** — Loads `maps/custom_map.json` and calls `SpawnTiles`.
- **InitMapForDesigner()** — Loads same map for designer mode (or empty map if file missing).
- **ErrMapNotFound** — Sentinel for missing map file.

### `spawn.go`

**Purpose:** **Create gameplay entities** and attach components.

- **SpawnPlayer(world, game)** — Creates one entity with tag `"player"`, adds Transform, Sprite (from game’s hero animation data), Collider, Input, Physics, Health. Returns the entity.
- **SpawnMob(world, game, x, y)** — Creates entity with tags `"enemy"`, `"mob"`; adds Transform, Sprite (mob walk), Collider, Physics, AIComponent (patrol path). Returns the entity.
- **ResetPlayerPosition(world)** — Finds player by tag and resets position to (100, 300) and velocity; clears `IsOnGround` if PhysicsComponent present.

---

## `ecs/`

Minimal ECS runtime: entities, component storage, world, and events.

### `entity.go`

**Purpose:** **Entity** and **entity management**.

- **EntityID** — Unique uint32.
- **Entity** — ID, Active flag, Tags slice; `HasTag()`, `AddTag()`, `RemoveTag()`.
- **EntityManager** — Next ID, map of entities; `CreateEntity(tags...)`, `GetEntity(id)`, `RemoveEntity(id)`, `GetAllEntities()`, `GetEntitiesWithTag(tag)`.

### `component.go`

**Purpose:** **Component storage** and registry.

- **Component** — Marker interface (empty).
- **ComponentStore[T]** — Map from EntityID to T; Add, Get, Remove, Has, All.
- **ComponentRegistry** — Map from `reflect.Type` to component store; used to hold one store per component type.
- **RegisterStore[T](cr)** — Get or create store for T, register it, return it.
- **GetStore[T](cr)** — Retrieve existing store for T.

Uses Go generics and reflection for type-safe, per-type stores.

### `world.go`

**Purpose:** **ECS world** and system execution.

- **System** — Interface with `Update(world *World, dt float32)`.
- **World** — EntityManager, ComponentRegistry, list of Systems, EventBus.
- **NewWorld()** — Creates world with new entity manager, component registry, empty systems, and event bus.
- **AddSystem(s)** — Appends a system.
- **Update(dt)** — Calls `Update(w, dt)` on each system in order.
- **CreateEntity**, **RemoveEntity**, **GetEntity**, **GetEntitiesWithTag**, **GetAllEntities** — Delegate to EntityManager.

### `event.go`

**Purpose:** **Event bus** for decoupled game events.

- **EventType** — Constants: EventPlayerJump, EventPlayerLand, EventCollision, EventDamage, EventDeath, EventCoinCollected.
- **Event** — Type, Source EntityID, Target EntityID, Data interface{}.
- **EventBus** — Handlers per event type, queue of events; `Subscribe(type, handler)`, `Publish(event)`, `Process()` (dispatch queue then clear), `PublishImmediate(event)`, `Clear()`.

Currently the bus is created and stored on the world; systems could publish/subscribe for future features (e.g. damage, death, SFX).

---

## `systems/`

ECS **systems** that run each frame during gameplay (in order: Input → Physics → Collision → Animation → Render). Each implements `System` and operates on the world’s entities and components.

### `input.go`

**Purpose:** **InputSystem** — Drives **InputComponent** from keyboard.

- Iterates entities with **InputComponent**.
- Sets `MoveX` from Left/Right; sets `JumpPressed` and `JumpHeld` from Up/Space.
- Resets `JumpPressed` each frame.

Only entities with InputComponent (e.g. player) are affected.

### `physics.go`

**Purpose:** **PhysicsSystem** — Movement and gravity.

- Iterates entities that have both **TransformComponent** and **PhysicsComponent**.
- Applies gravity to acceleration when not on ground; clears vertical acceleration when on ground.
- If entity also has **InputComponent**: applies horizontal velocity from `MoveX`, applies jump velocity when `JumpPressed` and on ground, updates `FacingRight`, uses `MoveSpeed`.
- If no InputComponent: only integrates acceleration into velocity (e.g. mobs).
- Integrates velocity into position.

### `collision.go`

**Purpose:** **CollisionSystem** — Resolves collisions with tiles.

- Gets all entities with **ColliderComponent** and **TransformComponent**; skips entities that have **TileComponent** (tiles are static).
- Builds list of tile entities (Transform + Collider + TileComponent).
- For each non-tile entity: resets `IsOnGround`; then two passes:
  - **Vertical:** Collides with each tile; resolves overlap (position + velocity); sets `IsOnGround` when landing on top.
  - **Horizontal:** Same for X overlap and velocity.
- **checkCollisionDirection()** — Returns unit vector of minimum penetration (left/right/top/bottom). **getEdges()** — Rectangle edges.

Uses AABB vs tile colliders only; no player–enemy or trigger logic yet.

### `animation.go`

**Purpose:** **AnimationSystem** — Animation state and frame advance.

- Iterates entities with **SpriteComponent**.
- Chooses **AnimationState** from Physics + Input: falling if velocity.Y > 0, jumping if in air and not falling, running if on ground and MoveX != 0, else idle. Writes to `SpriteComponent.CurrentAnim`.
- Advances current animation: increments frame counter; when past FrameDelay, advances CurrentFrame (loops). For non–sprite-sheet (GIF), updates texture from image frame data via `rl.UpdateTexture`.

Supports both sprite sheets and GIF-style frame buffers.

### `render.go`

**Purpose:** **RenderSystem** — Draws the game scene (and simple HUD).

- **RenderConfig** — Background texture, grass tile texture, health heart texture, pointer to HighlightBorders setting.
- **Update(world, dt)** — `BeginDrawing()`; draws background; **drawTiles()**; **drawSprites()**; **drawHealth()**; `EndDrawing()`.
- **drawTiles** — Entities with Transform + TileComponent (and optional Collider for debug); uses `textureForTile()` (currently always grass); optionally draws collider outline if HighlightBorders.
- **drawSprites** — Entities with Transform + SpriteComponent; current animation frame, sprite-sheet or full texture; flips source rect for FacingRight; optional collider outline.
- **drawHealth** — Draws 5 hearts in top-left (no binding to HealthComponent yet).
- **textureForTile** — Maps TileType to texture (only grass implemented).

Rendering runs as the last system so all simulation is done before draw.

---

## Summary

| Directory     | Role |
|--------------|------|
| **components** | ECS component definitions (transform, sprite, collider, input, physics, health, tile, AI). |
| **core**       | Game state, modes, assets, main menu, settings, designer, map load/save, player/mob/tile spawning. |
| **ecs**        | Entities, component stores/registry, world, system interface, event bus. |
| **systems**    | Input, physics, collision, animation, render — run in order each frame during gameplay. |

Gameplay uses **ECS** (world + entities + components + systems). Menus, designer, and settings use **traditional structs and Update/Draw** in `core`.
