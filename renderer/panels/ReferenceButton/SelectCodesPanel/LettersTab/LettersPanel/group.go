package letterspanel

import (
	"syscall/js"

	"github.com/pkg/errors"
)

/*

	DO NOT EDIT THIS FILE.

	This file is generated by kickasm and regenerated by rekickasm.

*/

// panelGroup is a group of 1 panel.
// It also has a show panel func for each panel in this panel group.
type panelGroup struct {
	lettersPanel js.Value
}

func (group *panelGroup) defineMembers() (err error) {

	defer func() {
		if err != nil {
			err = errors.WithMessage(err, "(group *panelGroup) defineMembers()")
		}
	}()

	if group.lettersPanel = notJS.GetElementByID("tabsMasterView_home_pad_ReferenceButton_SelectCodesPanel_tab_bar-LettersTabPanel-inner-LettersPanel"); group.lettersPanel == null {
		err = errors.New("unable to find #tabsMasterView_home_pad_ReferenceButton_SelectCodesPanel_tab_bar-LettersTabPanel-inner-LettersPanel")
		return
	}

	return
}

/*
	Show panel funcs.

	Call these from the controller, presenter and caller.
*/

// showLettersPanel shows the panel you named LettersPanel while hiding any other panels in this panel group.
// This panel will become visible only when this group of panels becomes visible.
/* Your note for this panel is:
This panel is displayed when the "Letters" tab button is clicked.
This panel is the only panel in it's panel group.

*/
func (group *panelGroup) showLettersPanel() {
	tools.ShowPanelInTabGroup(group.lettersPanel)
}
