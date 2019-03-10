package copypracticepanel

import (
	"github.com/pkg/errors"

	"github.com/josephbudd/cwt/domain/data/callids"
	"github.com/josephbudd/cwt/domain/interfaces/caller"
	"github.com/josephbudd/cwt/domain/types"
	"github.com/josephbudd/cwt/renderer/notjs"
	"github.com/josephbudd/cwt/renderer/viewtools"
)

/*

	Panel name: CopyPracticePanel

*/

// Caller communicates with the main process via an asynchrounous connection.
type Caller struct {
	panelGroup *PanelGroup
	presenter  *Presenter
	controler  *Controler
	quitCh     chan struct{} // send an empty struct to start the quit process.
	connection map[types.CallID]caller.Renderer
	tools      *viewtools.Tools // see /renderer/viewtools
	notJS      *notjs.NotJS

	/* NOTE TO DEVELOPER. Step 1 of 4.

	// 1: Declare your Caller members.

	*/

	// my new members

	state                   uint64
	getTextToCopyConnection caller.Renderer
	checkCopyConnection     caller.Renderer
	keyConnection           caller.Renderer
	updateCopyWPMConnection caller.Renderer
	getCopyWPMConnection    caller.Renderer
}

// addMainProcessCallBacks tells the main process what funcs to call back to.
func (panelCaller *Caller) addMainProcessCallBacks() (err error) {

	defer func() {
		if err != nil {
			err = errors.WithMessage(err, "(panelCaller *Caller) addMainProcessCallBacks()")
		}
	}()

	/* NOTE TO DEVELOPER. Step 2 of 4.

	// 2.1: Define each one of your Caller connection members as a conection to the main process.
	// 2.2: Tell the caller connection to the main processs to add a call back to each of your call back funcs.

	*/

	var found bool

	if panelCaller.checkCopyConnection, found = panelCaller.connection[callids.CheckCopyCallID]; !found {
		err = errors.New("unable to find panelCaller.connection[callids.CheckCopyCallID]")
		return
	}
	panelCaller.checkCopyConnection.AddCallBack(panelCaller.checkCopyCB)

	if panelCaller.getTextToCopyConnection, found = panelCaller.connection[callids.GetTextToCopyCallID]; !found {
		err = errors.New("unable to find panelCaller.connection[callids.GetTextToCopyCallID]")
		return
	}
	panelCaller.getTextToCopyConnection.AddCallBack(panelCaller.getTextToCopyCB)

	if panelCaller.keyConnection, found = panelCaller.connection[callids.KeyCallID]; !found {
		err = errors.New("unable to find panelCaller.connection[callids.KeyCallID]")
		return
	}
	panelCaller.keyConnection.AddCallBack(panelCaller.keyCB)

	// Define the update connection.
	if panelCaller.updateCopyWPMConnection, found = panelCaller.connection[callids.UpdateCopyWPMCallID]; !found {
		err = errors.New("unable to find panelCaller.connection[callids.UpdateCopyWPMCallID]")
		return
	}
	// Have the update connection call back to my call back handler.
	panelCaller.updateCopyWPMConnection.AddCallBack(panelCaller.updateCopyWPMCB)

	// Define the get connection.
	if panelCaller.getCopyWPMConnection, found = panelCaller.connection[callids.GetCopyWPMCallID]; !found {
		err = errors.New("unable to find panelCaller.connection[callids.GetCopyWPMCallID]")
		return
	}
	// Have the get connection call back to my call back handler.
	panelCaller.getCopyWPMConnection.AddCallBack(panelCaller.getCopyWPMCB)

	return
}

/* NOTE TO DEVELOPER. Step 3 of 4.

// 3.1: Define your funcs which call to the main process.
// 3.2: Define your funcs which the main process calls back to.

*/

// Check copy.

// checkCopy

func (panelCaller *Caller) checkCopy(copy []string, solution [][]*types.KeyCodeRecord, wpm uint64) {
	params := &types.RendererToMainProcessCheckCopyCallParams{
		Copy:     copy,
		Solution: solution,
		WPM:      wpm,
		State:    panelCaller.state,
	}
	panelCaller.checkCopyConnection.CallMainProcess(params)
}

func (panelCaller *Caller) checkCopyCB(params interface{}) {
	switch params := params.(type) {
	case *types.MainProcessToRendererCheckCopyCallParams:
		if params.State&panelCaller.state == panelCaller.state {
			if params.Error {
				panelCaller.tools.Error(params.ErrorMessage)
				return
			}
			// no errors
			panelCaller.presenter.showResult(params.CorrectCount, params.IncorrectCount, params.KeyedCount, params.TestResults)
		}
	}
}

// Get text to copy.

func (panelCaller *Caller) getTextToCopy() {
	params := &types.RendererToMainProcessGetTextToCopyCallParams{
		State: panelCaller.state,
	}
	panelCaller.getTextToCopyConnection.CallMainProcess(params)
}

func (panelCaller *Caller) getTextToCopyCB(params interface{}) {
	switch params := params.(type) {
	case *types.MainProcessToRendererGetTextToCopyCallParams:
		if params.State&panelCaller.state == panelCaller.state {
			if params.Error {
				panelCaller.tools.Error(params.ErrorMessage)
				return
			}
			// no errors
			panelCaller.controler.processTextToCopy(params.Solution)
		}
	}
}

// Key.

func (panelCaller *Caller) key(solution [][]*types.KeyCodeRecord, wpm, pause uint64) {
	params := &types.RendererToMainProcessKeyCallParams{
		Solution: solution,
		WPM:      wpm,
		Pause:    pause,
		State:    panelCaller.state,
		Run:      true,
	}
	panelCaller.keyConnection.CallMainProcess(params)
}

func (panelCaller *Caller) stopKeying() {
	params := &types.RendererToMainProcessKeyCallParams{
		State: panelCaller.state,
		Run:   false,
	}
	panelCaller.keyConnection.CallMainProcess(params)
}

func (panelCaller *Caller) keyCB(params interface{}) {
	switch params := params.(type) {
	case *types.MainProcessToRendererKeyCallParams:
		if params.State&panelCaller.state == panelCaller.state {
			if params.Error {
				panelCaller.tools.Error(params.ErrorMessage)
			}
			// no errors
			if params.Run {
				panelCaller.controler.processKeyFinished()
			} else {
				panelCaller.controler.processKeyStopped()
			}
		}
	}
}

// wpm

func (panelCaller *Caller) updateCopyWPMCB(params interface{}) {
	switch params := params.(type) {
	case *types.MainProcessToRendererUpdateCopyWPMCallParams:
		if params.Error {
			return
		}
		// no errors
		panelCaller.controler.processWPM(params.Record.WPM)
	}
}

func (panelCaller *Caller) getCopyWPMCB(params interface{}) {
	switch params := params.(type) {
	case *types.MainProcessToRendererGetCopyWPMCallParams:
		if params.Error {
			return
		}
		// no errors
		panelCaller.controler.processWPM(params.Record.WPM)
	}
}

// initialCalls makes the first calls to the main process.
func (panelCaller *Caller) initialCalls() {

	/* NOTE TO DEVELOPER. Step 4 of 4.

	//4: Make any initial calls to the main process that must be made when the app starts.

	*/

}
