// +build js, wasm

package keywpmpanel

import (
	"github.com/josephbudd/cwt/domain/lpc/message"
	"github.com/josephbudd/cwt/domain/store/record"
	"github.com/josephbudd/cwt/rendererprocess/api/display"
)

/*

	Panel name: KeyWPMPanel

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

// Set WPM.

func (messenger *panelMessenger) updateKeyWPM(r *record.WPM) {
	msg := &message.UpdateKeyWPMRendererToMainProcess{
		Record: r,
	}
	sendCh <- msg
}

func (messenger *panelMessenger) updateKeyWPMCB(msg *message.UpdateKeyWPMMainProcessToRenderer) {
	if msg.Error {
		display.Error(msg.ErrorMessage)
		return
	}
	messenger.controller.processWPM(msg.Record)
	display.Success("Key WPM Updated.")
}

// GetWPM

func (messenger *panelMessenger) getKeyWPMCB(msg *message.GetKeyWPMMainProcessToRenderer) {
	if msg.Error {
		display.Error(msg.ErrorMessage)
		return
	}
	messenger.controller.processWPM(msg.Record)
}

// dispatchMessages dispatches LPC messages from the main process.
// It stops when context is done.
func (messenger *panelMessenger) dispatchMessages() {
	go func() {
		for {
			select {
			case <-rendererProcessCtx.Done():
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

				case *message.UpdateKeyWPMMainProcessToRenderer:
					messenger.updateKeyWPMCB(msg)
				case *message.GetKeyWPMMainProcessToRenderer:
					messenger.getKeyWPMCB(msg)
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

	msg := &message.GetKeyWPMRendererToMainProcess{}
	sendCh <- msg
}
