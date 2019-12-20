// +build js, wasm

package widgets

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/josephbudd/cwt/domain/data"
	"github.com/josephbudd/cwt/domain/store/record"
	"github.com/josephbudd/cwt/rendererprocess/api/display"
	"github.com/josephbudd/cwt/rendererprocess/api/dom"
	"github.com/josephbudd/cwt/rendererprocess/api/event"
	"github.com/josephbudd/cwt/rendererprocess/api/markup"
)

const (
	initialInstructions           = "Click start to begin."
	mouseNotOverInstructions      = "Mouse over the area below to key."
	mouseNotOverAgainInstructions = "Mouse over the area below to key or click the [Check] button."
	mouseOverInstructions         = "Ready for you to key."
	noInstructions                = ""
)

// KeyWidget allows the user to key.
// KeyWidget is a widget of buttons and key areas which record keying.
type KeyWidget struct {
	document *dom.DOM

	heading           *markup.Element
	startButton       *markup.Element
	stopButton        *markup.Element
	metronomeCheckBox *markup.Element
	keyDiv            *markup.Element
	copyDiv           *markup.Element

	userKeyChecker UserKeyChecker
	metronomer     Metronomer

	userIsKeying bool
	metronomeOn  bool
	keyTime      time.Time
	durations    []time.Duration

	keyCodes [][]*record.KeyCode
	help     [][]data.HowTo
	wpm      uint64
	times    []time.Time
}

// UserKeyChecker processes keying millisecons.
type UserKeyChecker interface {
	GetKeyCodesWPM()
	CheckUserKey(milliSeconds []int64, solution [][]*record.KeyCode, wpm uint64)
}

// Metronomer starts and stops a metronome.
type Metronomer interface {
	StartMetronome(wpm uint64)
	StopMetronome()
}

// API
// * NewKeyWidget
// * SetKeyCodesWPM
// * ShowResults

// NewKeyWidget constructs a new KeyWidget
// Param startButton is the start button.
// Param stopButton is the stop button.
// Param keyDiv is the div moused over with the keyWidget-mouse when keying.
// Param copyDiv is where
// * the text is displayed that the user reads when keying.
// * the copy of the user's keying is displayed.
func NewKeyWidget(document *dom.DOM,
	heading *markup.Element,
	startButton, stopButton, metronomeCheckBox *markup.Element,
	keyDiv, copyDiv *markup.Element,
	userKeyChecker UserKeyChecker,
	metronomer Metronomer) (keyWidget *KeyWidget) {
	keyWidget = &KeyWidget{
		document:          document,
		heading:           heading,
		startButton:       startButton,
		stopButton:        stopButton,
		metronomeCheckBox: metronomeCheckBox,
		keyDiv:            keyDiv,
		copyDiv:           copyDiv,

		userKeyChecker: userKeyChecker,
		metronomer:     metronomer,

		times:     make([]time.Time, 0, 1024),
		durations: make([]time.Duration, 0, 1024),
	}

	if metronomeCheckBox != nil {
		keyWidget.metronomeOn = metronomeCheckBox.Checked()
		// Handle the keyWidget start button's click event.
		metronomeCheckBox.SetEventHandler(keyWidget.handleMetronomeCheck, "change", false)
	}

	heading.SetInnerText(initialInstructions)

	// Handle the keyWidget start button's click event.
	keyWidget.startButton.SetEventHandler(keyWidget.handleStart, "click", false)

	// Handle the keyWidget stop button's click event.
	keyWidget.stopButton.SetEventHandler(keyWidget.handleStop, "click", false)

	// Handle the keyWidget's key div mouse enter event.
	keyWidget.keyDiv.SetEventHandler(keyWidget.handleMouseEnter, "mouseenter", false)

	// Handle the keyWidget's key div mouse leave event.
	keyWidget.keyDiv.SetEventHandler(keyWidget.handleMouseLeave, "mouseleave", false)

	// Handle the keyWidget's key div mouse down event.
	// Mouse down is the morse code key down.
	keyWidget.keyDiv.SetEventHandler(keyWidget.handleMouseDown, "mousedown", false)

	// Handle the keyWidget's key div mouse up event.
	// Mouse down is the morse code key up.
	keyWidget.keyDiv.SetEventHandler(keyWidget.handleMouseUp, "mouseup", false)

	return
}

// SetKeyCodesHelpWPM sets the key codes, ( the solution ) that the user must key.
func (keyWidget *KeyWidget) SetKeyCodesHelpWPM(records [][]*record.KeyCode, help [][]data.HowTo, wpm uint64) {
	keyWidget.keyCodes = records
	keyWidget.help = help
	keyWidget.wpm = wpm
	keyWidget.tellUserToStart()
}

// SetKeyCodesWPM sets the key codes, ( the solution ) that the user must key.
func (keyWidget *KeyWidget) SetKeyCodesWPM(records [][]*record.KeyCode, wpm uint64) {
	keyWidget.keyCodes = records
	keyWidget.help = nil
	keyWidget.wpm = wpm
	keyWidget.tellUserToStart()
}

