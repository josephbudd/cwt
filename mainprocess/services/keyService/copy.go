package keyService

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
func Copy(milliSeconds []int64, wpm uint64, wPMStorer storer.WPMStorer, keyCodeStorer storer.KeyCodeStorer) (solution [][]*types.KeyCodeRecord, err error) {
	// get data from the stores.
	wpmRecord, err := wPMStorer.GetWPM(1)
	if err != nil {
		err = errors.New("Copy WPMStorer.GetWPM(1) error is " + err.Error())
		return
	}
	keyCodeRecords, err := keyCodeStorer.GetKeyCodes()
	if err != nil {
		err = errors.New("Copy keyCodeStorer.GetKeyCodes() error is " + err.Error())
		return
	}
	// define ms
	elementsPerMinute := (wpmRecord.WPM * 50) - 7
	elementMS := int64(60000 / elementsPerMinute)
	// allow dits to be a little long
	ditMS := elementMS + (elementMS / 2)
	fmt.Printf("ditMS is %d\n", ditMS)
	//dahMS := 3 * elementMS
	// pauses
	ditdahPauseMS := elementMS   // between dits and dahs
	charPauseMS := 3 * elementMS // between ".-" and "-..."
	//wordPauseMS := 7 * elementMS // between "._ -..." and "_... ._"
	// process stack ( pauseTime, keydownTime, ...)
	solution = make([][]*types.KeyCodeRecord, 0, 100)
	ditdahCharStack := make([]string, 0, 5)
	ditdahs := make([]string, 0, 5)
	for i, ms := range milliSeconds {
		if i%2 == 0 {
			// pause
			// if ms <= ditdahPauseMS continue to next dit or dah
			if ms > ditdahPauseMS {
				// pause between chars, between ".-" and "-..."
				// or pause between words, between "._ -..." and "_... ._"
				if len(ditdahs) > 0 {
					ditdahChar := strings.Join(ditdahs, "")
					ditdahCharStack = append(ditdahCharStack, ditdahChar)
					ditdahs = ditdahs[:0]
				}
			}
			if ms > charPauseMS {
				// pause between words, between "._ -..." and "_... ._"
				if len(ditdahCharStack) > 0 {
					rr := ditdahCharStackToText(ditdahCharStack, keyCodeRecords)
					solution = append(solution, rr)
					ditdahCharStack = ditdahCharStack[:0]
				}
			}
		} else {
			// key
			fmt.Printf("ms is %d\n", ms)
			if ms <= ditMS {
				ditdahs = append(ditdahs, ".")
			} else {
				ditdahs = append(ditdahs, "-")
			}
		}
	}
	if len(ditdahs) > 0 {
		ditdahCharStack = append(ditdahCharStack, strings.Join(ditdahs, ""))
		rr := ditdahCharStackToText(ditdahCharStack, keyCodeRecords)
		solution = append(solution, rr)
	}
	return
}

func ditdahCharStackToText(ditdahChars []string, records []*types.KeyCodeRecord) []*types.KeyCodeRecord {
	l := len(ditdahChars)
	rr := make([]*types.KeyCodeRecord, l, l)
	for i, d := range ditdahChars {
		rr[i] = ditDahToRecord(d, records)
	}
	return rr
}

func ditDahToRecord(ditdah string, records []*types.KeyCodeRecord) *types.KeyCodeRecord {
	for _, r := range records {
		if ditdah == r.DitDah {
			return r
		}
	}
	return nil
}
