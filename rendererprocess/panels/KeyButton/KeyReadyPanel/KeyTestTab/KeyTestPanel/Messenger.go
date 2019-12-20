// +build js, wasm

package keytestpanel

import (
	"github.com/josephbudd/cwt/domain/lpc/message"
	"github.com/josephbudd/cwt/domain/store/record"
	"github.com/josephbudd/cwt/rendererprocess/api/display"
)

/*

	Panel name: KeyTestPanel

*/

// panelMessenger communicates with the main process via an asynchrounous connection.
type panelMessenger struct {
	group      *panelGroup
	presenter  *panelPresenter
	controller *panelController

	/* NOTE TO DEVELOPER. Step 1 of 4.

	// 1.1: Declare your panelMessenger members.

	// example:

	state uint64

	*/
}

/* NOTE TO DEVELOPER. Step 2 of 4.

// 2.1: Define your funcs which send a message to the main process.
// 2.2: Define your funcs which receive a message from the main process.

// example:

import "github.com/josephbudd/cwt/domain/store/record"
import "github.com/josephbudd/cwt/domain/lpc/message"
import "github.com/josephbudd/cwt/rendererprocess/api/display"

// Add Customer.

func (messenger *panelMessenger) addCustomer(r *record.Customer) {
	msg := &message.AddCustomerRendererToMainProcess{
		UniqueID: messenger.uniqueID,
		Record:   record,
	}
	sendCh <- msg
}

func (messenger *panelMessenger) addCustomerRX(msg *message.AddCustomerMainProcessToRenderer) {
	if msg.UniqueID == messenger.uniqueID {
		if msg.Error {
			display.Error(msg.ErrorMessage)
			return
		}
		// no errors
		display.Success("Customer Added.")
	}
}

*/

// GetKeyCodesWPM gets new text for the user to key and the wpm required for keying.
func (messenger *panelMessenger) GetKeyCodesWPM() {
	msg := &message.GetTextToKeyRendererToMainProcess{
		State: state,
	}
	sendCh <- msg
}

func (messenger *panelMessenger) getTextWPMToKeyCB(msg *message.GetTextToKeyMainProcessToRenderer) {
	if msg.State&state != state {
		return
	}
	if msg.Error {
		display.Error(msg.ErrorMessage)
		return
	}
	messenger.controller.keyWidget.SetKeyCodesWPM(msg.Solution, msg.WPM)
}

// CheckUserKey checks the user's keying ability.
func (messenger *panelMessenger) CheckUserKey(milliSeconds []int64, solution [][]*record.KeyCode, wpm uint64) {
	msg := &message.CheckKeyRendererToMainProcess{
		MilliSeconds: milliSeconds,
		Solution:     solution,
		WPM:          wpm,
		State:        state,
		StoreResults: true,
	}
	sendCh <- msg
}

func (messenger *panelMessenger) checkKeyCB(msg *message.CheckKeyMainProcessToRenderer) {
	if msg.State&state != state {
		return
	}
	if msg.Error {
		display.Error(msg.ErrorMessage)
		return
	}
	messenger.controller.keyWidget.ShowResults(msg.CorrectCount, msg.IncorrectCount, msg.MaxCount, msg.TestResults)
}

/* Interface implementation:

panelMessenger implements widgets.Metronomer with funcs
 * StartMetronome
 * StopMetronome

*/

// StartMetronome starts the metronome.
func (messenger *panelMessenger) StartMetronome(wpm uint64) {
	msg := &message.MetronomeRendererToMainProcess{
		Run:   true,
		State: state,
		WPM:   wpm,
	}
	sendCh <- msg
}

// StopMetronome stops the metronome.
func (messenger *panelMessenger) StopMetronome() {
	msg := &message.MetronomeRendererToMainProcess{
		Run:   false,
		State: state,
	}
	sendCh <- msg
}

func (messenger *panelMessenger) metronomeCB(msg *message.MetronomeMainProcessToRenderer) {
	if msg.State&state != state {
		return
	}
	if msg.Error {
		display.Error(msg.ErrorMessage)
		return
	}
	if msg.Run {
		display.Success("Metronome Started")
	} else {
		display.Success("Metronome Stopped")
	}
}

// dispatchMessages dispatches LPC messages from the main process.
// It stops when it receives on the eoj channel.
func (messenger *panelMessenger) dispatchMessages() {
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

				// example:

				import "github.com/josephbudd/cwt/domain/lpc/message"

				case *message.AddCustomerMainProcessToRenderer:
					messenger.addCustomerRX(msg)

				*/

				default:
					_ = msg
				}
			}
		}
	}()

	return
}

// initialSends sends the first messages to the main process.
func (messenger *panelMessenger) initialSends() {

	/* NOTE TO DEVELOPER. Step 4 of 4.

	//4.1: Send messages to the main process right when the app starts.

	// example:

	import "github.com/josephbudd/cwt/domain/data/loglevels"
	import "github.com/josephbudd/cwt/domain/lpc/message"

	msg := &message.LogRendererToMainProcess{
		Level:   loglevels.LogLevelInfo,
		Message: "Started",
	}
	sendCh <- msg

	*/
}
