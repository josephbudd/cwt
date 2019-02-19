package types

// RendererToMainProcessCheckCopyCallParams is the CheckCopy function parameters that the renderer sends to the main process.
type RendererToMainProcessCheckCopyCallParams struct {
	Copy         []string
	Solution     [][]*KeyCodeRecord
	WPM          uint64
	StoreResults bool
	State        uint64
}

// MainProcessToRendererCheckCopyCallParams is the CheckCopy function parameters that the main process sends to the renderer.
type MainProcessToRendererCheckCopyCallParams struct {
	Error          bool
	ErrorMessage   string
	CorrectCount   uint64
	IncorrectCount uint64
	KeyedCount     uint64
	TestResults    [][]TestResult
	State          uint64
}
