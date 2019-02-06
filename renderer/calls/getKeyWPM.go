package calls

import (
	"encoding/json"

	"github.com/josephbudd/cwt/domain/data/callids"
	"github.com/josephbudd/cwt/domain/implementations/calling"
	"github.com/josephbudd/cwt/domain/types"
)

// newGetKeyWPMCall is the constructor for the GetKeyWPM Call.
// It should only receive the repos that are needed. In this case the keyCode repo.
// Param keyCodeStorer storer.KeyCodeStorer is the keyCode repo needed to get a keyCode record.
// Param rendererSendPayload: is a kickasm generated renderer func that sends data to the main process.
func newGetKeyWPMCall(rendererSendPayload func(payload []byte) error) *calling.Renderer {
	return calling.NewRenderer(
		callids.GetKeyWPMCallID,
		rendererReceiveAndDispatchGetKeyWPM,
		rendererSendPayload,
	)
}

// rendererReceiveAndDispatchGetKeyWPM is a renderer func.
// It receives and dispatches the params sent by the main process.
// Param params is a []byte of a MainProcessToRendererGetKeyWPMCallParams
// Param dispatch is a func that dispatches params to the renderer call backs.
// This func is simple.
// 1. Unmarshall params into a *MainProcessToRendererGetKeyWPMCallParams.
// 2. Dispatch the *MainProcessToRendererGetKeyWPMCallParams.
func rendererReceiveAndDispatchGetKeyWPM(params []byte, dispatch func(interface{})) {
	// 1. Unmarshall params into a *MainProcessToRendererGetKeyWPMCallParams.
	rxparams := &types.MainProcessToRendererGetKeyWPMCallParams{}
	if err := json.Unmarshal(params, rxparams); err != nil {
		// This error will only happend during the development stage.
		// It means a conflict with the txparams in func mainProcessReceiveGetKeyWPM defined about.
		rxparams = &types.MainProcessToRendererGetKeyWPMCallParams{
			Error:        true,
			ErrorMessage: err.Error(),
		}
	}
	// 2. Dispatch the *MainProcessToRendererGetKeyWPMCallParams to the renderer panel caller that want to handle the GetKeyWPM call backs.
	dispatch(rxparams)
}
