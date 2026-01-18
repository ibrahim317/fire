package systems

import (
	"fire/internal/components"
	"fire/internal/ecs"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// RenderConfig holds configuration for the render system.
type RenderConfig struct {
	Background       rl.Texture2D
	GrassTile        rl.Texture2D
	HealthHeart      rl.Texture2D
	HighlightBorders *bool
}

// RenderSystem draws all visible entities.
type RenderSystem struct {
	Config RenderConfig
}

// NewRenderSystem creates a new RenderSystem with the given configuration.
func NewRenderSystem(config RenderConfig) *RenderSystem {
	return &RenderSystem{Config: config}
}

// Update draws all entities (called during render phase).
func (s *RenderSystem) Update(world *ecs.World, dt float32) {
	rl.BeginDrawing()

	// Draw background
	rl.DrawTextureEx(s.Config.Background, rl.Vector2{X: 0, Y: 0}, 0.0, 2.7, rl.White)

	// Draw tiles first (background layer)
	s.drawTiles(world)

	// Draw characters/mobs (entity layer)
	s.drawSprites(world)

	// Draw health UI
	s.drawHealth()

	rl.EndDrawing()
}

// drawTiles draws all tile entities.
func (s *RenderSystem) drawTiles(world *ecs.World) {
	transformStore, ok1 := ecs.GetStore[*components.TransformComponent](world.Components)
	tileStore, ok2 := ecs.GetStore[*components.TileComponent](world.Components)
	colliderStore, _ := ecs.GetStore[*components.ColliderComponent](world.Components)

	if !ok1 || !ok2 {
		return
	}

	for _, id := range tileStore.All() {
		entity := world.GetEntity(id)
		if entity == nil || !entity.Active {
			continue
		}

		transform, ok := transformStore.Get(id)
		if !ok {
			continue
		}

		tile, _ := tileStore.Get(id)
		texture := s.textureForTile(tile.TileType)

		tileWidth := float32(texture.Width)
		tileHeight := float32(texture.Height)

		sourceRec := rl.Rectangle{X: 0, Y: 0, Width: tileWidth, Height: tileHeight}
		destRec := rl.Rectangle{X: transform.Position.X, Y: transform.Position.Y, Width: tileWidth, Height: tileHeight}
		rl.DrawTexturePro(texture, sourceRec, destRec, rl.Vector2{X: 0, Y: 0}, 0, rl.White)

		// Draw border if highlight is enabled
		if s.Config.HighlightBorders != nil && *s.Config.HighlightBorders {
			if colliderStore != nil {
				if collider, ok := colliderStore.Get(id); ok {
					bounds := collider.GetWorldBounds(transform.Position)
					rl.DrawRectangleLines(int32(bounds.X), int32(bounds.Y), int32(bounds.Width), int32(bounds.Height), rl.Green)
				}
			}
		}
	}
}

// drawSprites draws all entities with Transform and Sprite components.
func (s *RenderSystem) drawSprites(world *ecs.World) {
	transformStore, ok1 := ecs.GetStore[*components.TransformComponent](world.Components)
	spriteStore, ok2 := ecs.GetStore[*components.SpriteComponent](world.Components)
	colliderStore, _ := ecs.GetStore[*components.ColliderComponent](world.Components)

	if !ok1 || !ok2 {
		return
	}

	for _, id := range spriteStore.All() {
		entity := world.GetEntity(id)
		if entity == nil || !entity.Active {
			continue
		}

		transform, ok := transformStore.Get(id)
		if !ok {
			continue
		}

		sprite, _ := spriteStore.Get(id)
		animData := sprite.GetCurrentAnimation()
		if animData == nil {
			continue
		}

		frameWidth := float32(animData.Texture.Width)
		frameHeight := float32(animData.Texture.Height)

		var sourceRec rl.Rectangle
		if animData.IsSpriteSheet {
			// Sprite sheet: select frame by X offset
			sourceRec = rl.Rectangle{
				X:      float32(animData.CurrentFrame) * frameWidth / float32(animData.FrameCount),
				Y:      0,
				Width:  frameWidth / float32(animData.FrameCount),
				Height: frameHeight,
			}
			frameWidth = frameWidth / float32(animData.FrameCount)
		} else {
			// GIF animation: texture is already updated
			sourceRec = rl.Rectangle{X: 0, Y: 0, Width: frameWidth, Height: frameHeight}
		}

		// Flip horizontally if facing left
		if !transform.FacingRight {
			sourceRec.Width = -sourceRec.Width
		}

		destRec := rl.Rectangle{
			X:      transform.Position.X,
			Y:      transform.Position.Y,
			Width:  frameWidth * sprite.Scale,
			Height: frameHeight * sprite.Scale,
		}

		rl.DrawTexturePro(animData.Texture, sourceRec, destRec, rl.Vector2{X: 0, Y: 0}, 0, rl.White)

		// Draw border if highlight is enabled
		if s.Config.HighlightBorders != nil && *s.Config.HighlightBorders {
			if colliderStore != nil {
				if collider, ok := colliderStore.Get(id); ok {
					bounds := collider.GetWorldBounds(transform.Position)
					rl.DrawRectangleLines(int32(bounds.X), int32(bounds.Y), int32(bounds.Width), int32(bounds.Height), rl.Red)
				}
			}
		}
	}
}

// drawHealth draws the health UI.
func (s *RenderSystem) drawHealth() {
	for i := 0; i < 5; i++ {
		rl.DrawTextureEx(s.Config.HealthHeart, rl.Vector2{X: float32(10 + i*40), Y: 10}, 0.0, 0.02, rl.White)
	}
}

// textureForTile returns the appropriate texture for a tile type.
func (s *RenderSystem) textureForTile(tileType components.TileType) rl.Texture2D {
	// Currently all tiles use grass texture
	// Extend this when adding more tile types
	return s.Config.GrassTile
}
