package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/gen2brain/malgo"
	// "atomicgo.dev/keyboard"
	// "atomicgo.dev/keyboard/keys"
)

var green = "\033[32m"
var yellow = "\033[33m"
var blue = "\033[34m"
var red = "\033[31m"
var reset = "\033[0m"
var white = "\033[37m"

var width = 110
var height = 30
var impactFrames = 10
var frameRate = time.Millisecond * 13
var impactColor = white

var sfxData []byte
var sfxPos int
var sfxPlaying bool
var ctx *malgo.AllocatedContext

var freezeTime int
var frozenFrame frames

func initAudio() {
	f, _ := os.Open("hit.wav")
	defer f.Close()
	f.Seek(44, 0)
	sfxData, _ = io.ReadAll(f)

	ctx, _ = malgo.InitContext(nil, malgo.ContextConfig{}, nil)

	deviceConfig := malgo.DefaultDeviceConfig(malgo.Playback)
	deviceConfig.Playback.Format = malgo.FormatS16
	deviceConfig.Playback.Channels = 2
	deviceConfig.SampleRate = 44100

	sfxCallbacks := malgo.DeviceCallbacks{
		Data: func(pOutput, pInput []byte, frameCount uint32) {
			if !sfxPlaying {
				return
			}
			bytesNeeded := int(frameCount) * 4
			remaining := len(sfxData) - sfxPos
			if remaining <= 0 {
				sfxPlaying = false
				return
			}
			if bytesNeeded > remaining {
				bytesNeeded = remaining
			}
			copy(pOutput, sfxData[sfxPos:sfxPos+bytesNeeded])
			sfxPos += bytesNeeded
		},
	}

	sfxDevice, _ := malgo.InitDevice(ctx.Context, deviceConfig, sfxCallbacks)
	sfxDevice.Start()
}

func playSound() {
	sfxPos = 0
	sfxPlaying = true
}

type Timer struct {
	CurrentTime  int
	OriginalTime int
}

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

type TimedEvent struct {
	Timer Timer
	Event func()
}

type sprite struct {
	Velocity       velocity
	Charecters     []string // since a string is already a slice of charecters this way it's 2d
	HurtCharecters []string
	Hitbox         hitbox
	OriginPoint    point
	AbilityTimer   Timer
	Health         int
	Ability        func(s *sprite)
	CollisionFunc  func(meSprite *sprite, hitSprite *sprite)
	Alive          bool
	Name           string
	Color          string
	HurtColor      string
}

var sprites []sprite

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

func drawSprite(charecters []string, originPoint point, color string, frame frames) {

	midRowIndex := len(charecters) / 2
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
			frame[finalY][finalX] = color + string(charecter) + reset

		}
	}
}

func (s sprite) relativeHitbox() hitbox {
	return hitbox{
		TopLeft:     point{s.OriginPoint.X + s.Hitbox.TopLeft.X, s.OriginPoint.Y + s.Hitbox.TopLeft.Y},
		BottomRight: point{s.OriginPoint.X + s.Hitbox.BottomRight.X, s.OriginPoint.Y + s.Hitbox.BottomRight.Y},
	}
}

