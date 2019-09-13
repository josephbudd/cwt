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

// handleUpdateKeyCode is the *message.UpdateKeyCodeRendererToMainProcess handler.
// It's response back to the renderer is the *message.UpdateKeyCodeMainProcessToRenderer.
// Param rxmessage *message.UpdateKeyCodeRendererToMainProcess is the message received from the renderer.
// Param sending is the channel to use to send a *message.UpdateKeyCodeMainProcessToRenderer message back to the renderer.
// Param eojing lpc.EOJer ( End Of Job ) is an interface for your go routine to receive a stop signal.
//   It signals go routines that they must stop because the main process is ending.
//   So only use it inside a go routine if you have one.
//   In your go routine
//     1. Get a channel to listen to with eojing.NewEOJ().
//     2. Before your go routine returns, release that channel with eojing.Release().
// The func is simple:
// 1. Get the keyCode from the repo. Call back any errors or not found.
// 2. Send the keycode to the renderer.
func handleUpdateKeyCode(rxmessage *message.UpdateKeyCodeRendererToMainProcess, sending lpc.Sending, eojing lpc.EOJer, stores *store.Stores) {
	txmessage := &message.UpdateKeyCodeMainProcessToRenderer{
		Record: rxmessage.Record,
		State:  rxmessage.State,
	}
	// 1. Update the keyCode.
	if err := stores.KeyCode.Update(rxmessage.Record); err != nil {
		// Calling back the error.
		txmessage.Error = true
		txmessage.ErrorMessage = fmt.Sprintf("handleUpdateKeyCode: stores.KeyCode.Update(rxmessage.Record): error is %s\n", err.Error())
		sending <- txmessage
		log.Println(txmessage.ErrorMessage)
		return
	}
	// 3. Send the keycode to the renderer.
	sending <- txmessage
}
