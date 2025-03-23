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
