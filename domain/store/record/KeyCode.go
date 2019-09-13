package record

/*

	TODO:

	You need to complete this local bolt store record definition.

*/

// KeyCodeResult is a test result.
type KeyCodeResult struct {
	Correct, Attempts uint64
}

// KeyCode is the local bolt store KeyCode record.
type KeyCode struct {
	ID        uint64
	Name      string `json:"name"`
	Character string `json:"character"`
	DitDah    string `json:"ditdah"`
	Selected  bool   `json:"selected"`
	Type      uint64 `json:"type"`

	KeyWPMResults  map[uint64]KeyCodeResult
	CopyWPMResults map[uint64]KeyCodeResult
}

// NewKeyCode constructs a new local bolt store KeyCode.
func NewKeyCode() *KeyCode {
	return &KeyCode{
		KeyWPMResults: map[uint64]KeyCodeResult{
			5:  {0, 0},
			10: {0, 0},
			15: {0, 0},
			20: {0, 0},
			25: {0, 0},
			30: {0, 0},
			35: {0, 0},
		},
		CopyWPMResults: map[uint64]KeyCodeResult{
			5:  {0, 0},
			10: {0, 0},
			15: {0, 0},
			20: {0, 0},
			25: {0, 0},
			30: {0, 0},
			35: {0, 0},
		},
	}
}
