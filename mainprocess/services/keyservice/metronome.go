package keyservice

import "github.com/josephbudd/cwt/mainprocess/goalsa"

var metronomeRunning bool
var metronomeQuitChannel = make(chan struct{})

// StartMetronome starts the metronome.
func StartMetronome(wpm uint64, errCh chan error) {
	if !metronomeRunning {
		metronomeRunning = true
		go goalsa.Metronome(wpm, metronomeQuitChannel, errCh)
	}
}

// StopMetronome stops the metronome.
func StopMetronome() {
	if metronomeRunning {
		metronomeRunning = false
		metronomeQuitChannel <- struct{}{}
		return
	}
}
