package main

import (
	"fmt"
	"os"
	"time"

	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
)

var width = 80
var height = 40
var dirbomboclatX = 0
var dirbomboclatY = 0

type charecter struct {
	CurrentPosition []int
	Charecter       string
}

type frames [][]string

type velocity struct {
	X int
	Y int
}

type sprite struct {
	Charecters []charecter
	Velocity   velocity
}

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
		CurrentPosition: []int{1, 1},
		Charecter:       "# ",
	}

	charecter2 := charecter{
		CurrentPosition: []int{1, 2},
		Charecter:       "# ",
	}

	charecter3 := charecter{
		CurrentPosition: []int{2, 1},
		Charecter:       "# ",
	}
	charecter4 := charecter{
		CurrentPosition: []int{2, 2},
		Charecter:       "# ",
	}

	sprite1 := sprite{
		Charecters: []charecter{charecter1, charecter2, charecter3, charecter4},
		Velocity: velocity{
			X: 0,
			Y: 0,
		},
	}

	// input listener
	go keyboard.Listen(func(key keys.Key) (stop bool, err error) {
		switch key.Code {
		case keys.Up:
			dirbomboclatX = 0
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

		// composing the frame

		var frame frames
		for i := 0; i < height; i++ {

			var currentRow []string
			for j := 0; j < width; j++ {
				currentRow = append(currentRow, " ")
			}
			frame = append(frame, currentRow)
		}

		// out of bonds resetting
		for i := range sprite1.Charecters {

			futurePositionX := sprite1.Charecters[i].CurrentPosition[0] + sprite1.Velocity.X
			futurePositionY := sprite1.Charecters[i].CurrentPosition[1] + sprite1.Velocity.Y

			// bounces back if out of bonds for x
			if (futurePositionX > width-2) ||
				(futurePositionX < 0) {
				for j := 0; j < i; j++ {
					sprite1.Charecters[j].CurrentPosition[0] += sprite1.Velocity.X * -2
				}
				futurePositionX += sprite1.Velocity.X * -2
				sprite1.Velocity.X *= -1
			}
			// bounces back if out of bonds for y
			if (futurePositionY > height-2) ||
				(futurePositionY < 0) {

				for j := 0; j < i; j++ {

					sprite1.Charecters[j].CurrentPosition[1] += sprite1.Velocity.Y * -2
				}
				futurePositionY += sprite1.Velocity.Y * -2
				sprite1.Velocity.Y *= -1
			}
			sprite1.Charecters[i].CurrentPosition[0] = futurePositionX
			sprite1.Charecters[i].CurrentPosition[1] = futurePositionY
		}

		for j := range sprite1.Charecters {
			frame[sprite1.Charecters[j].CurrentPosition[1]][sprite1.Charecters[j].CurrentPosition[0]] = sprite1.Charecters[j].Charecter
		}

		drawFrame(frame)
		time.Sleep(time.Millisecond * 10)
		clearScreen()
	}
}
