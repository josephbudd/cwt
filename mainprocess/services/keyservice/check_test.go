package keyservice

import (
	"reflect"
	"testing"

	"github.com/boltdb/bolt"
	"github.com/josephbudd/cwt/domain/data/filepaths"
	"github.com/josephbudd/cwt/domain/data/keycodes"
	"github.com/josephbudd/cwt/domain/implementations/storing/boltstoring"
	"github.com/josephbudd/cwt/domain/interfaces/storer"
	"github.com/josephbudd/cwt/domain/types"
)

type checkKeyArgs struct {
	keyed         [][]*types.KeyCodeRecord
	solution      [][]*types.KeyCodeRecord
	keyCodeStorer storer.KeyCodeStorer
	wpm           uint64
	recordResults bool
}

type checkKeyData struct {
	name            string
	checkKeyArgs    checkKeyArgs
	wantNCorrect    uint64
	wantNIncorrect  uint64
	wantNRead       uint64
	wantTestResults [][]types.TestResult
	wantErr         bool
}

func TestCheck(t *testing.T) {
	var path string
	var db *bolt.DB
	var err error
	if path, db, err = buildBoltStores(); err != nil {
		t.Fatal(err)
	}
	keyCodeStore = boltstoring.NewKeyCodeBoltDB(db, path, filepaths.GetFmode())
	defer keyCodeStore.Close()
	var keyCodeWords [][]*types.KeyCodeRecord
	if keyCodeWords, err = getKeyCodes(keyCodeStore); err != nil {
		t.Fatal(err)
	}
	var keyCodeRecords []*types.KeyCodeRecord
	for _, keyCodeRecords = range keyCodeWords {
		okCheckTest(keyCodeRecords, t)
		missingKeysCheckTest(keyCodeRecords, t)
		extraKeysCheckTest(keyCodeRecords, t)
	}
}

func extraKeysCheckTest(keyCodeRecords []*types.KeyCodeRecord, t *testing.T) {
	wordLength := 5
	l := len(keyCodeRecords)
	if wordLength > l {
		wordLength = l
	}
	wordCount := 1
	solution := make([][]*types.KeyCodeRecord, wordCount, wordCount)
	keyed := make([][]*types.KeyCodeRecord, wordCount, wordCount)
	testResults := make([][]types.TestResult, wordCount, wordCount)
	solution[0] = keyCodeRecords[:wordLength-1]
	keyed[0] = keyCodeRecords[:wordLength]
	testResults[0] = makeResultsLine(solution[0], keyed[0])
	tests := []checkKeyData{
		// TODO: Add test cases.
		{
			name: "extra keys",
			checkKeyArgs: checkKeyArgs{
				keyed:         keyed,
				solution:      solution,
				keyCodeStorer: keyCodeStore,
			},
			wantNCorrect:    uint64(wordLength - 1),
			wantNIncorrect:  1,
			wantNRead:       uint64(wordLength),
			wantTestResults: testResults,
		},
	}
	testCheck(tests, t)
}

func missingKeysCheckTest(keyCodeRecords []*types.KeyCodeRecord, t *testing.T) {
	wordLength := 5
	l := len(keyCodeRecords)
	if wordLength > l {
		wordLength = l
	}
	wordCount := 1
	solution := make([][]*types.KeyCodeRecord, wordCount, wordCount)
	keyed := make([][]*types.KeyCodeRecord, wordCount, wordCount)
	testResults := make([][]types.TestResult, wordCount, wordCount)
	solution[0] = keyCodeRecords[:wordLength]
	keyed[0] = keyCodeRecords[:wordLength-1]
	testResults[0] = makeResultsLine(solution[0], keyed[0])
	tests := []checkKeyData{
		// TODO: Add test cases.
		{
			name: "missing keys",
			checkKeyArgs: checkKeyArgs{
				keyed:         keyed,
				solution:      solution,
				keyCodeStorer: keyCodeStore,
			},
			wantNCorrect:    uint64(wordLength - 1),
			wantNIncorrect:  1,
			wantNRead:       uint64(wordLength),
			wantTestResults: testResults,
		},
	}
	testCheck(tests, t)
}

