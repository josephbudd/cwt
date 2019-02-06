package calls

import (
	"encoding/json"

	"github.com/josephbudd/cwt/domain/data/callids"
	"github.com/josephbudd/cwt/domain/implementations/calling"
	"github.com/josephbudd/cwt/domain/types"
)

// newKeyCall is the constructor for the Key Call.
// Param rendererSendPayload: is a kickasm generated renderer func that sends data to the main process.
func newKeyCall(rendererSendPayload func(payload []byte) error) *calling.Renderer {
	return calling.NewRenderer(
		callids.KeyCallID,
		rendererReceiveAndDispatchKey,
		rendererSendPayload,
	)
}

// rendererReceiveAndDispatchKey is a renderer func.
// It receives and dispatches the params sent by the main process.
// Param params is a []byte of a MainProcessToRendererKeyCallParams
// Param dispatch is a func that dispatches params to the renderer call backs.
// This func is simple.
// 1. Unmarshall params into a *MainProcessToRendererKeyCallParams.
// 2. Dispatch the *MainProcessToRendererKeyCallParams.
func rendererReceiveAndDispatchKey(params []byte, dispatch func(interface{})) {
	// 1. Unmarshall params into a *MainProcessToRendererKeyCallParams.
	rxparams := &types.MainProcessToRendererKeyCallParams{}
	if err := json.Unmarshal(params, rxparams); err != nil {
		// This error will only happend during the development stage.
		// It means a conflict with the txparams in func mainProcessReceiveKey defined about.
		rxparams = &types.MainProcessToRendererKeyCallParams{
			Error:        true,
			ErrorMessage: err.Error(),
		}
	}
	// 2. Dispatch the *MainProcessToRendererKeyCallParams to the renderer panel caller that want to handle the Key call backs.
	dispatch(rxparams)
}