// ShowResults Processes and displays the checked results of what the user keyed.
func (keyWidget *KeyWidget) ShowResults(correctCount, incorrectCount, keyedCount uint64, testResults [][]data.TestResult) {
	// clear the div
	keyWidget.copyDiv.RemoveChildren()
	// heading
	h3 := keyWidget.document.NewH3()
	h3.SetInnerText(fmt.Sprintf("You got %d correct and %d incorrect out of %d.", correctCount, incorrectCount, keyedCount))
	keyWidget.copyDiv.AppendChild(h3)
	// details
	// put each word in its own table
	for _, mmWord := range testResults {
		// results
		controlChars := make([]string, 0, 20)
		copiedChars := make([]string, 0, 20)
		ditdahs := make([]string, 0, 20)
		for _, mmChar := range mmWord {
			controlChars = append(controlChars, mmChar.Control.Character)
			copiedChars = append(copiedChars, mmChar.Input.Character)
			ditdahs = append(ditdahs, mmChar.Input.DitDah)
		}
		controlWord := strings.Join(controlChars, "")
		ditDahWord := strings.Join(ditdahs, " ")
		copiedWord := strings.Join(copiedChars, "")
		keyWidget.addTable(controlWord, ditDahWord, copiedWord)
	}
	display.Resize()
	keyWidget.startButton.Show()
}

func (keyWidget *KeyWidget) addTable(controlWord, ditDahWord, copiedWord string) {
	var tdClass string
	if controlWord == copiedWord {
		tdClass = "match"
	} else {
		tdClass = "mismatch"
	}

	// table
	table := keyWidget.document.NewTABLE()
	table.AddClass("resize-me-width")
	tbody := keyWidget.document.NewTBODY()
	// control row
	tr := keyWidget.document.NewTR()
	th := keyWidget.document.NewTH()
	th.SetInnerText("Text")
	tr.AppendChild(th)
	// control characters
	td := keyWidget.document.NewTD()
	td.SetInnerText(controlWord)
	td.AddClasses("control", tdClass)
	tr.AppendChild(td)
	tbody.AppendChild(tr)

	// keyed row
	tr = keyWidget.document.NewTR()
	th = keyWidget.document.NewTH()
	th.SetInnerText("You Keyed")
	tr.AppendChild(th)
	// user keyed
	td = keyWidget.document.NewTD()
	td.SetInnerText(ditDahWord)
	td.AddClasses("ditdah", tdClass)
	tr.AppendChild(td)
	tbody.AppendChild(tr)

	// copied row
	tr = keyWidget.document.NewTR()
	th = keyWidget.document.NewTH()
	th.SetInnerText("App Copied")
	tr.AppendChild(th)
	// app copied
	td = keyWidget.document.NewTD()
	td.SetInnerText(copiedWord)
	td.AddClasses("copied", tdClass)
	tr.AppendChild(td)
	tbody.AppendChild(tr)

	table.AppendChild(tbody)
	keyWidget.copyDiv.AppendChild(table)
}

// handlers

func (keyWidget *KeyWidget) handleMetronomeCheck(e event.Event) (nilReturn interface{}) {
	if keyWidget.metronomeOn = keyWidget.metronomeCheckBox.Checked(); !keyWidget.metronomeOn {
		keyWidget.metronomer.StopMetronome()
	}
	log.Printf("keyWidget.metronomeOn is %#v", keyWidget.metronomeOn)
	return
}

func (keyWidget *KeyWidget) handleStart(e event.Event) (nilReturn interface{}) {
	keyWidget.userKeyChecker.GetKeyCodesWPM()
	keyWidget.startButton.Hide()
	keyWidget.heading.SetInnerText(mouseNotOverInstructions)
	keyWidget.copyDiv.RemoveChildren()
	keyWidget.keyTime = time.Now()
	keyWidget.times = make([]time.Time, 1, 200)
	keyWidget.times[0] = time.Now()
	return
}

func (keyWidget *KeyWidget) handleStop(e event.Event) (nilReturn interface{}) {
	keyWidget.setKeyingStopped()
	l := len(keyWidget.times)
	milliSeconds := make([]int64, 0, l)
	for i := 0; i < l; {
		start := keyWidget.times[i]
		i++
		if i < l {
			end := keyWidget.times[i]
			d := end.Sub(start)
			milliSeconds = append(milliSeconds, int64(d/1000000))
		}
	}
	keyWidget.userKeyChecker.CheckUserKey(milliSeconds, keyWidget.keyCodes, keyWidget.wpm)
	keyWidget.stopButton.Hide()
	keyWidget.heading.SetInnerText(noInstructions)
	return
}

