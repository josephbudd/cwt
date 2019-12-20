package message

import "github.com/josephbudd/cwt/domain/store/record"

// UpdateKeyCodeRendererToMainProcess is the UpdateKeyCode message that the renderer sends to the main process.
type UpdateKeyCodeRendererToMainProcess struct {
	Record *record.KeyCode
	State  uint64
}

// UpdateKeyCodeMainProcessToRenderer is the UpdateKeyCode message that the main process sends to the renderer.
type UpdateKeyCodeMainProcessToRenderer struct {
	Error        bool
	ErrorMessage string
	Fatal        bool

	Record *record.KeyCode
	State  uint64
}
