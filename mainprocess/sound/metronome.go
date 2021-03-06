package sound

import (
	"context"
	"math"
	"time"

	alsa "github.com/josephbudd/cwt/mainprocess/goalsa"
	"github.com/pkg/errors"
)

// Metronome clicks an element beat.
func Metronome(ctx context.Context, wpm uint64, errCh chan error) {

	var err error
	defer func() {
		if err != nil {
			err = errors.WithMessage(err, "Metronome(wpm uint64, quitCh chan struct{}, errCh chan error)")
		}
		// Always send err through the error channel even if its nil
		//   so that the go routine handling the error will stop.
		errCh <- err
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
	err = metronome(ctx, device, wpm)
}

func metronome(ctx context.Context, device *alsa.PlaybackDevice, wpm uint64) (err error) {
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
		case <-ctx.Done():
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
