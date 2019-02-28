package keyservice

import (
	"errors"
	"fmt"
	"strings"

	"github.com/josephbudd/cwt/domain/data/keycodes"
	"github.com/josephbudd/cwt/domain/interfaces/storer"
	"github.com/josephbudd/cwt/domain/types"
)

// Copy converts what the user keyed to text and ditdahs.
// It converts time of key-up and key-down durations to lines of text.
// Params milliSeconds is a slice of durations
//   0 is key-up duration, 1 is key-down duration,
//   2 is key-up duration, 3 is key-down duration, etc.
// Params wpm is the words per minute that the user is attempting to key at.
// Param wPMStorer is the wpm repository.Copy
// Param keyCodeStorer is the keycode test results repository.
func Copy(milliSeconds []int64, wpm uint64, keyCodeStorer storer.KeyCodeStorer) (solution [][]*types.KeyCodeRecord, err error) {
	fmt.Printf("%#v\n", milliSeconds)
	// the pause multiplier adjusts the pause times.
	pauseMultiplier := 1.5
	// the key multiplier adjusts the key times.
	keyMultiplier := 1.5
	// get data from the store.
	keyCodeRecords, err := keyCodeStorer.GetKeyCodes()
	if err != nil {
		err = errors.New("Copy keyCodeStorer.GetKeyCodes() error is " + err.Error())
		return
	}
	// 1. define the true milleseconds of a single element.
	elementsPerMinute := (wpm * 50) - 7
	elementMS := int64(60000 / elementsPerMinute)
	fmt.Printf("elementMS is %d\n", elementMS)
	// 2. define the maximum milliseconds allowed for a dit.
	//    a key down ms <= ditMaxMS is a dit.
	//    a key down ms >= ditMaxMS is a dah.
	ditMaxMS := int64(keyMultiplier * float64(elementMS))
	fmt.Printf("ditMaxMS is %d\n", ditMaxMS)
	// 3.a define the maximum millseconds allowed for a pause between dits and dahs
	// 3.b define the maximum millseconds allowed for a pause between characters
	betweenDitdahPauseMaxMS := int64(pauseMultiplier * float64(elementMS))
	betweenCharPauseMaxMS := int64(pauseMultiplier * float64(3*elementMS))
	// wordPauseMaxMS := int64(pauseMultiplier * float64(7*elementMS))
	fmt.Println("betweenDitdahPauseMaxMS is ", betweenDitdahPauseMaxMS)
	fmt.Println("betweenCharPauseMaxMS is ", betweenCharPauseMaxMS)
	// fmt.Println("wordPauseMaxMS is ", wordPauseMaxMS)
	// process stack ( pauseTime, keydownTime, ...)
	solution = make([][]*types.KeyCodeRecord, 0, 100)
	ditdahCharStack := make([]string, 0, 5)
	ditdahs := make([]string, 0, 5)
	for i, ms := range milliSeconds {
		// the first millisecond is a pause before keying so ignore it.
		if i > 0 {
			if i%2 == 0 {
				// this is a pause
				if ms <= betweenDitdahPauseMaxMS {
					// pause between dits and dahs inside a character.
					// between "." and "-" in ".-" ( "a" )
					fmt.Println("Pause between dits and dahs is ", ms)
				} else if ms < betweenCharPauseMaxMS {
					// pause between chars in a word.
					// between ".-" and "-." in ".- -." ( "an" )
					if len(ditdahs) > 0 {
						ditdahChar := strings.Join(ditdahs, "")
						ditdahCharStack = append(ditdahCharStack, ditdahChar)
						ditdahs = ditdahs[:0]
						fmt.Println("Pause between chars is ", ms)
					}
				} else {
					// pause between words in a phrase or sentence.
					// between ".- -." and ".- .--. .--. .-.. ." in ".- -.   .- .--. .--. .-.. ." ( "an apple" )
					fmt.Println("Pause between words is ", ms)
					if len(ditdahs) > 0 {
						ditdahChar := strings.Join(ditdahs, "")
						ditdahCharStack = append(ditdahCharStack, ditdahChar)
						ditdahs = ditdahs[:0]
					}
					if len(ditdahCharStack) > 0 {
						rr := ditdahCharStackToText(ditdahCharStack, keyCodeRecords)
						solution = append(solution, rr)
						ditdahCharStack = ditdahCharStack[:0]
						lrr := len(rr)
						chars := make([]string, lrr, lrr)
						for j, r := range rr {
							chars[j] = r.Character
						}
						fmt.Println("The word is ", strings.Join(chars, ""))
					}
				}
			} else {
				// key
				if ms <= ditMaxMS {
					ditdahs = append(ditdahs, ".")
					fmt.Println("Keyed a .")
				} else {
					ditdahs = append(ditdahs, "-")
					fmt.Println("Keyed a -")
				}
			}
		}
	}
	if len(ditdahs) > 0 {
		// the milliseconds did not end with an uint for a pause.
		keyed := strings.Join(ditdahs, "")
		ditdahCharStack = append(ditdahCharStack, keyed)
		rr := ditdahCharStackToText(ditdahCharStack, keyCodeRecords)
		solution = append(solution, rr)
	}
	return
}

func ditdahCharStackToText(ditdahChars []string, records []*types.KeyCodeRecord) (text []*types.KeyCodeRecord) {
	l := len(ditdahChars)
	text = make([]*types.KeyCodeRecord, l, l)
	for i, d := range ditdahChars {
		r := ditDahToRecord(d, records)
		text[i] = r
	}
	return
}

func ditDahToRecord(ditdah string, records []*types.KeyCodeRecord) (record *types.KeyCodeRecord) {
	for _, record = range records {
		if ditdah == record.DitDah {
			return
		}
	}
	// not found
	record = keycodes.UnknownKeyFromUser
	return
}
