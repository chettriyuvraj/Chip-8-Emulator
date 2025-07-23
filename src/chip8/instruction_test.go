package chip8

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInstruction_firstNibble(t *testing.T) {
	tests := []struct {
		name   string
		instr  instruction
		expect nibble
	}{
		{"all zeros", 0x0000, 0},
		{"first nibble set", 0xF000, 0xF},
		{"random value", 0xA123, 0xA},
		{"all ones", 0xFFFF, 0xF},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expect, tt.instr.firstNibble())
		})
	}
}

func TestInstruction_x(t *testing.T) {
	tests := []struct {
		name   string
		instr  instruction
		expect nibble
	}{
		{"all zeros", 0x0000, 0},
		{"x nibble set", 0x0F00, 0xF},
		{"random value", 0xA123, 0x1},
		{"all ones", 0xFFFF, 0xF},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expect, tt.instr.x())
		})
	}
}

func TestInstruction_y(t *testing.T) {
	tests := []struct {
		name   string
		instr  instruction
		expect nibble
	}{
		{"all zeros", 0x0000, 0},
		{"y nibble set", 0x00F0, 0xF},
		{"random value", 0xA123, 0x2},
		{"all ones", 0xFFFF, 0xF},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expect, tt.instr.y())
		})
	}
}

func TestInstruction_n(t *testing.T) {
	tests := []struct {
		name   string
		instr  instruction
		expect nibble
	}{
		{"all zeros", 0x0000, 0},
		{"n nibble set", 0x000F, 0xF},
		{"random value", 0xA123, 0x3},
		{"all ones", 0xFFFF, 0xF},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expect, tt.instr.n())
		})
	}
}

func TestInstruction_nn(t *testing.T) {
	tests := []struct {
		name   string
		instr  instruction
		expect byte
	}{
		{"all zeros", 0x0000, 0x00},
		{"nn set", 0x00FF, 0xFF},
		{"random value", 0xA123, 0x23},
		{"all ones", 0xFFFF, 0xFF},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expect, tt.instr.nn())
		})
	}
}

func TestInstruction_nnn(t *testing.T) {
	tests := []struct {
		name   string
		instr  instruction
		expect uint16
	}{
		{"all zeros", 0x0000, 0x000},
		{"nnn set", 0x0FFF, 0xFFF},
		{"random value", 0xA123, 0x123},
		{"all ones", 0xFFFF, 0xFFF},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expect, tt.instr.nnn())
		})
	}
}

func TestNibble_Equals(t *testing.T) {
	tests := []struct {
		name   string
		n      nibble
		val    uint8
		expect bool
	}{
		{"equal zero", 0, 0, true},
		{"equal max", 0xF, 0xF, true},
		{"not equal", 0xA, 0xB, false},
		{"not equal high", 0xF, 0x0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expect, tt.n.equals(tt.val))
		})
	}
}
