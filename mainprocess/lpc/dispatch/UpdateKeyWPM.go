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

// handleUpdateKeyWPM is the *message.UpdateKeyWPMRendererToMainProcess handler.
// It's response back to the renderer is the *message.UpdateKeyWPMMainProcessToRenderer.
// Param rxmessage *message.UpdateKeyWPMRendererToMainProcess is the message received from the renderer.
// Param sending is the channel to use to send a *message.UpdateKeyWPMMainProcessToRenderer message back to the renderer.
// Param eojing lpc.EOJer ( End Of Job ) is an interface for your go routine to receive a stop signal.
//   It signals go routines that they must stop because the main process is ending.
//   So only use it inside a go routine if you have one.
//   In your go routine
//     1. Get a channel to listen to with eojing.NewEOJ().
//     2. Before your go routine returns, release that channel with eojing.Release().
// The func is simple:
// 1. Update the wpm.
// 2. Send the wpm to the renderer.
func handleUpdateKeyWPM(rxmessage *message.UpdateKeyWPMRendererToMainProcess, sending lpc.Sending, eojing lpc.EOJer, stores *store.Stores) {
	txmessage := &message.UpdateKeyWPMMainProcessToRenderer{
		Record: rxmessage.Record,
	}
	// 1. Update the wpm.
	if err := stores.WPM.Update(rxmessage.Record); err != nil {
		txmessage.Error = true
		txmessage.ErrorMessage = fmt.Sprintf("handleUpdateKeyWPM: stores.WPM.Update(rxmessage.Record): error is %s\n", err.Error())
		sending <- txmessage
		log.Println(txmessage.ErrorMessage)
		return
	}
	// 2. Send the wpm to the renderer.
	sending <- txmessage
}
