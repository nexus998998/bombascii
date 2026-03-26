package main

import (
	"fmt"
	"time"
	"os"

	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
)

var width = 20
var height = 40
var dirbomboclatX = 0
var dirbomboclatY = 0

func drawFrame() {
	frame := ""
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			frame += "x"
		}
		frame += "\r\n"
	}
	fmt.Printf("dirbomboclatX: %d, dirbomboclatY: %d\r\n", dirbomboclatX, dirbomboclatY)
	fmt.Println(frame)
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func main() {


	for {
		go keyboard.Listen(func(key keys.Key) (stop bool, err error) {
		switch key.Code {
		case keys.Up:
			dirbomboclatX = -1
			dirbomboclatY = 0
		case keys.Down:
			dirbomboclatX = 1
			dirbomboclatY = 0
		case keys.Left:
			dirbomboclatX = 0
			dirbomboclatY = -1
		case keys.Right:
			dirbomboclatX = 0
			dirbomboclatY = 1
		case keys.CtrlC:
			fmt.Print("\033[?25h")
			os.Exit(0)
		}
		return false, nil
		})


		drawFrame()
		time.Sleep(time.Duration(0.05 * float64(time.Second)))
		clearScreen()
	}
}