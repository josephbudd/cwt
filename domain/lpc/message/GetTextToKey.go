package message

import (
	"github.com/josephbudd/cwt/domain/data"
	"github.com/josephbudd/cwt/domain/store/record"
)

// GetTextToKeyRendererToMainProcess is the GetTextToKey message that the renderer sends to the main process.
type GetTextToKeyRendererToMainProcess struct {
	State    uint64
	Practice bool
}

// GetTextToKeyMainProcessToRenderer is the GetTextToKey message that the main process sends to the renderer.
type GetTextToKeyMainProcessToRenderer struct {
	Error        bool
	ErrorMessage string
	Fatal        bool

	Solution [][]*record.KeyCode
	Help     [][]data.HowTo
	Practice bool
	WPM      uint64
	State    uint64
}
