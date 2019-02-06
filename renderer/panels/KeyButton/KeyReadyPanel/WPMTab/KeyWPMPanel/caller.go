package KeyWPMPanel

import (
	"github.com/pkg/errors"

	"github.com/josephbudd/cwt/domain/data/callids"
	"github.com/josephbudd/cwt/domain/interfaces/caller"
	"github.com/josephbudd/cwt/domain/types"
	"github.com/josephbudd/cwt/renderer/notjs"
	"github.com/josephbudd/cwt/renderer/viewtools"
)

/*

	Panel name: KeyWPMPanel

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

	// example:

	addCustomerConnection caller.Renderer

	*/

	updateKeyWPMConnection caller.Renderer
	getKeyWPMConnection    caller.Renderer
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

	// Define the update connection.
	if panelCaller.updateKeyWPMConnection, found = panelCaller.connection[callids.UpdateKeyWPMCallID]; !found {
		err = errors.New("unable to find panelCaller.connection[callids.UpdateKeyWPMCallID]")
		return
	}
	// Have the update connection call back to my call back handler.
	panelCaller.updateKeyWPMConnection.AddCallBack(panelCaller.updateKeyWPMCB)

	// Define the get connection.
	if panelCaller.getKeyWPMConnection, found = panelCaller.connection[callids.GetKeyWPMCallID]; !found {
		err = errors.New("unable to find panelCaller.connection[callids.GetKeyWPMCallID]")
		return
	}
	// Have the get connection call back to my call back handler.
	panelCaller.getKeyWPMConnection.AddCallBack(panelCaller.getKeyWPMCB)

	return
}

/* NOTE TO DEVELOPER. Step 3 of 4.

// 3.1: Define your funcs which call to the main process.
// 3.2: Define your funcs which the main process calls back to.


*/

// Set WPM.

func (panelCaller *Caller) updateKeyWPM(record *types.WPMRecord) {
	params := &types.RendererToMainProcessUpdateKeyWPMCallParams{
		Record: record,
	}
	panelCaller.updateKeyWPMConnection.CallMainProcess(params)
}

func (panelCaller *Caller) updateKeyWPMCB(params interface{}) {
	switch params := params.(type) {
	case *types.MainProcessToRendererUpdateKeyWPMCallParams:
		if params.Error {
			panelCaller.tools.Error(params.ErrorMessage)
			return
		}
		// no errors
		panelCaller.controler.processWPM(params.Record)
		panelCaller.tools.Success("Key WPM Updated.")
	}
}

// GetWPM

func (panelCaller *Caller) getKeyWPMCB(params interface{}) {
	switch params := params.(type) {
	case *types.MainProcessToRendererGetKeyWPMCallParams:
		if params.Error {
			panelCaller.tools.Error(params.ErrorMessage)
			return
		}
		// no errors
		panelCaller.controler.processWPM(params.Record)
	}
}

// initialCalls makes the first calls to the main process.
func (panelCaller *Caller) initialCalls() {

	/* NOTE TO DEVELOPER. Step 4 of 4.

	//4: Make any initial calls to the main process that must be made when the app starts.

	*/

	panelCaller.getKeyWPMConnection.CallMainProcess(nil)
}
