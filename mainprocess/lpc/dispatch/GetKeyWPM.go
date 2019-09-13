package dispatch

import (
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

// handleGetKeyWPM is the *message.GetKeyWPMRendererToMainProcess handler.
// It's response back to the renderer is the *message.GetKeyWPMMainProcessToRenderer.
// Param rxmessage *message.GetKeyWPMRendererToMainProcess is the message received from the renderer.
// Param sending is the channel to use to send a *message.GetKeyWPMMainProcessToRenderer message back to the renderer.
// Param eojing lpc.EOJer ( End Of Job ) is an interface for your go routine to receive a stop signal.
//   It signals go routines that they must stop because the main process is ending.
//   So only use it inside a go routine if you have one.
//   In your go routine
//     1. Get a channel to listen to with eojing.NewEOJ().
//     2. Before your go routine returns, release that channel with eojing.Release().
// The func is simple:
// 1. Get the wpm from the repo. Call back any errors or not found.
// 2. Send the wpm to the renderer.
func handleGetKeyWPM(rxmessage *message.GetKeyWPMRendererToMainProcess, sending lpc.Sending, eojing lpc.EOJer, stores *store.Stores) {
	txmessage := &message.GetKeyWPMMainProcessToRenderer{}
	// 1. Get the wpm from the repo.
	var err error
	if txmessage.Record, err = stores.WPM.GetKeyWPM(); err != nil {
		txmessage.Error = true
		txmessage.ErrorMessage = fmt.Sprintf("handleGetKeyWPM: stores.WPM.GetKeyWPM(): error is %s", err.Error())
		sending <- txmessage
		log.Println(txmessage.ErrorMessage)
		return
	}
	if txmessage.Record == nil {
		// Calling back the not found error.
		// This will only happen in development. It means that the data store is not getting initialized properly.
		txmessage.Error = true
		txmessage.ErrorMessage = "handleGetKeyWPM: stores.WPM.GetKeyWPM(): error is Not Found."
		sending <- txmessage
		log.Println(txmessage.ErrorMessage)
		return
	}
	// 2. Send the wpm to the renderer.
	sending <- txmessage
}
