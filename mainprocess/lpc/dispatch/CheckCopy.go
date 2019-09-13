package dispatch

import (
	"fmt"
	"log"
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
// Param rxmessage *message.CheckCopyRendererToMainProcess is the message received from the renderer.
// Param sending is the channel to use to send a *message.CheckCopyMainProcessToRenderer message back to the renderer.
// Param eojing lpc.EOJer ( End Of Job ) is an interface for your go routine to receive a stop signal.
//   It signals go routines that they must stop because the main process is ending.
//   So only use it inside a go routine if you have one.
//   In your go routine
//     1. Get a channel to listen to with eojing.NewEOJ().
//     2. Before your go routine returns, release that channel with eojing.Release().
// The func is simple:
// 1. Convert the copy string to key code records.
// 2. Check the copy.
// 3. Send the results to the renderer.
func handleCheckCopy(rxmessage *message.CheckCopyRendererToMainProcess, sending lpc.Sending, eojing lpc.EOJer, stores *store.Stores) {
	txmessage := &message.CheckCopyMainProcessToRenderer{
		State: rxmessage.State,
	}
	// 1. Convert the copy string to key code records.
	// convert copy [][]string to [][]*record.KeyCode
	var rr []*record.KeyCode
	var err error
	if rr, err = stores.KeyCode.GetAll(); err != nil {
		// Calling back the error.
		txmessage.Error = true
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
	if txmessage.CorrectCount, txmessage.IncorrectCount, txmessage.KeyedCount, txmessage.TestResults, err = copyservice.Check(copy, rxmessage.Solution, stores.KeyCode, rxmessage.WPM, rxmessage.StoreResults); err != nil {
		txmessage.Error = true
		txmessage.ErrorMessage = fmt.Sprintf("handleCheckCopy: copyservice.Check(copy, rxmessage.Solution, stores.KeyCode, rxmessage.WPM, rxmessage.StoreResults): error is %s", err.Error())
		sending <- txmessage
		log.Println(txmessage.ErrorMessage)
		return
	}
	// 3. Send the results to the renderer.
	sending <- txmessage
}
