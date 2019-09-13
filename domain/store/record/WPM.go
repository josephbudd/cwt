package record

/*

	TODO:

	You need to complete this local bolt store record definition.

*/

// WPM is the local bolt store WPM record.
type WPM struct {
	ID  uint64
	WPM uint64 `json:"wpm"`
}

// NewWPM constructs a new local bolt store WPM.
func NewWPM() *WPM {
	v := &WPM{}
	return v
}
