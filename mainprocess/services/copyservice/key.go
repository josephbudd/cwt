package copyservice

import (
	"time"

	"github.com/josephbudd/cwt/mainprocess/sound"
)

var keyRunning bool
var keyQuitChannel = make(chan struct{})

// Key plays the words (morse code) at wpm after pausing.
// Param words is an array of indivisual cw words.
//   A word is a combination of characters separated by single spaces.
//   A character is a combination of "-" and ".".
// Param wpm is the words per minute.
// Param pauseSecondsn is the amount of time in seconds to pause before playing the morse code.
func Key(words []string, wpm, pauseSeconds uint64) (err error) {
	keyRunning = true
	timeout := time.After(time.Second * time.Duration(pauseSeconds))
	select {
	case <-timeout:
	}
	err = sound.PlayCWWords(words, wpm, keyQuitChannel)
	keyRunning = false
	return
}

// StopKeying stops the keying.
func StopKeying() {
	if keyRunning {
		keyRunning = false
		keyQuitChannel <- struct{}{}
	}
}
