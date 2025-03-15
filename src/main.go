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
			display[i][j] = false
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
