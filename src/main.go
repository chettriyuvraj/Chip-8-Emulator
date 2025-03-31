package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
)

func main() {
	initialize()

	romPath := "/Users/yuvrajchettri/Desktop/Yuvi/Development/Chip-8/src/IBM_Logo.ch8"
	err := load(romPath)
	if err != nil {
		log.Fatal(err)
	}

	// set PC to start of rom
	PC = 0x200

	loop()

}

// ------------------------------------------------
// Loop for fetch-decode-execute cycle
// TODO: set speed for the cycle
// ------------------------------------------------
func loop() {

	for PC < RAM {
		instruction := fetch()
		PC += 2

		switch {
		case instruction == 0x00E0:
			clearDisplay()

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
			addToRegister(x, nn)

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
			clearConsole()
			printDisplay()

		// 2NNN
		case instruction.firstNibble().equals(0x2):
			nnn := instruction.nnn()
			pcToStack()
			jumpTo(nnn)

		// 00EE
		case instruction == 0x00E0:
			poppedInstruction := popStack()
			setPC(poppedInstruction)

		// 3XNN
		case instruction.firstNibble().equals(0x3):
			registerIdx := instruction.x()
			valToCheck := instruction.nn()
			skipInstructionIfRegisterEquals(registerIdx, valToCheck)

		// 4XNN
		case instruction.firstNibble().equals(0x4):
			registerIdx := instruction.x()
			valToCheck := instruction.nn()
			skipInstructionIfRegisterNotEquals(registerIdx, valToCheck)

		// 5XYO
		case instruction.firstNibble().equals(0x5):
			regXIdx := instruction.x()
			regYIdx := instruction.y()
			skipInstructionIfRegistersEqualEachOther(regXIdx, regYIdx)

		// 9XYO
		case instruction.firstNibble().equals(0x9):
			regXIdx := instruction.x()
			regYIdx := instruction.y()
			skipInstructionIfRegistersNotEqualEachOther(regXIdx, regYIdx)

		// 8X set of instructions
		case instruction.firstNibble().equals(0x8):
			logicalAndArithmetic(instruction)

		// BNNN or BXNN
		case instruction.firstNibble().equals(0xB):
			var offsetRegisterIdx nibble

			switch bnnn1 {
			case true: // BNNN
				offsetRegisterIdx = NIBBLE_0
			case false: // BXNN
				offsetRegisterIdx = instruction.x()
			}

			addr := instruction.nnn()
			jumpWithOffset(addr, offsetRegisterIdx)

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

func clearDisplay() {
	for i := range DISPLAY_ROWS {
		for j := range DISPLAY_COLS {
			display[i][j] = 0
		}
	}
}

func jumpTo(instruction uint16) {
	PC = instruction
}

func jumpWithOffset(addr uint16, offsetRegisterIdx nibble) {

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

	x %= DISPLAY_COLS
	y %= DISPLAY_ROWS

	originalX := x

	// VF is the collision register, set it to 0 initially
	registers[NIBBLE_F] = 0

	for n := NIBBLE_0; n < height; n++ {
		curSpritePosition := I + uint16(n)
		spriteVal := memory[curSpritePosition]

		byteIdx := 7
		for x < DISPLAY_COLS && byteIdx >= 0 {
			mask := isBitOn(spriteVal, byteIdx)

			// set collision register
			if mask != 0 && display[y][x] == 1 {
				registers[NIBBLE_F] = 1
			}

			display[y][x] ^= mask

			x += 1
			byteIdx -= 1
		}

		// Set x-coordinate back to its original value
		x = originalX

		// Move y coordinate forward to next row
		y += 1

		// Break if we have reached bottom of the screen i.e. last row
		if y >= DISPLAY_ROWS {
			break
		}

	}

}

func isBitOn(val uint8, idx int) int { // index must only be between between 0 and 7 otherwise modulo-ed
	idx %= 8

	mask := uint8(1) << idx
	return int(val & mask) // will return only 1 or 0
}

// ------------------------------------------------
// Load chip-8 'rom' to memory
// ------------------------------------------------

func load(filepath string) error {
	// this is where programs start in memory for chip-8
	start := 0x200

	f, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("error loading rom to memory: %w", err)
	}

	romData, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("error loading rom to memory: %w", err)
	}

	lenRom := len(romData)
	for i, j := start, 0; i < RAM && j < lenRom; i, j = i+1, j+1 {
		memory[i] = romData[j]
	}

	return nil
}

func printDisplay() {
	for i := 0; i < DISPLAY_ROWS; i++ {
		for j := 0; j < DISPLAY_COLS; j++ {
			if display[i][j] == 0 {
				fmt.Printf("%d", 0)
			} else {
				fmt.Printf("%d", 1)
			}

		}
		fmt.Println()
	}
}

// Function to clear the console screen
func clearConsole() {
	cmd := exec.Command("clear") // "cls" for Windows
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func pcToStack() {
	stack = append(stack, PC)
}

func popStack() uint16 {
	elemCount := len(stack)
	lastIdx := elemCount - 1
	lastElem := stack[lastIdx]

	// reduce stack length by one
	stack = stack[:elemCount-1]

	return lastElem
}

func setPC(instruction uint16) {
	PC = instruction
}

func skipInstructionIfRegisterEquals(registerIdx nibble, val uint8) {
	regVal := registers[registerIdx]
	if regVal == val {
		PC += 2
	}
}

func skipInstructionIfRegisterNotEquals(registerIdx nibble, val uint8) {
	regVal := registers[registerIdx]
	if regVal != val {
		PC += 2
	}
}

func skipInstructionIfRegistersEqualEachOther(regXIdx, regYIdx nibble) {
	regXVal := registers[regXIdx]
	regYVal := registers[regYIdx]
	if regXVal == regYVal {
		PC += 2
	}
}

func skipInstructionIfRegistersNotEqualEachOther(regXIdx, regYIdx nibble) {
	regXVal := registers[regXIdx]
	regYVal := registers[regYIdx]
	if regXVal != regYVal {
		PC += 2
	}
}

// 8X set of instructions
func logicalAndArithmetic(i instruction) {
	x := i.x()
	y := i.y()
	n := i.n()

	switch {
	// Set register vx to vy's val
	case n.equals(0x0):
		regYVal := registers[y]
		setRegister(x, regYVal)

	// VX = VX | VY
	case n.equals(0x1):
		regXVal := registers[x]
		regYVal := registers[y]
		registers[x] = regXVal | regYVal

	// VX = VX & VY
	case n.equals(0x2):
		regXVal := registers[x]
		regYVal := registers[y]
		registers[x] = regXVal & regYVal

	// VX = VX ^ VY
	case n.equals(0x3):
		regXVal := registers[x]
		regYVal := registers[y]
		registers[x] = regXVal ^ regYVal

	// VX = VX + VY and set carry flag if overflow
	case n.equals(0x4):
		regXVal := registers[x]
		regYVal := registers[y]
		newVal := regXVal + regYVal
		registers[x] = newVal
		if (newVal < regXVal) || (newVal < regYVal) { // overflow
			registers[NIBBLE_F] = 1
		} else {
			registers[NIBBLE_F] = 0
		}

	// VX = VX - VY and set carry flag if NO underflow
	case n.equals(0x5):
		regXVal := registers[x]
		regYVal := registers[y]
		registers[x] = regXVal - regYVal
		if regXVal > regYVal { // NO underflow
			registers[NIBBLE_F] = 1
		} else {
			registers[NIBBLE_F] = 0
		}

	// VX = VY - VX and set carry flag if NO underflow
	case n.equals(0x7):
		regXVal := registers[x]
		regYVal := registers[y]
		registers[x] = regYVal - regXVal
		if regYVal > regXVal { // NO underflow
			registers[NIBBLE_F] = 1
		} else {
			registers[NIBBLE_F] = 0
		}

	// Left and right shift
	case n.equals(0x6) || n.equals(0xE):
		if shift1 {
			regYVal := registers[y]
			registers[x] = regYVal
		}

		// right shift
		if n.equals(0x6) {
			// get rightmost bit and set carry flag
			regXVal := registers[x]
			rightMostBit := 0x1 & regXVal
			registers[NIBBLE_F] = rightMostBit
			// right shift
			registers[x] >>= 1
		} else { // left shift
			// get leftmost bit and set carry flag
			regXVal := registers[x]
			leftmostBit := 0x1 & (regXVal >> 7)
			registers[NIBBLE_F] = leftmostBit
			// left shift
			registers[x] <<= 1
		}

	}
}