func (keyWidget *KeyWidget) handleMouseEnter(e event.Event) (nilReturn interface{}) {
	e.PreventDefaultBehavior()
	e.StopAllPhasePropagation()
	if !keyWidget.userIsKeying {
		return
	}
	keyWidget.keyDiv.ReplaceClass("user-not-key-over", "user-key-over")
	keyWidget.heading.SetInnerText(mouseOverInstructions)
	if keyWidget.metronomeOn {
		keyWidget.metronomer.StartMetronome(keyWidget.wpm)
	}
	return
}

func (keyWidget *KeyWidget) handleMouseLeave(e event.Event) (nilReturn interface{}) {
	e.PreventDefaultBehavior()
	e.StopAllPhasePropagation()
	if !keyWidget.userIsKeying {
		return
	}
	keyWidget.keyDiv.ReplaceClass("user-key-over", "user-not-key-over")
	keyWidget.heading.SetInnerText(mouseNotOverAgainInstructions)
	if keyWidget.metronomer != nil {
		keyWidget.metronomer.StopMetronome()
	}
	return
}

func (keyWidget *KeyWidget) handleMouseDown(e event.Event) (nilReturn interface{}) {
	e.PreventDefaultBehavior()
	e.StopAllPhasePropagation()
	if !keyWidget.userIsKeying {
		return
	}
	keyWidget.times = append(keyWidget.times, time.Now())
	keyWidget.keyDiv.ReplaceClass("user-key-up", "user-key-down")
	return
}

func (keyWidget *KeyWidget) handleMouseUp(e event.Event) (nilReturn interface{}) {
	e.PreventDefaultBehavior()
	e.StopAllPhasePropagation()
	if !keyWidget.userIsKeying {
		return
	}
	keyWidget.times = append(keyWidget.times, time.Now())
	keyWidget.keyDiv.ReplaceClass("user-key-down", "user-key-up")
	return
}

// misc funcs

func (keyWidget *KeyWidget) tellUserToStart() {
	// clear durations.
	keyWidget.durations = keyWidget.durations[:0]
	// display the text to key.
	keyWidget.showTextToKey()
	keyWidget.setKeyingStarted()
}

func (keyWidget *KeyWidget) showTextToKey() {
	if keyWidget.help == nil {
		keyWidget.showTextOnlyToKey()
	} else {
		keyWidget.showTextHelpToKey()
	}
}

func (keyWidget *KeyWidget) showTextOnlyToKey() {
	// clear the div
	keyWidget.keyDiv.RemoveChildren()
	// put the lines in a paragraph.
	p := keyWidget.document.NewP()
	var b strings.Builder
	// the words each on a separate line.
	for _, wordCode := range keyWidget.keyCodes {
		for _, charCode := range wordCode {
			fmt.Fprint(&b, charCode.Character)
		}
		p.AppendText(b.String())
		p.AppendChild(keyWidget.document.NewBR())
		b.Reset()
	}
	keyWidget.keyDiv.AppendChild(p)
}

func (keyWidget *KeyWidget) showTextHelpToKey() {
	// clear the div
	keyWidget.keyDiv.RemoveChildren()
	for i, wordCode := range keyWidget.keyCodes {
		// the word
		l := len(wordCode)
		chars := make([]string, l, l)
		for j, charCode := range wordCode {
			chars[j] = charCode.Character
		}
		word := strings.Join(chars, "")
		table := keyWidget.document.NewTABLE()
		tbody := keyWidget.document.NewTBODY()
		table.AppendChild(tbody)
		tr := keyWidget.document.NewTR()
		td := keyWidget.document.NewTD()
		td.SetAttribute("colspan", 3)
		td.SetInnerText(word)
		tr.AppendChild(td)
		tbody.AppendChild(tr)
		// help
		wordHelp := keyWidget.help[i]
		for _, h := range wordHelp {
			tr = keyWidget.document.NewTR()
			// character
			td = keyWidget.document.NewTD()
			td.SetInnerText(h.Character)
			tr.AppendChild(td)
			// dit dah
			td = keyWidget.document.NewTD()
			td.SetInnerText(h.DitDah)
			tr.AppendChild(td)
			// instructions
			td = keyWidget.document.NewTD()
			td.SetInnerText(h.Instructions)
			tr.AppendChild(td)
			tbody.AppendChild(tr)
		}
		table.AppendChild(tbody)
		keyWidget.keyDiv.AppendChild(table)
	}
}

func (keyWidget *KeyWidget) setKeyingStarted() {
	keyWidget.userIsKeying = true
	keyWidget.keyDiv.ReplaceClass("user-not-keying", "user-keying")
	keyWidget.stopButton.Show()
	keyWidget.startButton.Hide()
}

func (keyWidget *KeyWidget) setKeyingStopped() {
	keyWidget.userIsKeying = false
	keyWidget.keyDiv.ReplaceClass("user-keying", "user-not-keying")
	keyWidget.stopButton.Hide()
	keyWidget.startButton.Show()
}
