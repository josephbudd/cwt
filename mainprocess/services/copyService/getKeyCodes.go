package copyService

import (
	"crypto/rand"
	"errors"
	"math/big"

	"github.com/josephbudd/cwt/domain/data/keyCodeTypes"
	"github.com/josephbudd/cwt/domain/interfaces/storer"
	"github.com/josephbudd/cwt/domain/types"
)

// GetKeyCodes builds key codes for keying to the user.
// Key codes for the user to copy.
// Returns a slice of key code words.
// A keycode word is a slice of key code records.
// Each word has a length of 5 key codes.
func GetKeyCodes(keyCodeStorer storer.KeyCodeStorer) (keyCodeWords [][]*types.KeyCodeRecord, err error) {
	rmap, err := getKeyCodeRecordsMap(keyCodeStorer)
	if err != nil {
		err = errors.New("GetKeyCodes getKeyCodeRecordsMap(keyCodeStorer) error is " + err.Error())
		return
	}
	// build the character/number list.
	var records []*types.KeyCodeRecord
	var ok bool
	if records, ok = rmap[keyCodeTypes.KeyCodeTypeLetter]; !ok {
		records = make([]*types.KeyCodeRecord, 0, 5)
	}
	if rr, ok := rmap[keyCodeTypes.KeyCodeTypeNumber]; ok {
		for _, r := range rr {
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
	shuffled = make([]*types.KeyCodeRecord, count, count)
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

func getKeyCodeRecordsMap(keyCodeStorer storer.KeyCodeStorer) (map[uint64][]*types.KeyCodeRecord, error) {
	rmap := make(map[uint64][]*types.KeyCodeRecord)
	records, err := keyCodeStorer.GetKeyCodes()
	if err != nil {
		return nil, err
	}
	for _, record := range records {
		if record.Selected {
			if _, ok := rmap[record.Type]; !ok {
				rmap[record.Type] = make([]*types.KeyCodeRecord, 0, 50)
			}
			rmap[record.Type] = append(rmap[record.Type], record)
		}
	}
	return rmap, nil
}
