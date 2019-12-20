package message

import "github.com/josephbudd/cwt/domain/store/record"

// UpdateCopyWPMRendererToMainProcess is the UpdateCopyWPM message that the renderer sends to the main process.
type UpdateCopyWPMRendererToMainProcess struct {
	Record *record.WPM
}

// UpdateCopyWPMMainProcessToRenderer is the UpdateCopyWPM message that the main process sends to the renderer.
type UpdateCopyWPMMainProcessToRenderer struct {
	Error        bool
	ErrorMessage string
	Fatal        bool

	Record *record.WPM
}
