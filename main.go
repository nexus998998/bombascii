package main

import (
	"fmt"
	"time"
)

var width = 20
var height = 40

func drawFrame() {
	frame := ""
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			frame += "x"

		}
		frame += "\n"
	}
	fmt.Println(frame)

}
func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func main() {
	for {

		drawFrame()
		time.Sleep(time.Duration(0.05 * float64(time.Second)))
		clearScreen()
	}
}
