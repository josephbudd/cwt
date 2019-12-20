// +build js, wasm

package letterspanel

import (
	"fmt"
	"sort"
	"strings"

	"github.com/josephbudd/cwt/domain/store/record"
	"github.com/josephbudd/cwt/rendererprocess/api/display"
	"github.com/josephbudd/cwt/rendererprocess/api/markup"
)

/*

	Panel name: LettersPanel

*/

// panelPresenter writes to the panel
type panelPresenter struct {
	group          *panelGroup
	controller     *panelController
	messenger      *panelMessenger
	tabPanelHeader *markup.Element
	tabButton      *markup.Element

	/* NOTE TO DEVELOPER: Step 1 of 3.

	// Declare your panelPresenter members here.

	// example:

	editCustomerName *markup.Element

	*/

	tableBody   *markup.Element
	initialized bool
}

// defineMembers defines the panelPresenter members by their html elements.
// Returns the error.
func (presenter *panelPresenter) defineMembers() (err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("(presenter *panelPresenter) defineMembers(): %w", err)
		}
	}()

	/* NOTE TO DEVELOPER. Step 2 of 3.

	// Define your panelPresenter members.

	// example:

	// Define the edit form's customer name input field.
	if presenter.editCustomerName = document.ElementByID("editCustomerName"); presenter.editCustomerName == nil {
		err = fmt.Errorf("unable to find #editCustomerName")
		return
	}

	*/

	if presenter.tableBody = document.ElementByID("lettersTableBody"); presenter.tableBody == nil {
		display.Error("unable to find lettersTableBody")
		return
	}

	return
}

// Tab button label.

func (presenter *panelPresenter) getTabLabel() (label string) {
	label = presenter.tabButton.InnerText()
	return
}

func (presenter *panelPresenter) setTabLabel(label string) {
	presenter.tabButton.SetInnerText(label)
}

// Tab panel heading.

func (presenter *panelPresenter) getTabPanelHeading() (heading string) {
	heading = presenter.tabPanelHeader.InnerText()
	return
}

func (presenter *panelPresenter) setTabPanelHeading(heading string) {
	heading = strings.TrimSpace(heading)
	if len(heading) == 0 {
		presenter.tabPanelHeader.Hide()
	} else {
		presenter.tabPanelHeader.Show()
	}
	presenter.tabPanelHeader.SetInnerText(heading)
}

/* NOTE TO DEVELOPER. Step 3 of 3.

// Define your panelPresenter functions.

// example:

// displayCustomer displays the customer in the edit customer form panel.
func (presenter *panelPresenter) displayCustomer(record *types.CustomerRecord) {
	presenter.editCustomerName.SetValue(record.Name)
}

*/

func (presenter *panelPresenter) recordIDToCheckBoxID(recordID uint64) string {
	return fmt.Sprintf("referenceKeyCodeCheckBox%d", recordID)
}

func (presenter *panelPresenter) recordIDToKeyPercentCorrectID(recordID, wpm uint64) string {
	return fmt.Sprintf("referenceKeyPercentCorrect%d-%d", wpm, recordID)
}

func (presenter *panelPresenter) fillTable(records []*record.KeyCode) {
	if !presenter.initialized {
		presenter.initialFillTable(records)
		presenter.initialized = true
	} else {
		presenter.refillTable(records)
	}
}

func (presenter *panelPresenter) refillTable(records []*record.KeyCode) {
	r := records[0]
	wpmSorted := make([]uint64, 0, 10)
	for wpm := range r.KeyWPMResults {
		wpmSorted = append(wpmSorted, wpm)
	}
	for _, r := range records {
		for _, wpm := range wpmSorted {
			id := presenter.recordIDToKeyPercentCorrectID(r.ID, wpm)
			td := document.ElementByID(id)
			td.RemoveChildren()
			copyResults := r.CopyWPMResults[wpm]
			keyResults := r.KeyWPMResults[wpm]
			// copying
			if copyResults.Attempts != 0 {
				perCentCopy := 100 * copyResults.Correct / copyResults.Attempts
				td.AppendText(fmt.Sprintf("%d%%", perCentCopy))
			} else {
				if keyResults.Attempts == 0 {
					td.AppendText("")
				} else {
					td.AppendText("NA")
				}
			}
			td.AppendChild(document.NewBR())
			// keying
			var keyTN *markup.Element
			if keyResults.Attempts != 0 {
				perCentKey := 100 * keyResults.Correct / keyResults.Attempts
				keyTN = document.NewText(fmt.Sprintf("%d%%", perCentKey))
			} else {
				if copyResults.Attempts == 0 {
					keyTN = document.NewText("")
				} else {
					keyTN = document.NewText("NA")
				}
			}
			td.AppendChild(keyTN)
		}
	}
}

