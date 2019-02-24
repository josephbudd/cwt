package calls

import (
	"github.com/josephbudd/cwt/domain/data/callids"
	"github.com/josephbudd/cwt/domain/interfaces/caller"
	"github.com/josephbudd/cwt/domain/types"
)

// TODO: Add your calls.
// Example:
//      callids.AddCustomerCallID: newAddCustomerCall(rendererSendPayload, customerStorer)

// GetCallMap returns a render call map.
func GetCallMap(rendererSendPayload func(payload []byte) error) map[types.CallID]caller.Renderer {
	return map[types.CallID]caller.Renderer{
		callids.LogCallID: newLogCall(rendererSendPayload),

		callids.GetKeyCodesCallID:   newGetKeyCodesCall(rendererSendPayload),
		callids.UpdateKeyCodeCallID: newUpdateKeyCodeCall(rendererSendPayload),

		callids.GetKeyWPMCallID:    newGetKeyWPMCall(rendererSendPayload),
		callids.UpdateKeyWPMCallID: newUpdateKeyWPMCall(rendererSendPayload),

		callids.GetCopyWPMCallID:    newGetCopyWPMCall(rendererSendPayload),
		callids.UpdateCopyWPMCallID: newUpdateCopyWPMCall(rendererSendPayload),

		callids.GetTextToCopyCallID: newGetTextToCopyCall(rendererSendPayload),
		callids.KeyCallID:           newKeyCall(rendererSendPayload),
		callids.CheckCopyCallID:     newCheckCopyCall(rendererSendPayload),

		callids.GetTextWPMToKeyCallID: newGetTextWPMToKeyCall(rendererSendPayload),
		callids.CheckKeyCallID:        newCheckKeyCall(rendererSendPayload),

		callids.MetronomeCallID: newMetronomeCall(rendererSendPayload),
	}
}
