package PunctuationPanel

import (
	"fmt"
	"sort"
	"syscall/js"

	"github.com/pkg/errors"

	"github.com/josephbudd/cwt/domain/types"
	"github.com/josephbudd/cwt/renderer/notjs"
	"github.com/josephbudd/cwt/renderer/viewtools"
)

/*

	Panel name: PunctuationPanel

*/

// Presenter writes to the panel
type Presenter struct {
	panelGroup *PanelGroup
	controler  *Controler
	caller     *Caller
	tools      *viewtools.Tools // see /renderer/viewtools
	notJS      *notjs.NotJS

	/* NOTE TO DEVELOPER: Step 1 of 3.

	// Declare your Presenter members here.

	*/

	tableBody   js.Value
	initialized bool
}

// defineMembers defines the Presenter members by their html elements.
// Returns the error.
func (panelPresenter *Presenter) defineMembers() (err error) {

	defer func() {
		if err != nil {
			err = errors.WithMessage(err, "(panelPresenter *Presenter) defineMembers()")
		}
	}()

	/* NOTE TO DEVELOPER. Step 2 of 3.

	// Define your Presenter members.

	*/

	notJS := panelPresenter.notJS
	tools := panelPresenter.tools
	null := js.Null()

	panelPresenter.tableBody = notJS.GetElementByID("punctuationsTableBody")
	if panelPresenter.tableBody == null {
		tools.Error("unable to find punctuationsTableBody")
		return
	}

	return
}

/* NOTE TO DEVELOPER. Step 3 of 3.

// Define your Presenter functions.

*/

func (panelPresenter *Presenter) recordIDToCheckBoxID(recordID uint64) string {
	return fmt.Sprintf("referenceKeyCodeCheckBox%d", recordID)
}

func (panelPresenter *Presenter) recordIDToKeyPercentCorrectID(recordID, wpm uint64) string {
	return fmt.Sprintf("referenceKeyPercentCorrect%d-%d", wpm, recordID)
}

func (panelPresenter *Presenter) fillTable(records []*types.KeyCodeRecord) {
	if !panelPresenter.initialized {
		panelPresenter.initialFillTable(records)
		panelPresenter.initialized = true
	} else {
		panelPresenter.refillTable(records)
	}
}

func (panelPresenter *Presenter) refillTable(records []*types.KeyCodeRecord) {
	notJS := panelPresenter.notJS
	r := records[0]
	wpmSorted := make([]uint64, 0, 10)
	for wpm := range r.KeyWPMResults {
		wpmSorted = append(wpmSorted, wpm)
	}
	for _, r := range records {
		for _, wpm := range wpmSorted {
			id := panelPresenter.recordIDToKeyPercentCorrectID(r.ID, wpm)
			td := notJS.GetElementByID(id)
			notJS.RemoveChildNodes(td)
			copyResults := r.CopyWPMResults[wpm]
			keyResults := r.KeyWPMResults[wpm]
			// copying
			var copyTN js.Value
			if copyResults.Attempts != 0 {
				perCentCopy := 100 * copyResults.Correct / copyResults.Attempts
				copyTN = notJS.CreateTextNode(fmt.Sprintf("%d%%", perCentCopy))
			} else {
				if keyResults.Attempts == 0 {
					copyTN = notJS.CreateTextNode("")
				} else {
					copyTN = notJS.CreateTextNode("NA")
				}
			}
			notJS.AppendChild(td, copyTN)
			notJS.AppendChild(td, notJS.CreateElementBR())
			// keying
			var keyTN js.Value
			if keyResults.Attempts != 0 {
				perCentKey := 100 * keyResults.Correct / keyResults.Attempts
				keyTN = notJS.CreateTextNode(fmt.Sprintf("%d%%", perCentKey))
			} else {
				if copyResults.Attempts == 0 {
					keyTN = notJS.CreateTextNode("")
				} else {
					keyTN = notJS.CreateTextNode("NA")
				}
			}
			notJS.AppendChild(td, keyTN)
		}
	}
}

