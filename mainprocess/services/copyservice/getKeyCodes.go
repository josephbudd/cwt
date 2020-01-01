package copyservice

import (
	"crypto/rand"
	"errors"
	"log"
	"math/big"
	"sort"

	"github.com/josephbudd/cwt/domain/store/record"
	"github.com/josephbudd/cwt/domain/store/storer"
)

// GetKeyCodes builds key codes for keying to the user.
// Key codes for the user to copy.
// Returns a slice of key code words.
// A keycode word is a slice of key code records.
// Each word has a length of 5 key codes.
func GetKeyCodes(keyCodeStorer storer.KeyCodeStorer) (keyCodeWords [][]*record.KeyCode, err error) {
	rr, err := getSelected(keyCodeStorer)
	log.Printf("GetKeyCodes: %d selected records", len(rr))
	if err != nil {
		err = errors.New("GetKeyCodes getSelected(keyCodeStorer) error is " + err.Error())
		return
	}
	keyCodeWords = buildKeyCodeWords(shuffleKeyCodeRecords(rr), 5, 5)
	return
}

// buildKeyCodeWords builds a slice of key code words.
// A keycode word is a slice of key code records.
func buildKeyCodeWords(rr []*record.KeyCode, wordSize, maxLines int) (keyCodeWords [][]*record.KeyCode) {
	keyCodeWords = make([][]*record.KeyCode, 0, maxLines)
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

func shuffleKeyCodeRecords(records []*record.KeyCode) (shuffled []*record.KeyCode) {
	count := len(records)
	shuffled = make([]*record.KeyCode, count, count)
	max := big.NewInt(int64(count))
	for i := 0; i < count; i++ {
		var bigJ *big.Int
		var err error
		var j int64
		if bigJ, err = rand.Int(rand.Reader, max); err != nil {
			j = 0
		} else {
			j = int64(bigJ.Int64())
		}
		shuffled[i] = records[j]
		shuffled[j] = records[i]
	}
	return
}

func getSelected(keyCodeStorer storer.KeyCodeStorer) (selected []*record.KeyCode, err error) {
	rr, err := keyCodeStorer.GetAll()
	if err != nil {
		return
	}
	selected = make([]*record.KeyCode, 0, len(rr))
	for _, r := range rr {
		if r.Selected {
			selected = append(selected, r)
		}
	}
	return
}

func getWorst(rr []*record.KeyCode, wpm uint64, count uint64) (worst []*record.KeyCode) {
	// map records to wpm
	var percent int
	percentRecord := make(map[int][]*record.KeyCode)
	for _, r := range rr {
		if r.Selected {
			if result, ok := r.CopyWPMResults[wpm]; ok {
				if result.Attempts == 0 {
					percent = 0
				} else {
					percent = int((result.Correct * 100) / result.Attempts)
				}
				if _, ok = percentRecord[percent]; !ok {
					percentRecord[percent] = make([]*record.KeyCode, 0, 50)
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
	sort.Sort(sort.IntSlice(sortedPercent))
	// collect the worst
	worst = make([]*record.KeyCode, 0, count)
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
