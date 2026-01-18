package core

import (
	"fmt"
	"path/filepath"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Designer handles the map design mode functionality
type Designer struct {
	MapPath       string
	TileWidth     float32
	TileHeight    float32
	StatusMessage string
	StatusColor   rl.Color
	StatusUntil   time.Time
	SaveButton    rl.Rectangle
	BackButton    rl.Rectangle
	Map           LevelMap
	GrassTile     rl.Texture2D
}

// NewDesigner creates a new designer instance
func NewDesigner(game *Game) *Designer {
	tileWidth := float32(game.GrassTile.Width)
	tileHeight := float32(game.GrassTile.Height)

	return &Designer{
		MapPath:    ResourcePath("maps/custom_map.json"),
		TileWidth:  tileWidth,
		TileHeight: tileHeight,
		SaveButton: rl.Rectangle{X: 20, Y: 20, Width: 140, Height: 44},
		BackButton: rl.Rectangle{X: 170, Y: 20, Width: 140, Height: 44},
		Map:        InitMapForDesigner(),
		GrassTile:  game.GrassTile,
	}
}

// Update handles designer input and returns the new game mode
func (d *Designer) Update(game *Game) GameMode {
	mouse := rl.GetMousePosition()

	// Handle escape key to go back to menu
	if rl.IsKeyPressed(rl.KeyEscape) {
		return ModeMainMenu
	}

	if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		// Check save button
		if rl.CheckCollisionPointRec(mouse, d.SaveButton) {
			if err := d.SaveMap(); err != nil {
				d.SetStatus(fmt.Sprintf("Save failed: %v", err), rl.Red)
			} else {
				d.SetStatus(
					fmt.Sprintf("Saved to %s", filepath.Base(d.MapPath)),
					rl.DarkGreen,
				)
			}
			return ModeDesigner
		}

		// Check back button
		if rl.CheckCollisionPointRec(mouse, d.BackButton) {
			return ModeMainMenu
		}

		// Add tile
		d.AddTileAt(mouse)
		d.SetStatus("Tile added", rl.DarkGreen)
	}

	if rl.IsMouseButtonPressed(rl.MouseRightButton) {
		if d.RemoveTileAt(mouse) {
			d.SetStatus("Tile removed", rl.Orange)
		}
	}

	return ModeDesigner
}

// Draw renders the designer UI
func (d *Designer) Draw(game *Game) {
	rl.ClearBackground(rl.RayWhite)
	rl.DrawTextureEx(game.Bg, rl.Vector2{X: 0, Y: 0}, 0, 2.7, rl.White)

	d.DrawGrid(game)
	d.DrawMap()
	d.DrawButtons()
	d.DrawInstructions()
	d.DrawStatus()
}

// DrawMap draws the tiles in the designer
func (d *Designer) DrawMap() {
	for _, tile := range d.Map.Tiles {
		texture := d.GrassTile
		tileWidth := float32(texture.Width)
		tileHeight := float32(texture.Height)

		sourceRec := rl.Rectangle{X: 0, Y: 0, Width: tileWidth, Height: tileHeight}
		destRec := rl.Rectangle{X: tile.X, Y: tile.Y, Width: tileWidth, Height: tileHeight}
		rl.DrawTexturePro(texture, sourceRec, destRec, rl.Vector2{X: 0, Y: 0}, 0, rl.White)
	}
}

// AddTileAt adds a tile at the given mouse position
func (d *Designer) AddTileAt(pos rl.Vector2) {
	snapped := d.SnapToGrid(pos)
	d.Map.AddTile(snapped.X, snapped.Y, 0) // 0 = Grass
}

// RemoveTileAt removes a tile at the given mouse position
func (d *Designer) RemoveTileAt(pos rl.Vector2) bool {
	snapped := d.SnapToGrid(pos)
	return d.Map.RemoveTileAt(snapped.X, snapped.Y)
}

// SnapToGrid snaps a position to the tile grid
func (d *Designer) SnapToGrid(pos rl.Vector2) rl.Vector2 {
	x := float32(int(pos.X/d.TileWidth)) * d.TileWidth
	y := float32(int(pos.Y/d.TileHeight)) * d.TileHeight
	return rl.Vector2{X: x, Y: y}
}

// DrawGrid draws the tile grid overlay
func (d *Designer) DrawGrid(game *Game) {
	gridColor := rl.Color{R: 255, G: 255, B: 255, A: 32}
	screenWidth := float32(game.ScreenWidth)
	screenHeight := float32(game.ScreenHeight)

	for x := float32(0); x <= screenWidth; x += d.TileWidth {
		rl.DrawLineV(rl.Vector2{X: x, Y: 0}, rl.Vector2{X: x, Y: screenHeight}, gridColor)
	}

	for y := float32(0); y <= screenHeight; y += d.TileHeight {
		rl.DrawLineV(rl.Vector2{X: 0, Y: y}, rl.Vector2{X: screenWidth, Y: y}, gridColor)
	}
}

// DrawButtons draws the save and back buttons
func (d *Designer) DrawButtons() {
	// Save button
	saveColor := rl.Color{R: 34, G: 139, B: 34, A: 220}
	if rl.CheckCollisionPointRec(rl.GetMousePosition(), d.SaveButton) {
		saveColor = rl.Color{R: 50, G: 205, B: 50, A: 220}
	}
	rl.DrawRectangleRec(d.SaveButton, saveColor)
	rl.DrawRectangleLinesEx(d.SaveButton, 2, rl.White)
	rl.DrawText("Save Map", int32(d.SaveButton.X)+20, int32(d.SaveButton.Y)+12, 16, rl.White)

	// Back button
	backColor := rl.Color{R: 70, G: 70, B: 90, A: 220}
	if rl.CheckCollisionPointRec(rl.GetMousePosition(), d.BackButton) {
		backColor = rl.Color{R: 100, G: 100, B: 130, A: 220}
	}
	rl.DrawRectangleRec(d.BackButton, backColor)
	rl.DrawRectangleLinesEx(d.BackButton, 2, rl.White)
	rl.DrawText("Back to Menu", int32(d.BackButton.X)+10, int32(d.BackButton.Y)+12, 16, rl.White)
}

// DrawInstructions draws the help text
func (d *Designer) DrawInstructions() {
	instruction := "Left-click: add tile | Right-click: remove tile | ESC: menu"
	rl.DrawText(instruction, 20, 70, 14, rl.DarkGray)
	rl.DrawText(fmt.Sprintf("Saving to %s", d.MapPath), 20, 90, 14, rl.DarkGray)
}

// DrawStatus draws the status message
func (d *Designer) DrawStatus() {
	if d.StatusMessage == "" || time.Now().After(d.StatusUntil) {
		return
	}
	rl.DrawText(d.StatusMessage, 20, 110, 16, d.StatusColor)
}

// SetStatus sets the status message to display
func (d *Designer) SetStatus(message string, color rl.Color) {
	d.StatusMessage = message
	d.StatusColor = color
	d.StatusUntil = time.Now().Add(2 * time.Second)
}

// SaveMap saves the current map to disk
func (d *Designer) SaveMap() error {
	return d.Map.Save(d.MapPath)
}

// ReloadMap reloads the map from disk
func (d *Designer) ReloadMap() {
	d.Map = InitMapForDesigner()
}
