package core

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"

	"fire/internal/components"
	"fire/internal/ecs"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// TileJSON represents a tile in the JSON format.
type TileJSON struct {
	X        float32 `json:"x"`
	Y        float32 `json:"y"`
	TileType int32   `json:"tileType"`
}

// LevelMapJSON represents the JSON structure for a level map.
type LevelMapJSON struct {
	Tiles []TileJSON `json:"tiles"`
}

// LevelMap holds the map data for the designer mode.
type LevelMap struct {
	Tiles []TileJSON `json:"tiles"`
}

// NewLevelMap creates an empty level map.
func NewLevelMap() LevelMap {
	return LevelMap{Tiles: make([]TileJSON, 0)}
}

// LoadLevelMap loads a map from a JSON file.
func LoadLevelMap(path string) (LevelMap, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return LevelMap{}, err
	}

	var levelMap LevelMap
	if err := json.Unmarshal(data, &levelMap); err != nil {
		return LevelMap{}, err
	}

	return levelMap, nil
}

// Save writes the map data to the given path in JSON format.
func (m LevelMap) Save(path string) error {
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	return os.WriteFile(path, data, 0o644)
}

// AddTile adds a tile to the map if it does not already exist at the
// provided coordinates. Returns true if the tile list changed.
func (m *LevelMap) AddTile(x, y float32, tileType int32) bool {
	for idx, existing := range m.Tiles {
		if existing.X == x && existing.Y == y {
			if existing.TileType == tileType {
				return false
			}
			m.Tiles[idx].TileType = tileType
			return true
		}
	}

	m.Tiles = append(m.Tiles, TileJSON{X: x, Y: y, TileType: tileType})
	return true
}

// RemoveTileAt removes the first tile found at the provided coordinates.
// Returns true if a tile was removed.
func (m *LevelMap) RemoveTileAt(x, y float32) bool {
	for idx, existing := range m.Tiles {
		if existing.X == x && existing.Y == y {
			m.Tiles = append(m.Tiles[:idx], m.Tiles[idx+1:]...)
			return true
		}
	}
	return false
}

// SpawnTiles creates tile entities in the ECS world from the loaded map.
func SpawnTiles(world *ecs.World, levelMap LevelMap, tileTexture rl.Texture2D) {
	// Get component stores
	transformStore := ecs.RegisterStore[*components.TransformComponent](world.Components)
	colliderStore := ecs.RegisterStore[*components.ColliderComponent](world.Components)
	tileStore := ecs.RegisterStore[*components.TileComponent](world.Components)

	tileWidth := float32(tileTexture.Width)
	tileHeight := float32(tileTexture.Height)

	for _, tile := range levelMap.Tiles {
		entity := world.CreateEntity("tile", "ground")

		// Add transform component
		transformStore.Add(entity.ID, &components.TransformComponent{
			Position:    rl.Vector2{X: tile.X, Y: tile.Y},
			FacingRight: true,
		})

		// Add collider component
		colliderStore.Add(entity.ID, &components.ColliderComponent{
			Bounds:    rl.Rectangle{X: 0, Y: 0, Width: tileWidth, Height: tileHeight},
			IsTrigger: false,
			Layer:     "ground",
		})

		// Add tile component
		tileStore.Add(entity.ID, &components.TileComponent{
			TileType: components.TileType(tile.TileType),
		})
	}
}

// LoadAndSpawnMap loads the map from disk and spawns tile entities.
func LoadAndSpawnMap(world *ecs.World, tileTexture rl.Texture2D) {
	mapPath := ResourcePath("maps/custom_map.json")
	levelMap, err := LoadLevelMap(mapPath)
	if err != nil {
		log.Printf("Unable to load map %s, starting empty: %v", mapPath, err)
		return
	}
	SpawnTiles(world, levelMap, tileTexture)
}

// InitMapForDesigner loads the map for the designer mode.
func InitMapForDesigner() LevelMap {
	mapPath := ResourcePath("maps/custom_map.json")
	levelMap, err := LoadLevelMap(mapPath)
	if err != nil {
		log.Printf("Unable to load map %s, starting empty: %v", mapPath, err)
		return NewLevelMap()
	}
	return levelMap
}

// ErrMapNotFound is returned when the requested map is missing on disk.
var ErrMapNotFound = errors.New("map not found")
