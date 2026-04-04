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
var frameRate = time.Millisecond * 16
var impactColor = white

var freezeTime int
var frozenFrame frames
var aliveTimers []*TimedEvent

var AllTimers []*TimedEvent

const numVoices = 4

type audioVoice struct {
	data    []byte
	pos     int
	playing bool
}

var voices [numVoices]audioVoice
var currentVoice int
var sfxFiles = map[string][]byte{}
var ctx *malgo.AllocatedContext

func loadSound(name string, file string) {
	f, _ := os.Open(file)
	defer f.Close()
	f.Seek(44, 0)
	data, _ := io.ReadAll(f)
	sfxFiles[name] = data
}

func playSound(name string) {
	voices[currentVoice].data = sfxFiles[name]
	voices[currentVoice].pos = 0
	voices[currentVoice].playing = true
	currentVoice = (currentVoice + 1) % numVoices
}

func initAudio() {
	loadSound("hit", "hit.wav")
	loadSound("timeStop", "timestop.wav")

	ctx, _ = malgo.InitContext(nil, malgo.ContextConfig{}, nil)

	deviceConfig := malgo.DefaultDeviceConfig(malgo.Playback)
	deviceConfig.Playback.Format = malgo.FormatS16
	deviceConfig.Playback.Channels = 2
	deviceConfig.SampleRate = 44100

	for i := 0; i < numVoices; i++ {
		i := i
		callbacks := malgo.DeviceCallbacks{
			Data: func(pOutput, pInput []byte, frameCount uint32) {
				if !voices[i].playing {
					return
				}
				bytesNeeded := int(frameCount) * 4
				remaining := len(voices[i].data) - voices[i].pos
				if remaining <= 0 {
					voices[i].playing = false
					return
				}
				if bytesNeeded > remaining {
					bytesNeeded = remaining
				}
				copy(pOutput, voices[i].data[voices[i].pos:voices[i].pos+bytesNeeded])
				voices[i].pos += bytesNeeded
			},
		}
		device, _ := malgo.InitDevice(ctx.Context, deviceConfig, callbacks)
		device.Start()
	}
}

var sprite1 = sprite{
	OriginPoint:    point{10, 10},
	NormalVelocity: velocity{1, 1},
	Velocity:       velocity{1, 1},
	Hitbox:         hitbox{point{-3, -3}, point{3, 3}},
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
	Health: 200,
	CollisionFunc: func(meSprite, hitSprite *sprite) {
		hitSprite.Health -= meSprite.Entity["damage"].(int)
		return
	},
	Color:     green,
	HurtColor: yellow,
	Name:      "slime",
	Entity: map[string]any{
		"damage": 15,
	},
}

var sprite2 = sprite{
	OriginPoint: point{30, 10},
	Charecters: []string{
		"   '.oOOOOOOOOo.'   ",
		"  oOOOOOOOOOOOOOOo  ",
		" oOOQQOOOOOOOOOOOOo ",
		"oOQQQQQOOOOOOOOOOOOo ",
		"oQQQQQQOOOOOOOOOOXXo",
		"oQQQQQQOOOOOOOOOOXXo",
		" ooOOOOOOOOOOOOOOXo ",
		"  .oOOOOOOOOOOOoo.  ",
		"    '.oOOOOOOOo.'   ",
	},
	HurtCharecters: []string{
		"   '.oOOOOOOOOo.'   ",
		"  o\\\\O//OO#OOOOOOo  ",
		" oO\\\\//OOO#OOOOOOOo ",
		"oOQQ//\\OO'''#O'OOOOo ",
		"oQQ//\\\\OO##'''OOOXXo",
		"oQQQQQQ\\\\'''OOOOOXXo",
		" ooOOOOO//OOOOOOOXo ",
		"  .oOOO//OOOOOOoo.  ",
		"    '.oOOOOOOOo.'   ",
	},
	Color:     yellow,
	HurtColor: red,
	Hitbox:    hitbox{point{-9, -4}, point{11, 5}},
	Velocity:  velocity{1, 1},
	Health:    100,
	CollisionFunc: func(meSprite *sprite, hitSprite *sprite) {
		return
	},
	Name: "sphere",
}

