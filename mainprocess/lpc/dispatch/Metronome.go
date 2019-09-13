package dispatch

import (
	"github.com/josephbudd/cwt/domain/lpc/message"
	"github.com/josephbudd/cwt/domain/store"
	"github.com/josephbudd/cwt/mainprocess/lpc"
	"github.com/josephbudd/cwt/mainprocess/services/keyservice"
)

/*
	YOU MAY EDIT THIS FILE.

	Rekickwasm will preserve this file for you.
	Kicklpc will not edit this file.

*/

// handleMetronome is the *message.MetronomeRendererToMainProcess handler.
// It's response back to the renderer is the *message.MetronomeMainProcessToRenderer.
// Param rxmessage *message.MetronomeRendererToMainProcess is the message received from the renderer.
// Param sending is the channel to use to send a *message.MetronomeMainProcessToRenderer message back to the renderer.
// Param eojing lpc.EOJer ( End Of Job ) is an interface for your go routine to receive a stop signal.
//   It signals go routines that they must stop because the main process is ending.
//   So only use it inside a go routine if you have one.
//   In your go routine
//     1. Get a channel to listen to with eojing.NewEOJ().
//     2. Before your go routine returns, release that channel with eojing.Release().
// The func is simple:
// 1. Turn the metronome on or off.
// 2. Let the renderer process know the job is running or stopped.
func handleMetronome(rxmessage *message.MetronomeRendererToMainProcess, sending lpc.Sending, eojing lpc.EOJer, stores *store.Stores) {
	txmessage := &message.MetronomeMainProcessToRenderer{
		Run:   rxmessage.Run,
		State: rxmessage.State,
	}
	// 1. Turn the metronome on or off.
	if rxmessage.Run {
		errCh := make(chan error)
		go handleMetronomeError(txmessage, sending, errCh, eojing)
		keyservice.StartMetronome(rxmessage.WPM, errCh)
	} else {
		keyservice.StopMetronome()
	}
	// 2. Let the renderer process know the job is running or stopped.
	sending <- txmessage
}

func handleMetronomeError(txmessage *message.MetronomeMainProcessToRenderer, sending lpc.Sending, errCh chan error, eojing lpc.EOJer) {
	eojCh := eojing.NewEOJ()
	for {
		select {
		case <-eojCh:
			eojing.Release()
			keyservice.StopMetronome()
			return
		case err := <-errCh:
			eojing.Release()
			if err != nil {
				txmessage.Error = true
				txmessage.ErrorMessage = err.Error()
				sending <- txmessage
			}
			return
		}
	}
}
