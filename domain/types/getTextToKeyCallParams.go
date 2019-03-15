package types

// HowTo is data for how to key a character.
type HowTo struct {
	Character    string
	DitDah       string
	Instructions string
}

// RendererToMainProcessGetTextWPMToKeyCallParams is the GetTextWPMToKey function parameters that the renderer sends to the main process.
type RendererToMainProcessGetTextWPMToKeyCallParams struct {
	State    uint64
	Practice bool
}

// MainProcessToRendererGetTextWPMToKeyCallParams is the GetTextWPMToKey function parameters that the main process sends to the renderer.
type MainProcessToRendererGetTextWPMToKeyCallParams struct {
	Solution     [][]*KeyCodeRecord
	Help         [][]HowTo
	Practice     bool
	WPM          uint64
	State        uint64
	Error        bool
	ErrorMessage string
}
