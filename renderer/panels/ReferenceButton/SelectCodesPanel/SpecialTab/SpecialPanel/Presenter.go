package specialpanel

import (
	"fmt"
	"sort"
	"strings"
	"syscall/js"

	"github.com/josephbudd/cwt/domain/store/record"
	"github.com/pkg/errors"
)

/*

	Panel name: SpecialPanel

*/

// panelPresenter writes to the panel
type panelPresenter struct {
	group          *panelGroup
	controller     *panelController
	caller         *panelCaller
	tabPanelHeader js.Value

	/* NOTE TO DEVELOPER: Step 1 of 3.

	// Declare your panelPresenter members here.

	*/

	tableBody   js.Value
	initialized bool
}

// defineMembers defines the panelPresenter members by their html elements.
// Returns the error.
func (presenter *panelPresenter) defineMembers() (err error) {

	defer func() {
		if err != nil {
			err = errors.WithMessage(err, "(presenter *panelPresenter) defineMembers()")
		}
	}()

	/* NOTE TO DEVELOPER. Step 2 of 3.

	// Define your panelPresenter members.

	*/

	presenter.tableBody = notJS.GetElementByID("specialsTableBody")
	if presenter.tableBody == null {
		tools.Error("unable to find specialsTableBody")
		return
	}

	return
}

// Tab panel heading.

func (presenter *panelPresenter) getTabPanelHeading() (heading string) {
	heading = notJS.GetInnerText(presenter.tabPanelHeader)
	return
}

func (presenter *panelPresenter) setTabPanelHeading(heading string) {
	heading = strings.TrimSpace(heading)
	if len(heading) == 0 {
		tools.ElementHide(presenter.tabPanelHeader)
	} else {
		tools.ElementShow(presenter.tabPanelHeader)
	}
	notJS.SetInnerText(presenter.tabPanelHeader, heading)
}

/* NOTE TO DEVELOPER. Step 3 of 3.

// Define your panelPresenter functions.

*/

func (presenter *panelPresenter) recordIDToCheckBoxID(recordID uint64) string {
	return fmt.Sprintf("referenceKeyCodeCheckBox%d", recordID)
}

func (presenter *panelPresenter) recordIDToKeyPercentCorrectID(recordID, wpm uint64) string {
	return fmt.Sprintf("referenceKeyPercentCorrect%d-%d", wpm, recordID)
}

func (presenter *panelPresenter) fillTable(rr []*record.KeyCode) {
	if !presenter.initialized {
		presenter.initialFillTable(rr)
		presenter.initialized = true
	} else {
		presenter.refillTable(rr)
	}
}

func (presenter *panelPresenter) refillTable(rr []*record.KeyCode) {
	r := rr[0]
	wpmSorted := make([]uint64, 0, 10)
	for wpm := range r.KeyWPMResults {
		wpmSorted = append(wpmSorted, wpm)
	}
	for _, r := range rr {
		for _, wpm := range wpmSorted {
			id := presenter.recordIDToKeyPercentCorrectID(r.ID, wpm)
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

func (presenter *panelPresenter) initialFillTable(rr []*record.KeyCode) {
	r := rr[0]
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
	notJS.AppendChild(presenter.tableBody, tr)
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
	notJS.AppendChild(presenter.tableBody, tr)
	for _, r := range rr {
		tr = notJS.CreateElementTR()
		// name column
		td := notJS.CreateElementTD()
		notJS.ClassListAddClass(td, "name")
		checkBox := notJS.CreateElementCheckBox()
		notJS.SetAttributeUint64(checkBox, recordIDAttribute, r.ID)
		cbID := presenter.recordIDToCheckBoxID(r.ID)
		notJS.SetID(checkBox, cbID)
		notJS.SetChecked(checkBox, r.Selected)
		notJS.AppendChild(td, checkBox)
		label := notJS.CreateElementLABEL()
		notJS.SetAttribute(label, "for", cbID)
		notJS.AppendChild(label, notJS.CreateTextNode(r.Name))
		notJS.AppendChild(td, label)
		notJS.AppendChild(tr, td)
		// ditdah column
		td = notJS.CreateElementTD()
		notJS.ClassListAddClass(td, "code")
		tn := notJS.CreateTextNode(r.DitDah)
		notJS.AppendChild(td, tn)
		notJS.AppendChild(tr, td)
		// stats
		for _, wpm := range wpmSorted {
			td = notJS.CreateElementTD()
			notJS.SetID(td, presenter.recordIDToKeyPercentCorrectID(r.ID, wpm))
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
		notJS.AppendChild(presenter.tableBody, tr)
	}
}
