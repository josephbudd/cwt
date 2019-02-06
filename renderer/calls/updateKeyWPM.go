package calls

import (
	"encoding/json"

	"github.com/josephbudd/cwt/domain/data/callids"
	"github.com/josephbudd/cwt/domain/implementations/calling"
	"github.com/josephbudd/cwt/domain/types"
)

// newUpdateKeyWPMCall is the constructor for the UpdateKeyWPM Call.
// It should only receive the repos that are needed. In this case the keyCode repo.
// Param keyCodeStorer storer.KeyWPMStorer is the keyCode repo needed to get a keyCode record.
// Param rendererSendPayload: is a kickasm generated renderer func that sends data to the main process.
func newUpdateKeyWPMCall(rendererSendPayload func(payload []byte) error) *calling.Renderer {
	return calling.NewRenderer(
		callids.UpdateKeyWPMCallID,
		rendererReceiveAndDispatchUpdateKeyWPM,
		rendererSendPayload,
	)
}

// rendererReceiveAndDispatchUpdateKeyWPM is a renderer func.
// It receives and dispatches the params sent by the main process.
// Param params is a []byte of a MainProcessToRendererUpdateKeyWPMCallParams
// Param dispatch is a func that dispatches params to the renderer call backs.
// This func is simple.
// 1. Unmarshall params into a *MainProcessToRendererUpdateKeyWPMCallParams.
// 2. Dispatch the *MainProcessToRendererUpdateKeyWPMCallParams.
func rendererReceiveAndDispatchUpdateKeyWPM(params []byte, dispatch func(interface{})) {
	// 1. Unmarshall params into a *MainProcessToRendererUpdateKeyWPMCallParams.
	rxparams := &types.MainProcessToRendererUpdateKeyWPMCallParams{}
	if err := json.Unmarshal(params, rxparams); err != nil {
		// This error will only happend during the development stage.
		// It means a conflict with the txparams in func mainProcessReceiveUpdateKeyWPM defined about.
		rxparams = &types.MainProcessToRendererUpdateKeyWPMCallParams{
			Error:        true,
			ErrorMessage: err.Error(),
		}
	}
	// 2. Dispatch the *MainProcessToRendererUpdateKeyWPMCallParams to the renderer panel caller that want to handle the UpdateKeyWPM call backs.
	dispatch(rxparams)
}
