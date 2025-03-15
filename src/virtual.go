package chip8

// ------------------------------------------------
// Virtual hardware used by the CHIP-8
// ------------------------------------------------

const (
	RAM          = 4096
	STACK_SIZE   = 100
	DISPLAY_COLS = 64
	DISPLAY_ROWS = 32
)

var memory = make([]byte, RAM)

var stack = make([]uint16, STACK_SIZE)

var display = make([][]int, DISPLAY_ROWS)

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

var PC uint16
var I uint16
var registers = map[nibble]uint8{
	NIBBLE_0: 0,
	NIBBLE_1: 0,
	NIBBLE_2: 0,
	NIBBLE_3: 0,
	NIBBLE_4: 0,
	NIBBLE_5: 0,
	NIBBLE_6: 0,
	NIBBLE_7: 0,
	NIBBLE_8: 0,
	NIBBLE_9: 0,
	NIBBLE_A: 0,
	NIBBLE_B: 0,
	NIBBLE_C: 0,
	NIBBLE_D: 0,
	NIBBLE_E: 0,
	NIBBLE_F: 0,
}

// TODO: add a ticker and subscription mechanism depending on how implementation pans out
type TimerRegister struct {
	val byte
}

// ------------------------------------------------
// 1. First 512 bytes in memory used to have the interpreter, that is no longer true as our interpreter runs in Go space. We can use first 512 for storing the font sprites. 60 bytes between 80-159 (0x050-0x09F)
// 2. Display is modelled as a 2D boolean array with 64 columns and 32 rows. To initialize the rows, we need a loop.
// ------------------------------------------------
func initialize() {

	// Initialize fonts in memory
	start := 0x50
	end := 0x9F
	for i, j := 0, start; i <= end; i, j = i+1, j+1 {
		val := font[i]
		location := j

		memory[location] = val
	}

	// Initialize the rows for display, each row has 'COLS' number of elems
	for i := 0; i < DISPLAY_ROWS; i++ {
		display[i] = make([]int, DISPLAY_COLS)
	}

}
