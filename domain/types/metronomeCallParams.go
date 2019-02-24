package types

// RendererToMainProcessMetronomeCallParams are the Metronome function parameters that the renderer sends to the main process.
type RendererToMainProcessMetronomeCallParams struct {
	Run   bool
	State uint64
	WPM   uint64
}

// MainProcessToRendererMetronomeCallParams are the Metronome function parameters that the main process sends to the renderer.
type MainProcessToRendererMetronomeCallParams struct {
	Error        bool
	ErrorMessage string
	Run          bool
	State        uint64
}
