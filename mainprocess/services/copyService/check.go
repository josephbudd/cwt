package copyService

import (
	"errors"
	"fmt"

	"github.com/josephbudd/cwt/domain/data/keycodes"
	"github.com/josephbudd/cwt/domain/interfaces/storer"
	"github.com/josephbudd/cwt/domain/types"
)

// Check checks the user's copy agains the solutionChars and returns results.
func Check(copy [][]*types.KeyCodeRecord, solution [][]*types.KeyCodeRecord, keyCodeStorer storer.KeyCodeStorer, wpm uint64, recordResults bool) (nCorrect, nIncorrect, nKeyed uint64, misMatches [][]types.MisMatch, err error) {
	defer func() {
		if err != nil {
			return
		}
		fmt.Printf("recordResults is %v", recordResults)
		if recordResults {
			err = recordCheckResults(keyCodeStorer, misMatches, wpm)
		}
	}()
	// now copy is a slice of strings and solutionChars is a slice of strings.
	// the 2 should match is the user copied correctly.
	misMatches = make([][]types.MisMatch, 0, 100)
	lc := len(copy)
	ls := len(solution)
	if ls > lc {
		// there are more solutions than there is copy.
		// iterate through the copy to find the mistakes.
		// after all the copy is checked iterate through the remaining solutions and mark those as mistakes.
		var i int
		misMatchesLine := make([]types.MisMatch, 0, ls)
		var copyLine []*types.KeyCodeRecord
		for i, copyLine = range copy {
			solutionLine := solution[i]
			lcl := len(copyLine)
			lsl := len(solutionLine)
			nKeyed += uint64(lsl)
			if lcl > lsl {
				var j int
				var sRecord *types.KeyCodeRecord
				for j, sRecord = range solutionLine {
					cRecord := copyLine[j]
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
				for j++; j < lcl; j++ {
					cRecord := copyLine[j]
					nIncorrect++
					m := types.MisMatch{
						Input:   cRecord,
						Control: keycodes.NotKeyed,
					}
					misMatchesLine = append(misMatchesLine, m)
				}
			} else {
				// lcl <= lsl
				var j int
				var cRecord *types.KeyCodeRecord
				for j, cRecord = range copyLine {
					// copied
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
				for j++; j < lsl; j++ {
					// missing copy
					sRecord := solutionLine[j]
					m := types.MisMatch{
						Input:   keycodes.NotCopied,
						Control: sRecord,
					}
					misMatchesLine = append(misMatchesLine, m)
					nIncorrect++
				}
			}
			misMatches = append(misMatches, misMatchesLine)
			misMatchesLine = make([]types.MisMatch, 0, ls)
		}
		// there is more solution but no more copy.
		// each remaining solution represents a mis match.
		for i++; i < ls; i++ {
			solutionLine := solution[i]
			for _, sRecord := range solutionLine {
				// not copied
				nIncorrect++
				nKeyed++
				m := types.MisMatch{
					Control: sRecord,
					Input:   keycodes.NotCopied,
				}
				misMatchesLine = append(misMatchesLine, m)
			}
			misMatches = append(misMatches, misMatchesLine)
			misMatchesLine = make([]types.MisMatch, 0, ls)
		}
	} else {
		// lc >= ls {
		// there is more copy than there is correct solutions.
		// iterate through the solutions and find the mistakes in the copy.
		// after the solutions are checked mark the rest of the copy as mistakes.
		var i int
		misMatchesLine := make([]types.MisMatch, 0, lc)
		var solutionLine []*types.KeyCodeRecord
		for i, solutionLine = range solution {
			copyLine := copy[i]
			lcl := len(copyLine)
			lsl := len(solutionLine)
			nKeyed += uint64(lsl)
			if lcl > lsl {
				var j int
				var sRecord *types.KeyCodeRecord
				for j, sRecord = range solutionLine {
					// copied
					cRecord := copyLine[j]
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
				for j++; j < lcl; j++ {
					// copied but nothing keyed
					nIncorrect++
					cRecord := copyLine[j]
					m := types.MisMatch{
						Input:   cRecord,
						Control: keycodes.NotKeyed,
					}
					misMatchesLine = append(misMatchesLine, m)
				}
			} else {
				// lc <= ls
				var j int
				var cRecord *types.KeyCodeRecord
				for j, cRecord = range copyLine {
					// keyed and copied
					sRecord := solutionLine[j]
					nKeyed++
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
				for j++; j < lsl; j++ {
					// keyed but never copied
					sRecord := solutionLine[j]
					nKeyed++
					m := types.MisMatch{
						Input:   keycodes.NotCopied,
						Control: sRecord,
					}
					misMatchesLine = append(misMatchesLine, m)
					nIncorrect++
				}
			}
			misMatches = append(misMatches, misMatchesLine)
			misMatchesLine = make([]types.MisMatch, 0, ls)
		}
		i++
		if i == lc {
			// no more copy
			return
		}
		// more copy and the rest of the copy is mistakes.
		for ; i < lc; i++ {
			copyLine := copy[i]
			for _, cRecord := range copyLine {
				nIncorrect++
				m := types.MisMatch{
					Input:   cRecord,
					Control: keycodes.NotKeyed,
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
			var results types.KeyCodeRecordResult
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
			if err = keyCodeStorer.UpdateKeyCode(m.Control); err != nil {
				return
			}
		}
	}
	return
}
