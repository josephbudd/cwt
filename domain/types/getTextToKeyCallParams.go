package types

// RendererToMainProcessGetTextWPMToKeyCallParams is the GetTextWPMToKey function parameters that the renderer sends to the main process.
type RendererToMainProcessGetTextWPMToKeyCallParams struct {
	State uint64
}

// MainProcessToRendererGetTextWPMToKeyCallParams is the GetTextWPMToKey function parameters that the main process sends to the renderer.
type MainProcessToRendererGetTextWPMToKeyCallParams struct {
	Solution     [][]*KeyCodeRecord
	WPM          uint64
	State        uint64
	Error        bool
	ErrorMessage string
}
