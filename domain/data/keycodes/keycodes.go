package keycodes

// IDs for valid key codes start at 100.
// These invalid key codes have an id < 100.

import "github.com/josephbudd/cwt/domain/store/record"

// NotKeyed represents
//  a code that was not keyed.
var NotKeyed *record.KeyCode

// NotKeyedByApp is when the user copies a key that never happened.
var NotKeyedByApp *record.KeyCode

// NotKeyedByUser is a char that the user never keyed.
var NotKeyedByUser *record.KeyCode

// NotCopiedByApp is ?
var NotCopiedByApp *record.KeyCode

// NotCopiedByUser is a key that the user missed.
var NotCopiedByUser *record.KeyCode

// NotCopied represents
//  or a key that was not copied.
var NotCopied *record.KeyCode

// NotInText represents non-existant text that the user keyed.
var NotInText *record.KeyCode

// UnknownKeyFromUser represents non-existant key that the user keyed.
var UnknownKeyFromUser *record.KeyCode

const nbsp = " "

func init() {
	NotCopied = record.NewKeyCode()
	NotCopied = &record.KeyCode{
		ID:        1,
		Character: "{NotCopied}",
		DitDah:    "{NotCopied}",
	}
	NotKeyed = record.NewKeyCode()
	NotKeyed = &record.KeyCode{
		ID:        2,
		Character: "{NotKeyed}",
		DitDah:    "{NotKeyed}",
	}
	NotKeyedByApp = record.NewKeyCode()
	NotKeyedByApp.ID = 3
	NotKeyedByApp.Character = "{Not Keyed BA}"
	NotKeyedByApp.DitDah = "{Not Keyed BA}"

	NotCopiedByUser = record.NewKeyCode()
	NotCopiedByUser.ID = 4
	NotCopiedByUser.Character = "{Not Copied BU}"
	NotCopiedByUser.DitDah = "{Not Copied BU}"

	NotCopiedByApp = record.NewKeyCode()
	NotCopiedByApp.ID = 5
	NotCopiedByApp.Character = "{Not Copied BA}"
	NotCopiedByApp.DitDah = "{Not Copied BA}"

	NotKeyedByUser = record.NewKeyCode()
	NotKeyedByUser.ID = 6
	NotKeyedByUser.Character = "{ }"
	NotKeyedByUser.DitDah = "{ }"

	NotInText = record.NewKeyCode()
	NotInText.ID = 7
	NotInText.Character = "{No Text}"
	NotInText.DitDah = "{No Text}"

	UnknownKeyFromUser = record.NewKeyCode()
	UnknownKeyFromUser.ID = 8
	UnknownKeyFromUser.Character = "{Unknown Key}"
	UnknownKeyFromUser.DitDah = "{Unknown Key}"
}
