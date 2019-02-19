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
	name           string
	checkKeyArgs   checkKeyArgs
	wantNCorrect   uint64
	wantNIncorrect uint64
	wantNRead      uint64
	wantTestResults [][]types.TestResult
	wantErr        bool
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
	var keyCodeRecords []*types.KeyCodeRecord
	if keyCodeRecords, err = getKeyCodes(); err != nil {
		t.Fatal(err)
	}
	okCheckTest(keyCodeRecords, t)
	missingKeysCheckTest(keyCodeRecords, t)
	extraKeysCheckTest(keyCodeRecords, t)
}

func extraKeysCheckTest(keyCodeRecords []*types.KeyCodeRecord, t *testing.T) {
	wordLength := 5
	wordCount := 1
	solution := make([][]*types.KeyCodeRecord, wordCount, wordCount)
	keyed := make([][]*types.KeyCodeRecord, wordCount, wordCount)
	testResults := make([][]types.TestResult, wordCount, wordCount)
	solution[0] = keyCodeRecords[:4]
	keyed[0] = keyCodeRecords[:5]
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
			wantNCorrect:   uint64(wordLength - 1),
			wantNIncorrect: 1,
			wantNRead:      uint64(wordLength),
			wantTestResults: testResults,
		},
	}
	testCheck(tests, t)
}

func missingKeysCheckTest(keyCodeRecords []*types.KeyCodeRecord, t *testing.T) {
	wordLength := 5
	wordCount := 1
	solution := make([][]*types.KeyCodeRecord, wordCount, wordCount)
	keyed := make([][]*types.KeyCodeRecord, wordCount, wordCount)
	testResults := make([][]types.TestResult, wordCount, wordCount)
	solution[0] = keyCodeRecords[:5]
	keyed[0] = keyCodeRecords[:4]
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
			wantNCorrect:   uint64(wordLength - 1),
			wantNIncorrect: 1,
			wantNRead:      uint64(wordLength),
			wantTestResults: testResults,
		},
	}
	testCheck(tests, t)
}

func okCheckTest(keyCodeRecords []*types.KeyCodeRecord, t *testing.T) {
	// pauses
	// make the words 5 characters long.
	wordLength := 5
	wordCount := len(keyCodeRecords) / wordLength
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
			wantNCorrect:   sizeSolution,
			wantNIncorrect: 0,
			wantNRead:      sizeSolution,
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
				Control: keycodes.NotInText,
			}
		}
	}
	return
}

func testCheck(tests []checkKeyData, t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNCorrect, gotNIncorrect, gotNRead, gotTestResults, err := Check(tt.checkKeyArgs.keyed, tt.checkKeyArgs.solution, tt.checkKeyArgs.keyCodeStorer, tt.checkKeyArgs.wpm, tt.checkKeyArgs.recordResults)
			if (err != nil) != tt.wantErr {
				t.Errorf("Check() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotNCorrect != tt.wantNCorrect {
				t.Errorf("Check() gotNCorrect = %v, want %v", gotNCorrect, tt.wantNCorrect)
			}
			if gotNIncorrect != tt.wantNIncorrect {
				t.Errorf("Check() gotNIncorrect = %v, want %v", gotNIncorrect, tt.wantNIncorrect)
			}
			if gotNRead != tt.wantNRead {
				t.Errorf("Check() gotNRead = %v, want %v", gotNRead, tt.wantNRead)
			}
			if !reflect.DeepEqual(gotTestResults, tt.wantTestResults) {
				t.Errorf("Check() gotTestResults = %v, want %v", gotTestResults, tt.wantTestResults)
			}
		})
	}
}
