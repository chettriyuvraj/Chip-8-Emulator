package chip8

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"golang.org/x/exp/rand"
)

func (chip8 *Chip8) Initialize() {
	// Delay timer and sound timer will keep running
	go chip8.initDelaySoundTimers()

	if chip8.speedHz <= 0 {
		chip8.speedHz = 700 // fallback default
	}
}

func (chip8 *Chip8) ExecuteInstruction(instruction instruction) {
	switch {
	case instruction == 0x00E0:
		chip8.clearDisplay()
		chip8.redraw = true

	// 1NNN
	case instruction.firstNibble().equals(0x01):
		nnn := instruction.nnn()
		chip8.jumpTo(nnn)

	// 6XNN
	case instruction.firstNibble().equals(0x06):
		nn := instruction.nn()
		x := instruction.x()
		chip8.setRegister(x, nn)

	// 7XNN
	case instruction.firstNibble().equals(0x07):
		nn := instruction.nn()
		x := instruction.x()
		chip8.addToRegister(x, nn)

	// ANNN
	case instruction.firstNibble().equals(0x0A):
		nnn := instruction.nnn()
		chip8.setIndexRegister(nnn)

	// DXYN
	case instruction.firstNibble().equals(0xD):
		x := instruction.x() // vx register contains the x coordinate
		y := instruction.y() // vy register contains the y coordinate
		n := instruction.n()
		chip8.draw(x, y, n)
		chip8.redraw = true

	// 2NNN
	case instruction.firstNibble().equals(0x2):
		nnn := instruction.nnn()
		// Push the address of the next instruction (PC already points to it)
		chip8.stack = append(chip8.stack, chip8.PC)
		chip8.jumpTo(nnn)

	// 00EE
	case instruction == 0x00EE:
		poppedInstruction := chip8.popStack()
		chip8.setPC(poppedInstruction)

	// 3XNN
	case instruction.firstNibble().equals(0x3):
		registerIdx := instruction.x()
		valToCheck := instruction.nn()
		chip8.skipInstructionIfRegisterEquals(registerIdx, valToCheck)

	// 4XNN
	case instruction.firstNibble().equals(0x4):
		registerIdx := instruction.x()
		valToCheck := instruction.nn()
		chip8.skipInstructionIfRegisterNotEquals(registerIdx, valToCheck)

	// 5XYO
	case instruction.firstNibble().equals(0x5):
		regXIdx := instruction.x()
		regYIdx := instruction.y()
		chip8.skipInstructionIfRegistersEqualEachOther(regXIdx, regYIdx)

	// 9XYO
	case instruction.firstNibble().equals(0x9):
		regXIdx := instruction.x()
		regYIdx := instruction.y()
		chip8.skipInstructionIfRegistersNotEqualEachOther(regXIdx, regYIdx)

	// 8X set of instructions
	case instruction.firstNibble().equals(0x8):
		chip8.logicalAndArithmetic(instruction)

	// BNNN or BXNN
	case instruction.firstNibble().equals(0xB):
		var offsetRegisterIdx nibble

		switch chip8.bnnn1 {
		case true: // BNNN
			offsetRegisterIdx = NIBBLE_0
		case false: // BXNN
			offsetRegisterIdx = instruction.x()
		}

		addr := instruction.nnn()
		chip8.jumpWithOffset(addr, offsetRegisterIdx)

	// CXNN: VX = random byte & NN
	case instruction.firstNibble().equals(0xC):
		nn := instruction.nn()
		x := instruction.x()
		randVal := byte(rand.Intn(256))
		chip8.setRegister(x, randVal&nn)

	// EX9E: Skip next instruction if key in VX is pressed
	case instruction.firstNibble().equals(0xE) && instruction.nn() == 0x9E:
		x := instruction.x()
		vx := chip8.registers[x]
		if chip8.isKeyPressed(vx) {
			chip8.PC += 2
		}

	// EXA1: Skip next instruction if key in VX is NOT pressed
	case instruction.firstNibble().equals(0xE) && instruction.nn() == 0xA1:
		x := instruction.x()
		vx := chip8.registers[x]
		if !chip8.isKeyPressed(vx) {
			chip8.PC += 2
		}

	// FX07: Set VX = delay timer
	case instruction.firstNibble().equals(0xF) && instruction.nn() == 0x07:
		x := instruction.x()
		chip8.setRegister(x, chip8.delayTimer)

	// FX15: Set delay timer = VX
	case instruction.firstNibble().equals(0xF) && instruction.nn() == 0x15:
		x := instruction.x()
		chip8.delayTimer = chip8.registers[x]

	// FX18: Set sound timer = VX
	case instruction.firstNibble().equals(0xF) && instruction.nn() == 0x18:
		x := instruction.x()
		chip8.soundTimer = chip8.registers[x]

	// FX1E: I += VX, set VF to 1 if overflow from 0x0FFF to >= 0x1000, else 0
	case instruction.firstNibble().equals(0xF) && instruction.nn() == 0x1E:
		x := instruction.x()
		vx := uint16(chip8.registers[x])
		oldI := chip8.I
		chip8.I += vx
		if oldI <= 0x0FFF && chip8.I > 0x0FFF {
			chip8.registers[NIBBLE_F] = 1
		} else {
			chip8.registers[NIBBLE_F] = 0
		}

	// FX0A: Wait for key press, store key value in VX
	case instruction.firstNibble().equals(0xF) && instruction.nn() == 0x0A:
		x := instruction.x()
		keyFound := false
		for _, chip8Key := range keys {
			chip8.keyboardMu.Lock()
			state := chip8.keyboardState[uint8(chip8Key)]
			chip8.keyboardMu.Unlock()
			if state {
				chip8.setRegister(x, chip8Key)
				keyFound = true
				break
			}
		}
		if !keyFound {
			// No key pressed, repeat this instruction next cycle
			chip8.PC -= 2
		}

	// FX29: Set I to the location of the sprite for the character in VX
	case instruction.firstNibble().equals(0xF) && instruction.nn() == 0x29:
		x := instruction.x()
		vx := chip8.registers[x] & 0xF // Only the lower 4 bits
		chip8.I = SPRITE_START_LOC + uint16(vx)*5

	// FX33: Store BCD representation of VX at I, I+1, I+2
	case instruction.firstNibble().equals(0xF) && instruction.nn() == 0x33:
		x := instruction.x()
		vx := chip8.registers[x]
		if chip8.I < RAM && chip8.I+1 < RAM && chip8.I+2 < RAM {
			chip8.memory[chip8.I] = vx / 100
			chip8.memory[chip8.I+1] = (vx / 10) % 10
			chip8.memory[chip8.I+2] = vx % 10
		} else {
			fmt.Printf("[WARN] FX33: I out of bounds: I=0x%X\n", chip8.I)
		}

	// FX55: Store V0 through VX in memory starting at I (modern: I unchanged)
	case instruction.firstNibble().equals(0xF) && instruction.nn() == 0x55:
		x := instruction.x()
		if chip8.I+uint16(x) < RAM {
			for i := nibble(0); i <= x; i++ {
				chip8.memory[chip8.I+uint16(i)] = chip8.registers[i]
			}
		} else {
			fmt.Printf("[WARN] FX55: I+X out of bounds: I=0x%X, X=0x%X\n", chip8.I, x)
		}

	// FX65: Load V0 through VX from memory starting at I (modern: I unchanged)
	case instruction.firstNibble().equals(0xF) && instruction.nn() == 0x65:
		x := instruction.x()
		if chip8.I+uint16(x) < RAM {
			for i := nibble(0); i <= x; i++ {
				chip8.registers[i] = chip8.memory[chip8.I+uint16(i)]
			}
		} else {
			fmt.Printf("[WARN] FX65: I+X out of bounds: I=0x%X, X=0x%X\n", chip8.I, x)
		}

	}
}

