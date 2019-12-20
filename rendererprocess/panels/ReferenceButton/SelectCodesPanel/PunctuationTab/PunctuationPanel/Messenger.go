// +build js, wasm

package punctuationpanel

import (
	"fmt"

	"github.com/josephbudd/cwt/domain/data"
	"github.com/josephbudd/cwt/domain/lpc/message"
	"github.com/josephbudd/cwt/domain/store/record"
	"github.com/josephbudd/cwt/rendererprocess/api/display"
)

/*

	Panel name: PunctuationPanel

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

// UpdateKeyCode

func (messenger *panelMessenger) updateKeyCode(r *record.KeyCode) {
	msg := &message.UpdateKeyCodeRendererToMainProcess{
		Record: r,
		State:  state,
	}
	sendCh <- msg
}

func (messenger *panelMessenger) updateKeyCodeCB(msg *message.UpdateKeyCodeMainProcessToRenderer) {
	if msg.State&state != state {
		return
	}
	if msg.Error {
		display.Error(msg.ErrorMessage)
		return
	}
	if msg.Record.Selected {
		display.Success(fmt.Sprintf("The punctuation %q is now selected.", msg.Record.Name))
	} else {
		display.Success(fmt.Sprintf("The punctuation %q is no longer selected.", msg.Record.Name))
	}
}

// GetKeyCodes

func (messenger *panelMessenger) getKeyCodesCB(msg *message.GetKeyCodesMainProcessToRenderer) {
	if msg.Error {
		// This panel will handle the GetKeyCode errors not the other panels.
		display.Error(msg.ErrorMessage)
		return
	}
	messenger.controller.setup(getPunctuationKeyCodes(msg.Records))
}

func getPunctuationKeyCodes(rr []*record.KeyCode) (group []*record.KeyCode) {
	group = make([]*record.KeyCode, 0, len(rr))
	for _, r := range rr {
		if r.Type == data.KeyCodeTypePunctuation {
			group = append(group, r)
		}
	}
	l := len(group) - 1
	for i := 0; i < l; i++ {
		for j := i + 1; j <= l; j++ {
			if group[i].Name > group[j].Name {
				temp := group[i]
				group[i] = group[j]
				group[j] = temp
			}
		}
	}
	return
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

				case *message.UpdateKeyCodeMainProcessToRenderer:
					messenger.updateKeyCodeCB(msg)
				case *message.GetKeyCodesMainProcessToRenderer:
					messenger.getKeyCodesCB(msg)
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