func (presenter *panelPresenter) initialFillTable(records []*record.KeyCode) {
	r := records[0]
	wpmSorted := make([]uint64, 0, 10)
	for wpm := range r.KeyWPMResults {
		wpmSorted = append(wpmSorted, wpm)
	}
	sort.Slice(wpmSorted, func(i, j int) bool { return wpmSorted[i] < wpmSorted[j] })
	// first heading row
	tr := document.NewTR()
	th := document.NewTH()
	th.SetInnerText("Name")
	tr.AppendChild(th)
	th = document.NewTH()
	th.SetInnerText("Dit Dah")
	tr.AppendChild(th)
	th = document.NewTH()
	th.SetInnerText("@ WPM: Copy / Key % Correct")
	th.SetAttribute("colspan", len(wpmSorted))
	tr.AppendChild(th)
	presenter.tableBody.AppendChild(tr)
	// second heading row
	tr = document.NewTR()
	th = document.NewTH()
	tr.AppendChild(th)
	th = document.NewTH()
	tr.AppendChild(th)
	for _, wpm := range wpmSorted {
		th = document.NewTH()
		th.AppendChild(document.NewText(fmt.Sprintf("%d", wpm)))
		th.AppendChild(document.NewBR())
		th.AppendChild(document.NewText("Copy"))
		th.AppendChild(document.NewBR())
		th.AppendChild(document.NewText("Key"))
		tr.AppendChild(th)
	}
	presenter.tableBody.AppendChild(tr)
	for _, r := range records {
		tr = document.NewTR()
		// name column
		td := document.NewTD()
		td.AddClass("name")
		checkBox := document.NewCheckBox()
		checkBox.SetAttribute(recordIDAttribute, r.ID)
		cbID := presenter.recordIDToCheckBoxID(r.ID)
		checkBox.SetID(cbID)
		checkBox.SetChecked(r.Selected)
		td.AppendChild(checkBox)
		label := document.NewLABEL()
		label.SetAttribute("for", cbID)
		label.SetInnerText(r.Name)
		td.AppendChild(label)
		tr.AppendChild(td)
		// ditdah column
		td = document.NewTD()
		td.AddClass("code")
		td.SetInnerText(r.DitDah)
		tr.AppendChild(td)
		for _, wpm := range wpmSorted {
			td = document.NewTD()
			td.SetID(presenter.recordIDToKeyPercentCorrectID(r.ID, wpm))
			copyResults := r.CopyWPMResults[wpm]
			keyResults := r.KeyWPMResults[wpm]
			// copying
			var copyTN *markup.Element
			if copyResults.Attempts != 0 {
				// copyTN = document.NewText(fmt.Sprintf("%d / %d", copyResults.Correct, copyResults.Attempts))
				perCentCopy := 100 * copyResults.Correct / copyResults.Attempts
				copyTN = document.NewText(fmt.Sprintf("%d%%", perCentCopy))
			} else {
				if keyResults.Attempts == 0 {
					copyTN = document.NewText("")
				} else {
					copyTN = document.NewText("NA")
				}
			}
			td.AppendChild(copyTN)
			td.AppendChild(document.NewBR())
			// keying
			var keyTN *markup.Element
			if keyResults.Attempts != 0 {
				// keyTN = document.NewText(fmt.Sprintf("%d / %d", keyResults.Correct, keyResults.Attempts))
				perCentKey := 100 * keyResults.Correct / keyResults.Attempts
				keyTN = document.NewText(fmt.Sprintf("%d%%", perCentKey))
			} else {
				if copyResults.Attempts == 0 {
					keyTN = document.NewText("")
				} else {
					keyTN = document.NewText("NA")
				}
			}
			td.AppendChild(keyTN)
			tr.AppendChild(td)
		}
		presenter.tableBody.AppendChild(tr)
	}
}
