package message

import "github.com/josephbudd/cwt/domain/store/record"

// KeyRendererToMainProcess is the Key message that the renderer sends to the main process.
type KeyRendererToMainProcess struct {
	Solution [][]*record.KeyCode
	WPM      uint64
	Pause    uint64
	Run      bool
	State    uint64
}

// KeyMainProcessToRenderer is the Key message that the main process sends to the renderer.
type KeyMainProcessToRenderer struct {
	Error        bool
	ErrorMessage string
	Fatal        bool

	Run   bool
	State uint64
}
