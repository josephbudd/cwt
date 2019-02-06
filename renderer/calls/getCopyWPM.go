package calls

import (
	"encoding/json"

	"github.com/josephbudd/cwt/domain/data/callids"
	"github.com/josephbudd/cwt/domain/implementations/calling"
	"github.com/josephbudd/cwt/domain/types"
)

// newGetCopyWPMCall is the constructor for the GetCopyWPM Call.
// It should only receive the repos that are needed. In this case the keyCode repo.
// Param keyCodeStorer storer.CopyCodeStorer is the keyCode repo needed to get a keyCode record.
// Param rendererSendPayload: is a kickasm generated renderer func that sends data to the main process.
func newGetCopyWPMCall(rendererSendPayload func(payload []byte) error) *calling.Renderer {
	return calling.NewRenderer(
		callids.GetCopyWPMCallID,
		rendererReceiveAndDispatchGetCopyWPM,
		rendererSendPayload,
	)
}

// rendererReceiveAndDispatchGetCopyWPM is a renderer func.
// It receives and dispatches the params sent by the main process.
// Param params is a []byte of a MainProcessToRendererGetCopyWPMCallParams
// Param dispatch is a func that dispatches params to the renderer call backs.
// This func is simple.
// 1. Unmarshall params into a *MainProcessToRendererGetCopyWPMCallParams.
// 2. Dispatch the *MainProcessToRendererGetCopyWPMCallParams.
func rendererReceiveAndDispatchGetCopyWPM(params []byte, dispatch func(interface{})) {
	// 1. Unmarshall params into a *MainProcessToRendererGetCopyWPMCallParams.
	rxparams := &types.MainProcessToRendererGetCopyWPMCallParams{}
	if err := json.Unmarshal(params, rxparams); err != nil {
		// This error will only happend during the development stage.
		// It means a conflict with the txparams in func mainProcessReceiveGetCopyWPM defined about.
		rxparams = &types.MainProcessToRendererGetCopyWPMCallParams{
			Error:        true,
			ErrorMessage: err.Error(),
		}
	}
	// 2. Dispatch the *MainProcessToRendererGetCopyWPMCallParams to the renderer panel caller that want to handle the GetCopyWPM call backs.
	dispatch(rxparams)
}
