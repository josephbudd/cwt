package sound

import (
	"math"
	"time"

	alsa "github.com/josephbudd/cwt/mainprocess/goalsa"
	"github.com/pkg/errors"
)

// Metronome clicks an element beat.
func Metronome(wpm uint64, quitCh chan struct{}, errCh chan error) {
	var err error
	defer func() {
		if err != nil {
			err = errors.WithMessage(err, "Metronome(wpm uint64, quitCh chan struct{}, errCh chan error)")
			errCh <- err
		} else {
			errCh <- nil
		}
	}()

	var device *alsa.PlaybackDevice
	device, err = alsa.NewPlaybackDevice(
		"default",
		2,
		alsa.FormatS16BE,
		44100,
		alsa.BufferParams{},
	)
	if err != nil {
		return
	}
	defer device.Close()
	err = metronome(device, wpm, quitCh)
}

func metronome(device *alsa.PlaybackDevice, wpm uint64, quitCh chan struct{}) (err error) {
	// word == "paris" == 50 elements.
	nElementsPerMinute := 50.0 * float64(wpm)
	nElementsPerSecond := math.Floor(nElementsPerMinute / 60.0)
	secondsPerElement := 1.0 / nElementsPerSecond
	duration := time.Duration(uint64(float64(time.Second) * secondsPerElement))
	// play sound and then wait for 1 element.
	data := buildMetronomeSound(device, wpm)
	timer := time.NewTimer(duration)
	for {
		// play the click sound
		if _, err = device.Write(data); err != nil {
			timer.Stop()
			return
		}
		select {
		case <-timer.C:
			timer = time.NewTimer(duration)
		case <-quitCh:
			return
		}
	}
}

func buildMetronomeSound(device *alsa.PlaybackDevice, wpm uint64) (data []int16) {
	nElementsPerMinute := 50.0 * float64(wpm)
	nElementsPerSecond := math.Floor(nElementsPerMinute / 60.0)
	samplesPerElement := device.Rate / int(nElementsPerSecond)
	// this metronome tick is 1/10 of an element.
	nSamples := samplesPerElement / 10
	dataSize := nSamples * device.Channels
	data = make([]int16, dataSize, dataSize)
	// each sample is a frame of device.Channels int16.
	for i := 0; i < nSamples; i++ {
		for j := 1; j <= device.Channels; j++ {
			index := (device.Channels * i) + j - 1
			data[index] = int16((i%(j*128))*100 - 1000)
		}
	}
	return
}
