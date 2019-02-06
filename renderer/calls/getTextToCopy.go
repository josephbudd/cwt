package calls

import (
	"encoding/json"

	"github.com/josephbudd/cwt/domain/data/callids"
	"github.com/josephbudd/cwt/domain/implementations/calling"
	"github.com/josephbudd/cwt/domain/types"
)

// newGetTextToCopyCall is the constructor for the GetTextToCopy Call.
// It should only receive the repos that are needed. In this case the customer repo.
// Param keyCodeStorer storer.CopyCodeStorer is the keycode store.
// Param rendererSendPayload: is a kickasm generated renderer func that sends data to the main process.
func newGetTextToCopyCall(rendererSendPayload func(payload []byte) error) *calling.Renderer {
	return calling.NewRenderer(
		callids.GetTextToCopyCallID,
		rendererReceiveAndDispatchGetTextToCopy,
		rendererSendPayload,
	)
}

// rendererReceiveAndDispatchGetTextToCopy is a renderer func.
// It receives and dispatches the params sent by the main process.
// Param params is a []byte of a MainProcessToRendererGetTextToCopyCallParams
// Param dispatch is a func that dispatches params to the renderer call backs.
// This func is simple.
// 1. Unmarshall params into a *MainProcessToRendererGetTextToCopyCallParams.
// 2. Dispatch the *MainProcessToRendererGetTextToCopyCallParams.
func rendererReceiveAndDispatchGetTextToCopy(params []byte, dispatch func(interface{})) {
	// 1. Unmarshall params into a *MainProcessToRendererGetTextToCopyCallParams.
	rxparams := &types.MainProcessToRendererGetTextToCopyCallParams{}
	if err := json.Unmarshal(params, rxparams); err != nil {
		// This error will only happend during the development stage.
		// It means a conflict with the txparams in func mainProcessReceiveGetTextToCopy defined about.
		rxparams = &types.MainProcessToRendererGetTextToCopyCallParams{
			Error:        true,
			ErrorMessage: err.Error(),
		}
	}
	// 2. Dispatch the *MainProcessToRendererGetTextToCopyCallParams to the renderer panel caller that want to handle the GetTextToCopy call backs.
	dispatch(rxparams)
}
