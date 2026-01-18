package core

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// Toggle represents a toggleable option
type Toggle struct {
	Rect        rl.Rectangle
	Label       string
	Value       *bool
	LabelColor  rl.Color
	OnColor     rl.Color
	OffColor    rl.Color
	KnobColor   rl.Color
	FontSize    int32
	IsHovered   bool
}

// NewToggle creates a new toggle with default styling
func NewToggle(x, y float32, label string, value *bool) Toggle {
	return Toggle{
		Rect:       rl.Rectangle{X: x, Y: y, Width: 60, Height: 30},
		Label:      label,
		Value:      value,
		LabelColor: rl.White,
		OnColor:    rl.Color{R: 76, G: 175, B: 80, A: 255},  // Green
		OffColor:   rl.Color{R: 100, G: 100, B: 100, A: 255}, // Gray
		KnobColor:  rl.White,
		FontSize:   20,
		IsHovered:  false,
	}
}

// Update checks for toggle interaction
func (t *Toggle) Update() {
	mouse := rl.GetMousePosition()
	t.IsHovered = rl.CheckCollisionPointRec(mouse, t.Rect)

	if t.IsHovered && rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		*t.Value = !*t.Value
	}
}

// Draw renders the toggle
func (t *Toggle) Draw() {
	// Draw label
	labelWidth := rl.MeasureText(t.Label, t.FontSize)
	rl.DrawText(t.Label, int32(t.Rect.X)-labelWidth-20, int32(t.Rect.Y)+5, t.FontSize, t.LabelColor)

	// Draw toggle background
	bgColor := t.OffColor
	if *t.Value {
		bgColor = t.OnColor
	}
	rl.DrawRectangleRounded(t.Rect, 0.5, 8, bgColor)

	// Draw toggle border on hover
	if t.IsHovered {
		rl.DrawRectangleRoundedLinesEx(t.Rect, 0.5, 8, 2, rl.White)
	}

	// Draw knob
	knobRadius := t.Rect.Height/2 - 4
	knobX := t.Rect.X + knobRadius + 4
	if *t.Value {
		knobX = t.Rect.X + t.Rect.Width - knobRadius - 4
	}
	knobY := t.Rect.Y + t.Rect.Height/2
	rl.DrawCircle(int32(knobX), int32(knobY), knobRadius, t.KnobColor)
}

// SettingsMenu holds the settings UI components
type SettingsMenu struct {
	Title            string
	TitleFontSize    int32
	BackButton       Button
	HighlightToggle  Toggle
}

// NewSettingsMenu creates a new settings menu
func NewSettingsMenu(screenWidth, screenHeight int32, highlightBorders *bool) *SettingsMenu {
	centerX := float32(screenWidth) / 2
	centerY := float32(screenHeight) / 2

	// Back button at the bottom
	buttonWidth := float32(200)
	buttonHeight := float32(50)
	backButton := NewButton(centerX-buttonWidth/2, float32(screenHeight)-100, buttonWidth, buttonHeight, "Back to Menu")

	// Toggle for highlight borders
	toggleX := centerX + 50
	toggleY := centerY - 20
	highlightToggle := NewToggle(toggleX, toggleY, "Highlight Object Borders", highlightBorders)

	return &SettingsMenu{
		Title:           "Settings",
		TitleFontSize:   48,
		BackButton:      backButton,
		HighlightToggle: highlightToggle,
	}
}

// Update updates the settings menu state and returns the new mode
func (s *SettingsMenu) Update() GameMode {
	s.BackButton.Update()
	s.HighlightToggle.Update()

	if s.BackButton.IsClicked() || rl.IsKeyPressed(rl.KeyEscape) {
		return ModeMainMenu
	}
	return ModeSettings
}

// Draw renders the settings menu
func (s *SettingsMenu) Draw(bg rl.Texture2D) {
	// Draw background with overlay
	rl.DrawTextureEx(bg, rl.Vector2{X: 0, Y: 0}, 0.0, 2.7, rl.White)

	// Semi-transparent overlay for better text readability
	screenWidth := int32(rl.GetScreenWidth())
	screenHeight := int32(rl.GetScreenHeight())
	rl.DrawRectangle(0, 0, screenWidth, screenHeight, rl.Color{R: 0, G: 0, B: 0, A: 180})

	// Draw title with shadow effect
	titleWidth := rl.MeasureText(s.Title, s.TitleFontSize)
	titleX := screenWidth/2 - titleWidth/2
	titleY := int32(60)

	// Shadow
	rl.DrawText(s.Title, titleX+2, titleY+2, s.TitleFontSize, rl.Color{R: 0, G: 0, B: 0, A: 180})
	// Main title
	rl.DrawText(s.Title, titleX, titleY, s.TitleFontSize, rl.Color{R: 100, G: 149, B: 237, A: 255})

	// Draw toggles
	s.HighlightToggle.Draw()

	// Draw back button
	s.BackButton.Draw()

	// Draw instructions
	instructions := "Click toggles to change settings"
	instructionSize := int32(16)
	instructionWidth := rl.MeasureText(instructions, instructionSize)
	rl.DrawText(instructions, screenWidth/2-instructionWidth/2, screenHeight-40, instructionSize, rl.Color{R: 150, G: 150, B: 150, A: 255})
}
