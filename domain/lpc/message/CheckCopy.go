package message

import (
	"github.com/josephbudd/cwt/domain/data"
	"github.com/josephbudd/cwt/domain/store/record"
)

// CheckCopyRendererToMainProcess is the CheckCopy message that the renderer sends to the main process.
type CheckCopyRendererToMainProcess struct {
	Copy         []string
	Solution     [][]*record.KeyCode
	WPM          uint64
	StoreResults bool
	State        uint64
}

// CheckCopyMainProcessToRenderer is the CheckCopy message that the main process sends to the renderer.
type CheckCopyMainProcessToRenderer struct {
	Error        bool
	ErrorMessage string

	CorrectCount   uint64
	IncorrectCount uint64
	KeyedCount     uint64
	TestResults    [][]data.TestResult
	State          uint64
}