func okCheckTest(keyCodeRecords []*types.KeyCodeRecord, t *testing.T) {
	// pauses
	// make the words 5 characters long.
	wordLength := 5
	l := len(keyCodeRecords)
	if wordLength > l {
		wordLength = l
	}
	wordCount := l / wordLength
	wantTestResults := make([][]types.TestResult, wordCount, wordCount)
	wantSolution := make([][]*types.KeyCodeRecord, wordCount, wordCount)
	var i int
	for i = 0; i < wordCount; i++ {
		start := i * wordLength
		end := start + wordLength
		solutionLine := keyCodeRecords[start:end]
		wantSolution[i] = solutionLine
		wantTestResults[i] = makeResultsLine(solutionLine, solutionLine)
	}
	sizeSolution := uint64(wordLength * i)

	tests := []checkKeyData{
		// TODO: Add test cases.
		{
			name: "ok",
			checkKeyArgs: checkKeyArgs{
				keyed:         wantSolution,
				solution:      wantSolution,
				keyCodeStorer: keyCodeStore,
			},
			wantNCorrect:    sizeSolution,
			wantNIncorrect:  0,
			wantNRead:       sizeSolution,
			wantTestResults: wantTestResults,
		},
	}
	testCheck(tests, t)
}

func makeResultsLine(solution, keyed []*types.KeyCodeRecord) (results []types.TestResult) {
	lenSolution := len(solution)
	lenKeyed := len(keyed)
	switch {
	case lenSolution == lenKeyed:
		results = make([]types.TestResult, lenSolution, lenSolution)
		for i := 0; i < lenKeyed; i++ {
			results[i] = types.TestResult{
				Input:   keyed[i],
				Control: solution[i],
			}
		}
	case lenSolution > lenKeyed:
		results = make([]types.TestResult, lenSolution, lenSolution)
		for i := 0; i < lenKeyed; i++ {
			results[i] = types.TestResult{
				Input:   keyed[i],
				Control: solution[i],
			}
		}
		// these unkeyed are bad
		for i := lenKeyed; i < lenSolution; i++ {
			results[i] = types.TestResult{
				Input:   keycodes.NotKeyedByUser,
				Control: solution[i],
			}
		}
	case lenSolution < lenKeyed:
		results = make([]types.TestResult, lenKeyed, lenKeyed)
		for i := 0; i < lenSolution; i++ {
			results[i] = types.TestResult{
				Input:   keyed[i],
				Control: solution[i],
			}
		}
		// these keyed are bad
		for i := lenSolution; i < lenKeyed; i++ {
			results[i] = types.TestResult{
				Input:   keyed[i],
				Control: keycodes.NoCopyToKey,
			}
		}
	}
	return
}

func testCheck(tests []checkKeyData, t *testing.T) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotNCorrect, gotNIncorrect, gotNRead, gotTestResults, err := Check(test.checkKeyArgs.keyed, test.checkKeyArgs.solution, test.checkKeyArgs.keyCodeStorer, test.checkKeyArgs.wpm, test.checkKeyArgs.recordResults)
			if (err != nil) != test.wantErr {
				t.Errorf("%s: Check() error = %v, wantErr %v", test.name, err, test.wantErr)
				return
			}
			if gotNCorrect != test.wantNCorrect {
				t.Errorf("%s: Check() gotNCorrect = %v, want %v", test.name, gotNCorrect, test.wantNCorrect)
			}
			if gotNIncorrect != test.wantNIncorrect {
				t.Errorf("%s: Check() gotNIncorrect = %v, want %v", test.name, gotNIncorrect, test.wantNIncorrect)
			}
			if gotNRead != test.wantNRead {
				t.Errorf("%s: Check() gotNRead = %v, want %v", test.name, gotNRead, test.wantNRead)
			}
			if !reflect.DeepEqual(gotTestResults, test.wantTestResults) {
				t.Errorf("%s: Check() gotTestResults = %v, want %v", test.name, gotTestResults, test.wantTestResults)
			}
		})
	}
}
