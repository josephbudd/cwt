package copywpmpanel

import (
	"github.com/josephbudd/cwt/domain/lpc/message"
	"github.com/josephbudd/cwt/domain/store/record"
)

/*

	Panel name: CopyWPMPanel

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

// Set WPM.

func (caller *panelCaller) updateCopyWPM(r *record.WPM) {
	msg := &message.UpdateCopyWPMRendererToMainProcess{
		Record: r,
	}
	sendCh <- msg
}

func (caller *panelCaller) updateCopyWPMCB(msg *message.UpdateCopyWPMMainProcessToRenderer) {
	if msg.Error {
		tools.Error(msg.ErrorMessage)
		return
	}
	caller.controller.processWPM(msg.Record)
	tools.Success("Copy WPM Updated.")
}

// GetWPM

func (caller *panelCaller) getCopyWPMCB(msg *message.GetCopyWPMMainProcessToRenderer) {
	if msg.Error {
		tools.Error(msg.ErrorMessage)
		return
	}
	// no errors
	caller.controller.processWPM(msg.Record)
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

				case *message.UpdateCopyWPMMainProcessToRenderer:
					caller.updateCopyWPMCB(msg)
				case *message.GetCopyWPMMainProcessToRenderer:
					caller.getCopyWPMCB(msg)
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

	sendCh <- &message.GetCopyWPMRendererToMainProcess{}
}
