package copyservice

import (
	"errors"
	"fmt"

	"github.com/josephbudd/cwt/domain/data"
	"github.com/josephbudd/cwt/domain/data/keycodes"
	"github.com/josephbudd/cwt/domain/store/record"
	"github.com/josephbudd/cwt/domain/store/storer"
	"github.com/josephbudd/cwt/domain/store/storing"
)

// Check checks the user's copy agains the solutionChars and returns results.
func Check(copy [][]*record.KeyCode, solution [][]*record.KeyCode, keyCodeStorer storer.KeyCodeStorer, wpm uint64, recordResults bool) (nCorrect, nIncorrect, nPossible uint64, testResults [][]data.TestResult, err error) {

	var total int
	defer func() {
		if err != nil {
			return
		}
		if recordResults {
			err = recordCheckResults(keyCodeStorer, testResults, wpm)
		}
		if err == nil {
			nIncorrect = uint64(total) - nCorrect
			nPossible = uint64(total)
		}
	}()

	// calc total
	if len(copy) > len(solution) {
		var i int
		var line []*record.KeyCode
		for i, line = range solution {
			ls := len(line)
			lc := len(copy[i])
			if ls > lc {
				total += ls
			} else {
				total += lc
			}
		}
		for i++; i < len(copy); i++ {
			total += len(copy[i])
		}
	} else {
		var i int
		var line []*record.KeyCode
		for i, line = range copy {
			lc := len(line)
			ls := len(solution[i])
			if ls > lc {
				total += ls
			} else {
				total += lc
			}
		}
		for i++; i < len(solution); i++ {
			total += len(solution[i])
		}
	}

	// get a list of record pointers
	var rr []*record.KeyCode
	if rr, err = keyCodeStorer.GetAll(); err != nil {
		return
	}
	controlIDRecord := make(map[uint64]*record.KeyCode, len(rr))
	for _, r := range rr {
		controlIDRecord[r.ID] = r
	}
	testResults = make([][]data.TestResult, 0, 100)
	lc := len(copy)
	ls := len(solution)
	if ls > lc {
		// there are more solutions than there is copy.
		// iterate through the copy to find the mistakes.
		// after all the copy is checked iterate through the remaining solutions and mark those as mistakes.
		var i int
		testResultsLine := make([]data.TestResult, 0, ls)
		var copyLine []*record.KeyCode
		for i, copyLine = range copy {
			// solutionLine must be true records
			tempSolutionLine := solution[i]
			l := len(tempSolutionLine)
			solutionLine := make([]*record.KeyCode, l, l)
			for j := range tempSolutionLine {
				r := tempSolutionLine[j]
				solutionLine[j] = controlIDRecord[r.ID]
			}
			lcl := len(copyLine)
			lsl := len(solutionLine)
			if lcl > lsl {
				var j int
				var sRecord *record.KeyCode
				for j, sRecord = range solutionLine {
					cRecord := copyLine[j]
					if cRecord.ID == sRecord.ID {
						nCorrect++
					}
					m := data.TestResult{
						Input:   cRecord,
						Control: sRecord,
					}
					testResultsLine = append(testResultsLine, m)
				}
				for j = lsl; j < lcl; j++ {
					cRecord := copyLine[j]
					m := data.TestResult{
						Input:   cRecord,
						Control: keycodes.NotKeyedByApp,
					}
					testResultsLine = append(testResultsLine, m)
				}
			} else {
				// lcl <= lsl
				var j int
				var cRecord *record.KeyCode
				for j, cRecord = range copyLine {
					// copied
					sRecord := solutionLine[j]
					if cRecord.ID == sRecord.ID {
						nCorrect++
					}
					m := data.TestResult{
						Input:   cRecord,
						Control: sRecord,
					}
					testResultsLine = append(testResultsLine, m)
				}
				for j = lcl; j < lsl; j++ {
					// missing copy
					sRecord := solutionLine[j]
					m := data.TestResult{
						Input:   keycodes.NotCopiedByUser,
						Control: sRecord,
					}
					testResultsLine = append(testResultsLine, m)
				}
			}
			testResults = append(testResults, testResultsLine)
			testResultsLine = make([]data.TestResult, 0, ls)
		}
		// there is more solution but no more copy.
		// each remaining solution represents a mis match.
		for i++; i < ls; i++ {
			// solutionLine must be true records
			tempSolutionLine := solution[i]
			l := len(tempSolutionLine)
			solutionLine := make([]*record.KeyCode, l, l)
			for i := range tempSolutionLine {
				r := tempSolutionLine[i]
				solutionLine[i] = controlIDRecord[r.ID]
			}
			for _, sRecord := range solutionLine {
				// not copied
				m := data.TestResult{
					Control: sRecord,
					Input:   keycodes.NotCopiedByUser,
				}
				testResultsLine = append(testResultsLine, m)
			}
			testResults = append(testResults, testResultsLine)
			testResultsLine = make([]data.TestResult, 0, ls)
		}
	} else {
		// lc >= ls {
		// there is more copy than there is correct solutions.
		// iterate through the solutions and find the mistakes in the copy.
		// after the solutions are checked mark the rest of the copy as mistakes.
		var i int
		testResultsLine := make([]data.TestResult, 0, lc)
		var tempSolutionLine []*record.KeyCode
		for i, tempSolutionLine = range solution {
			// solutionLine must be true records
			l := len(tempSolutionLine)
			solutionLine := make([]*record.KeyCode, l, l)
			for j := range tempSolutionLine {
				r := tempSolutionLine[j]
				solutionLine[j] = controlIDRecord[r.ID]
			}
			copyLine := copy[i]
			lcl := len(copyLine)
			lsl := len(solutionLine)
			if lcl > lsl {
				var j int
				var sRecord *record.KeyCode
				for j, sRecord = range solutionLine {
					// copied
					cRecord := copyLine[j]
					if cRecord.ID == sRecord.ID {
						nCorrect++
					}
					m := data.TestResult{
						Input:   cRecord,
						Control: sRecord,
					}
					testResultsLine = append(testResultsLine, m)
				}
				for j = lsl; j < lcl; j++ {
					// copied but nothing keyed
					cRecord := copyLine[j]
					m := data.TestResult{
						Input:   cRecord,
						Control: keycodes.NotKeyedByApp,
					}
					testResultsLine = append(testResultsLine, m)
				}
			} else {
				// lcl <= lsl
				var j int
				var cRecord *record.KeyCode
				for j, cRecord = range copyLine {
					// keyed and copied
					sRecord := solutionLine[j]
					if cRecord.Character == sRecord.Character {
						nCorrect++
					}
					m := data.TestResult{
						Input:   cRecord,
						Control: sRecord,
					}
					testResultsLine = append(testResultsLine, m)
				}
				for j = lcl; j < lsl; j++ {
					// keyed but never copied
					sRecord := solutionLine[j]
					m := data.TestResult{
						Input:   keycodes.NotCopiedByUser,
						Control: sRecord,
					}
					testResultsLine = append(testResultsLine, m)
				}
			}
			testResults = append(testResults, testResultsLine)
			testResultsLine = make([]data.TestResult, 0, ls)
		}
		if ls == lc {
			// no more copy
			return
		}
		// more copy and the rest of the copy is mistakes.
		for i = ls; i < lc; i++ {
			copyLine := copy[i]
			for _, cRecord := range copyLine {
				m := data.TestResult{
					Input:   cRecord,
					Control: keycodes.NotKeyedByApp,
				}
				testResultsLine = append(testResultsLine, m)
			}
		}
		testResults = append(testResults, testResultsLine)
	}
	return
}

func recordCheckResults(keyCodeStorer storer.KeyCodeStorer, testResults [][]data.TestResult, wpm uint64) (err error) {
	idRecord := make(map[uint64]*record.KeyCode, 100)
	for _, mm := range testResults {
		for _, m := range mm {
			if m.Control.ID >= storing.FirstValidID {
				var results record.KeyCodeResult
				var found bool
				if results, found = m.Control.CopyWPMResults[wpm]; !found {
					message := fmt.Sprintf("wpm %d is invalid", wpm)
					err = errors.New(message)
				}
				results.Attempts++
				if m.Control.ID == m.Input.ID {
					results.Correct++
				}
				m.Control.CopyWPMResults[wpm] = results
				idRecord[m.Control.ID] = m.Control
			}
		}
	}
	var r *record.KeyCode
	for _, r = range idRecord {
		if err = keyCodeStorer.Update(r); err != nil {
			return
		}
	}
	return
}