func (panelPresenter *Presenter) initialFillTable(records []*types.KeyCodeRecord) {
	notJS := panelPresenter.notJS
	r := records[0]
	wpmSorted := make([]uint64, 0, 10)
	for wpm := range r.KeyWPMResults {
		wpmSorted = append(wpmSorted, wpm)
	}
	sort.Slice(wpmSorted, func(i, j int) bool { return wpmSorted[i] < wpmSorted[j] })
	// first heading row
	tr := notJS.CreateElementTR()
	th := notJS.CreateElementTH()
	notJS.AppendChild(th, notJS.CreateTextNode("Name"))
	notJS.AppendChild(tr, th)
	th = notJS.CreateElementTH()
	notJS.AppendChild(th, notJS.CreateTextNode("Dit Dah"))
	notJS.AppendChild(tr, th)
	th = notJS.CreateElementTH()
	notJS.AppendChild(th, notJS.CreateTextNode("@ WPM: Copy / Key % Correct"))
	notJS.SetAttributeInt(th, "colspan", len(wpmSorted))
	notJS.AppendChild(tr, th)
	notJS.AppendChild(panelPresenter.tableBody, tr)
	// second heading row
	tr = notJS.CreateElementTR()
	th = notJS.CreateElementTH()
	notJS.AppendChild(tr, th)
	th = notJS.CreateElementTH()
	notJS.AppendChild(tr, th)
	for _, wpm := range wpmSorted {
		th = notJS.CreateElementTH()
		notJS.AppendChild(th, notJS.CreateTextNode(fmt.Sprintf("%d", wpm)))
		notJS.AppendChild(th, notJS.CreateElementBR())
		notJS.AppendChild(th, notJS.CreateTextNode("Copy"))
		notJS.AppendChild(th, notJS.CreateElementBR())
		notJS.AppendChild(th, notJS.CreateTextNode("Key"))
		notJS.AppendChild(tr, th)
	}
	notJS.AppendChild(panelPresenter.tableBody, tr)
	for _, r := range records {
		tr = notJS.CreateElementTR()
		// name column
		td := notJS.CreateElementTD()
		checkBox := notJS.CreateElementCheckBox()
		notJS.SetAttributeUint64(checkBox, recordIDAttribute, r.ID)
		cbID := panelPresenter.recordIDToCheckBoxID(r.ID)
		notJS.SetID(checkBox, cbID)
		notJS.SetChecked(checkBox, r.Selected)
		notJS.AppendChild(td, checkBox)
		label := notJS.CreateElementLabelForID(r.Name, cbID)
		notJS.AppendChild(td, label)
		notJS.AppendChild(tr, td)
		// ditdah column
		td = notJS.CreateElementTD()
		notJS.ClassListAddClass(td, "code")
		tn := notJS.CreateTextNode(r.DitDah)
		notJS.AppendChild(td, tn)
		notJS.AppendChild(tr, td)
		for _, wpm := range wpmSorted {
			td = notJS.CreateElementTD()
			notJS.SetID(td, panelPresenter.recordIDToKeyPercentCorrectID(r.ID, wpm))
			copyResults := r.CopyWPMResults[wpm]
			keyResults := r.KeyWPMResults[wpm]
			// copying
			var copyTN js.Value
			if copyResults.Attempts != 0 {
				perCentCopy := 100 * copyResults.Correct / copyResults.Attempts
				copyTN = notJS.CreateTextNode(fmt.Sprintf("%d%%", perCentCopy))
			} else {
				if keyResults.Attempts == 0 {
					copyTN = notJS.CreateTextNode("")
				} else {
					copyTN = notJS.CreateTextNode("NA")
				}
			}
			notJS.AppendChild(td, copyTN)
			notJS.AppendChild(td, notJS.CreateElementBR())
			// keying
			var keyTN js.Value
			if keyResults.Attempts != 0 {
				perCentKey := 100 * keyResults.Correct / keyResults.Attempts
				keyTN = notJS.CreateTextNode(fmt.Sprintf("%d%%", perCentKey))
			} else {
				if copyResults.Attempts == 0 {
					keyTN = notJS.CreateTextNode("")
				} else {
					keyTN = notJS.CreateTextNode("NA")
				}
			}
			notJS.AppendChild(td, keyTN)
			notJS.AppendChild(tr, td)
		}
		notJS.AppendChild(panelPresenter.tableBody, tr)
	}
}
