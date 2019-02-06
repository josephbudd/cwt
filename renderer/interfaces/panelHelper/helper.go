package panelHelper

// Helper is help for setting up the markup panels.
type Helper interface {

	// StateLetter returns the state for letters.
	StateLetter() uint64

	// StateNumber returns the state for numbers.
	StateNumber() uint64

	// StatePunctuation returns the state for punctuation.
	StatePunctuation() uint64

	// StateSpecial returns the state for special chars.
	StateSpecial() uint64

	// StatePractice returns the state for practicing.
	StatePractice() uint64

	// StateTest returns the state for testing.
	StateTest() uint64
}
