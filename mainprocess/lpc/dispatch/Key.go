package dispatch

import (
	"fmt"
	"log"
	"strings"

	"github.com/josephbudd/cwt/domain/lpc/message"
	"github.com/josephbudd/cwt/domain/store"
	"github.com/josephbudd/cwt/mainprocess/lpc"
	"github.com/josephbudd/cwt/mainprocess/services/copyservice"
)

/*
	YOU MAY EDIT THIS FILE.

	Rekickwasm will preserve this file for you.
	Kicklpc will not edit this file.

*/

// handleKey is the *message.KeyRendererToMainProcess handler.
// It's response back to the renderer is the *message.KeyMainProcessToRenderer.
// Param rxmessage *message.KeyRendererToMainProcess is the message received from the renderer.
// Param sending is the channel to use to send a *message.KeyMainProcessToRenderer message back to the renderer.
// Param eojing lpc.EOJer ( End Of Job ) is an interface for your go routine to receive a stop signal.
//   It signals go routines that they must stop because the main process is ending.
//   So only use it inside a go routine if you have one.
//   In your go routine
//     1. Get a channel to listen to with eojing.NewEOJ().
//     2. Before your go routine returns, release that channel with eojing.Release().
// The func is simple:
// 1. Turn off the keying if requested.
// 2. Build the morse code text.
// 3. Key the morse code. Call back any errors.
func handleKey(rxmessage *message.KeyRendererToMainProcess, sending lpc.Sending, eojing lpc.EOJer, stores *store.Stores) {
	txmessage := &message.KeyMainProcessToRenderer{
		Run:   rxmessage.Run,
		State: rxmessage.State,
	}
	// 1. Turn off the keying if requested.
	if !rxmessage.Run {
		copyservice.StopKeying()
		sending <- txmessage
		return
	}
	// 2. Build the morse code text.
	ditdahs := make([]string, 0, len(rxmessage.Solution))
	for _, line := range rxmessage.Solution {
		ditdahWord := make([]string, 0, len(line))
		for _, r := range line {
			ditdahWord = append(ditdahWord, r.DitDah)
		}
		ditdahs = append(ditdahs, strings.Join(ditdahWord, " "))
	}
	// 3. Key the morse code.
	if err := copyservice.Key(ditdahs, rxmessage.WPM, rxmessage.Pause); err != nil {
		txmessage.Error = true
		txmessage.ErrorMessage = fmt.Sprintf("mainProcessKey:  ditdah.Key(rxmessage.Ditdah, rxmessage.WPM, rxmessage.Delay): error is %s", err.Error())
		sending <- txmessage
		log.Println(txmessage.ErrorMessage)
		return
	}
	// no error so call back to the renderer.
	sending <- txmessage
	return
}
