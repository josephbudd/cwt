package types

/*

	TODO:

	You need to complete these record definitions.

*/

// WPMRecord is a WPM record.
type WPMRecord struct {
	ID  uint64
	WPM uint64 `json:"wpm"`
}

// KeyCodeRecordResult is a test result.
type KeyCodeRecordResult struct {
	Correct, Attempts uint64
}

// KeyCodeRecord is a KeyCode record.
type KeyCodeRecord struct {
	ID        uint64
	Name      string `json:"name"`
	Character string `json:"character"`
	DitDah    string `json:"ditdah"`
	Selected  bool   `json:"selected"`
	Type      uint64 `json:"type"`

	KeyWPMResults  map[uint64]KeyCodeRecordResult
	CopyWPMResults map[uint64]KeyCodeRecordResult
}

// NewWPMRecord constructs a new WPM record.
func NewWPMRecord() *WPMRecord {
	v := &WPMRecord{}
	return v
}

// NewKeyCodeRecord constructs a new KeyCode record.
func NewKeyCodeRecord() *KeyCodeRecord {
	return &KeyCodeRecord{
		KeyWPMResults: map[uint64]KeyCodeRecordResult{
			5:  {0, 0},
			10: {0, 0},
			15: {0, 0},
			20: {0, 0},
			25: {0, 0},
			30: {0, 0},
			35: {0, 0},
		},
		CopyWPMResults: map[uint64]KeyCodeRecordResult{
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
