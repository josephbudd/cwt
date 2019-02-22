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
	NotCopied = &types.KeyCodeRecord{
		ID:        1,
		Character: "{NotCopied}",
		DitDah:    "{NotCopied}",
	}
	NotKeyed = &types.KeyCodeRecord{
		ID:        2,
		Character: "{NotKeyed}",
		DitDah:    "{NotKeyed}",
	}
	NotKeyedByApp = &types.KeyCodeRecord{
		ID:        3,
		Character: "{Not Keyed BA}",
		DitDah:    "{Not Keyed BA}",
	}

	NotCopiedByUser = &types.KeyCodeRecord{
		ID:        4,
		Character: "{Not Copied BU}",
		DitDah:    "{Not Copied BU}",
	}

	NotCopiedByApp = &types.KeyCodeRecord{
		ID:        5,
		Character: "{Not Copied BA}",
		DitDah:    "{Not Copied BA}",
	}
	NotKeyedByUser = &types.KeyCodeRecord{
		ID:        6,
		Character: "{ }",
		DitDah:    "{ }",
	}
	NotInText = &types.KeyCodeRecord{
		ID:        7,
		Character: "{No Text}",
		DitDah:    "{No Text}",
	}
	UnknownKeyFromUser = &types.KeyCodeRecord{
		ID:        8,
		Character: "{Unknown Key}",
		DitDah:    "{Unknown Key}",
	}

}
