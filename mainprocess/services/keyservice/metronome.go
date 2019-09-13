package keyservice

import "github.com/josephbudd/cwt/mainprocess/sound"

var metronomeRunning bool
var metronomeQuitChannel = make(chan struct{})

/* TODO
make these funcs thread safe.
*/

// StartMetronome starts the metronome.
// If the metronome is already running
//   then it just sends a nil error through the errCh
//   to stop the error handler
//   and returns.
func StartMetronome(wpm uint64, errCh chan error) {
	if metronomeRunning {
		errCh <- nil
		return
	}
	// The metronome is not running so start it.
	metronomeRunning = true
	go sound.Metronome(wpm, metronomeQuitChannel, errCh)
}

// StopMetronome stops the metronome.
// After the metronome stops
//   it will send a nil error through the error channel
//   to stop the error handler.
func StopMetronome() {
	if metronomeRunning {
		metronomeRunning = false
		metronomeQuitChannel <- struct{}{}
	}
}
