package keycodes

// IDs for valid key codes start at 100.
// These invalid key codes have an id < 100.

import "github.com/josephbudd/cwt/domain/store/record"

// NotKeyedByApp is when the user copies a key that never happened.
var NotKeyedByApp *record.KeyCode

// NotKeyedByUser is a char that the user never keyed.
var NotKeyedByUser *record.KeyCode

// NotCopiedByUser is a key that the user missed.
var NotCopiedByUser *record.KeyCode

// NoCopyToKey represents non-existant text that the user keyed.
var NoCopyToKey *record.KeyCode

// UnknownKeyFromUser represents non-existant key that the user keyed.
var UnknownKeyFromUser *record.KeyCode

const (
	emptyBrackets = "{ }" // nbsp; in brackets
	doNotEnter    = "{🚫}"
	madFace       = "{😖}"
	trash         = "{🗑}"
	poop          = "{💩}"
)

func init() {
	NotKeyedByApp = record.NewKeyCode()
	NotKeyedByApp.ID = 3
	NotKeyedByApp.Character = emptyBrackets
	NotKeyedByApp.DitDah = emptyBrackets

	NotCopiedByUser = record.NewKeyCode()
	NotCopiedByUser.ID = 4
	NotCopiedByUser.Character = emptyBrackets
	NotCopiedByUser.DitDah = emptyBrackets

	NotKeyedByUser = record.NewKeyCode()
	NotKeyedByUser.ID = 6
	NotKeyedByUser.Character = emptyBrackets
	NotKeyedByUser.DitDah = emptyBrackets // doNotEnter

	NoCopyToKey = record.NewKeyCode()
	NoCopyToKey.ID = 7
	NoCopyToKey.Character = emptyBrackets
	NoCopyToKey.DitDah = emptyBrackets

	UnknownKeyFromUser = record.NewKeyCode()
	UnknownKeyFromUser.ID = 8
	UnknownKeyFromUser.Character = poop
	UnknownKeyFromUser.DitDah = poop
}
