package dispatch

import (
	"context"
	"fmt"
	"log"

	"github.com/josephbudd/cwt/domain/lpc/message"
	"github.com/josephbudd/cwt/domain/store"
	"github.com/josephbudd/cwt/mainprocess/lpc"
)

/*
	YOU MAY EDIT THIS FILE.

	Rekickwasm will preserve this file for you.
	Kicklpc will not edit this file.

*/

// handleGetCopyWPM is the *message.GetCopyWPMRendererToMainProcess handler.
// It's response back to the renderer is the *message.GetCopyWPMMainProcessToRenderer.
// Param ctx is the context. if <-ctx.Done() then the main process is shutting down.
// Param rxmessage *message.GetCopyWPMRendererToMainProcess is the message received from the renderer.
// Param sending is the channel to use to send a *message.GetCopyWPMMainProcessToRenderer message back to the renderer.
// Param stores is a struct the contains each of your stores.
// Param errChan is the channel to send the handler's error through since the handler does not return it's error.
func handleGetCopyWPM(ctx context.Context, rxmessage *message.GetCopyWPMRendererToMainProcess, sending lpc.Sending, stores *store.Stores, errChan chan error) {
	txmessage := &message.GetCopyWPMMainProcessToRenderer{}
	// 1. Get the keyCode from the repo.
	var err error
	if txmessage.Record, err = stores.WPM.GetCopyWPM(); err != nil {
		// Send the err to package main.
		errChan <- err
		// Send the error to the renderer.
		// A bolt database error is fatal.
		txmessage.Fatal = true
		txmessage.ErrorMessage = fmt.Sprintf("handleGetCopyWPM: stores.WPM.GetCopyWPM(): error is %s", err.Error())
		sending <- txmessage
		return
	}
	if txmessage.Record == nil {
		// Calling back the not found error.
		// This will only happen in development. It means that the data store is not getting initialized properly.
		// Send the err to package main.
		errChan <- err
		// Send the error to the renderer.
		// Not found is fatal.
		txmessage.Fatal = true
		txmessage.ErrorMessage = "handleGetCopyWPM: stores.WPM.GetCopyWPM(): error is Not Found."
		sending <- txmessage
		log.Println(txmessage.ErrorMessage)
		return
	}
	// 2. Send the wpm to the renderer.
	sending <- txmessage
	return
}
