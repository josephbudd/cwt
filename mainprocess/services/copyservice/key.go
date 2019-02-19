package copyservice

import (
	"time"

	"github.com/josephbudd/cwt/mainprocess/goalsa"
)

// Key plays the words (morse code) at wpm after pausing.
// Param words is an array of indivisual cw words.
//   A word is a combination of characters separated by single spaces.
//   A character is a combination of "-" and ".".
// Param wpm is the words per minute.
// Param pauseSecondsn is the amount of time in seconds to pause before playing the morse code.
func Key(words []string, wpm, pauseSeconds uint64) error {
	timeout := time.After(time.Second * time.Duration(pauseSeconds))
	select {
	case <-timeout:
	}
	return goalsa.PlayCWWords(words, wpm)
}
