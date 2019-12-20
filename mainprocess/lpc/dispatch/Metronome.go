package dispatch

import (
	"context"

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
// Param ctx is the context. if <-ctx.Done() then the main process is shutting down.
// Param rxmessage *message.MetronomeRendererToMainProcess is the message received from the renderer.
// Param sending is the channel to use to send a *message.MetronomeMainProcessToRenderer message back to the renderer.
// Param stores is a struct the contains each of your stores.
// Param errChan is the channel to send the handler's error through since the handler does not return it's error.
func handleMetronome(ctx context.Context, rxmessage *message.MetronomeRendererToMainProcess, sending lpc.Sending, stores *store.Stores, errChan chan error) {
	txmessage := &message.MetronomeMainProcessToRenderer{
		Run:   rxmessage.Run,
		State: rxmessage.State,
	}
	// 1. Turn the metronome on or off.
	if rxmessage.Run {
		metronomeErrCh := make(chan error)
		go handleMetronomeError(ctx, txmessage, sending, metronomeErrCh, errChan)
		keyservice.StartMetronome(ctx, rxmessage.WPM, metronomeErrCh)
	} else {
		keyservice.StopMetronome()
	}
	// 2. Let the renderer process know the job is running or stopped.
	sending <- txmessage
}

func handleMetronomeError(ctx context.Context, txmessage *message.MetronomeMainProcessToRenderer, sending lpc.Sending, metronomeErrCh, errChan chan error) {
	for {
		select {
		case <-ctx.Done():
			return
		case err := <-metronomeErrCh:
			if err != nil {
				// Send the err to package main.
				errChan <- err
				// Send the error to the renderer.
				// A sound err is not fatal.
				txmessage.Error = true
				txmessage.ErrorMessage = err.Error()
				sending <- txmessage
			}
			return
		}
	}
}