var sprite3 = sprite{
	OriginPoint: point{30, 10},
	Charecters: []string{
		"  _______  ",
		" /''12'''\\ ",
		"|''''|''''|",
		"|9'''|'''3|",
		"|'''''\\'''|",
		"|''''''\\''|",
		" \\___6___/ ",
	},
	HurtCharecters: []string{
		"  _______  ",
		" /''12''/  ",
		"|''''|'/   ",
		"|9###|''\\  ",
		"|####'\\'\\ ",
		"|''''''\\''|",
		" \\___6___/ ",
	},
	Hitbox:         hitbox{point{-5, -4}, point{5, 4}},
	Color:          blue,
	HurtColor:      red,
	Velocity:       velocity{1, 1},
	NormalVelocity: velocity{1, 1},
	Entity: map[string]any{
		"normalDamage":       5,
		"multipliedDamage":   25,
		"multipliedVelocity": velocity{3, 3},
		"normalVelocity":     velocity{1, 1},
		"timeFrozen":         false,
		"timeFreezeDuration": 60 * 2,
	},
	CollisionFunc: func(meSprite *sprite, hitSprite *sprite) {
		if meSprite.Entity["timeFrozen"].(bool) == true {
			hitSprite.Health -= meSprite.Entity["multipliedDamage"].(int)
			meSprite.Health += hitSprite.Entity["damage"].(int)
			return
		}
		hitSprite.Health -= meSprite.Entity["normalDamage"].(int)

	},
	Health: 200,
	Name:   "time stopper",
}

func timeContinue() {
	for x := range sprites {
		if sprites[x].Name == sprite3.Name {
			sprites[x].Entity["timeFrozen"] = false
			currVelocity := sprites[x].Velocity
			multipliedVelocity := sprites[x].Entity["normalVelocity"].(velocity)
			sprites[x].Velocity = velocity{
				currVelocity.X / multipliedVelocity.X,
				currVelocity.Y / multipliedVelocity.Y,
			}

		}
		sprites[x].Velocity = sprites[x].NormalVelocity
	}
}

func timeStop() {
	exists := false
	var TimeStopperIndex int
	for x := range sprites {
		if sprites[x].Name == sprite3.Name {
			exists = true
			TimeStopperIndex = x
		}
	}
	if !exists {
		return
	}
	playSound("timeStop")
	aliveTimers = append(aliveTimers,
		&TimedEvent{
			Repeat: false,
			Timer: Timer{
				OriginalTime: sprite3.Entity["timeFreezeDuration"].(int),
				CurrentTime:  sprite3.Entity["timeFreezeDuration"].(int),
			},
			Event: timeContinue,
		},
	)
	sprites[TimeStopperIndex].Entity["timeFrozen"] = true
	currVelocity := sprites[TimeStopperIndex].Velocity
	multipliedVelocity := sprites[TimeStopperIndex].Entity["multipliedVelocity"].(velocity)
	sprites[TimeStopperIndex].Velocity = velocity{
		currVelocity.X * multipliedVelocity.X,
		currVelocity.Y * multipliedVelocity.Y,
	}
	for x := range sprites {
		if sprites[x].Name == sprite3.Name {
			continue
		}
		sprites[x].Velocity = velocity{0, 0}

		continue
	}
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
	Timer  Timer
	Repeat bool
	Event  func()
}

type sprite struct {
	Velocity       velocity
	NormalVelocity velocity
	Charecters     []string // since a string is already a slice of charecters this way it's 2d
	HurtCharecters []string
	Hitbox         hitbox
	OriginPoint    point
	AbilityTimer   Timer
	Health         int
	CollisionFunc  func(meSprite *sprite, hitSprite *sprite)
	Name           string
	Color          string
	HurtColor      string
	Entity         map[string]any
}

var sprites []sprite

func sphereHealAbility() {
	exists := false
	for x := range sprites {
		if sprites[x].Name == "sphere" {
			exists = true
		}
	}

	if !exists {
		return
	}

	for x := range sprites {
		if sprites[x].Name == "sphere" {
			sprites[x].Health += 5
			continue
		}

		sprites[x].Health -= 5
	}
}

