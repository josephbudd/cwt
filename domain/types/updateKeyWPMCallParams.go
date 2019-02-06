package types

// RendererToMainProcessUpdateKeyWPMCallParams is the GetKeyWPM function parameters that the renderer sends to the main process.
type RendererToMainProcessUpdateKeyWPMCallParams struct {
	Record *WPMRecord
}

// MainProcessToRendererUpdateKeyWPMCallParams is the UpdateKeyWPM function parameters that the main process sends to the renderer.
type MainProcessToRendererUpdateKeyWPMCallParams struct {
	Error        bool
	ErrorMessage string
	Record       *WPMRecord
}
