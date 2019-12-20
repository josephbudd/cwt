package message

import "github.com/josephbudd/cwt/domain/store/record"

// GetTextToCopyRendererToMainProcess is the GetTextToCopy message that the renderer sends to the main process.
type GetTextToCopyRendererToMainProcess struct {
	State uint64
}

// GetTextToCopyMainProcessToRenderer is the GetTextToCopy message that the main process sends to the renderer.
type GetTextToCopyMainProcessToRenderer struct {
	Error        bool
	ErrorMessage string
	Fatal        bool

	Solution [][]*record.KeyCode
	WPM      uint64
	State    uint64
}
