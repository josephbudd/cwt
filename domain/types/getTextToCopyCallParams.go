package types

// RendererToMainProcessGetTextToCopyCallParams is the GetTextToCopy function parameters that the renderer sends to the main process.
type RendererToMainProcessGetTextToCopyCallParams struct {
	State uint64
}

// MainProcessToRendererGetTextToCopyCallParams is the GetTextToCopy function parameters that the main process sends to the renderer.
type MainProcessToRendererGetTextToCopyCallParams struct {
	Error        bool
	ErrorMessage string
	Solution     [][]*KeyCodeRecord
	WPM          uint64
	State        uint64
}
