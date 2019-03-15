package widgets

import (
	"fmt"
	"strings"
	"syscall/js"
	"time"

	"github.com/josephbudd/cwt/domain/types"
	"github.com/josephbudd/cwt/renderer/notjs"
	"github.com/josephbudd/cwt/renderer/viewtools"
)

const (
	initialInstructions           = "Click start to begin."
	mouseNotOverInstructions      = "Mouse over the area below to key."
	mouseNotOverAgainInstructions = "Mouse over the area below to key or click the [Check] button."
	mouseOverInstructions         = "Ready for you to key."
	noInstructions                = ""
)

// KeyWidget is a widget of buttons and key areas which record keying.
type KeyWidget struct {
	heading     js.Value
	startButton js.Value
	stopButton  js.Value
	keyDiv      js.Value
	copyDiv     js.Value

	userKeyChecker UserKeyChecker
	metronomer     Metronomer

	userIsKeying bool
	keyTime      time.Time
	durations    []time.Duration

	tools *viewtools.Tools // see /renderer/viewtools
	notJS *notjs.NotJS

	keyCodes [][]*types.KeyCodeRecord
	help     [][]types.HowTo
	wpm      uint64
	times    []time.Time
}

// UserKeyChecker processes keying millisecons.
type UserKeyChecker interface {
	GetKeyCodesWPM()
	CheckUserKey(milliSeconds []int64, solution [][]*types.KeyCodeRecord, wpm uint64)
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
func NewKeyWidget(heading js.Value,
	startButton, stopButton js.Value,
	keyDiv, copyDiv js.Value,
	userKeyChecker UserKeyChecker,
	metronomer Metronomer,
	tools *viewtools.Tools,
	notJS *notjs.NotJS) (keyWidget *KeyWidget) {
	keyWidget = &KeyWidget{
		heading:     heading,
		startButton: startButton,
		stopButton:  stopButton,
		keyDiv:      keyDiv,
		copyDiv:     copyDiv,

		tools: tools,
		notJS: notJS,

		userKeyChecker: userKeyChecker,
		metronomer:     metronomer,

		times:     make([]time.Time, 0, 1024),
		durations: make([]time.Duration, 0, 1024),
	}

	notJS.SetInnerText(heading, initialInstructions)

	cb := tools.RegisterEventCallBack(keyWidget.handleStart, true, true, true)
	notJS.SetOnClick(keyWidget.startButton, cb)

	cb = tools.RegisterEventCallBack(keyWidget.handleStop, true, true, true)
	notJS.SetOnClick(keyWidget.stopButton, cb)

	cb = tools.RegisterEventCallBack(keyWidget.handleMouseEnter, true, true, true)
	notJS.SetOnMouseEnter(keyWidget.keyDiv, cb)

	cb = tools.RegisterEventCallBack(keyWidget.handleMouseLeave, true, true, true)
	notJS.SetOnMouseLeave(keyWidget.keyDiv, cb)

	cb = tools.RegisterEventCallBack(keyWidget.handleMouseDown, true, true, true)
	notJS.SetOnMouseDown(keyWidget.keyDiv, cb)

	cb = tools.RegisterEventCallBack(keyWidget.handleMouseUp, true, true, true)
	notJS.SetOnMouseUp(keyWidget.keyDiv, cb)

	return
}

// SetKeyCodesHelpWPM sets the key codes, ( the solution ) that the user must key.
func (keyWidget *KeyWidget) SetKeyCodesHelpWPM(records [][]*types.KeyCodeRecord, help [][]types.HowTo, wpm uint64) {
	keyWidget.keyCodes = records
	keyWidget.help = help
	keyWidget.wpm = wpm
	keyWidget.tellUserToStart()
}

// SetKeyCodesWPM sets the key codes, ( the solution ) that the user must key.
func (keyWidget *KeyWidget) SetKeyCodesWPM(records [][]*types.KeyCodeRecord, wpm uint64) {
	keyWidget.keyCodes = records
	keyWidget.help = nil
	keyWidget.wpm = wpm
	keyWidget.tellUserToStart()
}

// ShowResults Processes and displays the checked results of what the user keyed.
func (keyWidget *KeyWidget) ShowResults(correctCount, incorrectCount, keyedCount uint64, testResults [][]types.TestResult) {
	copyDiv := keyWidget.copyDiv
	notJS := keyWidget.notJS
	// clear the div
	notJS.RemoveChildNodes(copyDiv)
	// heading
	h3 := notJS.CreateElementH3()
	tn := notJS.CreateTextNode(fmt.Sprintf("You got %d correct and %d incorrect out of %d.", correctCount, incorrectCount, keyedCount))
	notJS.AppendChild(h3, tn)
	notJS.AppendChild(copyDiv, h3)
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
	keyWidget.tools.SizeApp()
	keyWidget.tools.ElementShow(keyWidget.startButton)
}

func (keyWidget *KeyWidget) addTable(controlWord, ditDahWord, copiedWord string) {
	notJS := keyWidget.notJS

	var tdClass string
	if controlWord == copiedWord {
		tdClass = "match"
	} else {
		tdClass = "mismatch"
	}

	// table
	table := notJS.CreateElementTABLE()
	notJS.ClassListAddClass(table, "resize-me-width")
	tbody := notJS.CreateElementTBODY()
	// control row
	tr := notJS.CreateElementTR()
	th := notJS.CreateElementTH()
	tn := notJS.CreateTextNode("Actual Code")
	notJS.AppendChild(th, tn)
	notJS.AppendChild(tr, th)
	// control characters
	td := notJS.CreateElementTD()
	notJS.AppendChild(td, notJS.CreateTextNode(controlWord))
	notJS.ClassListAddClass(td, tdClass)
	notJS.ClassListAddClass(td, "control")
	notJS.AppendChild(tr, td)
	notJS.AppendChild(tbody, tr)

	// keyed row
	tr = notJS.CreateElementTR()
	th = notJS.CreateElementTH()
	tn = notJS.CreateTextNode("You Keyed")
	notJS.AppendChild(th, tn)
	notJS.AppendChild(tr, th)
	// user keyed
	td = notJS.CreateElementTD()
	notJS.AppendChild(td, notJS.CreateTextNode(ditDahWord))
	notJS.ClassListAddClass(td, tdClass)
	notJS.ClassListAddClass(td, "ditdah")
	notJS.AppendChild(tr, td)
	notJS.AppendChild(tbody, tr)

	// copied row
	tr = notJS.CreateElementTR()
	th = notJS.CreateElementTH()
	tn = notJS.CreateTextNode("App Copied")
	notJS.AppendChild(th, tn)
	notJS.AppendChild(tr, th)
	// app copied
	td = notJS.CreateElementTD()
	notJS.AppendChild(td, notJS.CreateTextNode(copiedWord))
	notJS.ClassListAddClass(td, tdClass)
	notJS.ClassListAddClass(td, "copied")
	notJS.AppendChild(tr, td)
	notJS.AppendChild(tbody, tr)

	notJS.AppendChild(table, tbody)
	notJS.AppendChild(keyWidget.copyDiv, table)
}

// handlers

func (keyWidget *KeyWidget) handleStart(event js.Value) interface{} {
	tools := keyWidget.tools
	notJS := keyWidget.notJS
	keyWidget.userKeyChecker.GetKeyCodesWPM()
	tools.ElementHide(keyWidget.startButton)
	notJS.SetInnerText(keyWidget.heading, mouseNotOverInstructions)
	notJS.RemoveChildNodes(keyWidget.copyDiv)
	keyWidget.keyTime = time.Now()
	keyWidget.times = make([]time.Time, 1, 200)
	keyWidget.times[0] = time.Now()
	return nil
}

func (keyWidget *KeyWidget) handleStop(event js.Value) interface{} {
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
	keyWidget.tools.ElementHide(keyWidget.stopButton)
	keyWidget.notJS.SetInnerText(keyWidget.heading, noInstructions)
	return nil
}

func (keyWidget *KeyWidget) handleMouseEnter(event js.Value) interface{} {
	if keyWidget.userIsKeying {
		notJS := keyWidget.notJS
		notJS.ClassListReplaceClass(keyWidget.keyDiv, "user-not-key-over", "user-key-over")
		keyWidget.notJS.SetInnerText(keyWidget.heading, mouseOverInstructions)
		if keyWidget.metronomer != nil {
			keyWidget.metronomer.StartMetronome(keyWidget.wpm)
		}
	}
	return nil
}

func (keyWidget *KeyWidget) handleMouseLeave(event js.Value) interface{} {
	if keyWidget.userIsKeying {
		keyWidget.notJS.ClassListReplaceClass(keyWidget.keyDiv, "user-key-over", "user-not-key-over")
		keyWidget.notJS.SetInnerText(keyWidget.heading, mouseNotOverAgainInstructions)
		if keyWidget.metronomer != nil {
			keyWidget.metronomer.StopMetronome()
		}
	}
	return nil
}

func (keyWidget *KeyWidget) handleMouseDown(event js.Value) interface{} {
	if keyWidget.userIsKeying {
		keyWidget.times = append(keyWidget.times, time.Now())
		keyWidget.notJS.ClassListReplaceClass(keyWidget.keyDiv, "user-key-up", "user-key-down")
	}
	return nil
}

func (keyWidget *KeyWidget) handleMouseUp(event js.Value) interface{} {
	if keyWidget.userIsKeying {
		keyWidget.times = append(keyWidget.times, time.Now())
		keyWidget.notJS.ClassListReplaceClass(keyWidget.keyDiv, "user-key-down", "user-key-up")
	}
	return nil
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
	notJS := keyWidget.notJS
	keyDiv := keyWidget.keyDiv
	// clear the div
	notJS.RemoveChildNodes(keyDiv)
	// put the lines in a paragraph.
	p := notJS.CreateElementP()
	var b strings.Builder
	// the words each on a separate line.
	for _, wordCode := range keyWidget.keyCodes {
		for _, charCode := range wordCode {
			fmt.Fprint(&b, charCode.Character)
		}
		tn := notJS.CreateTextNode(b.String())
		notJS.AppendChild(p, tn)
		notJS.AppendChild(p, notJS.CreateElementBR())
		b.Reset()
	}
	notJS.AppendChild(keyDiv, p)
}

func (keyWidget *KeyWidget) showTextHelpToKey() {
	notJS := keyWidget.notJS
	keyDiv := keyWidget.keyDiv
	// clear the div
	notJS.RemoveChildNodes(keyDiv)
	for i, wordCode := range keyWidget.keyCodes {
		// the word
		l := len(wordCode)
		chars := make([]string, l, l)
		for j, charCode := range wordCode {
			chars[j] = charCode.Character
		}
		word := strings.Join(chars, "")
		table := notJS.CreateElementTABLE()
		tbody := notJS.CreateElementTBODY()
		notJS.AppendChild(table, tbody)
		tr := notJS.CreateElementTR()
		td := notJS.CreateElementTD()
		notJS.SetAttributeInt(td, "colspan", 3)
		tn := notJS.CreateTextNode(word)
		notJS.AppendChild(td, tn)
		notJS.AppendChild(tr, td)
		notJS.AppendChild(tbody, tr)
		// help
		wordHelp := keyWidget.help[i]
		for _, h := range wordHelp {
			tr = notJS.CreateElementTR()
			// character
			td = notJS.CreateElementTD()
			tn = notJS.CreateTextNode(h.Character)
			notJS.AppendChild(td, tn)
			notJS.AppendChild(tr, td)
			// dit dah
			td = notJS.CreateElementTD()
			tn = notJS.CreateTextNode(h.DitDah)
			notJS.AppendChild(td, tn)
			notJS.AppendChild(tr, td)
			// instructions
			td = notJS.CreateElementTD()
			tn = notJS.CreateTextNode(h.Instructions)
			notJS.AppendChild(td, tn)
			notJS.AppendChild(tr, td)
			notJS.AppendChild(tbody, tr)
		}
		notJS.AppendChild(table, tbody)
		notJS.AppendChild(keyDiv, table)
	}
}

func (keyWidget *KeyWidget) setKeyingStarted() {
	keyWidget.userIsKeying = true
	keyWidget.notJS.ClassListReplaceClass(keyWidget.keyDiv, "user-not-keying", "user-keying")
	keyWidget.tools.ElementShow(keyWidget.stopButton)
	keyWidget.tools.ElementHide(keyWidget.startButton)
}

func (keyWidget *KeyWidget) setKeyingStopped() {
	keyWidget.userIsKeying = false
	keyWidget.notJS.ClassListReplaceClass(keyWidget.keyDiv, "user-keying", "user-not-keying")
	keyWidget.tools.ElementHide(keyWidget.stopButton)
	keyWidget.tools.ElementShow(keyWidget.startButton)
}
