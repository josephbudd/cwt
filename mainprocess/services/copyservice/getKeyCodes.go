package copyservice

import (
	"crypto/rand"
	"errors"
	"math/big"

	"github.com/josephbudd/cwt/domain/data"
	"github.com/josephbudd/cwt/domain/store/record"
	"github.com/josephbudd/cwt/domain/store/storer"
)

// GetKeyCodes builds key codes for keying to the user.
// Key codes for the user to copy.
// Returns a slice of key code words.
// A keycode word is a slice of key code records.
// Each word has a length of 5 key codes.
func GetKeyCodes(keyCodeStorer storer.KeyCodeStorer) (keyCodeWords [][]*record.KeyCode, err error) {
	rmap, err := getKeyCodeRecordsMap(keyCodeStorer)
	if err != nil {
		err = errors.New("GetKeyCodes getKeyCodeRecordsMap(keyCodeStorer) error is " + err.Error())
		return
	}
	// build the character/number list.
	var records []*record.KeyCode
	var ok bool
	if records, ok = rmap[data.KeyCodeTypeLetter]; !ok {
		records = make([]*record.KeyCode, 0, 5)
	}
	if rr, ok := rmap[data.KeyCodeTypeNumber]; ok {
		for _, r := range rr {
			records = append(records, r)
		}
	}
	keyCodeWords = buildKeyCodeWords(records, 5, 5)
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

func getKeyCodeRecordsMap(keyCodeStorer storer.KeyCodeStorer) (map[uint64][]*record.KeyCode, error) {
	rmap := make(map[uint64][]*record.KeyCode)
	rr, err := keyCodeStorer.GetAll()
	if err != nil {
		return nil, err
	}
	for _, r := range rr {
		if r.Selected {
			if _, ok := rmap[r.Type]; !ok {
				rmap[r.Type] = make([]*record.KeyCode, 0, 50)
			}
			rmap[r.Type] = append(rmap[r.Type], r)
		}
	}
	return rmap, nil
}
