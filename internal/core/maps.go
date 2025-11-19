package core

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
)

type TileType int32

const (
	Grass TileType = iota
	Stone
	Water
	Tree
	Rock
)

// Tile represents a single tile on the map grid.
type Tile struct {
	X        float32  `json:"x"`
	Y        float32  `json:"y"`
	TileType TileType `json:"tileType"`
}

type LevelMap struct {
	Tiles []Tile `json:"tiles"`
}

func NewLevelMap() LevelMap {
	return LevelMap{Tiles: make([]Tile, 0)}
}

func InitMap(g *Game) LevelMap {
	levelMap, err := LoadLevelMap(resourcePath("maps/custom_map.json"))
	if err != nil {
		log.Printf("unable to load map %s, starting empty: %v", resourcePath("maps/custom_map.json"), err)
		return NewLevelMap()
	}
	return levelMap
}

// AddTile adds a tile to the map if it does not already exist at the
// provided coordinates. Returns true if the tile list changed.
func (m *LevelMap) AddTile(tile Tile) bool {
	for idx, existing := range m.Tiles {
		if existing.X == tile.X && existing.Y == tile.Y {
			if existing.TileType == tile.TileType {
				return false
			}
			m.Tiles[idx] = tile
			return true
		}
	}

	m.Tiles = append(m.Tiles, tile)
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

// LoadLevelMap loads a map from the provided path.
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

// MustLoadLevelMap loads a map and panics if it fails. Mainly useful for tests.
func MustLoadLevelMap(path string) LevelMap {
	levelMap, err := LoadLevelMap(path)
	if err != nil {
		panic(err)
	}
	return levelMap
}

// ErrMapNotFound is returned when the requested map is missing on disk.
var ErrMapNotFound = errors.New("map not found")
