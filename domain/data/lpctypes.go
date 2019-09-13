package data

import "github.com/josephbudd/cwt/domain/store/record"

// TestResult represents a user input and the correct answer.
type TestResult struct {
	Input   *record.KeyCode
	Control *record.KeyCode
}

// HowTo is data for how to key a character.
type HowTo struct {
	Character    string
	DitDah       string
	Instructions string
}
