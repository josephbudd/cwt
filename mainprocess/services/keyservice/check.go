package keyservice

import (
	"github.com/josephbudd/cwt/domain/data/keycodes"
	"github.com/josephbudd/cwt/domain/implementations/storing/boltstoring"
	"github.com/josephbudd/cwt/domain/interfaces/storer"
	"github.com/josephbudd/cwt/domain/types"
)

// Check checks the user's keyed against the solutionChars and returns results.
func Check(keyed, solution [][]*types.KeyCodeRecord, keyCodeStorer storer.KeyCodeStorer, wpm uint64, recordResults bool) (nCorrect, nIncorrect, nRead uint64, testResults [][]types.TestResult, err error) {
	defer func() {
		if err != nil {
			return
		}
		if recordResults {
			err = recordCheckResults(keyCodeStorer, testResults, wpm)
		}
	}()

	var kRecord *types.KeyCodeRecord
	var sRecord *types.KeyCodeRecord
	var cRecord *types.KeyCodeRecord
	var testResultsLine []types.TestResult
	var solutionLine []*types.KeyCodeRecord
	var i, j int
	var m types.TestResult
	// get a list of record pointers
	var rr []*types.KeyCodeRecord
	if rr, err = keyCodeStorer.GetKeyCodes(); err != nil {
		return
	}
	controlIDRecord := make(map[uint64]*types.KeyCodeRecord, len(rr))
	for _, r := range rr {
		controlIDRecord[r.ID] = r
	}
	// keyed and solution should match if the user copied correctly.
	testResults = make([][]types.TestResult, 0, 100)
	lenKeyed := len(keyed)
	lenSolution := len(solution)
	if lenSolution > lenKeyed {
		// there are more solution lines than there is keyed lines.
		// iterate through the keyed lines to find the mistakes.
		// after all the keyed lines is checked iterate through the remaining solution lines and mark those as mistakes.
		testResultsLine = make([]types.TestResult, 0, lenSolution)
		var keyedLine []*types.KeyCodeRecord
		for i, keyedLine = range keyed {
			// solutionLine must be true records
			tempSolutionLine := solution[i]
			l := len(tempSolutionLine)
			solutionLine = make([]*types.KeyCodeRecord, l, l)
			for j = range tempSolutionLine {
				r := tempSolutionLine[j]
				solutionLine[j] = controlIDRecord[r.ID]
			}
			lenKeyedLine := len(keyedLine)
			lenSolutionLine := len(solutionLine)
			if lenKeyedLine > lenSolutionLine {
				nRead += uint64(lenKeyedLine)
				// the user added extra keys in this line.
				for j, sRecord = range solutionLine {
					kRecord = keyedLine[j]
					if kRecord.ID == sRecord.ID {
						nCorrect++
					} else {
						nIncorrect++
					}
					m = types.TestResult{
						Input:   kRecord,
						Control: sRecord,
					}
					testResultsLine = append(testResultsLine, m)
				}
				nIncorrect += uint64(lenKeyedLine - lenSolutionLine)
				for j = lenSolutionLine; j < lenKeyedLine; j++ {
					// extra keys in this line.
					cRecord = controlIDRecord[keyedLine[j].ID]
					m = types.TestResult{
						Input:   cRecord,
						Control: keycodes.NotInText,
					}
					testResultsLine = append(testResultsLine, m)
				}
			} else {
				nRead += uint64(lenSolutionLine)
				// the user did not add extra keys in this line.
				// there may not be enough keys.
				// lenKeyedLine <= lenSolutionLine
				// keyed
				for j, kRecord = range keyedLine {
					sRecord = solutionLine[j]
					if kRecord.ID == sRecord.ID {
						nCorrect++
					} else {
						nIncorrect++
					}
					m = types.TestResult{
						Input:   kRecord,
						Control: sRecord,
					}
					testResultsLine = append(testResultsLine, m)
				}
				nIncorrect += uint64(lenSolutionLine - lenKeyedLine)
				for j = lenKeyedLine; j < lenSolutionLine; j++ {
					// missing keys in this line.
					sRecord = solutionLine[j]
					m = types.TestResult{
						Input:   keycodes.NotKeyedByUser,
						Control: sRecord,
					}
					testResultsLine = append(testResultsLine, m)
				}
			}
			testResults = append(testResults, testResultsLine)
			testResultsLine = make([]types.TestResult, 0, lenSolution)
		}
		// there is more lines of solution but no more keyed.
		// each remaining solution represents a mis matched line.
		for i = lenKeyed; i < lenSolution; i++ {
			solutionLine = solution[i]
			lenSolutionLine := len(solutionLine)
			nRead += uint64(lenSolutionLine)
			nIncorrect += uint64(lenSolutionLine)
			for _, sRecord = range solutionLine {
				// not keyed by the user.
				m = types.TestResult{
					Control: sRecord,
					Input:   keycodes.NotKeyedByUser,
				}
				testResultsLine = append(testResultsLine, m)
			}
			testResults = append(testResults, testResultsLine)
			testResultsLine = make([]types.TestResult, 0, lenSolution)
		}
	} else {
		// lenSolution <= lenKeyed
		// there is more keyed lines than there are solution lines.
		// iterate through the solutions and find the mistakes in the keyed.
		// after the solutions are checked mark the rest of the keyed as mistakes.
		testResultsLine = make([]types.TestResult, 0, lenKeyed)
		for i, solutionLine = range solution {
			keyedLine := keyed[i]
			lenKeyedLine := len(keyedLine)
			lenSolutionLine := len(solutionLine)
			if lenKeyedLine > lenSolutionLine {
				nRead += uint64(lenKeyedLine)
				// the user keyed extra chars in this line.
				for j, sRecord = range solutionLine {
					// text and keys.
					kRecord := keyedLine[j]
					if kRecord.ID == sRecord.ID {
						nCorrect++
					} else {
						nIncorrect++
					}
					m = types.TestResult{
						Input:   kRecord,
						Control: sRecord,
					}
					testResultsLine = append(testResultsLine, m)
				}
				nIncorrect += uint64(lenKeyedLine - lenSolutionLine)
				for j = lenSolutionLine; j < lenKeyedLine; j++ {
					// keyed but there was no text to key in this line.
					kRecord = keyedLine[j]
					m = types.TestResult{
						Input:   kRecord,
						Control: keycodes.NotInText,
					}
					testResultsLine = append(testResultsLine, m)
				}
			} else {
				// lenKeyedLine <= lenSolutionLine
				// the user keyed the correct number or too few this line.
				nRead += uint64(lenSolutionLine)
				for j, kRecord = range keyedLine {
					// text and keys.
					sRecord = solutionLine[j]
					if kRecord.ID == sRecord.ID {
						nCorrect++
					} else {
						nIncorrect++
					}
					m = types.TestResult{
						Input:   kRecord,
						Control: sRecord,
					}
					testResultsLine = append(testResultsLine, m)
				}
				nIncorrect += uint64(lenSolutionLine - lenKeyedLine)
				for j = lenKeyedLine; j < lenSolutionLine; j++ {
					// text but the user missed keys for this line.
					sRecord = solutionLine[j]
					m = types.TestResult{
						Input:   keycodes.NotKeyedByUser,
						Control: sRecord,
					}
					testResultsLine = append(testResultsLine, m)
				}
			}
			testResults = append(testResults, testResultsLine)
			testResultsLine = make([]types.TestResult, 0, lenSolution)
		}
		i++
		if i == lenKeyed {
			// no more keyed
			return
		}
		// extra keyed lines are mistakes.
		for ; i < lenKeyed; i++ {
			keyedLine := keyed[i]
			nIncorrect += uint64(len(keyedLine))
			for _, kRecord = range keyedLine {
				m = types.TestResult{
					Input:   kRecord,
					Control: keycodes.NotInText,
				}
				testResultsLine = append(testResultsLine, m)
			}
		}
		testResults = append(testResults, testResultsLine)
	}
	return
}

func recordCheckResults(keyCodeStorer storer.KeyCodeStorer, testResults [][]types.TestResult, wpm uint64) (err error) {
	var wpmAC types.KeyCodeRecordResult
	var cRecord *types.KeyCodeRecord
	idRecord := make(map[uint64]*types.KeyCodeRecord, 100)
	for _, testResultsLine := range testResults {
		for _, testResult := range testResultsLine {
			if testResult.Control.ID >= boltstoring.FirstValidID {
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
	var r *types.KeyCodeRecord
	for _, r = range idRecord {
		if err = keyCodeStorer.UpdateKeyCode(r); err != nil {
			return
		}
	}
	return
}
