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

// handleCheckKey is the *message.CheckKeyRendererToMainProcess handler.
// It's response back to the renderer is the *message.CheckKeyMainProcessToRenderer.
// Param ctx is the context. if <-ctx.Done() then the main process is shutting down.
// Param rxmessage *message.CheckKeyRendererToMainProcess is the message received from the renderer.
// Param sending is the channel to use to send a *message.CheckKeyMainProcessToRenderer message back to the renderer.
// Param stores is a struct the contains each of your stores.
// Param errChan is the channel to send the handler's error through since the handler does not return it's error.
func handleCheckKey(ctx context.Context, rxmessage *message.CheckKeyRendererToMainProcess, sending lpc.Sending, stores *store.Stores, errChan chan error) {
	txmessage := &message.CheckKeyMainProcessToRenderer{
		State: rxmessage.State,
	}
	// 1. Copy: Convert the miliseconds to key code records.
	var copiedWords [][]*record.KeyCode
	var err error
	if copiedWords, err = keyservice.Copy(rxmessage.MilliSeconds, rxmessage.WPM, stores.KeyCode); err != nil {
		// Send the err to package main.
		errChan <- err
		// Send the error to the renderer.
		// A bolt database error is fatal.
		txmessage.Fatal = true
		txmessage.ErrorMessage = fmt.Sprintf("handleCheckKey: keyservice.Copy(rxmessage.MilliSeconds, rxmessage.WPM, stores.KeyCode): error is %s", err.Error())
		sending <- txmessage
		return
	}
	// 2. Check the copy against the solution.
	if txmessage.CorrectCount, txmessage.IncorrectCount, txmessage.MaxCount, txmessage.TestResults, err = keyservice.Check(copiedWords, rxmessage.Solution, stores.KeyCode, rxmessage.WPM, rxmessage.StoreResults); err != nil {
		// Send the err to package main.
		errChan <- err
		// Send the error to the renderer.
		// A bolt database error is fatal.
		txmessage.Fatal = true
		txmessage.ErrorMessage = fmt.Sprintf("handleCheckKey: keyservice.Check(copiedWords, rxmessage.Solution, stores.KeyCode, rxmessage.WPM, rxmessage.StoreResults): error is %s", err.Error())
		sending <- txmessage
		return
	}
	// 3. Send the results to the renderer.
	sending <- txmessage
}
