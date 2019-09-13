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

// handleGetCopyWPM is the *message.GetCopyWPMRendererToMainProcess handler.
// It's response back to the renderer is the *message.GetCopyWPMMainProcessToRenderer.
// Param rxmessage *message.GetCopyWPMRendererToMainProcess is the message received from the renderer.
// Param sending is the channel to use to send a *message.GetCopyWPMMainProcessToRenderer message back to the renderer.
// Param eojing lpc.EOJer ( End Of Job ) is an interface for your go routine to receive a stop signal.
//   It signals go routines that they must stop because the main process is ending.
//   So only use it inside a go routine if you have one.
//   In your go routine
//     1. Get a channel to listen to with eojing.NewEOJ().
//     2. Before your go routine returns, release that channel with eojing.Release().
// The func is simple:
// 1. Get the keyCodes from the repo. Call back any errors or not found.
// 2. Send the wpm to the renderer.
func handleGetCopyWPM(rxmessage *message.GetCopyWPMRendererToMainProcess, sending lpc.Sending, eojing lpc.EOJer, stores *store.Stores) {
	txmessage := &message.GetCopyWPMMainProcessToRenderer{}
	// 1. Get the keyCode from the repo.
	var err error
	if txmessage.Record, err = stores.WPM.GetCopyWPM(); err != nil {
		txmessage.Error = true
		txmessage.ErrorMessage = fmt.Sprintf("handleGetCopyWPM: stores.WPM.GetCopyWPM(): error is %s", err.Error())
		sending <- txmessage
		log.Println(txmessage.ErrorMessage)
		return
	}
	if txmessage.Record == nil {
		// Calling back the not found error.
		// This will only happen in development. It means that the data store is not getting initialized properly.
		txmessage.Error = true
		txmessage.ErrorMessage = "handleGetCopyWPM: stores.WPM.GetCopyWPM(): error is Not Found."
		sending <- txmessage
		log.Println(txmessage.ErrorMessage)
		return
	}
	// 2. Send the wpm to the renderer.
	sending <- txmessage
}
