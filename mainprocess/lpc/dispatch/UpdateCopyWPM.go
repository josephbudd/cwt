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

// handleUpdateCopyWPM is the *message.UpdateCopyWPMRendererToMainProcess handler.
// It's response back to the renderer is the *message.UpdateCopyWPMMainProcessToRenderer.
// Param ctx is the context. if <-ctx.Done() then the main process is shutting down.
// Param rxmessage *message.UpdateCopyWPMRendererToMainProcess is the message received from the renderer.
// Param sending is the channel to use to send a *message.UpdateCopyWPMMainProcessToRenderer message back to the renderer.
// Param stores is a struct the contains each of your stores.
// Param errChan is the channel to send the handler's error through since the handler does not return it's error.
func handleUpdateCopyWPM(ctx context.Context, rxmessage *message.UpdateCopyWPMRendererToMainProcess, sending lpc.Sending, stores *store.Stores, errChan chan error) {
	txmessage := &message.UpdateCopyWPMMainProcessToRenderer{
		Record: rxmessage.Record,
	}
	// 1. Update the wpm.
	if err := stores.WPM.Update(rxmessage.Record); err != nil {
		// Send the err to package main.
		errChan <- err
		// Send the error to the renderer.
		// A bolt database error is fatal.
		txmessage.Fatal = true
		txmessage.ErrorMessage = fmt.Sprintf("handleUpdateCopyWPM: stores.WPM.Update(rxmessage.Record): error is %s\n", err.Error())
		sending <- txmessage
		return
	}
	// 2. Send the renderer the updated record.
	sending <- txmessage
}
