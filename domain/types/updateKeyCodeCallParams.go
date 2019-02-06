package types

// RendererToMainProcessUpdateKeyCodeCallParams is the UpdateKeyCode function parameters that the renderer sends to the main process.
type RendererToMainProcessUpdateKeyCodeCallParams struct {
	Record *KeyCodeRecord
	State  uint64
}

// MainProcessToRendererUpdateKeyCodeCallParams is the UpdateKeyCode function parameters that the main process sends to the renderer.
type MainProcessToRendererUpdateKeyCodeCallParams struct {
	Error        bool
	ErrorMessage string
	Record       *KeyCodeRecord
	State        uint64
}
