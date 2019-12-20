package keyservice

import (
	"github.com/josephbudd/cwt/domain/data"
	"github.com/josephbudd/cwt/domain/data/keycodes"
	"github.com/josephbudd/cwt/domain/store/record"
	"github.com/josephbudd/cwt/domain/store/storer"
	"github.com/josephbudd/cwt/domain/store/storing"
)

// Check checks the user's keyed against the solutionChars and returns results.
func Check(keyed, solution [][]*record.KeyCode, keyCodeStorer storer.KeyCodeStorer, wpm uint64, recordResults bool) (nCorrect, nIncorrect, nPossible uint64, testResults [][]data.TestResult, err error) {

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
	if len(keyed) > len(solution) {
		var i int
		var line []*record.KeyCode
		for i, line = range solution {
			ls := len(line)
			lc := len(keyed[i])
			if ls > lc {
				total += ls
			} else {
				total += lc
			}
		}
		for i++; i < len(keyed); i++ {
			total += len(keyed[i])
		}
	} else {
		var i int
		var line []*record.KeyCode
		for i, line = range keyed {
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

	var kRecord *record.KeyCode
	var sRecord *record.KeyCode
	var cRecord *record.KeyCode
	var testResultsLine []data.TestResult
	var solutionLine []*record.KeyCode
	var i, j int
	var m data.TestResult
	// get a list of record pointers
	var rr []*record.KeyCode
	if rr, err = keyCodeStorer.GetAll(); err != nil {
		return
	}
	controlIDRecord := make(map[uint64]*record.KeyCode, len(rr))
	for _, r := range rr {
		controlIDRecord[r.ID] = r
	}
	// keyed and solution should match if the user copied correctly.
	testResults = make([][]data.TestResult, 0, 100)
	lenKeyed := len(keyed)
	lenSolution := len(solution)
	if lenSolution > lenKeyed {
		// there are more solution lines than there is keyed lines.
		// iterate through the keyed lines to find the mistakes.
		// after all the keyed lines is checked iterate through the remaining solution lines and mark those as mistakes.
		testResultsLine = make([]data.TestResult, 0, lenSolution)
		var keyedLine []*record.KeyCode
		for i, keyedLine = range keyed {
			// solutionLine must be true records
			tempSolutionLine := solution[i]
			l := len(tempSolutionLine)
			solutionLine = make([]*record.KeyCode, l, l)
			for j = range tempSolutionLine {
				r := tempSolutionLine[j]
				solutionLine[j] = controlIDRecord[r.ID]
			}
			lenKeyedLine := len(keyedLine)
			lenSolutionLine := len(solutionLine)
			if lenKeyedLine > lenSolutionLine {
				// the user added extra keys in this line.
				for j, sRecord = range solutionLine {
					kRecord = keyedLine[j]
					if kRecord.ID == sRecord.ID {
						nCorrect++
					}
					m = data.TestResult{
						Input:   kRecord,
						Control: sRecord,
					}
					testResultsLine = append(testResultsLine, m)
				}
				for j = lenSolutionLine; j < lenKeyedLine; j++ {
					// extra keys in this line.
					cRecord = controlIDRecord[keyedLine[j].ID]
					m = data.TestResult{
						Input:   cRecord,
						Control: keycodes.NoCopyToKey,
					}
					testResultsLine = append(testResultsLine, m)
				}
			} else {
				// the user did not add extra keys in this line.
				// there may not be enough keys.
				// lenKeyedLine <= lenSolutionLine
				// keyed
				for j, kRecord = range keyedLine {
					sRecord = solutionLine[j]
					if kRecord.ID == sRecord.ID {
						nCorrect++
					}
					m = data.TestResult{
						Input:   kRecord,
						Control: sRecord,
					}
					testResultsLine = append(testResultsLine, m)
				}
				for j = lenKeyedLine; j < lenSolutionLine; j++ {
					// missing keys in this line.
					sRecord = solutionLine[j]
					m = data.TestResult{
						Input:   keycodes.NotKeyedByUser,
						Control: sRecord,
					}
					testResultsLine = append(testResultsLine, m)
				}
			}
			testResults = append(testResults, testResultsLine)
			testResultsLine = make([]data.TestResult, 0, lenSolution)
		}
		// there is more lines of solution but no more keyed.
		// each remaining solution represents a mis matched line.
		for i = lenKeyed; i < lenSolution; i++ {
			solutionLine = solution[i]
			for _, sRecord = range solutionLine {
				// not keyed by the user.
				m = data.TestResult{
					Control: sRecord,
					Input:   keycodes.NotKeyedByUser,
				}
				testResultsLine = append(testResultsLine, m)
			}
			testResults = append(testResults, testResultsLine)
			testResultsLine = make([]data.TestResult, 0, lenSolution)
		}
	} else {
		// lenSolution <= lenKeyed
		// there is more keyed lines than there are solution lines.
		// iterate through the solutions and find the mistakes in the keyed.
		// after the solutions are checked mark the rest of the keyed as mistakes.
		testResultsLine = make([]data.TestResult, 0, lenKeyed)
		for i, solutionLine = range solution {
			keyedLine := keyed[i]
			lenKeyedLine := len(keyedLine)
			lenSolutionLine := len(solutionLine)
			if lenKeyedLine > lenSolutionLine {
				// the user keyed extra chars in this line.
				for j, sRecord = range solutionLine {
					// text and keys.
					kRecord := keyedLine[j]
					if kRecord.ID == sRecord.ID {
						nCorrect++
					}
					m = data.TestResult{
						Input:   kRecord,
						Control: sRecord,
					}
					testResultsLine = append(testResultsLine, m)
				}
				for j = lenSolutionLine; j < lenKeyedLine; j++ {
					// keyed but there was no text to key in this line.
					kRecord = keyedLine[j]
					m = data.TestResult{
						Input:   kRecord,
						Control: keycodes.NoCopyToKey,
					}
					testResultsLine = append(testResultsLine, m)
				}
			} else {
				// lenKeyedLine <= lenSolutionLine
				// the user keyed the correct number or too few this line.
				for j, kRecord = range keyedLine {
					// text and keys.
					sRecord = solutionLine[j]
					if kRecord.ID == sRecord.ID {
						nCorrect++
					}
					m = data.TestResult{
						Input:   kRecord,
						Control: sRecord,
					}
					testResultsLine = append(testResultsLine, m)
				}
				for j = lenKeyedLine; j < lenSolutionLine; j++ {
					// text but the user missed keys for this line.
					sRecord = solutionLine[j]
					m = data.TestResult{
						Input:   keycodes.NotKeyedByUser,
						Control: sRecord,
					}
					testResultsLine = append(testResultsLine, m)
				}
			}
			testResults = append(testResults, testResultsLine)
			testResultsLine = make([]data.TestResult, 0, lenSolution)
		}
		i++
		if i == lenKeyed {
			// no more keyed
			return
		}
		// extra keyed lines are mistakes.
		for ; i < lenKeyed; i++ {
			keyedLine := keyed[i]
			for _, kRecord = range keyedLine {
				m = data.TestResult{
					Input:   kRecord,
					Control: keycodes.NoCopyToKey,
				}
				testResultsLine = append(testResultsLine, m)
			}
		}
		testResults = append(testResults, testResultsLine)
	}
	return
}

func recordCheckResults(keyCodeStorer storer.KeyCodeStorer, testResults [][]data.TestResult, wpm uint64) (err error) {
	var wpmAC record.KeyCodeResult
	var cRecord *record.KeyCode
	idRecord := make(map[uint64]*record.KeyCode, 100)
	for _, testResultsLine := range testResults {
		for _, testResult := range testResultsLine {
			if testResult.Control.ID >= storing.FirstValidID {
				if _, found := idRecord[testResult.Control.ID]; !found {
					idRecord[testResult.Control.ID] = testResult.Control
				}
				cRecord = idRecord[testResult.Control.ID]
				wpmAC = cRecord.KeyWPMResults[wpm]
				wpmAC.Attempts++
				if testResult.Control.ID == testResult.Input.ID {
					wpmAC.Correct++
				}
				cRecord.KeyWPMResults[wpm] = wpmAC
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
