package dispatch

import (
	"context"
	"fmt"

	"github.com/josephbudd/cwt/domain/lpc/message"
	"github.com/josephbudd/cwt/domain/store"
	"github.com/josephbudd/cwt/mainprocess/lpc"
)

/*
	YOU MAY EDIT THIS FILE.

	Rekickwasm will preserve this file for you.
	Kicklpc will not edit this file.

*/

// handleGetKeyCodes is the *message.GetKeyCodesRendererToMainProcess handler.
// It's response back to the renderer is the *message.GetKeyCodesMainProcessToRenderer.
// Param ctx is the context. if <-ctx.Done() then the main process is shutting down.
// Param rxmessage *message.GetKeyCodesRendererToMainProcess is the message received from the renderer.
// Param sending is the channel to use to send a *message.GetKeyCodesMainProcessToRenderer message back to the renderer.
// Param stores is a struct the contains each of your stores.
// Param errChan is the channel to send the handler's error through since the handler does not return it's error.
func handleGetKeyCodes(ctx context.Context, rxmessage *message.GetKeyCodesRendererToMainProcess, sending lpc.Sending, stores *store.Stores, errChan chan error) {
	txmessage := &message.GetKeyCodesMainProcessToRenderer{
		State: rxmessage.State,
	}
	// 1. Get the keyCodes from the repo.
	var err error
	if txmessage.Records, err = stores.KeyCode.GetAll(); err != nil {
		// Send the err to package main.
		errChan <- err
		// Send the error to the renderer.
		// A bolt database error is fatal.
		txmessage.Fatal = true
		txmessage.ErrorMessage = fmt.Sprintf("handleGetKeyCodes: stores.KeyCode.GetAll(): error is %s", err.Error())
		sending <- txmessage
		return
	}
	// 2. Send the key codes to the renderer.
	sending <- txmessage
}
