package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/yuvrajchettri/chip-8-emulator/chip8"

	"github.com/veandco/go-sdl2/sdl"
)

var keyMap = map[uint8]sdl.Scancode{
	0x1: sdl.SCANCODE_1,
	0x2: sdl.SCANCODE_2,
	0x3: sdl.SCANCODE_3,
	0xC: sdl.SCANCODE_4,
	0x4: sdl.SCANCODE_Q,
	0x5: sdl.SCANCODE_W,
	0x6: sdl.SCANCODE_E,
	0xD: sdl.SCANCODE_R,
	0x7: sdl.SCANCODE_A,
	0x8: sdl.SCANCODE_S,
	0x9: sdl.SCANCODE_D,
	0xE: sdl.SCANCODE_F,
	0xA: sdl.SCANCODE_Z,
	0x0: sdl.SCANCODE_X,
	0xB: sdl.SCANCODE_C,
	0xF: sdl.SCANCODE_V,
}

func main() {
	// Default ROM filename
	romFile := "TANK"

	// If a filename is passed as an argument, use it
	if len(os.Args) > 1 {
		romFile = os.Args[1]
	}

	// Construct the path relative to ../roms
	romPath := filepath.Join("..", "roms", romFile)

	// Initialize SDL
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		log.Fatalf("Failed to initialize SDL: %v", err)
	}
	defer sdl.Quit()

	// Create a window
	var modifier = 10
	window, windowErr := sdl.CreateWindow(
		"Chip 8",
		sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		int32(chip8.DISPLAY_COLS*modifier), int32(chip8.DISPLAY_ROWS*modifier),
		sdl.WINDOW_SHOWN,
	)
	if windowErr != nil {
		panic(windowErr)
	}
	defer window.Destroy()

	// Create render surface - canvas
	canvas, canvasErr := sdl.CreateRenderer(window, -1, 0)
	if canvasErr != nil {
		panic(canvasErr)
	}
	defer canvas.Destroy()

	// Create a new chip-8 instance
	emulator := chip8.NewChip8(false, false, 700)

	// Load ROM
	if err := emulator.Load(romPath); err != nil {
		log.Fatalf("Failed to load ROM: %v", err)
	}

	// Set PC to start of ROM
	emulator.PC = 0x200

	loop(emulator, canvas, int32(modifier))
}

// ------------------------------------------------
// Loop for fetch-decode-execute cycle
// ------------------------------------------------
func loop(emulator *chip8.Chip8, canvas *sdl.Renderer, modifier int32) {
	emulator.Initialize()

	go updateKeyboardState(emulator)

	speedHz := emulator.Speed()
	instructionDelay := time.Second / time.Duration(speedHz)
	ticker := time.NewTicker(instructionDelay)
	defer ticker.Stop()

	for emulator.ProgramCounter() < chip8.RAM {
		// Pump events to update keyboard state only from main thread
		sdl.PumpEvents()

		// Render display if redraw is true
		if emulator.ShouldRedraw() {
			emulator.ResetRedraw()
			renderDisplay(emulator, canvas, modifier)
		}

		// Main instruction loop
		select {
		case <-ticker.C:
			instruction := emulator.Fetch()
			emulator.NextInstruction()

			emulator.ExecuteInstruction(instruction)
		}
	}
}

func renderDisplay(emulator *chip8.Chip8, canvas *sdl.Renderer, modifier int32) {
	canvas.SetDrawColor(255, 0, 0, 255)
	canvas.Clear()

	// Get the display buffer and render
	vector := emulator.GetDisplay()
	for j := 0; j < len(vector); j++ {
		for i := 0; i < len(vector[j]); i++ {
			// Values of pixel are stored in 1D array of size 64 * 32
			if vector[j][i] != 0 {
				canvas.SetDrawColor(255, 255, 0, 255)
			} else {
				canvas.SetDrawColor(255, 0, 0, 255)
			}
			canvas.FillRect(&sdl.Rect{
				Y: int32(j) * modifier,
				X: int32(i) * modifier,
				W: modifier,
				H: modifier,
			})
		}
	}

	canvas.Present()
}

func updateKeyboardState(emulator *chip8.Chip8) {
	ticker := time.NewTicker(time.Millisecond * 16) // ~60Hz refresh rate
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			keys := sdl.GetKeyboardState()

			// Update internal keyboard state
			for chip8Key, scancode := range keyMap {
				state := keys[scancode] != 0
				emulator.UpdateKeyboardState(chip8Key, state)
			}
		}
	}
}
