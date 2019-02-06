package types

// MainProcessToRendererGetCopyWPMCallParams is the GetCopyWPM function parameters that the main process sends to the renderer.
type MainProcessToRendererGetCopyWPMCallParams struct {
	Error        bool
	ErrorMessage string
	Record       *WPMRecord
}
