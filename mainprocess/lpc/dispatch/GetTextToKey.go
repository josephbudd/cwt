package dispatch

import (
	"context"
	"fmt"

	"github.com/josephbudd/cwt/domain/lpc/message"
	"github.com/josephbudd/cwt/domain/store"
	"github.com/josephbudd/cwt/domain/store/record"
	"github.com/josephbudd/cwt/mainprocess/lpc"
	"github.com/josephbudd/cwt/mainprocess/services/keyservice"
)

/*
	YOU MAY EDIT THIS FILE.

	Rekickwasm will preserve this file for you.
	Kicklpc will not edit this file.

*/

// handleGetTextToKey is the *message.GetTextToKeyRendererToMainProcess handler.
// It's response back to the renderer is the *message.GetTextToKeyMainProcessToRenderer.
// Param ctx is the context. if <-ctx.Done() then the main process is shutting down.
// Param rxmessage *message.GetTextToKeyRendererToMainProcess is the message received from the renderer.
// Param sending is the channel to use to send a *message.GetTextToKeyMainProcessToRenderer message back to the renderer.
// Param stores is a struct the contains each of your stores.
// Param errChan is the channel to send the handler's error through since the handler does not return it's error.
func handleGetTextToKey(ctx context.Context, rxmessage *message.GetTextToKeyRendererToMainProcess, sending lpc.Sending, stores *store.Stores, errChan chan error) {
	txmessage := &message.GetTextToKeyMainProcessToRenderer{
		State:    rxmessage.State,
		Practice: rxmessage.Practice,
	}
	// 1 Get the wpm to key.
	var r *record.WPM
	var err error
	if r, err = stores.WPM.GetKeyWPM(); err != nil {
		// Send the err to package main.
		errChan <- err
		// Send the error to the renderer.
		// A bolt database error is fatal.
		txmessage.Fatal = true
		txmessage.ErrorMessage = fmt.Sprintf("mainProcessGetTextWPMToKey: keyCodeStorer.GetKeyWPM(): error is %s", err.Error())
		sending <- txmessage
		return
	}
	txmessage.WPM = r.WPM
	// 3. Get the text for the user to key from the repo and the help.
	if rxmessage.Practice {
		// practicing
		if txmessage.Solution, txmessage.Help, err = keyservice.GetPracticeKeyCodes(stores.KeyCode, r.WPM); err != nil {
			// Send the err to package main.
			errChan <- err
			// Send the error to the renderer.
			// A bolt database error is fatal.
			txmessage.Fatal = true
			txmessage.ErrorMessage = fmt.Sprintf("mainProcessGetTextWPMToKey: keyservice.GetPracticeKeyCodes(stores.KeyCode, r.WPM): error is %s", err.Error())
			sending <- txmessage
			return
		}
	} else {
		// testing not practicing
		if txmessage.Solution, err = keyservice.GetTestKeyCodes(stores.KeyCode); err != nil {
			// Send the err to package main.
			errChan <- err
			// Send the error to the renderer.
			// A bolt database error is fatal.
			txmessage.Fatal = true
			txmessage.ErrorMessage = fmt.Sprintf("mainProcessGetTextWPMToKey: keyservice.GetTestKeyCodes(keyCodeStorer): error is %s\n", err.Error())
			sending <- txmessage
			return
		}
	}
	sending <- txmessage
}
