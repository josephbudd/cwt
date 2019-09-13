package message

import "github.com/josephbudd/cwt/domain/store/record"

// GetKeyCodesRendererToMainProcess is the GetKeyCodes message that the renderer sends to the main process.
type GetKeyCodesRendererToMainProcess struct {
	State uint64
}

// GetKeyCodesMainProcessToRenderer is the GetKeyCodes message that the main process sends to the renderer.
type GetKeyCodesMainProcessToRenderer struct {
	Error        bool
	ErrorMessage string

	Records []*record.KeyCode
	State   uint64
}