func AbsInt(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func moveOriginPoints() {
	for x := range sprites {
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

}

func drawFrame(frame frames) {
	frameString := "\033[H\033[2J"
	for _, row := range frame {
		for _, charecter := range row {
			frameString += charecter
		}
		frameString += "\r\n"
	}
	fmt.Print(frameString)
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

func goRight(index int) {
	sprites[index].Velocity.X = AbsInt(sprites[index].Velocity.X)
}

func goLeft(index int) {
	sprites[index].Velocity.X = AbsInt(sprites[index].Velocity.X) * -1
}

func goUp(index int) {
	sprites[index].Velocity.Y = AbsInt(sprites[index].Velocity.Y) * -1
}

func goDown(index int) {
	sprites[index].Velocity.Y = AbsInt(sprites[index].Velocity.Y)
}

func (s sprite) moveBy(x int, y int) {
	s.OriginPoint = point{
		s.OriginPoint.X + x,
		s.OriginPoint.Y + y,
	}
}

func (s sprite) relativeHitbox() hitbox {
	return hitbox{
		TopLeft:     point{s.OriginPoint.X + s.Hitbox.TopLeft.X, s.OriginPoint.Y + s.Hitbox.TopLeft.Y},
		BottomRight: point{s.OriginPoint.X + s.Hitbox.BottomRight.X, s.OriginPoint.Y + s.Hitbox.BottomRight.Y},
	}
}

func main() {
	sprites = []sprite{
		sprite1,
		sprite3,
	}

	AllTimers = []*TimedEvent{
		&TimedEvent{
			Timer:  Timer{60 * 6, 60 * 6},
			Repeat: true,
			Event:  timeStop,
		},
	}

	initAudio()

	for {

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

		sprites = alive // removing dead sprites

		// ticks all timers
		aliveTimers = []*TimedEvent{}
		for _, timedEvent := range AllTimers {
			timedEvent.Timer.CurrentTime--
			if timedEvent.Timer.CurrentTime == 0 {
				timedEvent.Event()
				if timedEvent.Repeat == true {
					timedEvent.Timer.CurrentTime = timedEvent.Timer.OriginalTime
					aliveTimers = append(aliveTimers, timedEvent)
					continue
				}
				continue
			}
			aliveTimers = append(aliveTimers, timedEvent)
		}

		AllTimers = aliveTimers
		moveOriginPoints() // moves the ascii's origin points by their velocity
		for x := range sprites {

			fmt.Printf("%s's health : %d\n", sprites[x].Name, sprites[x].Health)

			sprites[x].AbilityTimer.CurrentTime--
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
				alreadyCollidedSprites[j] = x
				playSound("hit") // hit sound
				freezeTime = impactFrames
				sprites[x].CollisionFunc(&sprites[x], &sprites[j])
				sprites[j].CollisionFunc(&sprites[j], &sprites[x])

				var rightMostAsciiIndex int
				var leftMostAsciiIndex int
				var upperAsciiIndex int
				var bottomAsciiIndex int

				if sprites[x].OriginPoint.X > sprites[j].OriginPoint.X {
					rightMostAsciiIndex = x
					leftMostAsciiIndex = j
				} else {
					rightMostAsciiIndex = j
					leftMostAsciiIndex = x
				}
				if sprites[x].OriginPoint.Y > sprites[j].OriginPoint.Y {
					upperAsciiIndex = j
					bottomAsciiIndex = x
				} else {
					upperAsciiIndex = x
					bottomAsciiIndex = j
				}

				if collisionBoxHeight == collisionBoxWidth {

					goLeft(leftMostAsciiIndex)
					goRight(rightMostAsciiIndex)

					sprites[x].moveBy(sprites[x].Velocity.X*2, sprites[x].Velocity.Y*2)
					sprites[j].moveBy(sprites[j].Velocity.X*2, sprites[j].Velocity.Y*2)
					alreadyCollidedSprites[j] = x
					continue
				}
				if collisionBoxHeight > collisionBoxWidth {
					sprites[rightMostAsciiIndex].Velocity.X = AbsInt(sprites[rightMostAsciiIndex].Velocity.X)
					sprites[leftMostAsciiIndex].Velocity.X = AbsInt(sprites[leftMostAsciiIndex].Velocity.X) * -1

					sprites[x].moveBy(sprites[j].Velocity.X*2, 0)
					sprites[j].moveBy(sprites[x].Velocity.X*2, 0)
					continue
				}
				sprites[upperAsciiIndex].Velocity.Y = AbsInt(sprites[upperAsciiIndex].Velocity.Y) * -1
				sprites[bottomAsciiIndex].Velocity.Y = AbsInt(sprites[bottomAsciiIndex].Velocity.Y)

				sprites[x].moveBy(0, sprites[x].Velocity.Y*2)
				sprites[j].moveBy(0, sprites[j].Velocity.Y*2)
				fmt.Println(collisionBoxHeight, collisionBoxWidth)

			}
		}

		for x := range sprites {
			// find hitbox and draw sprite
			color := sprites[x].Color
			charecters := sprites[x].Charecters
			if sprites[x].Health < 20 {
				color = sprites[x].HurtColor
				charecters = sprites[x].HurtCharecters
			}
			if freezeTime > 0 {
				color = impactColor
			}
			drawSprite(charecters, sprites[x].OriginPoint, color, frame)

		}

		if freezeTime == impactFrames { // which means impact happened in this frame
			frozenFrame = frame
		}
		drawFrame(frame)
		time.Sleep(frameRate)
	}
}
