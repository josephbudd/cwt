package sound

import (
	"fmt"
	"math"
	"strings"
	"time"

	alsa "github.com/josephbudd/cwt/mainprocess/goalsa"
	"github.com/pkg/errors"
)

// PlayCWWords plays an array of indivisual cw words.
// A word is a combination of characters separated by single spaces.
// A character is a combination of "-" and ".".
// PlayCWWords joins the words with "\t"
// PlayCWWords plays them at wpm words per minute.
func PlayCWWords(words []string, wpm uint64, quitCh chan struct{}) (err error) {
	err = PlayCW(strings.Join(words, "\t"), wpm, quitCh)
	return
}

// PlayCW plays "-. ._\t-. ._" at wpm words per minute.
// "." == dit.
// "-" == dah.
// " " == character separator. A character is a string of dits and dahs.
// "\t" == word separator. A word is a string of characters joined by character separators.
func PlayCW(ditdah string, wpm uint64, quitCh chan struct{}) (err error) {
	for i, r := range ditdah {
		if r != '.' && r != '-' && r != ' ' && r != '\t' {
			err = fmt.Errorf("%q is not a valid dit dah character at position %d. It should be a \".\", a \"-\", a space or a tab", r, i)
			return
		}
	}
	err = playCW(ditdah, wpm, quitCh)
	return
}

func playCW(ditdah string, wpm uint64, quitCh chan struct{}) (err error) {
	defer func() {
		if err != nil {
			err = errors.WithMessage(err, "playCW(ditdah string, wpm uint64)")
		}
	}()

	var device *alsa.PlaybackDevice
	device, err = alsa.NewPlaybackDevice(
		"default",
		2,
		alsa.FormatS16LE,
		44100,
		alsa.BufferParams{},
	)
	if err != nil {
		err = errors.WithMessage(err, "alsa.NewPlaybackDevice(...)")
		return
	}
	defer device.Close()
	err = cw(device, wpm, ditdah, quitCh)
	return
}

// cw converts ".- -."  to sound
func cw(device *alsa.PlaybackDevice, wpm uint64, ditdah string, quitCh chan struct{}) (err error) {
	defer func() {
		if err != nil {
			err = errors.WithMessage(err, "cw(device *alsa.PlaybackDevice, wpm uint64, ditdah string)")
		}
	}()
	// word == "paris" == 50 elements.
	nElementsPerMinute := 50.0 * float64(wpm)
	nElementsPerSecond := math.Floor(nElementsPerMinute / 60.0)
	nElementsPerHalfMinute := math.Floor(nElementsPerMinute / 2.0)
	secondsPerElement := 1.0 / nElementsPerSecond

	// build the cw data to be played
	ditSound := buildCWSound(1, device, wpm)
	dahSound := buildCWSound(3, device, wpm)
	oneSecondSound := buildCWSound(uint64(nElementsPerSecond), device, wpm)
	thirySecondSound := buildCWSound(uint64(nElementsPerHalfMinute), device, wpm)
	// loop through the runes in the ditdah string.
	var soundCount int64
	var silenceCount int64
	for _, r := range ditdah {
		switch r {
		case ' ':
			// space is used as a character separator.
			// char separator : 3 elements of silence.
			soundCount = 0
			silenceCount = 2 // 3 - the 1 silence that followed the last dit or dah.
		case '\t':
			// tab is used as a word separator.
			// word separator : 7 elements of silence.
			soundCount = 0
			silenceCount = 6 // 7 - the 1 silence that followed the last dit or dah.
		case '.':
			// period is used as a dit.
			// dit : 1 element of sound followed by 1 element of silence
			// play the dit.
			if _, err = device.Write(ditSound); err != nil {
				return
			}
			soundCount = 1
			silenceCount = 1
		case '-':
			// dash is used as a dah.
			// dah : 3 elements of sound followed by 1 element of silence.
			// play the dah.
			if _, err = device.Write(dahSound); err != nil {
				return
			}
			soundCount = 3
			silenceCount = 1
		case 's':
			// 1 second dit for testing timing.
			if _, err = device.Write(oneSecondSound); err != nil {
				return
			}
			soundCount = int64(nElementsPerSecond)
			silenceCount = 0
		case 'm':
			// 30 second dah for testing timing.
			if _, err = device.Write(thirySecondSound); err != nil {
				return
			}
			soundCount = int64(nElementsPerHalfMinute)
			silenceCount = 0
		}
		pauseFor := secondsPerElement * float64(time.Second) * float64(silenceCount+soundCount)
		if pauseFor > 0 {
			timeout := time.After(time.Duration(int64(pauseFor)))
			select {
			case <-timeout:
			}
		}
		select {
		case <-quitCh:
			return
		default:
			break
		}
	}
	return
}

func buildCWSound(nElements uint64, device *alsa.PlaybackDevice, wpm uint64) (data []int16) {
	nElementsPerMinute := 50 * wpm // 50 == # elements in "paris"
	nElementsPerSecond := nElementsPerMinute / 60
	samplesPerElement := device.Rate / int(nElementsPerSecond)
	nSamples := samplesPerElement * int(nElements)
	dataSize := nSamples * device.Channels
	data = make([]int16, dataSize, dataSize)
	// each sample is a frame of device.Channels int16.
	for i := 0; i < nSamples; i++ {
		for j := 1; j <= device.Channels; j++ {
			data[device.Channels*i] = int16((i%(j*128))*100 - 1000)
		}
	}
	return
}
