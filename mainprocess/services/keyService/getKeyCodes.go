package keyService

import (
	"crypto/rand"
	"errors"
	"math/big"

	"github.com/josephbudd/cwt/domain/interfaces/storer"
	"github.com/josephbudd/cwt/domain/types"
)

// GetKeyCodes builds key codes for displaying text to the user.
// Key codes for the user to key.
// Returns a slice of key code words.
// A keycode word is a slice of key code records.
// Each word has a length of 5 key codes.
func GetKeyCodes(keyCodeStorer storer.KeyCodeStorer) (keyCodeWords [][]*types.KeyCodeRecord, err error) {
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
