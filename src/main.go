package chip8

func main() {
	loop()
}

// ------------------------------------------------
// Loop for fetch-decode-execute cycle
// TODO: set speed for the cycle
// ------------------------------------------------
func loop() {

	for {
		instruction := fetch()

		switch {
		case instruction == 0x00E0:
			clearScreen()

		// 1NNN
		case instruction.firstNibble().equals(0x01):
			nnn := instruction.nnn()
			jumpTo(nnn)

		// 6XNN
		case instruction.firstNibble().equals(0x06):
			nn := instruction.nn()
			x := instruction.x()
			setRegister(x, nn)

		// 7XNN
		case instruction.firstNibble().equals(0x07):
			nn := instruction.nn()
			x := instruction.x()
			setRegister(x, nn)

		// ANNN
		case instruction.firstNibble().equals(0x0A):
			nnn := instruction.nnn()
			setIndexRegister(nnn)

		// DXYN
		case instruction.firstNibble().equals(0xD):
			x := instruction.x() // vx register contains the x coordinate
			y := instruction.y() // vy register contains the y coordinate
			n := instruction.n()
			draw(x, y, n)

		}
	}

}

// ------------------------------------------------
// Fetches the 16-byte instruction
// TODO: What happens if PC overshoots the cycle?
// ------------------------------------------------
func fetch() instruction {
	var rawInstruction uint16

	firstByte := memory[PC]
	secondByte := memory[PC+1]

	firstByteExtended := uint16(firstByte)
	rawInstruction |= firstByteExtended
	rawInstruction <<= 8

	secondByteExtended := uint16(secondByte)
	rawInstruction |= secondByteExtended

	instruction := instruction(rawInstruction)

	return instruction
}

func clearScreen() {
	for i := range DISPLAY_ROWS {
		for j := range DISPLAY_COLS {
			display[i][j] = 0
		}
	}
}

func jumpTo(instruction uint16) {
	PC = instruction
}

func setRegister(registerNum nibble, val byte) {
	_, exists := registers[registerNum]
	if !exists {
		panic("trying to set invalid vx register")
	}

	registers[registerNum] = val
}

func addToRegister(registerNum nibble, val byte) {
	_, exists := registers[registerNum]
	if !exists {
		panic("trying to add val to invalid vx register")
	}

	registers[registerNum] += val
}

func setIndexRegister(val uint16) {
	I = val
}

func draw(registerXNo, registerYNo nibble, height nibble) {
	// Get x and y coordinate where sprite will start in display
	x, registerExists := registers[registerXNo]
	if !registerExists {
		panic("trying to get draw x coordinate from invalid register")
	}

	y, registerExists := registers[registerYNo]
	if !registerExists {
		panic("trying to get draw y coordinate from invalid register")
	}

	x %= DISPLAY_ROWS
	y %= DISPLAY_COLS

	// VF is the collision register, set it to 0 initially
	registers[NIBBLE_F] = 0

	for n := NIBBLE_0; n < height; n++ {
		curSpritePosition := nibble(I) + n
		spriteVal := memory[curSpritePosition+n]

		byteIdx := 7
		for y < DISPLAY_COLS {
			mask := isBitOn(spriteVal, byteIdx)

			// set collision register
			if mask == 1 && display[x][y] == 1 {
				registers[NIBBLE_F] = 1
			}

			display[x][y] ^= mask

			y += 1
			byteIdx -= 1
		}

		// Set y to the 0th column
		y %= DISPLAY_COLS

		// Move x forward to next row
		x += 1

		// Break if we have reached bottom of the screen i.e. last row
		if x%DISPLAY_ROWS == 0 {
			break
		}

	}

}

func isBitOn(val uint8, idx int) int { // index must only be between between 0 and 7 otherwise modulo-ed
	idx %= 8

	mask := uint8(1) << idx
	return int(val & mask) // will return only 1 or 0
}