// ------------------------------------------------
// Fetches the 16-byte instruction
// TODO: What happens if PC overshoots the cycle?
// ------------------------------------------------
func (chip8 *Chip8) Fetch() instruction {
	var rawInstruction uint16

	firstByte := chip8.memory[chip8.PC]
	secondByte := chip8.memory[chip8.PC+1]

	firstByteExtended := uint16(firstByte)
	rawInstruction |= firstByteExtended
	rawInstruction <<= 8

	secondByteExtended := uint16(secondByte)
	rawInstruction |= secondByteExtended

	inst := instruction(rawInstruction)

	return inst
}

func (chip8 *Chip8) clearDisplay() {
	for i := 0; i < DISPLAY_ROWS; i++ {
		for j := 0; j < DISPLAY_COLS; j++ {
			chip8.display[i][j] = 0
		}
	}
}

func (chip8 *Chip8) jumpTo(instruction uint16) {
	chip8.PC = instruction
}

func (chip8 *Chip8) jumpWithOffset(addr uint16, offsetRegisterIdx nibble) {
	offset := chip8.registers[offsetRegisterIdx]
	chip8.PC = addr + uint16(offset)
}

func (chip8 *Chip8) setRegister(registerNum nibble, val byte) {
	_, exists := chip8.registers[registerNum]
	if !exists {
		panic("trying to set invalid vx register")
	}

	chip8.registers[registerNum] = val
}

