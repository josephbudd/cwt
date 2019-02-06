package keyService

import (
	"fmt"

	"github.com/josephbudd/cwt/domain/data/keycodes"
	"github.com/josephbudd/cwt/domain/interfaces/storer"
	"github.com/josephbudd/cwt/domain/types"
)

// Check checks the user's keyed against the solutionChars and returns results.
func Check(keyed [][]*types.KeyCodeRecord, solution [][]*types.KeyCodeRecord, keyCodeStorer storer.KeyCodeStorer, wpm uint64, recordResults bool) (nCorrect, nIncorrect, nRead uint64, misMatches [][]types.MisMatch, err error) {
	defer func() {
		if err != nil {
			return
		}
		if recordResults {
			err = recordCheckResults(keyCodeStorer, misMatches, wpm)
		}
	}()
	// now keyed is a slice of strings and solutionChars is a slice of strings.
	// the 2 should match is the user copied correctly.
	misMatches = make([][]types.MisMatch, 0, 100)
	lenKeyed := len(keyed)
	lenSolution := len(solution)
	fmt.Printf("lenKeyed is %d. lenSolution is %d", lenKeyed, lenSolution)
	if lenSolution > lenKeyed {
		// there are more solution lines than there is keyed lines.
		// iterate through the keyed lines to find the mistakes.
		// after all the keyed lines is checked iterate through the remaining solution lines and mark those as mistakes.
		var i int
		misMatchesLine := make([]types.MisMatch, 0, lenSolution)
		var keyedLine []*types.KeyCodeRecord
		for i, keyedLine = range keyed {
			solutionLine := solution[i]
			lenKeyedLine := len(keyedLine)
			lenSolutionLine := len(solutionLine)
			nRead += uint64(lenSolutionLine)
			if lenKeyedLine > lenSolutionLine {
				// the user added extra keys in this line.
				var j int
				var sRecord *types.KeyCodeRecord
				for j, sRecord = range solutionLine {
					cRecord := keyedLine[j]
					if cRecord.Character == sRecord.Character {
						nCorrect++
					} else {
						nIncorrect++
					}
					m := types.MisMatch{
						Input:   cRecord,
						Control: sRecord,
					}
					misMatchesLine = append(misMatchesLine, m)
				}
				for j++; j < lenKeyedLine; j++ {
					// extra keys in this line.
					cRecord := keyedLine[j]
					nIncorrect++
					m := types.MisMatch{
						Input:   cRecord,
						Control: keycodes.NotCopiedByApp,
					}
					misMatchesLine = append(misMatchesLine, m)
				}
			} else {
				// the user did not add extra keys in this line.
				// there may not be enough keys.
				// lenKeyedLine <= lenSolutionLine
				var j int
				var cRecord *types.KeyCodeRecord
				// keyed
				for j, cRecord = range keyedLine {
					sRecord := solutionLine[j]
					if cRecord.Character == sRecord.Character {
						nCorrect++
					} else {
						nIncorrect++
					}
					m := types.MisMatch{
						Input:   cRecord,
						Control: sRecord,
					}
					misMatchesLine = append(misMatchesLine, m)
				}
				for j++; j < lenSolutionLine; j++ {
					// missing keys in this line.
					sRecord := solutionLine[j]
					m := types.MisMatch{
						Input:   keycodes.NotKeyedByUser,
						Control: sRecord,
					}
					misMatchesLine = append(misMatchesLine, m)
					nIncorrect++
				}
			}
			misMatches = append(misMatches, misMatchesLine)
			misMatchesLine = make([]types.MisMatch, 0, lenSolution)
		}
		// there is more lines of solution but no more keyed.
		// each remaining solution represents a mis matched line.
		for i = len(keyed); i < lenSolution; i++ {
			solutionLine := solution[i]
			for _, sRecord := range solutionLine {
				// not keyed by the user.
				nIncorrect++
				nRead++
				m := types.MisMatch{
					Control: sRecord,
					Input:   keycodes.NotKeyedByUser,
				}
				misMatchesLine = append(misMatchesLine, m)
			}
			misMatches = append(misMatches, misMatchesLine)
			misMatchesLine = make([]types.MisMatch, 0, lenSolution)
		}
	} else {
		// lenSolution <= lenKeyed
		// there is more keyed lines than there are solution lines.
		// iterate through the solutions and find the mistakes in the keyed.
		// after the solutions are checked mark the rest of the keyed as mistakes.
		var i int
		misMatchesLine := make([]types.MisMatch, 0, lenKeyed)
		var solutionLine []*types.KeyCodeRecord
		for i, solutionLine = range solution {
			keyedLine := keyed[i]
			lenKeyedLine := len(keyedLine)
			lenSolutionLine := len(solutionLine)
			nRead += uint64(lenSolutionLine)
			if lenKeyedLine > lenSolutionLine {
				// the user keyed extra chars in this line.
				var j int
				var sRecord *types.KeyCodeRecord
				for j, sRecord = range solutionLine {
					// text and keys.
					cRecord := keyedLine[j]
					if cRecord.Character == sRecord.Character {
						nCorrect++
					} else {
						nIncorrect++
					}
					m := types.MisMatch{
						Input:   cRecord,
						Control: sRecord,
					}
					misMatchesLine = append(misMatchesLine, m)
				}
				for j++; j < lenKeyedLine; j++ {
					// keyed but there was no text to key in this line.
					nIncorrect++
					cRecord := keyedLine[j]
					m := types.MisMatch{
						Input:   cRecord,
						Control: keycodes.NotCopiedByApp,
					}
					misMatchesLine = append(misMatchesLine, m)
				}
			} else {
				// lenKeyed <= lenSolution
				// the user keyed the correct number or too few this line.
				var j int
				var cRecord *types.KeyCodeRecord
				for j, cRecord = range keyedLine {
					// text and keys.
					sRecord := solutionLine[j]
					nRead++
					if cRecord.Character == sRecord.Character {
						nCorrect++
					} else {
						nIncorrect++
					}
					m := types.MisMatch{
						Input:   cRecord,
						Control: sRecord,
					}
					misMatchesLine = append(misMatchesLine, m)
				}
				for j++; j < lenSolutionLine; j++ {
					// text but the user missed keys for this line.
					sRecord := solutionLine[j]
					nRead++
					m := types.MisMatch{
						Input:   keycodes.NotKeyedByUser,
						Control: sRecord,
					}
					misMatchesLine = append(misMatchesLine, m)
					nIncorrect++
				}
			}
			misMatches = append(misMatches, misMatchesLine)
			misMatchesLine = make([]types.MisMatch, 0, lenSolution)
		}
		i++
		if i == lenKeyed {
			// no more keyed
			return
		}
		// extra keyed lines are mistakes.
		for ; i < lenKeyed; i++ {
			keyedLine := keyed[i]
			for _, cRecord := range keyedLine {
				nIncorrect++
				m := types.MisMatch{
					Input:   cRecord,
					Control: keycodes.NotCopiedByApp,
				}
				misMatchesLine = append(misMatchesLine, m)
			}
		}
		misMatches = append(misMatches, misMatchesLine)
	}
	return
}

func recordCheckResults(keyCodeStorer storer.KeyCodeStorer, misMatches [][]types.MisMatch, wpm uint64) (err error) {
	for _, mm := range misMatches {
		for _, m := range mm {
			results := m.Control.KeyWPMResults[wpm]
			results.Attempts++
			if m.Control.ID == m.Input.ID {
				results.Correct++
			}
			m.Control.KeyWPMResults[wpm] = results
			keyCodeStorer.UpdateKeyCode(m.Control)
		}
	}
	return
}
