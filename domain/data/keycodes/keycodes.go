package keycodes

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

const nbsp = " "

func init() {
	NotCopied = &types.KeyCodeRecord{
		Character: "{NotCopied}",
		DitDah:    "{NotCopied}",
	}
	NotKeyed = &types.KeyCodeRecord{
		Character: "{NotKeyed}",
		DitDah:    "{NotKeyed}",
	}
	NotKeyedByApp = &types.KeyCodeRecord{
		Character: "{Not Keyed BA}",
		DitDah:    "{Not Keyed BA}",
	}

	NotCopiedByUser = &types.KeyCodeRecord{
		Character: "{Not Copied BU}",
		DitDah:    "{Not Copied BU}",
	}

	NotCopiedByApp = &types.KeyCodeRecord{
		Character: "{Not Copied BA}",
		DitDah:    "{Not Copied BA}",
	}
	NotKeyedByUser = &types.KeyCodeRecord{
		Character: "{ }",
		DitDah:    "{ }",
	}
	NotInText = &types.KeyCodeRecord{
		Character: "{No Text}",
		DitDah:    "{No Text}",
	}

}
