package main

import "embed"

//go:embed roms/PONG roms/TANK roms/TETRIS
var embeddedROMs embed.FS

// GetROMBytes returns the bytes of the specified ROM file
func GetROMBytes(romName string) ([]byte, error) {
	return embeddedROMs.ReadFile("roms/" + romName)
}

// DefaultROM is the ROM that will be loaded if no argument is provided
const DefaultROM = "PONG"

// ValidROMs contains the list of available ROMs
var ValidROMs = []string{"PONG", "TANK", "TETRIS"}