func (chip8 *Chip8) addToRegister(registerNum nibble, val byte) {
	_, exists := chip8.registers[registerNum]
	if !exists {
		panic("trying to add val to invalid vx register")
	}

	chip8.registers[registerNum] += val
}

func (chip8 *Chip8) setIndexRegister(val uint16) {
	chip8.I = val
}

func (chip8 *Chip8) draw(registerXNo, registerYNo nibble, height nibble) {
	// Get x and y coordinate where sprite will start in display
	x, registerExists := chip8.registers[registerXNo]
	if !registerExists {
		panic("trying to get draw x coordinate from invalid register")
	}

	y, registerExists := chip8.registers[registerYNo]
	if !registerExists {
		panic("trying to get draw y coordinate from invalid register")
	}

	x %= DISPLAY_COLS
	y %= DISPLAY_ROWS

	originalX := x

	// VF is the collision register, set it to 0 initially
	chip8.registers[NIBBLE_F] = 0

	for n := NIBBLE_0; n < height; n++ {
		curSpritePosition := chip8.I + uint16(n)
		spriteVal := chip8.memory[curSpritePosition]

		byteIdx := 7
		for x < DISPLAY_COLS && byteIdx >= 0 {
			mask := isBitOn(spriteVal, byteIdx)

			// set collision register
			if mask != 0 && chip8.display[y][x] == 1 {
				chip8.registers[NIBBLE_F] = 1
			}

			chip8.display[y][x] ^= mask

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

func isBitOn(val uint8, idx int) int {
	idx %= 8
	mask := uint8(1) << idx
	if val&mask != 0 {
		return 1
	}
	return 0
}

func (chip8 *Chip8) printDisplay() {
	for i := 0; i < DISPLAY_ROWS; i++ {
		for j := 0; j < DISPLAY_COLS; j++ {
			if chip8.display[i][j] == 0 {
				fmt.Printf("%d", 0)
			} else {
				fmt.Printf("%d", 1)
			}

		}
		fmt.Println()
	}
}

// Function to clear the console screen
func (chip8 *Chip8) clearConsole() {
	cmd := exec.Command("clear") // "cls" for Windows
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func (chip8 *Chip8) pcToStack() {
	chip8.stack = append(chip8.stack, chip8.PC)
}

func (chip8 *Chip8) popStack() uint16 {
	elemCount := len(chip8.stack)
	lastIdx := elemCount - 1
	lastElem := chip8.stack[lastIdx]

	// reduce stack length by one
	chip8.stack = chip8.stack[:elemCount-1]

	return lastElem
}

func (chip8 *Chip8) setPC(instruction uint16) {
	chip8.PC = instruction
}

func (chip8 *Chip8) skipInstructionIfRegisterEquals(registerIdx nibble, val uint8) {
	regVal := chip8.registers[registerIdx]
	if regVal == val {
		chip8.PC += 2
	}
}

func (chip8 *Chip8) skipInstructionIfRegisterNotEquals(registerIdx nibble, val uint8) {
	regVal := chip8.registers[registerIdx]
	if regVal != val {
		chip8.PC += 2
	}
}

func (chip8 *Chip8) skipInstructionIfRegistersEqualEachOther(regXIdx, regYIdx nibble) {
	regXVal := chip8.registers[regXIdx]
	regYVal := chip8.registers[regYIdx]
	if regXVal == regYVal {
		chip8.PC += 2
	}
}

func (chip8 *Chip8) skipInstructionIfRegistersNotEqualEachOther(regXIdx, regYIdx nibble) {
	regXVal := chip8.registers[regXIdx]
	regYVal := chip8.registers[regYIdx]
	if regXVal != regYVal {
		chip8.PC += 2
	}
}

// 8X set of instructions
func (chip8 *Chip8) logicalAndArithmetic(i instruction) {
	x := i.x()
	y := i.y()
	n := i.n()

	switch {
	// Set register vx to vy's val
	case n.equals(0x0):
		regYVal := chip8.registers[y]
		chip8.setRegister(x, regYVal)

	// VX = VX | VY
	case n.equals(0x1):
		regXVal := chip8.registers[x]
		regYVal := chip8.registers[y]
		chip8.registers[x] = regXVal | regYVal

	// VX = VX & VY
	case n.equals(0x2):
		regXVal := chip8.registers[x]
		regYVal := chip8.registers[y]
		chip8.registers[x] = regXVal & regYVal

	// VX = VX ^ VY
	case n.equals(0x3):
		regXVal := chip8.registers[x]
		regYVal := chip8.registers[y]
		chip8.registers[x] = regXVal ^ regYVal

	// VX = VX + VY and set carry flag if overflow
	case n.equals(0x4):
		regXVal := chip8.registers[x]
		regYVal := chip8.registers[y]
		newVal := regXVal + regYVal
		chip8.registers[x] = newVal
		if (newVal < regXVal) || (newVal < regYVal) { // overflow
			chip8.registers[NIBBLE_F] = 1
		} else {
			chip8.registers[NIBBLE_F] = 0
		}

	// VX = VX - VY and set carry flag if NO underflow
	case n.equals(0x5):
		regXVal := chip8.registers[x]
		regYVal := chip8.registers[y]
		chip8.registers[x] = regXVal - regYVal
		if regXVal > regYVal { // NO underflow
			chip8.registers[NIBBLE_F] = 1
		} else {
			chip8.registers[NIBBLE_F] = 0
		}

	// VX = VY - VX and set carry flag if NO underflow
	case n.equals(0x7):
		regXVal := chip8.registers[x]
		regYVal := chip8.registers[y]
		chip8.registers[x] = regYVal - regXVal
		if regYVal > regXVal { // NO underflow
			chip8.registers[NIBBLE_F] = 1
		} else {
			chip8.registers[NIBBLE_F] = 0
		}

	// Left and right shift
	case n.equals(0x6) || n.equals(0xE):
		if chip8.shift1 {
			regYVal := chip8.registers[y]
			chip8.registers[x] = regYVal
		}

		// right shift
		if n.equals(0x6) {
			// get rightmost bit and set carry flag
			regXVal := chip8.registers[x]
			rightMostBit := 0x1 & regXVal
			chip8.registers[NIBBLE_F] = rightMostBit
			// right shift
			chip8.registers[x] >>= 1
		} else { // left shift
			// get leftmost bit and set carry flag
			regXVal := chip8.registers[x]
			leftmostBit := 0x1 & (regXVal >> 7)
			chip8.registers[NIBBLE_F] = leftmostBit
			// left shift
			chip8.registers[x] <<= 1
		}

	}
}

// Add new function for keyboard state management
func (chip8 *Chip8) UpdateKeyboardState(key uint8, state bool) {
	chip8.keyboardMu.Lock()
	chip8.keyboardState[key] = state
	chip8.keyboardMu.Unlock()
}

// Update isKeyPressed to use internal state
func (chip8 *Chip8) isKeyPressed(key uint8) bool {
	chip8.keyboardMu.Lock()
	state := chip8.keyboardState[key]
	chip8.keyboardMu.Unlock()
	return state
}

func (chip8 *Chip8) initDelaySoundTimers() {
	delaySpeedHz := 60
	delayTicker := time.NewTicker(time.Second / time.Duration(delaySpeedHz))
	defer delayTicker.Stop()
	for {
		select {
		case <-delayTicker.C:
			if chip8.delayTimer > 0 {
				chip8.delayTimer -= 1
			}
			if chip8.soundTimer > 0 {
				chip8.soundTimer -= 1
				fmt.Print("\a") // ASCII Bell character - make computer beep as long as > 0
			}
		}
	}
}
