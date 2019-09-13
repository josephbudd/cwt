package paneling

// Help helps initialize panels.
type Help struct{}

// NewHelp constructs a new *Help.
func NewHelp() *Help {
	return &Help{}
}

// StateLetter returns the state for letters.
func (help *Help) StateLetter() uint64 { return uint64(1) }

// StateNumber returns the state for numbers.
func (help *Help) StateNumber() uint64 { return uint64(1 << 1) }

// StatePunctuation returns the state for punctuation.
func (help *Help) StatePunctuation() uint64 { return uint64(1 << 2) }

// StateSpecial returns the state for special chars.
func (help *Help) StateSpecial() uint64 { return uint64(1 << 3) }

// StatePractice returns the state for practice.
func (help *Help) StatePractice() uint64 { return uint64(1 << 4) }

// StateTest returns the state for test.
func (help *Help) StateTest() uint64 { return uint64(1 << 5) }
