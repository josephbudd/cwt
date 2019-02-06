package types

// MisMatch represents a user input and the correct answer.
type MisMatch struct {
	Input   *KeyCodeRecord
	Control *KeyCodeRecord
}
