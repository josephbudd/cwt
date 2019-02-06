package goalsa

/*

#cgo CFLAGS: -Iinclude
#include "alsa.h"

*/
import "C"
import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// PlayCWWords plays an array of indivisual cw words.
// A word is a combination of characters separated by single spaces.
// A character is a combination of "-" and ".".
// PlayCWWords joins the words with "\t"
// PlayCWWords plays them at wpm words per minute.
func PlayCWWords(words []string, wpm uint64) (err error) {
	err = playCW(strings.Join(words, "\t"), wpm)
	return
}

// PlayCW plays "-. ._\t-. ._" at wpm words per minute.
// "." == dit.
// "-" == dah.
// " " == character separator. A character is a string of dits and dahs.
// "\t" == word separator. A word is a string of characters joined by character separators.
func PlayCW(ditdah string, wpm uint64) (err error) {
	for i, r := range ditdah {
		if r != '.' && r != '-' && r != ' ' && r != '\t' {
			err = fmt.Errorf("%q is not a valid dit dah character at position %d. It should be a \".\", a \"-\", a space or a tab", r, i)
			return
		}
	}
	err = playCW(ditdah, wpm)
	return
}

func playCW(ditdah string, wpm uint64) (err error) {
	defer func() {
		if err != nil {
			err = errors.WithMessage(err, "PlayCW(ditdah string, wpm uint64)")
		}
	}()
	player := newPlayer(440, wpm)
	if err = player.open(); err != nil {
		return
	}
	defer player.close()
	err = player.ditDahToFF(ditdah)
	return
}

// AlsaInfo supplies alsa information for this device.
type AlsaInfo struct {
	Version     string
	StreamTypes map[int]string
	AccessTypes map[int]string
	Formats     map[int]struct{ Name, Description string }
	SubFormats  map[int]struct{ Name, Description string }
	States      map[int]string
}

// GetInfo returns the alsa information for this device.
func GetInfo() *AlsaInfo {
	return getAlsaInfo()
}
