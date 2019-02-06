package goalsa

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

const paris = ".--. .- .-. .. ..."

func TestCW(t *testing.T) {

	// goAlsaTestFail(t, "s", 0)
	// goAlsaTestFail(t, ".- \tm.- \t", 4)
	// goAlsaTestSecond(t, 30)
	// goAlsaTestSecond(t, 5)
	// goAlsaTestMinute(t, 30)
	// goAlsaTestMinute(t, 5)
	fmt.Println("Testing ditdah without word separators.")
	goAlsaTest5(t, false)
	fmt.Println("Testing ditdah with word separators.")
	goAlsaTest5(t, true)
}

func goAlsaTest5(t *testing.T, separateWords bool) {
	// pp is 5 words. At 5wpm should take 1 minute to play.
	fmt.Println("Testing ditdahs. This should take 1 minute.")
	pp := []string{paris, paris, paris, paris, paris}
	var ditdah string
	if separateWords {
		ditdah = strings.Join(pp, "\t")
	} else {
		ditdah = strings.Join(pp, "")
	}
	startTime := time.Now()
	err := PlayCW(ditdah, 5)
	elapsedTime := time.Since(startTime)
	seconds := elapsedTime.Seconds()
	fmt.Printf(" That took %f seconds.", seconds)
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println("")
}

func goAlsaTestFail(t *testing.T, errorString string, errorIndex int) {
	err := PlayCW(errorString, 5)
	if err != nil {
		got := err.Error()
		want := fmt.Sprintf("'%s' is not a valid dit dah character at position %d. It should be a \".\", a \"-\", a space or a tab", errorString[errorIndex:errorIndex+1], errorIndex)
		if got != want {
			t.Errorf("got %q\nwant %q\n", got, want)
		}
	}
}

func goAlsaTestSecond(t *testing.T, wpm uint64) {
	fmt.Printf("1 second at %d wpm.", wpm)
	startTime := time.Now()
	err := playCW("s", wpm)
	if err != nil {
		t.Error("\nerr is ", err.Error())
		return
	}
	elapsedTime := time.Since(startTime)
	seconds := elapsedTime.Seconds()
	fmt.Printf(" That took %f seconds.", seconds)
	if err != nil {
		t.Error("\nerr is ", err.Error())
		return
	}
	if seconds > 1.05 {
		off := seconds - float64(1.0)
		t.Errorf(" That was %f seconds too long.\n", off)
		return
	}
	if seconds < 1.0 {
		off := float64(1.0) - seconds
		t.Errorf(" That was %f seconds too short.\n", off)
		return
	}
	fmt.Print("\n")
}

func goAlsaTestMinute(t *testing.T, wpm uint64) {
	fmt.Printf("30 seconds at %d wpm.", wpm)
	startTime := time.Now()
	err := playCW("m", wpm)
	if err != nil {
		t.Error("\nerr is ", err.Error())
		return
	}
	elapsedTime := time.Since(startTime)
	seconds := elapsedTime.Seconds()
	fmt.Printf(" That took %f seconds.", seconds)
	if err != nil {
		t.Error("\nerr is ", err.Error())
		return
	}
	if seconds > 31.0 {
		off := seconds - float64(1.0)
		t.Errorf(" That was %f seconds too long.\n", off)
		return
	}
	if seconds < 29.9 {
		off := float64(60.0) - seconds
		t.Errorf(" That was %f seconds too short.\n", off)
		return
	}
	fmt.Print("\n")
}
