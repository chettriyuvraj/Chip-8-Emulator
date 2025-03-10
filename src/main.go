package chip8

func main() {

}

// ------------------------------------------------
// Loop for fetch-decode-execute cycle
// TODO: set speed for the cycle
// ------------------------------------------------
func loop() {

	for {
		instruction := fetch()
	}

}

// ------------------------------------------------
// Fetches the 16-byte instruction
// TODO: What happens if PC overshoots the cycle?
// ------------------------------------------------
func fetch() (instruction uint16) {
	firstByte := memory[PC]
	secondByte := memory[PC+1]

	firstByteExtended := uint16(firstByte)
	instruction |= firstByteExtended
	instruction <<= 8

	secondByteExtended := uint16(secondByte)
	instruction |= secondByteExtended

	return instruction
}