func main() {

	initAudio()

	SlimeSprite := sprite{
		OriginPoint: point{10, 10},
		Velocity:    velocity{1, 1},
		Hitbox:      hitbox{point{-3, -3}, point{3, 3}},
		Charecters: []string{
			"small slime",
			"+---------+",
			"| o     o |",
			"|  \\___/  |",
			"+---------+",
		},
		HurtCharecters: []string{
			"small slime",
			"+---------+",
			"| o ___ o |",
			"|  /   \\  |",
			"+---------+",
		},
		AbilityTimer: Timer{OriginalTime: 30, CurrentTime: 30},
		Health:       200,
		Ability: func(s *sprite) {
			return
		},
		CollisionFunc: func(meSprite, hitSprite *sprite) {
			hitSprite.Health -= 20
			meSprite.Health += 5
			return
		},
		Color:     green,
		HurtColor: yellow,
		Alive:     true,
		Name:      "slime",
	}

	qwqSprite := sprite{
		OriginPoint: point{30, 10},
		Charecters: []string{
			"  .oOOOOo.  ",
			" oOOOOOOOOo ",
			"oOOQQOOOOOOo",
			"oQQQQOOOOOOo",
			"oQQQQOOOOOXo",
			"oQQQQOOOOOXo",
			" oOOOOOOOXo ",
			"  'oOOOOo'  ",
		},
		HurtCharecters: []string{
			"  .oOOOo.  ",
			" oOOOO/OOo ",
			"oO\\QQ/OOOOo",
			"oQQQ'O#OOOo",
			"o#;#''OOOXo",
			"o#####O\\\\Xo",
			" oO''#///o ",
			"  'oOOOo'  ",
		},
		Color:        blue,
		HurtColor:    red,
		Hitbox:       hitbox{point{-6, -4}, point{6, 4}},
		Velocity:     velocity{1, 1},
		AbilityTimer: Timer{OriginalTime: 40, CurrentTime: 40},
		Health:       100,
		Ability: func(s *sprite) {
			s.Health += 5
			for i := range sprites {
				if sprites[i].Name == "qwq" {
					continue
				}
				sprites[i].Health -= 5
			}
		},
		CollisionFunc: func(meSprite *sprite, hitSprite *sprite) {
			return
		},
		Alive: true,
		Name:  "qwq",
	}

	sprites = []sprite{
		SlimeSprite,
		qwqSprite,
	}

	// input listener
	// go keyboard.Listen(func(key keys.Key) (stop bool, err error) {
	// 	switch key.Code {
	// 	case keys.Up:
	// 		sprites[0].Velocity.Y += -1
	// 	case keys.Down:
	// 		sprites[0].Velocity.Y += 1
	// 	case keys.Left:
	// 		sprites[0].Velocity.X += -1
	// 	case keys.Right:
	// 		sprites[0].Velocity.X += 1
	// 	case keys.CtrlC:
	// 		fmt.Print("\033[?25h")
	// 		os.Exit(0)
	// 	}
	// 	return false, nil
	// })

	for {
		clearScreen()

		if freezeTime > 0 {
			freezeTime--
			drawFrame(frozenFrame)
			time.Sleep(frameRate)
			continue
		}

		// composing the frame
		var frame frames
		for i := 0; i < height; i++ {
			var currentRow []string
			for j := 0; j < width; j++ {
				currentRow = append(currentRow, " ")
			}
			frame = append(frame, currentRow)

		}

		alive := []sprite{}

		for x := range sprites {
			if sprites[x].Health <= 0 {
				continue
			}

			alive = append(alive, sprites[x])
		}

		sprites = alive

		for x := range sprites {

			fmt.Printf("%s's health : %d\n", sprites[x].Name, sprites[x].Health)

			if sprites[x].AbilityTimer.CurrentTime == 0 {

				sprites[x].Ability(&sprites[x])
				sprites[x].AbilityTimer.CurrentTime = sprites[x].AbilityTimer.OriginalTime
			}

			sprites[x].AbilityTimer.CurrentTime--

			futurePositionX := sprites[x].OriginPoint.X + sprites[x].Velocity.X
			futurePositionY := sprites[x].OriginPoint.Y + sprites[x].Velocity.Y
			TopLeftFuturePosition := point{futurePositionX + sprites[x].Hitbox.TopLeft.X, futurePositionY + sprites[x].Hitbox.TopLeft.Y}
			BottomRightFuturePosition := point{futurePositionX + sprites[x].Hitbox.BottomRight.X, futurePositionY + sprites[x].Hitbox.BottomRight.Y}

			if (TopLeftFuturePosition.X < 0) || (BottomRightFuturePosition.X > width-1) {
				futurePositionX += sprites[x].Velocity.X * -2
				sprites[x].Velocity.X *= -1
			}

			if (BottomRightFuturePosition.Y > height-1) || (TopLeftFuturePosition.Y < 0) {
				futurePositionY += sprites[x].Velocity.Y * -2
				sprites[x].Velocity.Y *= -1
			}

			sprites[x].OriginPoint = point{futurePositionX, futurePositionY}

		}

		alreadyCollidedSprites := make(map[int]int)

		for x := range sprites {
			// checking for collisions

			checkerHitbox := sprites[x].relativeHitbox()
			for j := range sprites {
				if alreadyCollidedSprites[j] == x {
					continue
				}
				if j == x {
					continue
				}
				checkedHitbox := sprites[j].relativeHitbox()
				collisionBoxWidth := min(checkerHitbox.BottomRight.X, checkedHitbox.BottomRight.X) - max(checkerHitbox.TopLeft.X, checkedHitbox.TopLeft.X)
				if collisionBoxWidth < 0 || collisionBoxWidth > checkerHitbox.BottomRight.X-checkerHitbox.TopLeft.X {
					continue
				}
				collisionBoxHeight := min(checkerHitbox.BottomRight.Y, checkedHitbox.BottomRight.Y) - max(checkerHitbox.TopLeft.Y, checkedHitbox.TopLeft.Y)
				if collisionBoxHeight < 0 || collisionBoxHeight > checkerHitbox.BottomRight.Y-checkerHitbox.TopLeft.Y {
					continue
				}
				// impact between 2 asciis here
				playSound()
				freezeTime = impactFrames
				sprites[x].CollisionFunc(&sprites[x], &sprites[j])
				sprites[j].CollisionFunc(&sprites[j], &sprites[x])

				if collisionBoxHeight == collisionBoxWidth {
					sprites[x].Velocity.Y *= -1
					sprites[j].Velocity.Y *= -1
					sprites[x].Velocity.X *= -1
					sprites[j].Velocity.X *= -1

					sprites[x].OriginPoint = point{
						sprites[x].OriginPoint.X + sprites[x].Velocity.X*2,
						sprites[x].OriginPoint.Y + sprites[x].Velocity.Y*2,
					}
					sprites[j].OriginPoint = point{
						sprites[j].OriginPoint.X + sprites[j].Velocity.X*2,
						sprites[j].OriginPoint.Y + sprites[j].Velocity.Y*2,
					}
					fmt.Println(collisionBoxHeight, collisionBoxWidth)
					alreadyCollidedSprites[j] = x
					continue
				}
				if collisionBoxHeight > collisionBoxWidth {
					sprites[x].Velocity.X *= -1
					sprites[j].Velocity.X *= -1

					sprites[x].OriginPoint = point{
						sprites[x].OriginPoint.X + sprites[x].Velocity.X*2,
						sprites[x].OriginPoint.Y,
					}
					sprites[j].OriginPoint = point{
						sprites[j].OriginPoint.X + sprites[j].Velocity.X*2,
						sprites[j].OriginPoint.Y,
					}
					continue
				}
				sprites[x].Velocity.Y *= -1
				sprites[j].Velocity.Y *= -1

				sprites[x].OriginPoint = point{
					sprites[x].OriginPoint.X,
					sprites[x].OriginPoint.Y + sprites[x].Velocity.Y*2,
				}
				sprites[j].OriginPoint = point{
					sprites[j].OriginPoint.X,
					sprites[j].OriginPoint.Y + sprites[j].Velocity.Y*2,
				}
				fmt.Println(collisionBoxHeight, collisionBoxWidth)
				alreadyCollidedSprites[j] = x

			}
		}

		for x := range sprites {
			// find hitbox and draw sprite
			color := sprites[x].Color
			if sprites[x].Health < 20 {
				color = sprites[x].HurtColor
			}
			if freezeTime > 0 {
				color = impactColor
			}
			drawSprite(sprites[x].Charecters, sprites[x].OriginPoint, color, frame)

		}

		if freezeTime == impactFrames { // which means impact happened in this frame
			frozenFrame = frame
		}
		drawFrame(frame)
		time.Sleep(frameRate)
	}
}
