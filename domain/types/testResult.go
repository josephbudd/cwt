package types

// TestResult represents a user input and the correct answer.
type TestResult struct {
	Input   *KeyCodeRecord
	Control *KeyCodeRecord
}
