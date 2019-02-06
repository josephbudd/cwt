package calls

import (
	"encoding/json"

	"github.com/josephbudd/cwt/domain/data/callids"
	"github.com/josephbudd/cwt/domain/implementations/calling"
	"github.com/josephbudd/cwt/domain/types"
)

// newUpdateKeyCodeCall is the constructor for the UpdateKeyCode Call.
// It should only receive the repos that are needed. In this case the keyCode repo.
// Param keyCodeStorer storer.KeyCodeStorer is the keyCode repo needed to get a keyCode record.
// Param rendererSendPayload: is a kickasm generated renderer func that sends data to the main process.
func newUpdateKeyCodeCall(rendererSendPayload func(payload []byte) error) *calling.Renderer {
	return calling.NewRenderer(
		callids.UpdateKeyCodeCallID,
		rendererReceiveAndDispatchUpdateKeyCode,
		rendererSendPayload,
	)
}

// rendererReceiveAndDispatchUpdateKeyCode is a renderer func.
// It receives and dispatches the params sent by the main process.
// Param params is a []byte of a MainProcessToRendererUpdateKeyCodeCallParams
// Param dispatch is a func that dispatches params to the renderer call backs.
// This func is simple.
// 1. Unmarshall params into a *MainProcessToRendererUpdateKeyCodeCallParams.
// 2. Dispatch the *MainProcessToRendererUpdateKeyCodeCallParams.
func rendererReceiveAndDispatchUpdateKeyCode(params []byte, dispatch func(interface{})) {
	// 1. Unmarshall params into a *MainProcessToRendererUpdateKeyCodeCallParams.
	rxparams := &types.MainProcessToRendererUpdateKeyCodeCallParams{}
	if err := json.Unmarshal(params, rxparams); err != nil {
		// This error will only happend during the development stage.
		// It means a conflict with the txparams in func mainProcessReceiveUpdateKeyCode defined about.
		rxparams = &types.MainProcessToRendererUpdateKeyCodeCallParams{
			Error:        true,
			ErrorMessage: err.Error(),
		}
	}
	// 2. Dispatch the *MainProcessToRendererUpdateKeyCodeCallParams to the renderer panel caller that want to handle the UpdateKeyCode call backs.
	dispatch(rxparams)
}
