//go:build js && wasm

package main

import (
	"fmt"
	"log"
	"os"
	"syscall/js"

	"github.com/yuvrajchettri/chip-8-emulator/chip8"
)

var keyMap = map[string]uint8{
	"1": 0x1, "2": 0x2, "3": 0x3, "4": 0xC,
	"q": 0x4, "w": 0x5, "e": 0x6, "r": 0xD,
	"a": 0x7, "s": 0x8, "d": 0x9, "f": 0xE,
	"z": 0xA, "x": 0x0, "c": 0xB, "v": 0xF,
}

var (
	ctx       js.Value
	keyStates = make(map[uint8]bool)
)

func main() {
	// Get the canvas context from JavaScript
	doc := js.Global().Get("document")
	canvas := doc.Call("getElementById", "chip8-canvas")
	ctx = canvas.Call("getContext", "2d")

	// Default ROM filename
	romName := DefaultROM

	// If a filename is passed as an argument, use it
	if len(os.Args) > 1 {
		romName = os.Args[1]
	}

	// Validate ROM name
	isValid := false
	for _, validROM := range ValidROMs {
		if romName == validROM {
			isValid = true
			break
		}
	}
	if !isValid {
		fmt.Printf("Invalid ROM name. Available ROMs: %v\n", ValidROMs)
		os.Exit(1)
	}

	// Get ROM bytes
	romBytes, err := GetROMBytes(romName)
	if err != nil {
		log.Fatalf("Failed to read ROM: %v", err)
	}

	// Create a new chip-8 instance
	emulator := chip8.NewChip8(false, false, 1400)

	// Load ROM bytes
	if err := emulator.LoadBytes(romBytes); err != nil {
		log.Fatalf("Failed to load ROM: %v", err)
	}

	// Set PC to start of ROM
	emulator.PC = 0x200

	// Setup keyboard event listeners
	setupKeyboardHandlers()

	// Start the emulation loop
	loop(emulator, 10) // modifier of 10 like in SDL version

	// Keep the program running
	select {}
}

func renderDisplay(emulator *chip8.Chip8, modifier int32) {
	// Clear the canvas
	ctx.Set("fillStyle", "#FF0000") // Red background
	ctx.Call("fillRect", 0, 0, chip8.DISPLAY_COLS*modifier, chip8.DISPLAY_ROWS*modifier)

	// Get the display buffer and render
	vector := emulator.GetDisplay()
	ctx.Set("fillStyle", "#FFFF00") // Yellow pixels

	for j := 0; j < len(vector); j++ {
		for i := 0; i < len(vector[j]); i++ {
			if vector[j][i] != 0 {
				ctx.Call("fillRect",
					i*int(modifier), // x
					j*int(modifier), // y
					int(modifier),   // width
					int(modifier))   // height
			}
		}
	}
}

func setupKeyboardHandlers() {
	// Create keydown handler
	keydownHandler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		event := args[0]
		key := event.Get("key").String()
		if chip8Key, ok := keyMap[key]; ok {
			keyStates[chip8Key] = true
		}
		return nil
	})

	// Create keyup handler
	keyupHandler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		event := args[0]
		key := event.Get("key").String()
		if chip8Key, ok := keyMap[key]; ok {
			keyStates[chip8Key] = false
		}
		return nil
	})

	// Register the event listeners
	js.Global().Get("document").Call("addEventListener", "keydown", keydownHandler)
	js.Global().Get("document").Call("addEventListener", "keyup", keyupHandler)
}

func updateKeyboardState(emulator *chip8.Chip8) {
	// Update the emulator's keyboard state based on our keyStates map
	for chip8Key, isPressed := range keyStates {
		emulator.UpdateKeyboardState(chip8Key, isPressed)
	}
}

func loop(emulator *chip8.Chip8, modifier int32) {
	emulator.Initialize()

	// Rendering loop runs at 60 FPS because of requestAnimationFrame.
	// This means you have 60 “slots” per second to both run the CPU and update the display.
	// This means for each of those 60 frames, you execute ~12 CHIP-8 instructions, so that by the end of 1 second:
	// 12 instructions/frame * 60 frames/second = 720 instructions/second ~ nearly the correct number of instructions
	instrPerFrame := emulator.Speed() / 60 // e.g. 700/60 ≈ 12

	var renderFrame js.Func
	renderFrame = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// Update keyboard state
		updateKeyboardState(emulator)

		// 1. Run CPU
		for i := 0; i < instrPerFrame; i++ {
			instr := emulator.Fetch()
			emulator.NextInstruction()
			emulator.ExecuteInstruction(instr)
		}

		// 2. If the VRAM changed, paint it
		if emulator.ShouldRedraw() {
			emulator.ResetRedraw()
			renderDisplay(emulator, modifier)
		}

		js.Global().Call("requestAnimationFrame", renderFrame)
		return nil
	})
	js.Global().Call("requestAnimationFrame", renderFrame)
}
