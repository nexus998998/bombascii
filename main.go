package main

import (
	"fmt"
	"os"
	"time"

	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
)

var width = 50
var height = 25

type charecter struct {
	PreviousPosition []int
	CurrentPosition  []int
	Charecter        string
}
type frames [][]string
type velocity struct {
	X int
	Y int
}
type point struct {
	X int
	Y int
}
type hitbox struct {
	TopLeft     point
	BottomRight point
}
type sprite struct {
	Velocity    velocity
	Charecters  []string // since a string is already a slice of charecters this way it's 2d
	Hitbox      hitbox
	OriginPoint point
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

	sprites := []sprite{
		{
			OriginPoint: point{10, 10},
			Velocity:    velocity{1, 1},
			Hitbox:      hitbox{point{-6, 3}, point{5, -3}},
			Charecters: []string{
				"small slime",
				"+---------+",
				"| o     o |",
				"|  \\___/  |",
				"+---------+",
			},
		},
		{
			OriginPoint: point{10, 10},
			Charecters: []string{
				"  /\\/\\  ",
				" |  | ",
				"  \\  /  ",
				"   \\/   ",
			},
			Hitbox:   hitbox{point{-2, 1}, point{2, -3}},
			Velocity: velocity{1, 1},
		},
	}

	// input listener
	go keyboard.Listen(func(key keys.Key) (stop bool, err error) {
		switch key.Code {
		case keys.Up:
			sprites[0].Velocity.Y += -1
		case keys.Down:
			sprites[0].Velocity.Y += 1
		case keys.Left:
			sprites[0].Velocity.X += -1
		case keys.Right:
			sprites[0].Velocity.X += 1
		case keys.CtrlC:
			fmt.Print("\033[?25h")
			os.Exit(0)
		}
		return false, nil
	})

	for {
		clearScreen()
		// composing the frame
		var frame frames
		for i := 0; i < height; i++ {
			var currentRow []string
			for j := 0; j < width; j++ {
				currentRow = append(currentRow, "  ")
			}
			frame = append(frame, currentRow)

		}

		for x := range sprites {
			futurePositionX := sprites[x].OriginPoint.X + sprites[x].Velocity.X
			futurePositionY := sprites[x].OriginPoint.Y + sprites[x].Velocity.Y
			TopLeftFuturePosition := point{futurePositionX + sprites[x].Hitbox.TopLeft.X, futurePositionY + sprites[x].Hitbox.TopLeft.Y}
			BottomRightFuturePosition := point{futurePositionX + sprites[x].Hitbox.BottomRight.X, futurePositionY + sprites[x].Hitbox.BottomRight.Y}

			if (TopLeftFuturePosition.X < 0) || (BottomRightFuturePosition.X > width-1) {
				futurePositionX += sprites[x].Velocity.X * -2
				sprites[x].Velocity.X *= -1
			}

			if (BottomRightFuturePosition.Y < 0) || (TopLeftFuturePosition.Y > height-1) {
				futurePositionY += sprites[x].Velocity.Y * -2
				sprites[x].Velocity.Y *= -1
			}

			sprites[x].OriginPoint = point{futurePositionX, futurePositionY}

		}

		for x := range sprites {
			originPoint := sprites[x].OriginPoint
			charecters := sprites[x].Charecters
			midRowIndex := len(sprites[x].Charecters) / 2

			for rowNumber, row := range charecters {
				relativeY := rowNumber - midRowIndex
				midCharecterIndex := len(row) / 2
				finalY := relativeY + originPoint.Y
				if (finalY < 0) || (finalY > height-1) {
					continue
				}
				for charecterIndex, charecter := range row {
					relativeX := charecterIndex - midCharecterIndex
					finalX := relativeX + originPoint.X
					if (finalX > width-1) || (finalX < 0) {
						continue
					}
					frame[finalY][finalX] = string(charecter)

				}
			}
		}

		drawFrame(frame)
		time.Sleep(time.Millisecond * 40)
	}
}
