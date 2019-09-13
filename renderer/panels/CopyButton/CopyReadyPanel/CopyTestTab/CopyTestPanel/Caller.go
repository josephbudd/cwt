package copytestpanel

/*

	Panel name: CopyTestPanel

*/

import (
	"github.com/josephbudd/cwt/domain/lpc/message"
	"github.com/josephbudd/cwt/domain/store/record"
)

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

// Check copy.

// checkCopy

func (caller *panelCaller) checkCopy(copy []string, solution [][]*record.KeyCode, wpm uint64) {
	msg := message.CheckCopyRendererToMainProcess{
		Copy:         copy,
		Solution:     solution,
		WPM:          wpm,
		State:        state,
		StoreResults: true,
	}
	sendCh <- msg
}

func (caller *panelCaller) checkCopyCB(msg *message.CheckCopyMainProcessToRenderer) {
	if msg.State&state != state {
		return
	}
	if msg.Error {
		tools.Error(msg.ErrorMessage)
		return
	}
	caller.presenter.showResult(msg.CorrectCount, msg.IncorrectCount, msg.KeyedCount, msg.TestResults)
}

// Get text to copy.

func (caller *panelCaller) getTextToCopy() {
	msg := &message.GetTextToCopyRendererToMainProcess{
		State: state,
	}
	sendCh <- msg
}

func (caller *panelCaller) getTextToCopyCB(msg *message.GetTextToCopyMainProcessToRenderer) {
	if msg.State&state != state {
		return
	}
	if msg.Error {
		tools.Error(msg.ErrorMessage)
		return
	}
	caller.controller.processTextToCopy(msg.Solution)
}

// Key.

func (caller *panelCaller) key(solution [][]*record.KeyCode, wpm, pause uint64) {
	msg := &message.KeyRendererToMainProcess{
		Solution: solution,
		WPM:      wpm,
		Pause:    pause,
		State:    state,
		Run:      true,
	}
	sendCh <- msg
}

func (caller *panelCaller) keyCB(msg *message.KeyMainProcessToRenderer) {
	if msg.State&state != state {
		return
	}
	if msg.Error {
		tools.Error(msg.ErrorMessage)
	}
	caller.controller.processKeyFinished()
}

// WPM

func (caller *panelCaller) updateCopyWPMCB(msg *message.UpdateCopyWPMMainProcessToRenderer) {
	if msg.Error {
		return
	}
	caller.controller.processWPM(msg.Record.WPM)
}

func (caller *panelCaller) getCopyWPMCB(msg *message.GetCopyWPMMainProcessToRenderer) {
	if msg.Error {
		return
	}
	caller.controller.processWPM(msg.Record.WPM)
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

				case *message.CheckCopyMainProcessToRenderer:
					caller.checkCopyCB(msg)
				case *message.GetTextToCopyMainProcessToRenderer:
					caller.getTextToCopyCB(msg)
				case *message.KeyMainProcessToRenderer:
					caller.keyCB(msg)
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
}
