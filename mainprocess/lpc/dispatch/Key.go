package dispatch

import (
	"context"
	"fmt"
	"strings"

	"github.com/josephbudd/cwt/domain/lpc/message"
	"github.com/josephbudd/cwt/domain/store"
	"github.com/josephbudd/cwt/mainprocess/lpc"
	"github.com/josephbudd/cwt/mainprocess/services/copyservice"
)

/*
	YOU MAY EDIT THIS FILE.

	Rekickwasm will preserve this file for you.
	Kicklpc will not edit this file.

*/

// handleKey is the *message.KeyRendererToMainProcess handler.
// It's response back to the renderer is the *message.KeyMainProcessToRenderer.
// Param ctx is the context. if <-ctx.Done() then the main process is shutting down.
// Param rxmessage *message.KeyRendererToMainProcess is the message received from the renderer.
// Param sending is the channel to use to send a *message.KeyMainProcessToRenderer message back to the renderer.
// Param stores is a struct the contains each of your stores.
// Param errChan is the channel to send the handler's error through since the handler does not return it's error.
func handleKey(ctx context.Context, rxmessage *message.KeyRendererToMainProcess, sending lpc.Sending, stores *store.Stores, errChan chan error) {
	txmessage := &message.KeyMainProcessToRenderer{
		Run:   rxmessage.Run,
		State: rxmessage.State,
	}
	// 1. Turn off the keying if requested.
	if !rxmessage.Run {
		copyservice.StopKeying()
		sending <- txmessage
		return
	}
	// 2. Build the morse code text.
	ditdahs := make([]string, 0, len(rxmessage.Solution))
	for _, line := range rxmessage.Solution {
		ditdahWord := make([]string, 0, len(line))
		for _, r := range line {
			ditdahWord = append(ditdahWord, r.DitDah)
		}
		ditdahs = append(ditdahs, strings.Join(ditdahWord, " "))
	}
	// 3. Key the morse code.
	if err := copyservice.Key(ditdahs, rxmessage.WPM, rxmessage.Pause); err != nil {
		// Send the err to package main.
		errChan <- err
		// Send the error to the renderer.
		// A bolt database error is fatal.
		txmessage.Fatal = true
		txmessage.ErrorMessage = fmt.Sprintf("mainProcessKey:  ditdah.Key(rxmessage.Ditdah, rxmessage.WPM, rxmessage.Delay): error is %s", err.Error())
		sending <- txmessage
		return
	}
	// no error so call back to the renderer.
	sending <- txmessage
	return
}
