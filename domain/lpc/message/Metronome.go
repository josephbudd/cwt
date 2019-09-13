package message

// MetronomeRendererToMainProcess is the Metronome message that the renderer sends to the main process.
type MetronomeRendererToMainProcess struct {
	Run   bool
	State uint64
	WPM   uint64
}

// MetronomeMainProcessToRenderer is the Metronome message that the main process sends to the renderer.
type MetronomeMainProcessToRenderer struct {
	Error        bool
	ErrorMessage string

	Run   bool
	State uint64
}
