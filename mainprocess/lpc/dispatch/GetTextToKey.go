package dispatch

import (
	"fmt"
	"log"

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
// Param rxmessage *message.GetTextToKeyRendererToMainProcess is the message received from the renderer.
// Param sending is the channel to use to send a *message.GetTextToKeyMainProcessToRenderer message back to the renderer.
// Param eojing lpc.EOJer ( End Of Job ) is an interface for your go routine to receive a stop signal.
//   It signals go routines that they must stop because the main process is ending.
//   So only use it inside a go routine if you have one.
//   In your go routine
//     1. Get a channel to listen to with eojing.NewEOJ().
//     2. Before your go routine returns, release that channel with eojing.Release().
func handleGetTextToKey(rxmessage *message.GetTextToKeyRendererToMainProcess, sending lpc.Sending, eojing lpc.EOJer, stores *store.Stores) {
	txmessage := &message.GetTextToKeyMainProcessToRenderer{
		State:    rxmessage.State,
		Practice: rxmessage.Practice,
	}
	// 1 Get the wpm to key.
	var r *record.WPM
	var err error
	if r, err = stores.WPM.GetKeyWPM(); err != nil {
		// Calling back the error.
		txmessage.Error = true
		txmessage.ErrorMessage = fmt.Sprintf("mainProcessGetTextWPMToKey: keyCodeStorer.GetKeyWPM(): error is %s", err.Error())
		sending <- txmessage
		log.Println(txmessage.ErrorMessage)
		return
	}
	txmessage.WPM = r.WPM
	// 3. Get the text for the user to key from the repo and the help.
	if rxmessage.Practice {
		// practicing
		if txmessage.Solution, txmessage.Help, err = keyservice.GetPracticeKeyCodes(stores.KeyCode, r.WPM); err != nil {
			// Calling back the error.
			txmessage.Error = true
			txmessage.ErrorMessage = fmt.Sprintf("mainProcessGetTextWPMToKey: keyservice.GetPracticeKeyCodes(stores.KeyCode, r.WPM): error is %s", err.Error())
			sending <- txmessage
			log.Println(txmessage.ErrorMessage)
			return
		}
	} else {
		// testing not practicing
		if txmessage.Solution, err = keyservice.GetTestKeyCodes(stores.KeyCode); err != nil {
			// Calling back the error.
			txmessage.Error = true
			txmessage.ErrorMessage = fmt.Sprintf("mainProcessGetTextWPMToKey: keyservice.GetTestKeyCodes(keyCodeStorer): error is %s\n", err.Error())
			sending <- txmessage
			log.Println(txmessage.ErrorMessage)
			return
		}
	}
	sending <- txmessage
}
