package types

// RendererToMainProcessCheckKeyCallParams is the CheckKey function parameters that the renderer sends to the main process.
type RendererToMainProcessCheckKeyCallParams struct {
	Solution     [][]*KeyCodeRecord
	MilliSeconds []int64
	WPM          uint64
	StoreResults bool
	State        uint64
}

// MainProcessToRendererCheckKeyCallParams is the CheckKey function parameters that the main process sends to the renderer.
type MainProcessToRendererCheckKeyCallParams struct {
	Error          bool
	ErrorMessage   string
	CorrectCount   uint64
	IncorrectCount uint64
	KeyedCount     uint64
	TestResults     [][]TestResult
	State          uint64
}
