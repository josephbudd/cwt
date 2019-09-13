package punctuationpanel

import (
	"fmt"

	"github.com/josephbudd/cwt/domain/data"
	"github.com/josephbudd/cwt/domain/lpc/message"
	"github.com/josephbudd/cwt/domain/store/record"
)

/*

	Panel name: PunctuationPanel

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

// example:

// import "github.com/josephbudd/cwt/domain/store/record"
// import "github.com/josephbudd/cwt/domain/lpc/message"


// Add Customer.

func (caller *panelCaller) addCustomer(r *record.Customer) {
	msg := &message.AddCustomerRendererToMainProcess{
		UniqueID: caller.uniqueID,
		Record:   record,
	}
	sendCh <- msg
}

func (caller *panelCaller) addCustomerRX(msg *message.AddCustomerMainProcessToRenderer) {
	if msg.UniqueID == caller.uniqueID {
		if msg.Error {
			tools.Error(msg.ErrorMessage)
			return
		}
		// no errors
		tools.Success("Customer Added.")
	}
}

*/

// UpdateKeyCode

func (caller *panelCaller) updateKeyCode(r *record.KeyCode) {
	msg := &message.UpdateKeyCodeRendererToMainProcess{
		Record: r,
		State:  state,
	}
	sendCh <- msg
}

func (caller *panelCaller) updateKeyCodeCB(msg *message.UpdateKeyCodeMainProcessToRenderer) {
	if msg.State&state != state {
		return
	}
	if msg.Error {
		tools.Error(msg.ErrorMessage)
		return
	}
	if msg.Record.Selected {
		tools.Success(fmt.Sprintf("The punctuation %q is now selected.", msg.Record.Name))
	} else {
		tools.Success(fmt.Sprintf("The punctuation %q is no longer selected.", msg.Record.Name))
	}
}

// GetKeyCodes

func (caller *panelCaller) getKeyCodesCB(msg *message.GetKeyCodesMainProcessToRenderer) {
	if msg.Error {
		// This panel will handle the GetKeyCode errors not the other panels.
		tools.Error(msg.ErrorMessage)
		return
	}
	caller.controller.setup(getPunctuationKeyCodes(msg.Records))
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

				case *message.UpdateKeyCodeMainProcessToRenderer:
					caller.updateKeyCodeCB(msg)
				case *message.GetKeyCodesMainProcessToRenderer:
					caller.getKeyCodesCB(msg)
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

	*/
}
