package types

// RendererToMainProcessKeyCallParams is the Key function parameters that the renderer sends to the main process.
type RendererToMainProcessKeyCallParams struct {
	Solution [][]*KeyCodeRecord
	WPM      uint64
	Pause    uint64
	Run      bool
	State    uint64
}

// MainProcessToRendererKeyCallParams is the Key function parameters that the main process sends to the renderer.
type MainProcessToRendererKeyCallParams struct {
	Error        bool
	ErrorMessage string
	Run          bool
	State        uint64
}
