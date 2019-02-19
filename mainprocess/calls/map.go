package calls

import (
	"github.com/josephbudd/cwt/domain/data/callids"
	"github.com/josephbudd/cwt/domain/interfaces/caller"
	"github.com/josephbudd/cwt/domain/interfaces/storer"
	"github.com/josephbudd/cwt/domain/types"
)

// TODO: Add your calls.
// Example:
//      callids.AddConactCallID: newAddContactCall(contactStorer)

// GetCallMap returns a map of each mainprocess call.
func GetCallMap(wPMStore storer.WPMStorer, keyCodeStore storer.KeyCodeStorer) map[types.CallID]caller.MainProcesser {
	return map[types.CallID]caller.MainProcesser{
		callids.LogCallID: newLogCall(),

		callids.GetKeyCodesCallID:   newGetKeyCodesCall(keyCodeStore),
		callids.UpdateKeyCodeCallID: newUpdateKeyCodeCall(keyCodeStore),

		callids.GetKeyWPMCallID:    newGetKeyWPMCall(wPMStore),
		callids.UpdateKeyWPMCallID: newUpdateKeyWPMCall(wPMStore),

		callids.GetCopyWPMCallID:    newGetCopyWPMCall(wPMStore),
		callids.UpdateCopyWPMCallID: newUpdateCopyWPMCall(wPMStore),

		callids.GetTextToCopyCallID: newGetTextToCopyCall(keyCodeStore, wPMStore),
		callids.KeyCallID:           newKeyCall(),
		callids.CheckCopyCallID:     newCheckCopyCall(keyCodeStore),

		callids.GetTextWPMToKeyCallID: newGetTextWPMToKeyCall(keyCodeStore, wPMStore),
		callids.CheckKeyCallID:        newCheckKeyCall(keyCodeStore),
	}
}
