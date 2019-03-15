package keyservice

import (
	"fmt"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/josephbudd/cwt/domain/data/keycodes"

	"github.com/boltdb/bolt"
	"github.com/pkg/errors"

	"github.com/josephbudd/cwt/domain/data/filepaths"
	"github.com/josephbudd/cwt/domain/implementations/storing/boltstoring"
	"github.com/josephbudd/cwt/domain/interfaces/storer"
	"github.com/josephbudd/cwt/domain/types"
)

var (
	keyCodeStore storer.KeyCodeStorer
)

func getElementMS(wpm uint64) (elementMS int64) {
	elementsPerMinute := (wpm * 50) - 7
	elementMS = int64(60000 / elementsPerMinute)
	return
}

func getCopyTestKeyCodes() (keyCodeRecords []*types.KeyCodeRecord, err error) {
	if keyCodeRecords, err = keyCodeStore.GetKeyCodes(); err != nil {
		err = errors.WithMessage(err, `keyCodeStore.GetKeyCodes()`)
	}
	return
}

// buildBoltStores makes bolt data stores.
// Each store is the implementation of an interface defined in package repoi.
// Each store uses the same bolt database so closing one will close all.
func buildBoltStores() (path string, db *bolt.DB, err error) {
	if path, err = filepaths.BuildUserSubFoldersPath("boltdb"); err != nil {
		err = errors.WithMessage(err, `filepaths.BuildUserSubFoldersPath("boltdb")`)
		return
	}
	path = filepath.Join(path, "allstores.nosql")
	if db, err = bolt.Open(path, filepaths.GetFmode(), nil); err != nil {
		err = errors.WithMessage(err, `bolt.Open(path, filepaths.GetFmode(), nil)`)
	}
	return
}

func TestCopy(t *testing.T) {
	var path string
	var db *bolt.DB
	var err error
	if path, db, err = buildBoltStores(); err != nil {
		t.Fatal(err)
	}
	keyCodeStore = boltstoring.NewKeyCodeBoltDB(db, path, filepaths.GetFmode())
	defer keyCodeStore.Close()
	var keyCodeRecords []*types.KeyCodeRecord
	if keyCodeRecords, err = getCopyTestKeyCodes(); err != nil {
		t.Fatal(err)
	}
	l := len(keyCodeRecords)
	wpm := uint64(20)
	// pauses
	ditMS := getElementMS(wpm)
	dahMS := 3 * ditMS
	ddPauseMS := ditMS
	charPauseMS := ditMS * 3
	wordPauseMS := ditMS * 7
	// build the milliseconds and the solution
	milliSeconds := make([]int64, 0, (l*2)+30)
	milliSeconds = append(milliSeconds, wordPauseMS)
	// make the words 5 characters long.
	wordLength := 5
	wordCount := len(keyCodeRecords) / wordLength
	lastI := wordCount - 1
	wordCount++
	wantSolution := make([][]*types.KeyCodeRecord, wordCount, wordCount)
	for i := 0; i <= lastI; i++ {
		start := i * wordLength
		end := start + wordLength
		wantSolution[i] = keyCodeRecords[start:end]
		var j int
		for j = start; j < end; j++ {
			ditdah := keyCodeRecords[j].DitDah
			fmt.Println("ditdah is ", ditdah, ": char is ", keyCodeRecords[j].Character)
			lditdah := len(ditdah) - 1
			for k, dd := range ditdah {
				// dit or dah
				if dd == '.' {
					milliSeconds = append(milliSeconds, ditMS)
				} else {
					milliSeconds = append(milliSeconds, dahMS)
				}
				if k < lditdah {
					// new dit or dah will follow
					milliSeconds = append(milliSeconds, ddPauseMS)
				}
			}
			if j < end-1 {
				// new character in word will follow
				milliSeconds = append(milliSeconds, charPauseMS)
			}
		}
		// still more words
		milliSeconds = append(milliSeconds, wordPauseMS)
	}
	// final word is an unknown char: 10 dahs
	for i := 0; i < 10; i++ {
		milliSeconds = append(milliSeconds, dahMS)
		if i < 9 {
			milliSeconds = append(milliSeconds, ddPauseMS)
		}
	}
	// final solution word
	finalSolutionWord := []*types.KeyCodeRecord{keycodes.UnknownKeyFromUser}
	//wantSolution = append(wantSolution, finalSolutionWord)
	wantSolution[wordCount-1] = finalSolutionWord
	// end with a long pause
	//milliSeconds[len(milliSeconds)-1] = wordPauseMS
	milliSeconds = append(milliSeconds, wordPauseMS)

	fmt.Printf("milliSeconds is %#v\n", milliSeconds)

	type args struct {
		milliSeconds  []int64
		wpm           uint64
		keyCodeStorer storer.KeyCodeStorer
	}
	tests := []struct {
		name         string
		args         args
		wantSolution [][]*types.KeyCodeRecord
		wantErr      bool
	}{
		{
			// TODO: Add test cases.
			name: "first",
			args: args{
				milliSeconds:  milliSeconds,
				wpm:           wpm,
				keyCodeStorer: keyCodeStore,
			},
			wantSolution: wantSolution,
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSolution, err := Copy(tt.args.milliSeconds, tt.args.wpm, tt.args.keyCodeStorer)
			if (err != nil) != tt.wantErr {
				t.Errorf("Copy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotSolution, tt.wantSolution) {
				t.Errorf("Copy() = %v, want %v", gotSolution, tt.wantSolution)
			}
		})
	}
}
