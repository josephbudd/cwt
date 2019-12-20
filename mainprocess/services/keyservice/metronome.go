package keyservice

import (
	"context"

	"github.com/josephbudd/cwt/mainprocess/sound"
)

var (
	metronomeRunning bool
	metronomeCtx     context.Context
	metronomeCancel  context.CancelFunc
)

/* TODO
make these funcs thread safe.
*/

// StartMetronome starts the metronome.
// If the metronome is already running
//   then it just sends a nil error through the errCh
//   to stop the error handler
//   and returns.
func StartMetronome(ctx context.Context, wpm uint64, errCh chan error) {
	if metronomeRunning {
		errCh <- nil
		return
	}
	metronomeCtx, metronomeCancel = context.WithCancel(ctx)
	// The metronome is not running so start it.
	metronomeRunning = true
	go sound.Metronome(metronomeCtx, wpm, errCh)
}

// StopMetronome stops the metronome.
func StopMetronome() {
	if metronomeRunning {
		metronomeRunning = false
		metronomeCancel()
	}
}
