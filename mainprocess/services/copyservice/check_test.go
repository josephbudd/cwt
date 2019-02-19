package copyservice

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/boltdb/bolt"
	"github.com/josephbudd/cwt/domain/data/filepaths"
	"github.com/josephbudd/cwt/domain/data/keycodes"
	"github.com/josephbudd/cwt/domain/implementations/storing/boltstoring"
	"github.com/josephbudd/cwt/domain/interfaces/storer"
	"github.com/josephbudd/cwt/domain/types"
	"github.com/pkg/errors"
)

type checkCopyArgs struct {
	copy          [][]*types.KeyCodeRecord
	solution      [][]*types.KeyCodeRecord
	keyCodeStorer storer.KeyCodeStorer
	wpm           uint64
	recordResults bool
}

type checkCopyData struct {
	name           string
	checkCopyArgs  checkCopyArgs
	wantNCorrect   uint64
	wantNIncorrect uint64
	wantNKeyed     uint64
	wantTestResults [][]types.TestResult
	wantErr        bool
}

var (
	keyCodeStore storer.KeyCodeStorer
)

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
	wordCount := 1
	solution := make([][]*types.KeyCodeRecord, wordCount, wordCount)
	copied := make([][]*types.KeyCodeRecord, wordCount, wordCount)
	testResults := make([][]types.TestResult, wordCount, wordCount)
	solution[0] = keyCodeRecords[:4]
	copied[0] = keyCodeRecords[:5]
	testResults[0] = makeResultsLine(solution[0], copied[0])
	tests := []checkCopyData{
		// TODO: Add test cases.
		{
			name: "extra keys",
			checkCopyArgs: checkCopyArgs{
				copy:          copied,
				solution:      solution,
				keyCodeStorer: keyCodeStore,
			},
			wantNCorrect:   uint64(len(solution[0])),
			wantNIncorrect: 1,
			wantNKeyed:     uint64(len(solution[0])),
			wantTestResults: testResults,
		},
	}
	testCheck(tests, t)
}

func missingKeysCheckTest(keyCodeRecords []*types.KeyCodeRecord, t *testing.T) {
	wordCount := 1
	solution := make([][]*types.KeyCodeRecord, wordCount, wordCount)
	copied := make([][]*types.KeyCodeRecord, wordCount, wordCount)
	testResults := make([][]types.TestResult, wordCount, wordCount)
	solution[0] = keyCodeRecords[:5]
	copied[0] = keyCodeRecords[:4]
	testResults[0] = makeResultsLine(solution[0], copied[0])
	tests := []checkCopyData{
		// TODO: Add test cases.
		{
			name: "missing keys",
			checkCopyArgs: checkCopyArgs{
				copy:          copied,
				solution:      solution,
				keyCodeStorer: keyCodeStore,
			},
			wantNCorrect:   uint64(len(copied[0])),
			wantNIncorrect: 1,
			wantNKeyed:     uint64(len(solution[0])),
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
	testResults := make([][]types.TestResult, wordCount, wordCount)
	solution := make([][]*types.KeyCodeRecord, wordCount, wordCount)
	var i int
	for i = 0; i < wordCount; i++ {
		start := i * wordLength
		end := start + wordLength
		solutionLine := keyCodeRecords[start:end]
		solution[i] = solutionLine
		testResults[i] = makeResultsLine(solutionLine, solutionLine)
	}
	sizeSolution := uint64(wordLength * i)

	tests := []checkCopyData{
		// TODO: Add test cases.
		{
			name: "ok",
			checkCopyArgs: checkCopyArgs{
				copy:          solution,
				solution:      solution,
				keyCodeStorer: keyCodeStore,
			},
			wantNCorrect:   sizeSolution,
			wantNIncorrect: 0,
			wantNKeyed:     sizeSolution,
			wantTestResults: testResults,
		},
	}
	testCheck(tests, t)
}

func testCheck(tests []checkCopyData, t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNCorrect, gotNIncorrect, gotNKeyed, gotTestResults, err := Check(tt.checkCopyArgs.copy, tt.checkCopyArgs.solution, tt.checkCopyArgs.keyCodeStorer, tt.checkCopyArgs.wpm, tt.checkCopyArgs.recordResults)
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
			if gotNKeyed != tt.wantNKeyed {
				t.Errorf("Check() gotNKeyed = %v, want %v", gotNKeyed, tt.wantNKeyed)
			}
			if !reflect.DeepEqual(gotTestResults, tt.wantTestResults) {
				t.Errorf("Check() gotTestResults = %v, want %v", gotTestResults, tt.wantTestResults)
			}
		})
	}
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

func getKeyCodes() (keyCodeRecords []*types.KeyCodeRecord, err error) {
	if keyCodeRecords, err = keyCodeStore.GetKeyCodes(); err != nil {
		err = errors.WithMessage(err, `keyCodeStore.GetKeyCodes()`)
	}
	return
}

func makeResultsLine(solution, copied []*types.KeyCodeRecord) (results []types.TestResult) {
	lenSolution := len(solution)
	lenKeyed := len(copied)
	switch {
	case lenSolution == lenKeyed:
		results = make([]types.TestResult, lenSolution, lenSolution)
		for i := 0; i < lenKeyed; i++ {
			results[i] = types.TestResult{
				Input:   copied[i],
				Control: solution[i],
			}
		}
	case lenSolution > lenKeyed:
		results = make([]types.TestResult, lenSolution, lenSolution)
		for i := 0; i < lenKeyed; i++ {
			results[i] = types.TestResult{
				Input:   copied[i],
				Control: solution[i],
			}
		}
		// these uncopied are bad
		for i := lenKeyed; i < lenSolution; i++ {
			results[i] = types.TestResult{
				Input:   keycodes.NotCopiedByUser,
				Control: solution[i],
			}
		}
	case lenSolution < lenKeyed:
		results = make([]types.TestResult, lenKeyed, lenKeyed)
		for i := 0; i < lenSolution; i++ {
			results[i] = types.TestResult{
				Input:   copied[i],
				Control: solution[i],
			}
		}
		// these copied are bad
		for i := lenSolution; i < lenKeyed; i++ {
			results[i] = types.TestResult{
				Input:   copied[i],
				Control: keycodes.NotKeyedByApp,
			}
		}
	}
	return
}
