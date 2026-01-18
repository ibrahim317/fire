package core

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// Button represents a clickable UI button
type Button struct {
	Rect         rl.Rectangle
	Text         string
	NormalColor  rl.Color
	HoverColor   rl.Color
	TextColor    rl.Color
	BorderColor  rl.Color
	BorderWidth  float32
	FontSize     int32
	IsHovered    bool
}

// NewButton creates a new button with default styling
func NewButton(x, y, width, height float32, text string) Button {
	return Button{
		Rect:         rl.Rectangle{X: x, Y: y, Width: width, Height: height},
		Text:         text,
		NormalColor:  rl.Color{R: 45, G: 45, B: 58, A: 240},
		HoverColor:   rl.Color{R: 65, G: 105, B: 225, A: 255},
		TextColor:    rl.White,
		BorderColor:  rl.Color{R: 100, G: 100, B: 120, A: 255},
		BorderWidth:  2,
		FontSize:     24,
		IsHovered:    false,
	}
}

// Update checks if the button is being hovered
func (b *Button) Update() {
	mouse := rl.GetMousePosition()
	b.IsHovered = rl.CheckCollisionPointRec(mouse, b.Rect)
}

// IsClicked returns true if the button was clicked this frame
func (b *Button) IsClicked() bool {
	return b.IsHovered && rl.IsMouseButtonPressed(rl.MouseLeftButton)
}

// Draw renders the button
func (b *Button) Draw() {
	color := b.NormalColor
	if b.IsHovered {
		color = b.HoverColor
	}

	// Draw button background with rounded corners effect
	rl.DrawRectangleRec(b.Rect, color)
	rl.DrawRectangleLinesEx(b.Rect, b.BorderWidth, b.BorderColor)

	// Center text in button
	textWidth := rl.MeasureText(b.Text, b.FontSize)
	textX := int32(b.Rect.X) + (int32(b.Rect.Width)-textWidth)/2
	textY := int32(b.Rect.Y) + (int32(b.Rect.Height)-b.FontSize)/2

	rl.DrawText(b.Text, textX, textY, b.FontSize, b.TextColor)
}

// MainMenu holds the main menu UI components
type MainMenu struct {
	Title          string
	PlayButton     Button
	DesignButton   Button
	SettingsButton Button
	TitleFontSize  int32
}

// NewMainMenu creates a new main menu centered on screen
func NewMainMenu(screenWidth, screenHeight int32) *MainMenu {
	buttonWidth := float32(280)
	buttonHeight := float32(60)
	buttonSpacing := float32(20)
	centerX := float32(screenWidth)/2 - buttonWidth/2
	centerY := float32(screenHeight)/2 - buttonHeight

	return &MainMenu{
		Title:          "FIRE",
		TitleFontSize:  72,
		PlayButton:     NewButton(centerX, centerY, buttonWidth, buttonHeight, "Play Game"),
		DesignButton:   NewButton(centerX, centerY+buttonHeight+buttonSpacing, buttonWidth, buttonHeight, "Design Mode"),
		SettingsButton: NewButton(centerX, centerY+2*(buttonHeight+buttonSpacing), buttonWidth, buttonHeight, "Settings"),
	}
}

// Update updates the main menu state
func (m *MainMenu) Update() GameMode {
	m.PlayButton.Update()
	m.DesignButton.Update()
	m.SettingsButton.Update()

	if m.PlayButton.IsClicked() {
		return ModeGame
	}
	if m.DesignButton.IsClicked() {
		return ModeDesigner
	}
	if m.SettingsButton.IsClicked() {
		return ModeSettings
	}
	return ModeMainMenu
}

// Draw renders the main menu
func (m *MainMenu) Draw(bg rl.Texture2D) {
	// Draw background with overlay
	rl.DrawTextureEx(bg, rl.Vector2{X: 0, Y: 0}, 0.0, 2.7, rl.White)

	// Semi-transparent overlay for better text readability
	screenWidth := int32(rl.GetScreenWidth())
	screenHeight := int32(rl.GetScreenHeight())
	rl.DrawRectangle(0, 0, screenWidth, screenHeight, rl.Color{R: 0, G: 0, B: 0, A: 150})

	// Draw title with shadow effect
	titleWidth := rl.MeasureText(m.Title, m.TitleFontSize)
	titleX := screenWidth/2 - titleWidth/2
	titleY := int32(100)

	// Shadow
	rl.DrawText(m.Title, titleX+3, titleY+3, m.TitleFontSize, rl.Color{R: 0, G: 0, B: 0, A: 180})
	// Main title with gradient-like effect
	rl.DrawText(m.Title, titleX, titleY, m.TitleFontSize, rl.Color{R: 255, G: 165, B: 0, A: 255})

	// Subtitle
	subtitle := "A Platformer Adventure"
	subtitleSize := int32(20)
	subtitleWidth := rl.MeasureText(subtitle, subtitleSize)
	rl.DrawText(subtitle, screenWidth/2-subtitleWidth/2, titleY+m.TitleFontSize+10, subtitleSize, rl.Color{R: 200, G: 200, B: 200, A: 255})

	// Draw buttons
	m.PlayButton.Draw()
	m.DesignButton.Draw()
	m.SettingsButton.Draw()

	// Draw instructions
	instructions := "Use arrow keys to move, Up to jump"
	instructionSize := int32(16)
	instructionWidth := rl.MeasureText(instructions, instructionSize)
	rl.DrawText(instructions, screenWidth/2-instructionWidth/2, screenHeight-50, instructionSize, rl.Color{R: 150, G: 150, B: 150, A: 255})
}
