package main

import (
	"machine"
	"math/rand"
	"time"

	"tinygo.org/x/drivers/buzzer"
)

const (
	RED = iota
	GREEN
	BLUE
)

const (
	IDLE = iota
	START_GAME
	GENERATE_SEQUENCE
	PLAYER_INPUT
	GAME_OVER
	PLAYER_WINS
)

// Game holds the pins for buttons and leds, and other game information
type Game struct {
	leds       [3]machine.Pin
	buttons    [3]machine.Pin
	tones      [3]float64
	bzr        buzzer.Device
	soundONOFF machine.Pin
	sequence   [20]uint8
	round      uint8
	state      uint8
}

func main() {
	var i uint8
	var k uint8
	game := Game{}

	// Set up the pins for the leds
	game.leds[RED] = machine.A1
	game.leds[GREEN] = machine.A2
	game.leds[BLUE] = machine.A3

	// Set up the pins for the buttons
	game.buttons[RED] = machine.A5
	game.buttons[GREEN] = machine.A6
	game.buttons[BLUE] = machine.A7

	// Configure the LEDs pins as output, the buttons as input
	// set the leds off
	for i = 0; i < 3; i++ {
		game.leds[i].Configure(machine.PinConfig{Mode: machine.PinOutput})
		game.buttons[i].Configure(machine.PinConfig{Mode: machine.PinInput})

		game.leds[i].Low()
	}

	// Configure the buzzer pin with the buzzer driver
	bzrPin := machine.A0
	bzrPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	game.bzr = buzzer.New(bzrPin)

	// Assign each color/button a different tone
	game.tones[RED] = buzzer.G4
	game.tones[GREEN] = buzzer.C4
	game.tones[BLUE] = buzzer.E4

	// Use the slide switch on the CPE to disable sound
	game.soundONOFF = machine.D7
	game.soundONOFF.Configure(machine.PinConfig{Mode: machine.PinInput})

	// Play a happy sound
	game.happySound()

	// Start the game in IDLE mode
	game.state = IDLE

	for {
		switch game.state {
		case IDLE:
			for game.state == IDLE {
				// Check if any button is pressed
				for i = 0; i < 3; i++ {
					game.leds[i].Low()
					if !game.buttons[i].Get() {
						game.state = START_GAME
						break
					}
				}

				game.leds[k].High()
				k = (k + 1) % 3

				time.Sleep(500 * time.Millisecond)
			}
			break
		case START_GAME:
			for i = 0; i < 3; i++ {
				game.leds[i].Low()
			}
			game.round = 0
			game.state = GENERATE_SEQUENCE
			break
		case GENERATE_SEQUENCE:
			// play existing sequence of color/sounds
			if game.round > 0 {
				for i = 0; i < game.round; i++ {
					game.playTune(game.sequence[i])
					time.Sleep(100 * time.Millisecond)
				}
			}
			// generate new step in the sequence
			game.sequence[game.round] = uint8(rand.Intn(3))
			game.playTune(game.sequence[game.round])
			time.Sleep(100 * time.Millisecond)
			game.state = PLAYER_INPUT
			break
		case PLAYER_INPUT:
			for game.state == PLAYER_INPUT {
				for i = 0; i < 3; i++ {
					game.leds[i].Low()
					if !game.buttons[i].Get() {
						if i != game.sequence[game.round] {
							game.state = GAME_OVER
						} else {
							game.state = GENERATE_SEQUENCE
						}
						break
					}
				}
				time.Sleep(100 * time.Millisecond)
			}
			break
		case GAME_OVER:
			game.sadSound()
			time.Sleep(3 * time.Second)
			break
		case PLAYER_WINS:
			game.happySound()
			time.Sleep(3 * time.Second)
			break
		default:
			break
		}
	}

}

func (game *Game) happySound() {
	game.bzr.Tone(buzzer.G3, 0.5)
	time.Sleep(100 * time.Millisecond)
	game.bzr.Tone(buzzer.A3, 0.5)
	time.Sleep(100 * time.Millisecond)
	game.bzr.Tone(buzzer.B4, 0.5)
	time.Sleep(100 * time.Millisecond)
}

func (game *Game) sadSound() {
	game.bzr.Tone(buzzer.B4, 0.5)
	time.Sleep(100 * time.Millisecond)
	game.bzr.Tone(buzzer.A3, 0.5)
	time.Sleep(100 * time.Millisecond)
	game.bzr.Tone(buzzer.G3, 0.5)
	time.Sleep(100 * time.Millisecond)
}

func (game *Game) playTune(color uint8) {
	game.leds[color].High()
	game.bzr.Tone(game.tones[color], 0.5)
	time.Sleep(100 * time.Millisecond)
	game.leds[color].Low()
}
