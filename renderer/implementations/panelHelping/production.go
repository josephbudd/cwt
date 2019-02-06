package panelHelping

// NewProductionHelper constructs a new  production helper.
func NewProductionHelper() *ProductionHelper {
	v := &ProductionHelper{}
	return v
}

// ProductionHelper is the helper for this app.
type ProductionHelper struct{}

// StateLetter returns the state for letters.
func (helper *ProductionHelper) StateLetter() uint64 { return uint64(1) }

// StateNumber returns the state for numbers.
func (helper *ProductionHelper) StateNumber() uint64 { return uint64(1 << 1) }

// StatePunctuation returns the state for punctuation.
func (helper *ProductionHelper) StatePunctuation() uint64 { return uint64(1 << 2) }

// StateSpecial returns the state for special chars.
func (helper *ProductionHelper) StateSpecial() uint64 { return uint64(1 << 3) }

// StatePractice returns the state for practice.
func (helper *ProductionHelper) StatePractice() uint64 { return uint64(1 << 4) }

// StateTest returns the state for test.
func (helper *ProductionHelper) StateTest() uint64 { return uint64(1 << 5) }
