package chip8

// ------------------------------------------------
// Methods for extracting emulator instruction
// into constituent parts
// ------------------------------------------------

type nibble uint8 // leftmost 4 bits will always be zero

type instruction uint16

// ------------------------------------------------
// First nibble doesn't seem to have any specific name unlike other nibbles
// ------------------------------------------------
func (i instruction) firstNibble() nibble {
	mask := uint16(0xF000)
	val := uint16(i)
	res := val & mask
	res >>= 12
	firstNibble := nibble(res)
	return firstNibble
}

// ------------------------------------------------
// The second nibble is called x
// ------------------------------------------------
func (i instruction) x() nibble {
	mask := uint16(0x0F00)
	val := uint16(i)
	res := val & mask
	res >>= 8
	x := nibble(res)
	return x
}

// ------------------------------------------------
// The third nibble is called y
// ------------------------------------------------
func (i instruction) y() nibble {
	mask := uint16(0x00F0)
	val := uint16(i)
	res := val & mask
	res >>= 4
	y := nibble(res)
	return y
}

// ------------------------------------------------
// The fourth nibble is called n
// ------------------------------------------------
func (i instruction) n() nibble {
	mask := uint16(0x000F)
	val := uint16(i)
	res := val & mask
	n := nibble(res)
	return n
}

// ------------------------------------------------
// The second byte is called nn
// ------------------------------------------------
func (i instruction) nn() byte {
	mask := uint16(0x0FF0)
	val := uint16(i)
	res := val & mask
	res >>= 4
	nn := byte(res)
	return nn
}

// ------------------------------------------------
// The second, third and fourth nibble - leftmost nibble is always 0
// ------------------------------------------------
func (i instruction) nnn() uint16 {
	mask := uint16(0x0FFF)
	val := uint16(i)
	res := val & mask
	nnn := res
	return nnn
}
