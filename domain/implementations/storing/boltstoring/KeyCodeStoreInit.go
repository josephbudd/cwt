package boltstoring

import (
	"github.com/josephbudd/cwt/domain/data/keyCodeTypes"
	"github.com/josephbudd/cwt/domain/types"
)

var (
	cwLetters = map[string]struct {
		Char   string
		DitDah string
	}{
		"A": {Char: "A", DitDah: ".-"},
		"B": {Char: "B", DitDah: "-..."},
		"C": {Char: "C", DitDah: "-.-."},
		"D": {Char: "D", DitDah: "-.."},
		"E": {Char: "E", DitDah: "."},
		"F": {Char: "F", DitDah: "..-."},
		"G": {Char: "G", DitDah: "--."},
		"H": {Char: "H", DitDah: "...."},
		"I": {Char: "I", DitDah: ".."},
		"J": {Char: "J", DitDah: ".---"},
		"K": {Char: "K", DitDah: "-.-"},
		"L": {Char: "L", DitDah: ".-.."},
		"M": {Char: "M", DitDah: "--"},
		"N": {Char: "N", DitDah: "-."},
		"O": {Char: "O", DitDah: "---"},
		"P": {Char: "P", DitDah: ".--."},
		"Q": {Char: "Q", DitDah: "--.-"},
		"R": {Char: "R", DitDah: ".-."},
		"S": {Char: "S", DitDah: "..."},
		"T": {Char: "T", DitDah: "-"},
		"U": {Char: "U", DitDah: "..-"},
		"V": {Char: "V", DitDah: "...-"},
		"W": {Char: "W", DitDah: ".--"},
		"X": {Char: "X", DitDah: "-..-"},
		"Y": {Char: "Y", DitDah: "-.--"},
		"Z": {Char: "Z", DitDah: "--.."},
	}

	cwNumbers = map[string]struct {
		Char   string
		DitDah string
	}{
		"1": {Char: "1", DitDah: ".----"},
		"2": {Char: "2", DitDah: "..---"},
		"3": {Char: "3", DitDah: "...--"},
		"4": {Char: "4", DitDah: "....-"},
		"5": {Char: "5", DitDah: "....."},
		"6": {Char: "6", DitDah: "-...."},
		"7": {Char: "7", DitDah: "--..."},
		"8": {Char: "8", DitDah: "---.."},
		"9": {Char: "9", DitDah: "----."},
		"0": {Char: "0", DitDah: "-----"},
	}

	cwPunctuation = map[string]struct {
		Char   string
		DitDah string
	}{
		"Period":        {Char: ".", DitDah: ".-.-.-"},
		"Comma":         {Char: ",", DitDah: "--..--"},
		"Slash":         {Char: "/", DitDah: "-..-."},
		"Plus":          {Char: "+", DitDah: ".-.-."},
		"Equals":        {Char: "=", DitDah: "-...-"},
		"Question Mark": {Char: "?", DitDah: "..--.."},
		"Open Paren":    {Char: "(", DitDah: "-.--."},
		"Close Paren":   {Char: ")", DitDah: "-.--.-"},
		"Dash":          {Char: "-", DitDah: "-....-"},
		"Double Quote":  {Char: "\"", DitDah: ".-..-."},
		"Underline":     {Char: "_", DitDah: "..--.-"},
		"Single Quote":  {Char: "'", DitDah: ".----."},
		"Colon":         {Char: ":", DitDah: "---..."},
		"Semicolon":     {Char: ";", DitDah: "-.-.-."},
		"Dollar Sign":   {Char: "$", DitDah: "...-..-"},
		"At Sign":       {Char: "@", DitDah: ".--.-."},
	}

	cwSpecial = map[string]struct {
		Char   string
		DitDah string
	}{
		"Warning":    {Char: "Warning", DitDah: ".-..-"},
		"Error":      {Char: "Error", DitDah: "........"},
		"Repetition": {Char: "Repetition", DitDah: ".. .."},
	}
)

// initialize is only useful if you want to add the default records to the db.
// otherwise you don't need it to do anything.
func (keyCodedb *KeyCodeBoltDB) initialize() (err error) {
	for name, cd := range cwLetters {
		record := types.NewKeyCodeRecord()
		record.Name = name
		record.Character = cd.Char
		record.DitDah = cd.DitDah
		record.Type = keyCodeTypes.KeyCodeTypeLetter
		if err = keyCodedb.updateKeyCodeBucket(record); err != nil {
			return
		}
	}
	for name, cd := range cwNumbers {
		record := types.NewKeyCodeRecord()
		record.Name = name
		record.Character = cd.Char
		record.DitDah = cd.DitDah
		record.Type = keyCodeTypes.KeyCodeTypeNumber
		if err = keyCodedb.updateKeyCodeBucket(record); err != nil {
			return err
		}
	}
	for name, cd := range cwPunctuation {
		record := types.NewKeyCodeRecord()
		record.Name = name
		record.Character = cd.Char
		record.DitDah = cd.DitDah
		record.Type = keyCodeTypes.KeyCodeTypePunctuation
		if err = keyCodedb.updateKeyCodeBucket(record); err != nil {
			return err
		}
	}
	for name, cd := range cwSpecial {
		record := types.NewKeyCodeRecord()
		record.Name = name
		record.Character = cd.Char
		record.DitDah = cd.DitDah
		record.Type = keyCodeTypes.KeyCodeTypeSpecial
		if err = keyCodedb.updateKeyCodeBucket(record); err != nil {
			return err
		}
	}
	return nil
}
