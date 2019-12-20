package dispatch

import (
	"context"
	"fmt"
	"strings"

	"github.com/josephbudd/cwt/domain/lpc/message"
	"github.com/josephbudd/cwt/domain/store"
	"github.com/josephbudd/cwt/domain/store/record"
	"github.com/josephbudd/cwt/mainprocess/lpc"
	"github.com/josephbudd/cwt/mainprocess/services/copyservice"
)

/*
	YOU MAY EDIT THIS FILE.

	Rekickwasm will preserve this file for you.
	Kicklpc will not edit this file.

*/

// handleCheckCopy is the *message.CheckCopyRendererToMainProcess handler.
// It's response back to the renderer is the *message.CheckCopyMainProcessToRenderer.
// Param ctx is the context. if <-ctx.Done() then the main process is shutting down.
// Param rxmessage *message.CheckCopyRendererToMainProcess is the message received from the renderer.
// Param sending is the channel to use to send a *message.CheckCopyMainProcessToRenderer message back to the renderer.
// Param stores is a struct the contains each of your stores.
// Param errChan is the channel to send the handler's error through since the handler does not return it's error.
func handleCheckCopy(ctx context.Context, rxmessage *message.CheckCopyRendererToMainProcess, sending lpc.Sending, stores *store.Stores, errChan chan error) {
	txmessage := &message.CheckCopyMainProcessToRenderer{
		State: rxmessage.State,
	}
	// 1. Convert the copy string to key code records.
	// convert copy [][]string to [][]*record.KeyCode
	var rr []*record.KeyCode
	var err error
	if rr, err = stores.KeyCode.GetAll(); err != nil {
		// Send the err to package main.
		errChan <- err
		// Send the error to the renderer.
		// A bolt database error is fatal.
		txmessage.Fatal = true
		txmessage.ErrorMessage = fmt.Sprintf("handleCheckCopy: stores.KeyCode.GetAll(): error is %s", err.Error())
		sending <- txmessage
		return
	}
	copy := make([][]*record.KeyCode, 0, len(rxmessage.Solution))
	for _, rxCopyLine := range rxmessage.Copy {
		copyLine := make([]*record.KeyCode, 0, len(rxCopyLine))
		for _, rxCopyChar := range rxCopyLine {
			uc := strings.ToUpper(string(rxCopyChar))
			found := false
			for _, r := range rr {
				if r.Character == uc {
					copyLine = append(copyLine, r)
					found = true
					break
				}
			}
			if !found {
				copyLine = append(copyLine, nil)
			}
		}
		copy = append(copy, copyLine)
	}
	// 2. Check the copy.
	if txmessage.CorrectCount, txmessage.IncorrectCount, txmessage.MaxCount, txmessage.TestResults, err = copyservice.Check(copy, rxmessage.Solution, stores.KeyCode, rxmessage.WPM, rxmessage.StoreResults); err != nil {
		// Send the err to package main.
		errChan <- err
		// Send the error to the renderer.
		// A bolt database error is fatal.
		txmessage.Fatal = true
		txmessage.ErrorMessage = fmt.Sprintf("handleCheckCopy: copyservice.Check(copy, rxmessage.Solution, stores.KeyCode, rxmessage.WPM, rxmessage.StoreResults): error is %s", err.Error())
		sending <- txmessage
		return
	}
	// 3. Send the results to the renderer.
	sending <- txmessage
}
