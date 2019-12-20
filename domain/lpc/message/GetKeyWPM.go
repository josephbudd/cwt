package message

import "github.com/josephbudd/cwt/domain/store/record"

// GetKeyWPMRendererToMainProcess is the GetKeyWPM message that the renderer sends to the main process.
type GetKeyWPMRendererToMainProcess struct {
}

// GetKeyWPMMainProcessToRenderer is the GetKeyWPM message that the main process sends to the renderer.
type GetKeyWPMMainProcessToRenderer struct {
	Error        bool
	ErrorMessage string
	Fatal        bool

	Record *record.WPM
}
