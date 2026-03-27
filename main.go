package main

import (
	"fmt"
	"os"
	"time"

	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
)

var width = 40
var height = 20
var dirbomboclatX int
var dirbomboclatY int

type charecter struct {
	CurrentPosition []int
	Velocity        []int
	Charecter       string
}

type frames [][]string
type sprite []charecter

func drawFrame(frame frames) {
	var frameString string
	for _, row := range frame {
		for _, charecter := range row {
			frameString += charecter
		}
		frameString += "\r\n"
	}
	fmt.Print(frameString)
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func main() {

	charecter1 := charecter{
		CurrentPosition: []int{1, 0},
		Velocity:        []int{1, 0},
	}

	charecter2 := charecter{
		CurrentPosition: []int{1, 1},
		Velocity:        []int{1, 0},
	}

	charecter3 := charecter{
		CurrentPosition: []int{1, 2},
		Velocity:        []int{1, 0},
	}

	sprite1 := sprite{charecter1, charecter2, charecter3}

	// input listener
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

		// composing the frame

		var frame frames
		for i := 0; i < height; i++ {

			var currentRow []string
			for j := 0; j < width; j++ {
				currentRow = append(currentRow, " ")
			}
			frame = append(frame, currentRow)
		}

		for i := range sprite1 {
			if (sprite1[i].CurrentPosition[0] >= width-1 && sprite1[i].Velocity[0] > 0) ||
				(sprite1[i].CurrentPosition[0] <= 0 && sprite1[i].Velocity[0] < 0) {
				sprite1[i].Velocity[0] *= -1
			}
			sprite1[i].CurrentPosition[0] += sprite1[i].Velocity[0]
			sprite1[i].CurrentPosition[1] += sprite1[i].Velocity[1]

			frame[sprite1[i].CurrentPosition[1]][sprite1[i].CurrentPosition[0]] = "#"
		}

		drawFrame(frame)
		time.Sleep(time.Millisecond * 10)
		clearScreen()
	}
}
