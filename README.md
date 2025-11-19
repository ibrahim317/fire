## Fire (raylib-go)

A tiny raylib-go project showcasing a character with GIF-based animations.

### Requirements
- **Go** 1.20+ installed and on your PATH

### Run
```bash
go run main.go
```

### Designer mode
An interactive map designer is available to sketch platform layouts with the mouse.

```bash
go run ./cmd/designer -map maps/my_level.json
```

- Left click: add a grass tile snapped to the grid
- Right click: remove a tile
- `Save Map` button: exports the current layout to JSON (default `maps/custom_map.json`)
- Saved files can be loaded in-game via `core.LoadLevelMap`

### Controls
- Left/Right Arrow: move the character (switches to Running)
- Release keys: character returns to Idle

### Project layout
- `main.go`: program entry point
- `internal/core/`: core game types and asset loading (GIF animation support)
- `internal/logic/`: input handling and animation state updates
- `internal/render/`: drawing
- `resources/`: images, fonts, and other assets

### Notes
- Assets are loaded from the `resources/` directory; keep paths intact when running.
