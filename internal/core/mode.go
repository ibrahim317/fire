package core

// GameMode represents the current mode of the application
type GameMode int

const (
	// ModeMainMenu shows the main menu with options to play or design
	ModeMainMenu GameMode = iota
	// ModeGame is the actual gameplay mode
	ModeGame
	// ModeDesigner is the map design mode
	ModeDesigner
	// ModeSettings is the settings/options mode
	ModeSettings
)
