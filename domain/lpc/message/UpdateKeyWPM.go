package message

import "github.com/josephbudd/cwt/domain/store/record"

// UpdateKeyWPMRendererToMainProcess is the UpdateKeyWPM message that the renderer sends to the main process.
type UpdateKeyWPMRendererToMainProcess struct {
	Record *record.WPM
}

// UpdateKeyWPMMainProcessToRenderer is the UpdateKeyWPM message that the main process sends to the renderer.
type UpdateKeyWPMMainProcessToRenderer struct {
	Error        bool
	ErrorMessage string
	Fatal        bool

	Record *record.WPM
}
