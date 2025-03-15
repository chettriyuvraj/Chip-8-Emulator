package chip8

// ------------------------------------------------
// Methods for extracting emulator instruction
// into constituent parts
// ------------------------------------------------

type nibble uint8 // leftmost 4 bits will always be zero

type instruction uint16

const (
	NIBBLE_0 = nibble(0)
	NIBBLE_1 = nibble(1)
	NIBBLE_2 = nibble(2)
	NIBBLE_3 = nibble(3)
	NIBBLE_4 = nibble(4)
	NIBBLE_5 = nibble(5)
	NIBBLE_6 = nibble(6)
	NIBBLE_7 = nibble(7)
	NIBBLE_8 = nibble(8)
	NIBBLE_9 = nibble(9)
	NIBBLE_A = nibble(10)
	NIBBLE_B = nibble(11)
	NIBBLE_C = nibble(12)
	NIBBLE_D = nibble(13)
	NIBBLE_E = nibble(14)
	NIBBLE_F = nibble(15)
)

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

	// or
	// y := i.y() // third nibble
	// n := i.n() // fourth nibble
	// y <<= 4
	// nn := byte(y | n)
	// return nn
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

	// or
	// x := i.x()
	// y := i.y() // third nibble
	// n := i.n() // fourth nibble
	// xExtended := uint16(x)
	// yExtended := uint16(y)
	// nExtended := uint16(n)
	// xExtended <<= 8
	// yExtended <<= 4
	// nnn := xExtended | yExtended | nExtended
	// return nnn
}

func (n nibble) equals(val uint8) bool {
	nVal := uint8(n)
	return nVal == val
}
