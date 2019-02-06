package types

// RendererToMainProcessGetKeyCodesCallParams is the GetKeyCodes function parameters that the renderer sends to the main process.
type RendererToMainProcessGetKeyCodesCallParams struct {
	State uint64
}

// MainProcessToRendererGetKeyCodesCallParams is the GetKeyCodes function parameters that the main process sends to the renderer.
type MainProcessToRendererGetKeyCodesCallParams struct {
	Error        bool
	ErrorMessage string
	Records      []*KeyCodeRecord
	State        uint64
}
