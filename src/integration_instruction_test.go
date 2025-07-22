package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInstruction_00E0_ClearScreen(t *testing.T) {
	chip8 := NewChip8(false, false, 700)
	// Fill display with 1s
	for i := 0; i < DISPLAY_ROWS; i++ {
		for j := 0; j < DISPLAY_COLS; j++ {
			chip8.display[i][j] = 1
		}
	}
	// Execute 00E0
	chip8.clearDisplay()
	// Check all display is 0
	for i := 0; i < DISPLAY_ROWS; i++ {
		for j := 0; j < DISPLAY_COLS; j++ {
			require.Equal(t, 0, chip8.display[i][j], "display[%d][%d] should be 0", i, j)
		}
	}
}

func TestInstruction_1NNN_Jump(t *testing.T) {
	chip8 := NewChip8(false, false, 700)
	chip8.PC = 0
	addr := uint16(0x345)
	chip8.jumpTo(addr)
	require.Equal(t, addr, chip8.PC)
}

func TestInstruction_6XNN_SetRegisterVX(t *testing.T) {
	chip8 := NewChip8(false, false, 700)
	x := nibble(0xA)
	val := byte(0x77)
	chip8.setRegister(x, val)
	require.Equal(t, val, chip8.registers[x])
}

func TestInstruction_7XNN_AddToRegisterVX(t *testing.T) {
	chip8 := NewChip8(false, false, 700)
	x := nibble(0x3)
	chip8.setRegister(x, 5)
	chip8.addToRegister(x, 7)
	require.Equal(t, uint8(12), chip8.registers[x])
}

func TestInstruction_ANNN_SetIndexRegister(t *testing.T) {
	chip8 := NewChip8(false, false, 700)
	chip8.I = 0
	addr := uint16(0x2AB)
	chip8.setIndexRegister(addr)
	require.Equal(t, addr, chip8.I)
}

func TestInstruction_DXYN_Draw(t *testing.T) {
	chip8 := NewChip8(false, false, 700)
	// Place a sprite in memory at I
	chip8.I = 0x300
	chip8.memory[0x300] = 0b10000000 // Only leftmost pixel on
	// Set Vx and Vy to (0,0)
	chip8.setRegister(0, 0)
	chip8.setRegister(1, 0)
	// Draw 1 row sprite at (0,0)
	chip8.draw(0, 1, 1)
	// Check that display[0][0] is 1
	require.Equal(t, 1, chip8.display[0][0])
	// All other pixels in first row should be 0
	for j := 1; j < DISPLAY_COLS; j++ {
		require.Equal(t, 0, chip8.display[0][j])
	}
}
