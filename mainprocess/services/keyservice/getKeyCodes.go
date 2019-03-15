package keyservice

import (
	"crypto/rand"
	"errors"
	"math/big"
	"sort"

	"github.com/josephbudd/cwt/domain/interfaces/storer"
	"github.com/josephbudd/cwt/domain/types"
	"github.com/josephbudd/cwt/mainprocess/howto"
)

// GetTestKeyCodes returns
// Key codes for the user to key in a test.
// Help
func GetTestKeyCodes(keyCodeStorer storer.KeyCodeStorer) (keyCodeWords [][]*types.KeyCodeRecord, err error) {
	if keyCodeWords, err = getKeyCodes(keyCodeStorer); err != nil {
		return
	}
	return
}

// GetPracticeKeyCodes returns
// Key codes for the user to key in practice.
// Help
func GetPracticeKeyCodes(keyCodeStorer storer.KeyCodeStorer, wpm uint64) (keyCodeWords [][]*types.KeyCodeRecord, help [][]types.HowTo, err error) {
	if keyCodeWords, err = getPracticeKeyCodes(keyCodeStorer, wpm); err != nil {
		return
	}
	help = howto.KeyCodesToHelp(keyCodeWords)
	return
}

// getPracticeKeyCodes builds key codes for displaying text to the user.
// Key codes for the user to key.
// Returns a slice of key code words.
// A keycode word is a slice of key code records.
// Each word has a length of 5 key codes.
func getPracticeKeyCodes(keyCodeStorer storer.KeyCodeStorer, wpm uint64) (keyCodeWords [][]*types.KeyCodeRecord, err error) {
	rr, err := keyCodeStorer.GetKeyCodes()
	if err != nil {
		err = errors.New("GetKeyCodes keyCodeStorer.GetKeyCodes() error is " + err.Error())
		return
	}
	worst := getWorst(rr, wpm, 5)
	keyCodeWords = buildKeyCodeWords(worst, 5, 1)
	return
}

func getWorst(rr []*types.KeyCodeRecord, wpm uint64, count uint64) (worst []*types.KeyCodeRecord) {
	// map records to wpm
	var percent int
	percentRecord := make(map[int][]*types.KeyCodeRecord)
	for _, r := range rr {
		if r.Selected {
			if result, ok := r.KeyWPMResults[wpm]; ok {
				if result.Attempts == 0 {
					percent = 0
				} else {
					percent = int((result.Correct * 100) / result.Attempts)
				}
				if _, ok = percentRecord[percent]; !ok {
					percentRecord[percent] = make([]*types.KeyCodeRecord, 0, 50)
				}
				percentRecord[percent] = append(percentRecord[percent], r)
			}
		}
	}
	// sort the percent
	l := len(percentRecord)
	sortedPercent := make([]int, 0, l)
	for k := range percentRecord {
		sortedPercent = append(sortedPercent, k)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(sortedPercent)))
	// collect the worst
	worst = make([]*types.KeyCodeRecord, 0, count)
	max := int(count)
	for _, percent = range sortedPercent {
		for _, r := range percentRecord[percent] {
			worst = append(worst, r)
		}
		if len(worst) >= max {
			break
		}
	}
	if len(worst) <= max {
		return
	}
	// shuffle worst and resize to max
	worst = shuffleKeyCodeRecords(worst)
	worst = worst[:max]
	return
}

// getKeyCodes builds key codes for displaying text to the user.
// Key codes for the user to key.
// Returns a slice of key code words.
// A keycode word is a slice of key code records.
// Each word has a length of 5 key codes.
func getKeyCodes(keyCodeStorer storer.KeyCodeStorer) (keyCodeWords [][]*types.KeyCodeRecord, err error) {
	rr, err := keyCodeStorer.GetKeyCodes()
	if err != nil {
		err = errors.New("GetKeyCodes keyCodeStorer.GetKeyCodes() error is " + err.Error())
		return
	}
	records := make([]*types.KeyCodeRecord, 0, len(rr))
	for _, r := range rr {
		if r.Selected {
			records = append(records, r)
		}
	}
	keyCodeWords = buildKeyCodeWords(records, 5, 5)
	return
}

// buildKeyCodeWords builds a slice of key code words.
// A keycode word is a slice of key code records.
func buildKeyCodeWords(rr []*types.KeyCodeRecord, wordSize, maxLines int) (keyCodeWords [][]*types.KeyCodeRecord) {
	keyCodeWords = make([][]*types.KeyCodeRecord, 0, maxLines)
	recordCount := len(rr)
	if recordCount < wordSize {
		wordSize = recordCount
	}
	for i := 0; i < maxLines; i++ {
		shuffled := shuffleKeyCodeRecords(rr)
		keyCodeWords = append(keyCodeWords, shuffled[:wordSize])
	}
	return
}

func shuffleKeyCodeRecords(records []*types.KeyCodeRecord) (shuffled []*types.KeyCodeRecord) {
	count := len(records)
	temp := make(map[int64]*types.KeyCodeRecord)
	max := big.NewInt(int64(count))
	for i := 0; i < count; i++ {
		var bigJ *big.Int
		var err error
		var j, k int64
		var used bool
		if bigJ, err = rand.Int(rand.Reader, max); err != nil {
			j = 0
		} else {
			j = int64(bigJ.Int64())
		}
		last := int64(count) - 1
		if _, used = temp[j]; used {
			// temp[j] is already filled.
			for k = j; k < last; {
				k++
				if _, used = temp[k]; !used {
					break
				}
			}
			if used {
				for k = j; k > 0; {
					k--
					if _, used = temp[k]; !used {
						break
					}
				}
			}
			temp[k] = records[i]
		} else {
			// temp[j] is empty.
			temp[j] = records[i]
		}
	}
	// build shuffled
	shuffled = make([]*types.KeyCodeRecord, count, count)
	for k, v := range temp {
		shuffled[k] = v
	}
	return
}
