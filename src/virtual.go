package main

import "github.com/veandco/go-sdl2/sdl"

// ------------------------------------------------
// Virtual hardware used by the CHIP-8
// ------------------------------------------------

// ------------------------------------------------
// Constants
// ------------------------------------------------

const (
	RAM          = 4096
	STACK_SIZE   = 100
	DISPLAY_COLS = 64
	DISPLAY_ROWS = 32
)

var font = []uint8{
	0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
	0x20, 0x60, 0x20, 0x20, 0x70, // 1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
	0x90, 0x90, 0xF0, 0x10, 0x10, // 4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
	0xF0, 0x10, 0x20, 0x40, 0x40, // 7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
	0xF0, 0x90, 0xF0, 0x90, 0x90, // A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
	0xF0, 0x80, 0x80, 0x80, 0xF0, // C
	0xE0, 0x90, 0x90, 0x90, 0xE0, // D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
	0xF0, 0x80, 0xF0, 0x80, 0x80, // F
}

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

// ------------------------------------------------
// Chip8 struct
// ------------------------------------------------
type Chip8 struct {
	memory    []byte
	stack     []uint16
	display   [][]int
	registers map[nibble]uint8
	PC        uint16
	I         uint16
	speedHz   int  // Instructions per second
	shift1    bool // Configurable behaviour for shift instructions (8XY6 and 8XYE) - consider Y register or not
	bnnn1     bool // Configurable behaviour for BNNN instruction - BNNN or not (if not then BXNN)
}

func NewChip8(shift1, bnnn1 bool, speedHz int) *Chip8 {
	chip8 := &Chip8{
		memory:  make([]byte, RAM),
		stack:   make([]uint16, STACK_SIZE),
		display: make([][]int, DISPLAY_ROWS),
		speedHz: speedHz,
		shift1:  shift1,
		bnnn1:   bnnn1,
	}
	chip8.initialize()
	return chip8
}

// ------------------------------------------------
// 1. First 512 bytes in memory used to have the interpreter, that is no longer true as our interpreter runs in Go space. We can use first 512 for storing the font sprites. 60 bytes between 80-159 (0x050-0x09F)
// 2. Display is modelled as a 2D boolean array with 64 columns and 32 rows. To initialize the rows, we need a loop.
// 3. Initialize registers
// ------------------------------------------------
func (chip8 *Chip8) initialize() {
	// Initialize fonts in memory
	memory := chip8.memory
	start := 0x50
	end := 0x9F
	for i, j := 0, start; j <= end; i, j = i+1, j+1 {
		val := font[i]
		location := j

		memory[location] = val
	}

	// Initialize the rows for display, each row has 'COLS' number of elems
	display := chip8.display
	for i := 0; i < DISPLAY_ROWS; i++ {
		display[i] = make([]int, DISPLAY_COLS)
	}

	// Initialize registers
	chip8.registers = make(map[nibble]uint8)
	registers := chip8.registers
	for i := 0; i < 16; i++ {
		registers[nibble(i)] = 0
	}

}

// TODO: add a ticker and subscription mechanism depending on how implementation pans out
type TimerRegister struct {
	val byte
}
