package calls

import (
	"encoding/json"

	"github.com/josephbudd/cwt/domain/data/callids"
	"github.com/josephbudd/cwt/domain/implementations/calling"
	"github.com/josephbudd/cwt/domain/types"
)

// newCheckKeyCall is the constructor for the CheckKey Call.
// Param rendererSendPayload: is a kickasm generated renderer func that sends data to the main process.
// Param keyCodeStorer is the key code storer.
func newCheckKeyCall(rendererSendPayload func(payload []byte) error) *calling.Renderer {
	return calling.NewRenderer(
		callids.CheckKeyCallID,
		rendererReceiveAndDispatchCheckKey,
		rendererSendPayload,
	)
}

// rendererReceiveAndDispatchCheckKey is a renderer func.
// It receives and dispatches the params sent by the main process.
// Param params is a []byte of a MainProcessToRendererCheckKeyCallParams
// Param dispatch is a func that dispatches params to the renderer call backs.
// This func is simple.
// 1. Unmarshall params into a *MainProcessToRendererCheckKeyCallParams.
// 2. Dispatch the *MainProcessToRendererCheckKeyCallParams.
func rendererReceiveAndDispatchCheckKey(params []byte, dispatch func(interface{})) {
	// 1. Unmarshall params into a *MainProcessToRendererCheckKeyCallParams.
	rxparams := &types.MainProcessToRendererCheckKeyCallParams{}
	if err := json.Unmarshal(params, rxparams); err != nil {
		// This error will only happend during the development stage.
		// It means a conflict with the txparams in func mainProcessReceiveCheckKey defined about.
		rxparams = &types.MainProcessToRendererCheckKeyCallParams{
			Error:        true,
			ErrorMessage: err.Error(),
		}
	}
	// 2. Dispatch the *MainProcessToRendererCheckKeyCallParams to the renderer panel caller that want to handle the CheckKey call backs.
	dispatch(rxparams)
}
