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
)

var leds [3]machine.Pin
var buttons [3]machine.Pin
var tones [3]float64
var bzr buzzer.Device
var soundONOFF machine.Pin

func main() {
	var i uint8

	leds[RED] = machine.A1
	leds[GREEN] = machine.A2
	leds[BLUE] = machine.A3

	buttons[RED] = machine.A5
	buttons[GREEN] = machine.A6
	buttons[BLUE] = machine.A7

	for i = 0; i < 3; i++ {
		leds[i].Configure(machine.PinConfig{Mode: machine.PinOutput})
		buttons[i].Configure(machine.PinConfig{Mode: machine.PinInput})

		leds[i].Low()
	}

	bzrPin := machine.A0
	bzrPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	bzr = buzzer.New(bzrPin)

	tones[RED] = buzzer.G4
	tones[GREEN] = buzzer.C4
	tones[BLUE] = buzzer.E4

	soundONOFF = machine.D7
	soundONOFF.Configure(machine.PinConfig{Mode: machine.PinInput})

	var state uint8
	happySound()
	var round uint8

	for {
		switch state {
		case IDLE:
			state = idle()
			break
		case START_GAME:
			for i = 0; i < 3; i++ {
				leds[i].Low()
			}
			round = 0
			state = GENERATE_SEQUENCE
			break
		case GENERATE_SEQUENCE:
			if round > 0 {

			}
			for {
				println(rand.Intn(3))
				time.Sleep(1 * time.Second)
			}
			break
		default:
			break
		}
	}

	/*for {
		if !btnR.Get() {
			ledR.High()
		} else {
			ledR.Low()
		}
		if !btnG.Get() {
			ledG.High()
		} else {
			ledG.Low()
		}
		if !btnB.Get() {
			ledB.High()
		} else {
			ledB.Low()
		}
		time.Sleep(100 * time.Millisecond)
	}*/
}

func idle() uint8 {
	var k uint8
	var i uint8

	for {
		for i = 0; i < 3; i++ {
			leds[i].Low()
			if !buttons[i].Get() {
				return START_GAME
			}
		}

		leds[k].High()
		k = (k + 1) % 3

		time.Sleep(500 * time.Millisecond)
	}
	return IDLE
}

func happySound() {
	bzr.Tone(buzzer.G3, 0.5)
	time.Sleep(100 * time.Millisecond)
	bzr.Tone(buzzer.A3, 0.5)
	time.Sleep(100 * time.Millisecond)
	bzr.Tone(buzzer.B4, 0.5)
	time.Sleep(100 * time.Millisecond)
}
