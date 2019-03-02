package sound

import (
	"strings"
	"testing"
	"time"
)

const paris = ".--. .- .-. .. ..."

func TestPlayCW(t *testing.T) {
	cwTest15Seconds(t)
}

func cwTest15Seconds(t *testing.T) {
	// pp is 6 words. At 24wpm should take 1/4 minute to play.
	pp := []string{paris, paris, paris, paris, paris, paris}
	ditdah := strings.Join(pp, "\t")
	startTime := time.Now()
	if err := PlayCW(ditdah, 24); err != nil {
		t.Error(err.Error())
		return
	}
	elapsedTime := time.Since(startTime)
	seconds := elapsedTime.Seconds()
	if seconds < 14.0 || seconds > 16.0 {
		t.Errorf(" That took %f seconds instead of 15.", seconds)
	}
}
