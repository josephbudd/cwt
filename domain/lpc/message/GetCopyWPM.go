package message

import "github.com/josephbudd/cwt/domain/store/record"

// GetCopyWPMRendererToMainProcess is the GetCopyWPM message that the renderer sends to the main process.
type GetCopyWPMRendererToMainProcess struct {
}

// GetCopyWPMMainProcessToRenderer is the GetCopyWPM message that the main process sends to the renderer.
type GetCopyWPMMainProcessToRenderer struct {
	Error        bool
	ErrorMessage string
	Fatal        bool

	Record *record.WPM
}
