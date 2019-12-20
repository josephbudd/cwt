package message

import (
	"github.com/josephbudd/cwt/domain/data"
	"github.com/josephbudd/cwt/domain/store/record"
)

// CheckKeyRendererToMainProcess is the CheckKey message that the renderer sends to the main process.
type CheckKeyRendererToMainProcess struct {
	Solution     [][]*record.KeyCode
	MilliSeconds []int64
	WPM          uint64
	StoreResults bool
	State        uint64
}

// CheckKeyMainProcessToRenderer is the CheckKey message that the main process sends to the renderer.
type CheckKeyMainProcessToRenderer struct {
	Error        bool
	ErrorMessage string
	Fatal        bool

	CorrectCount   uint64
	IncorrectCount uint64
	MaxCount       uint64
	TestResults    [][]data.TestResult
	State          uint64
}
