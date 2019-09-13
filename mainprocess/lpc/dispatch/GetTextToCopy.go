package dispatch

import (
	"fmt"
	"log"

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

// handleGetTextToCopy is the *message.GetTextToCopyRendererToMainProcess handler.
// It's response back to the renderer is the *message.GetTextToCopyMainProcessToRenderer.
// Param rxmessage *message.GetTextToCopyRendererToMainProcess is the message received from the renderer.
// Param sending is the channel to use to send a *message.GetTextToCopyMainProcessToRenderer message back to the renderer.
// Param eojing lpc.EOJer ( End Of Job ) is an interface for your go routine to receive a stop signal.
//   It signals go routines that they must stop because the main process is ending.
//   So only use it inside a go routine if you have one.
//   In your go routine
//     1. Get a channel to listen to with eojing.NewEOJ().
//     2. Before your go routine returns, release that channel with eojing.Release().
// The func is simple:
// 1. Get the text and ditdah to copy.
// 2. Get the copy WPM for rendering the ditdahs.
// 3. Call the renderer back with the text.
func handleGetTextToCopy(rxmessage *message.GetTextToCopyRendererToMainProcess, sending lpc.Sending, eojing lpc.EOJer, stores *store.Stores) {
	txmessage := &message.GetTextToCopyMainProcessToRenderer{
		State: rxmessage.State,
	}
	// 1. Get the text and ditdah to copy.
	var err error
	if txmessage.Solution, err = copyservice.GetKeyCodes(stores.KeyCode); err != nil {
		// Calling back the error.
		txmessage.Error = true
		txmessage.ErrorMessage = fmt.Sprintf("mainProcessGetTextToCopy: keyservice.GetTextToCopy(stores.KeyCode): error is %s", err.Error())
		sending <- txmessage
		log.Println(txmessage.ErrorMessage)
		return
	}
	// 3. Get the copy WPM for rendering the ditdahs.
	var r *record.WPM
	if r, err = stores.WPM.GetCopyWPM(); err != nil {
		// Calling back the error.
		txmessage.Error = true
		txmessage.ErrorMessage = fmt.Sprintf("mainProcessGetTextToCopy: wPMStorer.GetCopyWPM(): error is %s\n", err.Error())
		sending <- txmessage
		log.Println(txmessage.ErrorMessage)
		return
	}
	// 4. Call the renderer back with the text.
	txmessage.WPM = r.WPM
	sending <- txmessage
}
