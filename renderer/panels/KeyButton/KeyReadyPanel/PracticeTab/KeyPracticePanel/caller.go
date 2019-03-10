package keypracticepanel

import (
	"github.com/pkg/errors"

	"github.com/josephbudd/cwt/domain/data/callids"
	"github.com/josephbudd/cwt/domain/interfaces/caller"
	"github.com/josephbudd/cwt/domain/types"
	"github.com/josephbudd/cwt/renderer/notjs"
	"github.com/josephbudd/cwt/renderer/viewtools"
)

/*

	Panel name: KeyPracticePanel

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

	state                     uint64
	getTextWPMToKeyConnection caller.Renderer
	checkKeyConnection        caller.Renderer
	metronomeConnection       caller.Renderer
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

	// Get Text.
	// Define the connection.
	if panelCaller.getTextWPMToKeyConnection, found = panelCaller.connection[callids.GetTextWPMToKeyCallID]; !found {
		err = errors.New("unable to find panelCaller.connection[callids.GetTextWPMToKeyCallID]")
		return
	}
	// Have the connection call back to my call back handler.
	panelCaller.getTextWPMToKeyConnection.AddCallBack(panelCaller.getTextWPMToKeyCB)

	// check key
	// Define the connection.
	if panelCaller.checkKeyConnection, found = panelCaller.connection[callids.CheckKeyCallID]; !found {
		err = errors.New("unable to find panelCaller.connection[callids.CheckKeyCallID]")
		return
	}
	// Have the connection call back to my call back handler.
	panelCaller.checkKeyConnection.AddCallBack(panelCaller.checkKeyCB)

	// metronome
	if panelCaller.metronomeConnection, found = panelCaller.connection[callids.MetronomeCallID]; !found {
		err = errors.New("unable to find panelCaller.connection[callids.MetronomeCallID]")
		return
	}
	// Have the connection call back to my call back handler.
	panelCaller.metronomeConnection.AddCallBack(panelCaller.metronomeCB)

	return
}

/* NOTE TO DEVELOPER. Step 3 of 4.

// 3.1: Define your funcs which call to the main process.
// 3.2: Define your funcs which the main process calls back to.

*/

// Caller implements widgets.UserKeyChecker with funcs GetKeyCodesWPM and CheckUserKey.
// Caller implements widgets.Metronomer with funcs StartMetronome and StopMetronome.

// Get text to key.

// GetKeyCodesWPM

// GetKeyCodesWPM gets new text for the user to key and the wpm required for keying.
func (panelCaller *Caller) GetKeyCodesWPM() {
	params := &types.RendererToMainProcessGetTextWPMToKeyCallParams{
		State: panelCaller.state,
	}
	panelCaller.getTextWPMToKeyConnection.CallMainProcess(params)
}

func (panelCaller *Caller) getTextWPMToKeyCB(params interface{}) {
	switch params := params.(type) {
	case *types.MainProcessToRendererGetTextWPMToKeyCallParams:
		if params.State&panelCaller.state == panelCaller.state {
			if params.Error {
				panelCaller.tools.Error(params.ErrorMessage)
				return
			}
			// no errors
			panelCaller.controler.keyWidget.SetKeyCodesWPM(params.Solution, params.WPM)
		}
	}
}

// check key

// CheckUserKey checks the user's keying ability.
func (panelCaller *Caller) CheckUserKey(milliSeconds []int64, solution [][]*types.KeyCodeRecord, wpm uint64) {
	params := &types.RendererToMainProcessCheckKeyCallParams{
		MilliSeconds: milliSeconds,
		Solution:     solution,
		WPM:          wpm,
		State:        panelCaller.state,
	}
	panelCaller.checkKeyConnection.CallMainProcess(params)
}

func (panelCaller *Caller) checkKeyCB(params interface{}) {
	switch params := params.(type) {
	case *types.MainProcessToRendererCheckKeyCallParams:
		if params.State&panelCaller.state == panelCaller.state {
			if params.Error {
				panelCaller.tools.Error(params.ErrorMessage)
				return
			}
			// no errors
			panelCaller.controler.keyWidget.ShowResults(params.CorrectCount, params.IncorrectCount, params.KeyedCount, params.TestResults)
		}
	}
}

// metronome

// StartMetronome starts the metronome.
func (panelCaller *Caller) StartMetronome(wpm uint64) {
	params := &types.RendererToMainProcessMetronomeCallParams{
		Run:   true,
		State: panelCaller.state,
		WPM:   wpm,
	}
	panelCaller.metronomeConnection.CallMainProcess(params)
}

// StopMetronome stops the metronome.
func (panelCaller *Caller) StopMetronome() {
	params := &types.RendererToMainProcessMetronomeCallParams{
		Run:   false,
		State: panelCaller.state,
	}
	panelCaller.metronomeConnection.CallMainProcess(params)
}

func (panelCaller *Caller) metronomeCB(params interface{}) {
	switch params := params.(type) {
	case *types.MainProcessToRendererMetronomeCallParams:
		if params.State&panelCaller.state == panelCaller.state {
			if params.Error {
				panelCaller.tools.Error(params.ErrorMessage)
				return
			}
		}
	}
}

// initialCalls makes the first calls to the main process.
func (panelCaller *Caller) initialCalls() {

	/* NOTE TO DEVELOPER. Step 4 of 4.

	//4: Make any initial calls to the main process that must be made when the app starts.

	*/

}
