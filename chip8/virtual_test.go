package chip8

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewChip8Initialization(t *testing.T) {
	chip8 := NewChip8(false, false, 700)

	// 1. Fonts are initialized correctly in memory from 0x50 to 0x9F
	for i, j := 0, SPRITE_START_LOC; j <= SPRITE_END_LOC; i, j = i+1, j+1 {
		require.Equal(t, font[i], chip8.memory[j])

	}

	// 2. Display has 32 rows and 64 columns
	require.Equal(t, DISPLAY_ROWS, len(chip8.display))
	for _, row := range chip8.display {
		require.Equal(t, DISPLAY_COLS, len(row))
	}

	// 3. Registers are initialized to zero for all 16 registers
	for i := 0; i < 16; i++ {
		require.Equal(t, uint8(0), chip8.registers[nibble(i)])
	}
}
