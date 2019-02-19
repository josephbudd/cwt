package keyservice

import (
	"errors"
	"fmt"
	"strings"

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
	// get data from the stores.
	keyCodeRecords, err := keyCodeStorer.GetKeyCodes()
	if err != nil {
		err = errors.New("Copy keyCodeStorer.GetKeyCodes() error is " + err.Error())
		return
	}
	// define ms
	elementsPerMinute := (wpm * 50) - 7
	elementMS := int64(60000 / elementsPerMinute)
	// allow dits to be a little long
	ditMS := elementMS + (elementMS / 2)
	fmt.Printf("ditMS is %d\n", ditMS)
	//dahMS := 3 * elementMS
	// pauses
	//ditdahPauseMS := elementMS   // between dits and dahs
	ditdahPauseMS := int64(1.5 * float64(elementMS))
	fmt.Println("ditdahPauseMS is ", ditdahPauseMS)
	charPauseMS := int64(4.5 * float64(elementMS)) // 3 * elementMS // between ".-" and "-..."
	fmt.Println("charPauseMS is ", charPauseMS)
	//wordPauseMS := 7 * elementMS // between "._ -..." and "_... ._"
	// process stack ( pauseTime, keydownTime, ...)
	solution = make([][]*types.KeyCodeRecord, 0, 100)
	ditdahCharStack := make([]string, 0, 5)
	ditdahs := make([]string, 0, 5)
	for i, ms := range milliSeconds {
		// the first millisecond is a pause before keying so ignore it.
		if i > 0 {
			if i%2 == 0 {
				// pause
				// if ms <= ditdahPauseMS continue to next dit or dah
				//if ms > ditdahPauseMS {
				fmt.Printf("ms is %d, ditdahPauseMS is %d, charPauseMS is %d\n", ms, ditdahPauseMS, charPauseMS)
				if ms <= ditdahPauseMS {
					// pause between dits and dahs inside a word
				} else if ms < charPauseMS {
					// pause between chars, between ".-" and "-..."
					// or pause between words, between "._ -..." and "_... ._"
					if len(ditdahs) > 0 {
						ditdahChar := strings.Join(ditdahs, "")
						ditdahCharStack = append(ditdahCharStack, ditdahChar)
						ditdahs = ditdahs[:0]
						fmt.Println("copied ditdahChar is ", ditdahChar)
					}
				} else {
					// word pause
					// pause between words, between "._ -..." and "_... ._"
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
					}
				}
			} else {
				// key
				if ms <= ditMS {
					ditdahs = append(ditdahs, ".")
				} else {
					ditdahs = append(ditdahs, "-")
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
	record = nil
	return
}
