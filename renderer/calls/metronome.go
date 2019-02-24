package calls

import (
	"encoding/json"

	"github.com/josephbudd/cwt/domain/data/callids"
	"github.com/josephbudd/cwt/domain/implementations/calling"
	"github.com/josephbudd/cwt/domain/types"
)

// newMetronomeCall is the constructor for the Metronome Call.
// It should only receive the repos that are needed. In this case the keyCode repo.
// Param keyCodeStorer storer.CopyWPMStorer is the keyCode repo needed to get a keyCode record.
// Param rendererSendPayload: is a kickasm generated renderer func that sends data to the main process.
func newMetronomeCall(rendererSendPayload func(payload []byte) error) *calling.Renderer {
	return calling.NewRenderer(
		callids.MetronomeCallID,
		rendererReceiveAndDispatchMetronome,
		rendererSendPayload,
	)
}

// rendererReceiveAndDispatchMetronome is a renderer func.
// It receives and dispatches the params sent by the main process.
// Param params is a []byte of a MainProcessToRendererMetronomeCallParams
// Param dispatch is a func that dispatches params to the renderer call backs.
// This func is simple.
// 1. Unmarshall params into a *MainProcessToRendererMetronomeCallParams.
// 2. Dispatch the *MainProcessToRendererMetronomeCallParams.
func rendererReceiveAndDispatchMetronome(params []byte, dispatch func(interface{})) {
	// 1. Unmarshall params into a *MainProcessToRendererMetronomeCallParams.
	rxparams := &types.MainProcessToRendererMetronomeCallParams{}
	if err := json.Unmarshal(params, rxparams); err != nil {
		// This error will only happend during the development stage.
		// It means a conflict with the txparams in func mainProcessReceiveMetronome defined about.
		rxparams = &types.MainProcessToRendererMetronomeCallParams{
			Error:        true,
			ErrorMessage: err.Error(),
		}
	}
	// 2. Dispatch the *MainProcessToRendererMetronomeCallParams to the renderer panel caller that want to handle the Metronome call backs.
	dispatch(rxparams)
}
