package calls

import (
	"encoding/json"

	"github.com/josephbudd/cwt/domain/data/callids"
	"github.com/josephbudd/cwt/domain/implementations/calling"
	"github.com/josephbudd/cwt/domain/types"
)

// newCheckCopyCall is the constructor for the CheckCopy Call.
// Param rendererSendPayload: is a kickasm generated renderer func that sends data to the main process.
func newCheckCopyCall(rendererSendPayload func(payload []byte) error) *calling.Renderer {
	return calling.NewRenderer(
		callids.CheckCopyCallID,
		rendererReceiveAndDispatchCheckCopy,
		rendererSendPayload,
	)
}

// rendererReceiveAndDispatchCheckCopy is a renderer func.
// It receives and dispatches the params sent by the main process.
// Param params is a []byte of a MainProcessToRendererCheckCopyCallParams
// Param dispatch is a func that dispatches params to the renderer call backs.
// This func is simple.
// 1. Unmarshall params into a *MainProcessToRendererCheckCopyCallParams.
// 2. Dispatch the *MainProcessToRendererCheckCopyCallParams.
func rendererReceiveAndDispatchCheckCopy(params []byte, dispatch func(interface{})) {
	// 1. Unmarshall params into a *MainProcessToRendererCheckCopyCallParams.
	rxparams := &types.MainProcessToRendererCheckCopyCallParams{}
	if err := json.Unmarshal(params, rxparams); err != nil {
		// This error will only happend during the development stage.
		// It means a conflict with the txparams in func mainProcessReceiveCheckCopy defined about.
		rxparams = &types.MainProcessToRendererCheckCopyCallParams{
			Error:        true,
			ErrorMessage: err.Error(),
		}
	}
	// 2. Dispatch the *MainProcessToRendererCheckCopyCallParams to the renderer panel caller that want to handle the CheckCopy call backs.
	dispatch(rxparams)
}
