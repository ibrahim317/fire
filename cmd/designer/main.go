package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"fire/internal/core"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type designerApp struct {
	game          *core.Game
	mapPath       string
	buttonRect    rl.Rectangle
	tileWidth     float32
	tileHeight    float32
	statusMessage string
	statusColor   rl.Color
	statusUntil   time.Time
}

func main() {
	var mapPath string
	flag.StringVar(&mapPath, "map", "maps/custom_map.json", "Path where the designer will save the map JSON")
	flag.Parse()

	var game core.Game
	game.Init()

	rl.InitWindow(game.ScreenWidth, game.ScreenHeight, "Fire - Map Designer")
	defer rl.CloseWindow()

	rl.SetTargetFPS(core.FPS)

	game.LoadAssets()
	defer game.UnloadAssets()

	game.Map = maybeLoadExistingMap(mapPath)

	app := newDesignerApp(&game, mapPath)

	for !rl.WindowShouldClose() {
		app.update()
		app.draw()
	}
}

func maybeLoadExistingMap(path string) core.LevelMap {
	levelMap, err := core.LoadLevelMap(path)
	if err == nil {
		return levelMap
	}

	if os.IsNotExist(err) {
		return core.NewLevelMap()
	}

	log.Printf("unable to load map %s, starting empty: %v", path, err)
	return core.NewLevelMap()
}

func newDesignerApp(game *core.Game, mapPath string) *designerApp {
	tileWidth := float32(game.GrassTile.Width)
	tileHeight := float32(game.GrassTile.Height)

	return &designerApp{
		game:       game,
		mapPath:    mapPath,
		tileWidth:  tileWidth,
		tileHeight: tileHeight,
		buttonRect: rl.Rectangle{X: 20, Y: 20, Width: 220, Height: 44},
	}
}

func (d *designerApp) update() {
	mouse := rl.GetMousePosition()

	if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		if rl.CheckCollisionPointRec(mouse, d.buttonRect) {
			if err := d.saveMap(); err != nil {
				d.setStatus(fmt.Sprintf("Save failed: %v", err), rl.Red)
			} else {
				d.setStatus(
					fmt.Sprintf("Saved map to %s", filepath.Base(d.mapPath)),
					rl.DarkGreen,
				)
			}
			return
		}

		d.addTileAt(mouse)
		d.setStatus("Tile added", rl.DarkGreen)
	}

	if rl.IsMouseButtonPressed(rl.MouseRightButton) {
		if d.removeTileAt(mouse) {
			d.setStatus("Tile removed", rl.Orange)
		}
	}
}

func (d *designerApp) draw() {
	rl.BeginDrawing()
	defer rl.EndDrawing()

	rl.ClearBackground(rl.RayWhite)
	rl.DrawTextureEx(d.game.Bg, rl.Vector2{X: 0, Y: 0}, 0, 2.7, rl.White)

	d.drawGrid()
	d.game.DrawMap()
	d.drawButton()
	d.drawInstructions()
	d.drawStatus()
}

func (d *designerApp) addTileAt(pos rl.Vector2) {
	snapped := d.snapToGrid(pos)
	d.game.Map.AddTile(core.Tile{
		X:        snapped.X,
		Y:        snapped.Y,
		TileType: core.Grass,
	})
}

func (d *designerApp) removeTileAt(pos rl.Vector2) bool {
	snapped := d.snapToGrid(pos)
	return d.game.Map.RemoveTileAt(snapped.X, snapped.Y)
}

func (d *designerApp) snapToGrid(pos rl.Vector2) rl.Vector2 {
	x := float32(int(pos.X/d.tileWidth)) * d.tileWidth
	y := float32(int(pos.Y/d.tileHeight)) * d.tileHeight
	return rl.Vector2{X: x, Y: y}
}

func (d *designerApp) drawGrid() {
	gridColor := rl.Color{R: 255, G: 255, B: 255, A: 32}
	screenWidth := float32(d.game.ScreenWidth)
	screenHeight := float32(d.game.ScreenHeight)

	for x := float32(0); x <= screenWidth; x += d.tileWidth {
		rl.DrawLineV(rl.Vector2{X: x, Y: 0}, rl.Vector2{X: x, Y: screenHeight}, gridColor)
	}

	for y := float32(0); y <= screenHeight; y += d.tileHeight {
		rl.DrawLineV(rl.Vector2{X: 0, Y: y}, rl.Vector2{X: screenWidth, Y: y}, gridColor)
	}
}

func (d *designerApp) drawButton() {
	buttonColor := rl.Color{R: 34, G: 139, B: 34, A: 220}
	if rl.CheckCollisionPointRec(rl.GetMousePosition(), d.buttonRect) {
		buttonColor = rl.Color{R: 50, G: 205, B: 50, A: 220}
	}

	rl.DrawRectangleRec(d.buttonRect, buttonColor)
	rl.DrawRectangleLinesEx(d.buttonRect, 2, rl.White)
	rl.DrawText("Save Map", int32(d.buttonRect.X)+20, int32(d.buttonRect.Y)+12, 16, rl.White)
}

func (d *designerApp) drawInstructions() {
	instruction := "Left-click: add tile | Right-click: remove tile"
	rl.DrawText(instruction, 20, 70, 14, rl.DarkGray)
	rl.DrawText(fmt.Sprintf("Saving to %s", d.mapPath), 20, 90, 14, rl.DarkGray)
}

func (d *designerApp) drawStatus() {
	if d.statusMessage == "" || time.Now().After(d.statusUntil) {
		return
	}
	rl.DrawText(d.statusMessage, 20, 110, 16, d.statusColor)
}

func (d *designerApp) setStatus(message string, color rl.Color) {
	d.statusMessage = message
	d.statusColor = color
	d.statusUntil = time.Now().Add(2 * time.Second)
}

func (d *designerApp) saveMap() error {
	if err := d.game.Map.Save(d.mapPath); err != nil {
		return err
	}
	return nil
}
