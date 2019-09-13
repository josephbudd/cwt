package keytestpanel

import (
	"github.com/josephbudd/cwt/domain/lpc/message"
	"github.com/josephbudd/cwt/domain/store/record"
)

/*

	Panel name: KeyTestPanel

*/

// panelCaller communicates with the main process via an asynchrounous connection.
type panelCaller struct {
	group      *panelGroup
	presenter  *panelPresenter
	controller *panelController

	/* NOTE TO DEVELOPER. Step 1 of 4.

	// 1.1: Declare your panelCaller members.

	*/
}

/* NOTE TO DEVELOPER. Step 2 of 4.

// 2.1: Define your funcs which send a message to the main process.
// 2.2: Define your funcs which receive a message from the main process.

*/

/* Interface implementation:

panelCaller implements widgets.UserKeyChecker with funcs
 * GetKeyCodesWPM
 * CheckUserKey

*/

// GetKeyCodesWPM gets new text for the user to key and the wpm required for keying.
func (caller *panelCaller) GetKeyCodesWPM() {
	msg := &message.GetTextToKeyRendererToMainProcess{
		State: state,
	}
	sendCh <- msg
}

func (caller *panelCaller) getTextWPMToKeyCB(msg *message.GetTextToKeyMainProcessToRenderer) {
	if msg.State&state != state {
		return
	}
	if msg.Error {
		tools.Error(msg.ErrorMessage)
		return
	}
	caller.controller.keyWidget.SetKeyCodesWPM(msg.Solution, msg.WPM)
}

// CheckUserKey checks the user's keying ability.
func (caller *panelCaller) CheckUserKey(milliSeconds []int64, solution [][]*record.KeyCode, wpm uint64) {
	msg := &message.CheckKeyRendererToMainProcess{
		MilliSeconds: milliSeconds,
		Solution:     solution,
		WPM:          wpm,
		State:        state,
		StoreResults: true,
	}
	sendCh <- msg
}

func (caller *panelCaller) checkKeyCB(msg *message.CheckKeyMainProcessToRenderer) {
	if msg.State&state != state {
		return
	}
	if msg.Error {
		tools.Error(msg.ErrorMessage)
		return
	}
	caller.controller.keyWidget.ShowResults(msg.CorrectCount, msg.IncorrectCount, msg.KeyedCount, msg.TestResults)
}

/* Interface implementation:

panelCaller implements widgets.Metronomer with funcs
 * StartMetronome
 * StopMetronome

*/

// StartMetronome starts the metronome.
func (caller *panelCaller) StartMetronome(wpm uint64) {
	msg := &message.MetronomeRendererToMainProcess{
		Run:   true,
		State: state,
		WPM:   wpm,
	}
	sendCh <- msg
}

// StopMetronome stops the metronome.
func (caller *panelCaller) StopMetronome() {
	msg := &message.MetronomeRendererToMainProcess{
		Run:   false,
		State: state,
	}
	sendCh <- msg
}

func (caller *panelCaller) metronomeCB(msg *message.MetronomeMainProcessToRenderer) {
	if msg.State&state != state {
		return
	}
	if msg.Error {
		tools.Error(msg.ErrorMessage)
		return
	}
	if msg.Run {
		tools.Success("Metronome Started")
	} else {
		tools.Success("Metronome Stopped")
	}
}

// dispatchMessages dispatches LPC messages from the main process.
// It stops when it receives on the eoj channel.
func (caller *panelCaller) dispatchMessages() {
	go func() {
		for {
			select {
			case <-eojCh:
				return
			case msg := <-receiveCh:
				// A message sent from the main process to the renderer.
				switch msg := msg.(type) {

				/* NOTE TO DEVELOPER. Step 3 of 4.

				// 3.1:   Remove the default clause below.
				// 3.2.a: Add a case for each of the messages
				//          that you are expecting from the main process.
				// 3.2.b: In that case statement, pass the message to your message receiver func.

				*/

				case *message.GetTextToKeyMainProcessToRenderer:
					caller.getTextWPMToKeyCB(msg)
				case *message.CheckKeyMainProcessToRenderer:
					caller.checkKeyCB(msg)
				case *message.MetronomeMainProcessToRenderer:
					caller.metronomeCB(msg)
				}
			}
		}
	}()

	return
}

// initialCalls makes the first calls to the main process.
func (caller *panelCaller) initialCalls() {

	/* NOTE TO DEVELOPER. Step 4 of 4.

	//4.1: Make any initial calls to the main process that must be made when the app starts.

	// example:

	// import "github.com/josephbudd/cwt/domain/data/loglevels"
	// import "github.com/josephbudd/cwt/domain/lpc/message"

	msg := &message.LogRendererToMainProcess{
		Level:   loglevels.LogLevelInfo,
		Message: "Started",
	}
	sendCh <- msg

	*/
}
