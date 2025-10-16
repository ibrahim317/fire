## Fire (raylib-go)

A tiny raylib-go project showcasing a character with GIF-based animations.

### Requirements
- **Go** 1.20+ installed and on your PATH

### Run
```bash
go run main.go
```

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
