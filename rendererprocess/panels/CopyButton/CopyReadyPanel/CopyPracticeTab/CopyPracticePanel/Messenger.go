// +build js, wasm

package copypracticepanel

import (
	"github.com/josephbudd/cwt/domain/lpc/message"
	"github.com/josephbudd/cwt/domain/store/record"
	"github.com/josephbudd/cwt/rendererprocess/api/display"
)

/*

	Panel name: CopyPracticePanel

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

// Check copy.

func (messenger *panelMessenger) checkCopy(copy []string, solution [][]*record.KeyCode, wpm uint64) {
	msg := &message.CheckCopyRendererToMainProcess{
		Copy:     copy,
		Solution: solution,
		WPM:      wpm,
		State:    state,
	}
	sendCh <- msg
}

func (messenger *panelMessenger) checkCopyCB(msg *message.CheckCopyMainProcessToRenderer) {
	if msg.State&state != state {
		return
	}
	if msg.Error {
		display.Error(msg.ErrorMessage)
		return
	}
	messenger.presenter.showResult(msg.CorrectCount, msg.IncorrectCount, msg.MaxCount, msg.TestResults)
}

// Get text to copy.

func (messenger *panelMessenger) getTextToCopy() {
	msg := &message.GetTextToCopyRendererToMainProcess{
		State: state,
	}
	sendCh <- msg
}

func (messenger *panelMessenger) getTextToCopyCB(msg *message.GetTextToCopyMainProcessToRenderer) {
	if msg.State&state != state {
		return
	}
	if msg.Error {
		display.Error(msg.ErrorMessage)
		return
	}
	messenger.controller.processTextToCopy(msg.Solution)
}

// Key.

func (messenger *panelMessenger) key(solution [][]*record.KeyCode, wpm, pause uint64) {
	msg := &message.KeyRendererToMainProcess{
		Solution: solution,
		WPM:      wpm,
		Pause:    pause,
		State:    state,
		Run:      true,
	}
	sendCh <- msg
}

func (messenger *panelMessenger) stopKeying() {
	msg := &message.KeyRendererToMainProcess{
		State: state,
	}
	sendCh <- msg
}

func (messenger *panelMessenger) keyCB(msg *message.KeyMainProcessToRenderer) {
	if msg.State&state != state {
		return
	}
	if msg.Error {
		display.Error(msg.ErrorMessage)
		return
	}
	if msg.Run {
		messenger.controller.processKeyFinished()
	} else {
		messenger.controller.processKeyStopped()
	}
}

// wpm

func (messenger *panelMessenger) updateCopyWPMCB(msg *message.UpdateCopyWPMMainProcessToRenderer) {
	if msg.Error {
		display.Error(msg.ErrorMessage)
		return
	}
	messenger.controller.processWPM(msg.Record.WPM)
}

func (messenger *panelMessenger) getCopyWPMCB(msg *message.GetCopyWPMMainProcessToRenderer) {
	if msg.Error {
		display.Error(msg.ErrorMessage)
		return
	}
	messenger.controller.processWPM(msg.Record.WPM)
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

				case *message.CheckCopyMainProcessToRenderer:
					messenger.checkCopyCB(msg)
				case *message.GetTextToCopyMainProcessToRenderer:
					messenger.getTextToCopyCB(msg)
				case *message.KeyMainProcessToRenderer:
					messenger.keyCB(msg)
				case *message.UpdateCopyWPMMainProcessToRenderer:
					messenger.updateCopyWPMCB(msg)
				case *message.GetCopyWPMMainProcessToRenderer:
					messenger.getCopyWPMCB(msg)
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
