package howto

import (
	"fmt"
	"strings"

	"github.com/josephbudd/cwt/domain/types"
)

const (
	wordSeparaterInstructions = "2 3 4 5 6 7"
	charSeparaterInstructions = "2 3"
	separaterHowTo            = ", "
)

// KeyCodesToHelp turns key code records that form a word into help.
func KeyCodesToHelp(words [][]*types.KeyCodeRecord) (howtos [][]types.HowTo) {
	lws := len(words)
	howtos = make([][]types.HowTo, lws, lws)
	for i, word := range words {
		//lw := len(word)
		helps := keyCodesWordToHelp(word)
		l := len(helps)
		if i > 0 {
			// new word
			var howto []types.HowTo
			howto = make([]types.HowTo, 0, l+1)
			howto = append(howto, types.HowTo{Instructions: wordSeparaterInstructions})
			for _, h := range helps {
				howto = append(howto, h)
			}
			howtos[i] = howto
		} else {
			howtos[i] = helps
		}
	}
	return
}

func keyCodesWordToHelp(word []*types.KeyCodeRecord) (howto []types.HowTo) {
	l := len(word)
	howto = make([]types.HowTo, 0, l*2)
	for i, char := range word {
		if i > 0 {
			howto = append(howto, types.HowTo{
				Instructions: charSeparaterInstructions,
			})
		}
		howto = append(howto, types.HowTo{
			Character:    char.Character,
			DitDah:       char.DitDah,
			Instructions: ditDahToHowTo(char.DitDah),
		})
	}
	return
}

func ditDahToHowTo(ditdah string) (howto string) {
	var b strings.Builder
	for i, r := range ditdah {
		switch r {
		case '.':
			// dit in character
			if i > 0 {
				fmt.Fprint(&b, separaterHowTo)
			}
			fmt.Fprint(&b, "down up")
		case '-':
			// dah in character
			if i > 0 {
				fmt.Fprint(&b, separaterHowTo)
			}
			fmt.Fprint(&b, "down 2 3 up")
		case ' ':
			// space between character in word
			if i > 0 {
				fmt.Fprint(&b, separaterHowTo)
			}
			fmt.Fprint(&b, "2 3")
		}
	}
	// howto = strings.Join(steps, ", ")
	howto = b.String()
	return
}
