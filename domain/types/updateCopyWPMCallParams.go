package types

// RendererToMainProcessUpdateCopyWPMCallParams is the GetCopyWPM function parameters that the renderer sends to the main process.
type RendererToMainProcessUpdateCopyWPMCallParams struct {
	Record *WPMRecord
}

// MainProcessToRendererUpdateCopyWPMCallParams is the UpdateCopyWPM function parameters that the main process sends to the renderer.
type MainProcessToRendererUpdateCopyWPMCallParams struct {
	Error        bool
	ErrorMessage string
	Record       *WPMRecord
}
