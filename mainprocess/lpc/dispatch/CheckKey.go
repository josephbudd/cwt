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

// handleCheckKey is the *message.CheckKeyRendererToMainProcess handler.
// It's response back to the renderer is the *message.CheckKeyMainProcessToRenderer.
// Param rxmessage *message.CheckKeyRendererToMainProcess is the message received from the renderer.
// Param sending is the channel to use to send a *message.CheckKeyMainProcessToRenderer message back to the renderer.
// Param eojing lpc.EOJer ( End Of Job ) is an interface for your go routine to receive a stop signal.
//   It signals go routines that they must stop because the main process is ending.
//   So only use it inside a go routine if you have one.
//   In your go routine
//     1. Get a channel to listen to with eojing.NewEOJ().
//     2. Before your go routine returns, release that channel with eojing.Release().
// The func is simple:
// 1. Copy: Convert the miliseconds to key code records.
// 2. Check the copy against the solution.
// 3. Send the results to the renderer.
func handleCheckKey(rxmessage *message.CheckKeyRendererToMainProcess, sending lpc.Sending, eojing lpc.EOJer, stores *store.Stores) {
	txmessage := &message.CheckKeyMainProcessToRenderer{
		State: rxmessage.State,
	}
	// 1. Copy: Convert the miliseconds to key code records.
	var copiedWords [][]*record.KeyCode
	var err error
	if copiedWords, err = keyservice.Copy(rxmessage.MilliSeconds, rxmessage.WPM, stores.KeyCode); err != nil {
		txmessage.Error = true
		txmessage.ErrorMessage = fmt.Sprintf("handleCheckKey: keyservice.Copy(rxmessage.MilliSeconds, rxmessage.WPM, stores.KeyCode): error is %s", err.Error())
		sending <- txmessage
		log.Println(txmessage.ErrorMessage)
		return
	}
	// 2. Check the copy against the solution.
	if txmessage.CorrectCount, txmessage.IncorrectCount, txmessage.KeyedCount, txmessage.TestResults, err = keyservice.Check(copiedWords, rxmessage.Solution, stores.KeyCode, rxmessage.WPM, rxmessage.StoreResults); err != nil {
		txmessage.Error = true
		txmessage.ErrorMessage = fmt.Sprintf("handleCheckKey: keyservice.Check(copiedWords, rxmessage.Solution, stores.KeyCode, rxmessage.WPM, rxmessage.StoreResults): error is %s", err.Error())
		sending <- txmessage
		log.Println(txmessage.ErrorMessage)
		return
	}
	// 3. Send the results to the renderer.
	sending <- txmessage
}
