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

// handleGetKeyCodes is the *message.GetKeyCodesRendererToMainProcess handler.
// It's response back to the renderer is the *message.GetKeyCodesMainProcessToRenderer.
// Param rxmessage *message.GetKeyCodesRendererToMainProcess is the message received from the renderer.
// Param sending is the channel to use to send a *message.GetKeyCodesMainProcessToRenderer message back to the renderer.
// Param eojing lpc.EOJer ( End Of Job ) is an interface for your go routine to receive a stop signal.
//   It signals go routines that they must stop because the main process is ending.
//   So only use it inside a go routine if you have one.
//   In your go routine
//     1. Get a channel to listen to with eojing.NewEOJ().
//     2. Before your go routine returns, release that channel with eojing.Release().
// The func is simple:
// 1. Get the keyCodes from the repo. Call back any errors or not found.
// 2. Call the renderer back with the keyCode records.
func handleGetKeyCodes(rxmessage *message.GetKeyCodesRendererToMainProcess, sending lpc.Sending, eojing lpc.EOJer, stores *store.Stores) {
	txmessage := &message.GetKeyCodesMainProcessToRenderer{
		State: rxmessage.State,
	}
	// 1. Get the keyCodes from the repo.
	var err error
	if txmessage.Records, err = stores.KeyCode.GetAll(); err != nil {
		txmessage.Error = true
		txmessage.ErrorMessage = fmt.Sprintf("handleGetKeyCodes: stores.KeyCode.GetAll(): error is %s", err.Error())
		sending <- txmessage
		log.Println(txmessage.ErrorMessage)
		return
	}
	// 2. Send the key codes to the renderer.
	sending <- txmessage
}
