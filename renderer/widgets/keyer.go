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

	userIsKeying bool
	keyTime      time.Time
	durations    []time.Duration

	tools *viewtools.Tools // see /renderer/viewtools
	notJS *notjs.NotJS

	keyCodes [][]*types.KeyCodeRecord
	wpm      uint64
}

// UserKeyChecker processes keying millisecons.
type UserKeyChecker interface {
	GetKeyCodesWPM()
	CheckUserKey(milliSeconds []int64, solution [][]*types.KeyCodeRecord, wpm uint64)
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
		durations:      make([]time.Duration, 0, 1024),
	}

	notJS.SetInnerText(heading, initialInstructions)

	cb := notJS.RegisterEventCallBack(false, false, false, keyWidget.handleStart)
	notJS.SetOnClick(keyWidget.startButton, cb)

	cb = notJS.RegisterEventCallBack(false, false, false, keyWidget.handleStop)
	notJS.SetOnClick(keyWidget.stopButton, cb)

	cb = notJS.RegisterEventCallBack(false, false, false, keyWidget.handleMouseEnter)
	notJS.SetOnMouseEnter(keyWidget.keyDiv, cb)

	cb = notJS.RegisterEventCallBack(false, false, false, keyWidget.handleMouseLeave)
	notJS.SetOnMouseLeave(keyWidget.keyDiv, cb)

	cb = notJS.RegisterEventCallBack(false, false, false, keyWidget.handleMouseDown)
	notJS.SetOnMouseDown(keyWidget.keyDiv, cb)

	cb = notJS.RegisterEventCallBack(false, false, false, keyWidget.handleMouseUp)
	notJS.SetOnMouseUp(keyWidget.keyDiv, cb)

	return
}

// SetKeyCodesWPM sets the key codes, ( the solution ) that the user must key.
func (keyWidget *KeyWidget) SetKeyCodesWPM(records [][]*types.KeyCodeRecord, wpm uint64) {
	keyWidget.keyCodes = records
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

func (keyWidget *KeyWidget) handleStart(event js.Value) {
	tools := keyWidget.tools
	notJS := keyWidget.notJS
	keyWidget.userKeyChecker.GetKeyCodesWPM()
	tools.ElementHide(keyWidget.startButton)
	notJS.SetInnerText(keyWidget.heading, mouseNotOverInstructions)
	notJS.RemoveChildNodes(keyWidget.copyDiv)
	keyWidget.keyTime = time.Now()
}

func (keyWidget *KeyWidget) handleStop(event js.Value) {
	keyWidget.setKeyingStopped()
	l := len(keyWidget.durations)
	milliSeconds := make([]int64, l, l)
	for i, d := range keyWidget.durations {
		milliSeconds[i] = int64(d) / 1000000
	}
	keyWidget.userKeyChecker.CheckUserKey(milliSeconds, keyWidget.keyCodes, keyWidget.wpm)
	keyWidget.tools.ElementHide(keyWidget.stopButton)
	keyWidget.notJS.SetInnerText(keyWidget.heading, noInstructions)
}

func (keyWidget *KeyWidget) handleMouseEnter(event js.Value) {
	if keyWidget.userIsKeying {
		notJS := keyWidget.notJS
		notJS.ClassListReplaceClass(keyWidget.keyDiv, "user-not-key-over", "user-key-over")
		keyWidget.notJS.SetInnerText(keyWidget.heading, mouseOverInstructions)
	}
}

func (keyWidget *KeyWidget) handleMouseLeave(event js.Value) {
	if keyWidget.userIsKeying {
		keyWidget.notJS.ClassListReplaceClass(keyWidget.keyDiv, "user-key-over", "user-not-key-over")
		keyWidget.notJS.SetInnerText(keyWidget.heading, mouseNotOverAgainInstructions)
	}
}

func (keyWidget *KeyWidget) handleMouseDown(event js.Value) {
	if keyWidget.userIsKeying {
		keyWidget.durations = append(keyWidget.durations, time.Since(keyWidget.keyTime))
		keyWidget.keyTime = time.Now()
	}
}

func (keyWidget *KeyWidget) handleMouseUp(event js.Value) {
	if keyWidget.userIsKeying {
		since := time.Since(keyWidget.keyTime)
		keyWidget.durations = append(keyWidget.durations, since)
		keyWidget.keyTime = time.Now()
	}
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
	notJS := keyWidget.notJS
	keyDiv := keyWidget.keyDiv
	// clear the div
	notJS.RemoveChildNodes(keyDiv)
	p := notJS.CreateElementP()
	lineLength := 0
	for _, wordCode := range keyWidget.keyCodes {
		l := len(wordCode)
		word := make([]string, l, l)
		for j, charCode := range wordCode {
			word[j] = charCode.Character
		}
		tn := notJS.CreateTextNode(strings.Join(word, ""))
		notJS.AppendChild(p, tn)
		lineLength += l
		if lineLength >= 40 {
			lineLength = 0
			notJS.AppendChild(p, notJS.CreateElementBR())
		} else {
			notJS.AppendChild(p, notJS.CreateTextNode(" "))
		}
	}
	notJS.AppendChild(keyDiv, p)
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
