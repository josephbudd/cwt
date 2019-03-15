package keycodes

// IDs for valid key codes start at 100.
// These invalid key codes have an id < 100.

import "github.com/josephbudd/cwt/domain/types"

// NotKeyed represents
//  a code that was not keyed.
var NotKeyed *types.KeyCodeRecord

// NotKeyedByApp is when the user copies a key that never happened.
var NotKeyedByApp *types.KeyCodeRecord

// NotKeyedByUser is a char that the user never keyed.
var NotKeyedByUser *types.KeyCodeRecord

// NotCopiedByApp is ?
var NotCopiedByApp *types.KeyCodeRecord

// NotCopiedByUser is a key that the user missed.
var NotCopiedByUser *types.KeyCodeRecord

// NotCopied represents
//  or a key that was not copied.
var NotCopied *types.KeyCodeRecord

// NotInText represents non-existant text that the user keyed.
var NotInText *types.KeyCodeRecord

// UnknownKeyFromUser represents non-existant key that the user keyed.
var UnknownKeyFromUser *types.KeyCodeRecord

const nbsp = " "

func init() {
	NotCopied = types.NewKeyCodeRecord()
	NotCopied = &types.KeyCodeRecord{
		ID:        1,
		Character: "{NotCopied}",
		DitDah:    "{NotCopied}",
	}
	NotKeyed = types.NewKeyCodeRecord()
	NotKeyed = &types.KeyCodeRecord{
		ID:        2,
		Character: "{NotKeyed}",
		DitDah:    "{NotKeyed}",
	}
	NotKeyedByApp = types.NewKeyCodeRecord()
	NotKeyedByApp.ID = 3
	NotKeyedByApp.Character = "{Not Keyed BA}"
	NotKeyedByApp.DitDah = "{Not Keyed BA}"

	NotCopiedByUser = types.NewKeyCodeRecord()
	NotCopiedByUser.ID = 4
	NotCopiedByUser.Character = "{Not Copied BU}"
	NotCopiedByUser.DitDah = "{Not Copied BU}"

	NotCopiedByApp = types.NewKeyCodeRecord()
	NotCopiedByApp.ID = 5
	NotCopiedByApp.Character = "{Not Copied BA}"
	NotCopiedByApp.DitDah = "{Not Copied BA}"

	NotKeyedByUser = types.NewKeyCodeRecord()
	NotKeyedByUser.ID = 6
	NotKeyedByUser.Character = "{ }"
	NotKeyedByUser.DitDah = "{ }"

	NotInText = types.NewKeyCodeRecord()
	NotInText.ID = 7
	NotInText.Character = "{No Text}"
	NotInText.DitDah = "{No Text}"

	UnknownKeyFromUser = types.NewKeyCodeRecord()
	UnknownKeyFromUser.ID = 8
	UnknownKeyFromUser.Character = "{Unknown Key}"
	UnknownKeyFromUser.DitDah = "{Unknown Key}"
}
