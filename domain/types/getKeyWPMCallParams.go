package types

// MainProcessToRendererGetKeyWPMCallParams is the GetKeyWPM function parameters that the main process sends to the renderer.
type MainProcessToRendererGetKeyWPMCallParams struct {
	Error        bool
	ErrorMessage string
	Record       *WPMRecord
}
