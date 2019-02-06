package NumbersPanel

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/josephbudd/cwt/domain/data/callids"
	"github.com/josephbudd/cwt/domain/data/keyCodeTypes"
	"github.com/josephbudd/cwt/domain/interfaces/caller"
	"github.com/josephbudd/cwt/domain/types"
	"github.com/josephbudd/cwt/renderer/notjs"
	"github.com/josephbudd/cwt/renderer/viewtools"
)

/*

	Panel name: NumbersPanel

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

	// my added members
	state, testState        uint64
	getKeyCodesConnection   caller.Renderer
	updateKeyCodeConnection caller.Renderer
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

	if panelCaller.updateKeyCodeConnection, found = panelCaller.connection[callids.UpdateKeyCodeCallID]; !found {
		err = errors.New("unable to find panelCaller.connection[callids.UpdateKeyCodeCallID]")
		return
	}
	panelCaller.updateKeyCodeConnection.AddCallBack(panelCaller.updateKeyCodeCB)

	if panelCaller.getKeyCodesConnection, found = panelCaller.connection[callids.GetKeyCodesCallID]; !found {
		err = errors.New("unable to find panelCaller.connection[callids.GetKeyCodesCallID]")
		return
	}
	panelCaller.getKeyCodesConnection.AddCallBack(panelCaller.getKeyCodesCB)

	return
}

/* NOTE TO DEVELOPER. Step 3 of 4.

// 3.1: Define your funcs which call to the main process.
// 3.2: Define your funcs which the main process calls back to.

*/

// UpdateKeyCode

func (panelCaller *Caller) updateKeyCode(record *types.KeyCodeRecord) {
	params := &types.RendererToMainProcessUpdateKeyCodeCallParams{
		Record: record,
		State:  panelCaller.state,
	}
	panelCaller.updateKeyCodeConnection.CallMainProcess(params)
}

func (panelCaller *Caller) updateKeyCodeCB(params interface{}) {
	switch params := params.(type) {
	case *types.MainProcessToRendererUpdateKeyCodeCallParams:
		if params.State&panelCaller.state == panelCaller.state {
			if params.Error {
				panelCaller.tools.Error(params.ErrorMessage)
				return
			}
			// no errors
			if params.Record.Selected {
				panelCaller.tools.Success(fmt.Sprintf("The number %q is now selected.", params.Record.Name))
			} else {
				panelCaller.tools.Success(fmt.Sprintf("The number %q is no longer selected.", params.Record.Name))
			}
		}
	}
}

// GetKeyCodes

func (panelCaller *Caller) getKeyCodesCB(params interface{}) {
	switch params := params.(type) {
	case *types.MainProcessToRendererGetKeyCodesCallParams:
		if params.Error {
			// This panel will handle the GetKeyCode errors not the other panels.
			panelCaller.tools.Error(params.ErrorMessage)
			return
		}
		// no errors
		panelCaller.controler.setup(getNumberKeyCodes(params.Records))
	default:
		panelCaller.tools.Error("getKeyCodesCB: Unknown type.")
	}
}

func getNumberKeyCodes(records []*types.KeyCodeRecord) []*types.KeyCodeRecord {
	group := make([]*types.KeyCodeRecord, 0, len(records))
	for _, r := range records {
		if r.Type == keyCodeTypes.KeyCodeTypeNumber {
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
	return group
}

// initialCalls makes the first calls to the main process.
func (panelCaller *Caller) initialCalls() {

	/* NOTE TO DEVELOPER. Step 4 of 4.

	//4: Make any initial calls to the main process that must be made when the app starts.

	*/

}
