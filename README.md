# CHIP-8

For nostalgic reasons, I want to create a Gameboy emulator.

Bad news: I don't know the heads or tails of how to program an emulator.

r/EmuDev recommends the [Chip-8](https://en.wikipedia.org/wiki/CHIP-8) as the gateway drug of choice for emulator development, so that is what I have built here.

This emulator has 3 ROMs embedded into it - mostly to avoid the hassle of opening a file and loading the ROM in WASM.

You can modify the code to use other Chip-8 ROMs.

## Usage

## Native Go build

```
go build -o emulator

./emulator [PONG | TANK | TETRIS]
```

## WASM build

```
$GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o chip8.wasm

# Start a server which serves the index.html file
$python3 -m http.server 8080

# Go to localhost:8080 where you will find the emulator running on a web-frontend

https://github.com/user-attachments/assets/23af413b-d5ab-4bf2-a86f-49296ea57c62


```

## Troubleshooting

### Tetris

I have attachmed a _wasm_exec.js_ file - you might have to use your own one for the WASM build.

Copy your wasm_exec.js into the base directory using:

```
cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" .
```

### Pong
## In action

https://github.com/user-attachments/assets/0f5f6c04-7aa6-4071-94ee-b1e3623333f2


https://github.com/user-attachments/assets/9f40bfbe-8299-49d9-9486-8da590fab02d





