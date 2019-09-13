package keypracticepanel

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
	keyPracticePanel js.Value
}

func (group *panelGroup) defineMembers() (err error) {

	defer func() {
		if err != nil {
			err = errors.WithMessage(err, "(group *panelGroup) defineMembers()")
		}
	}()

	if group.keyPracticePanel = notJS.GetElementByID("tabsMasterView_home_pad_KeyButton_KeyReadyPanel_tab_bar-KeyPracticeTabPanel-inner-KeyPracticePanel"); group.keyPracticePanel == null {
		err = errors.New("unable to find #tabsMasterView_home_pad_KeyButton_KeyReadyPanel_tab_bar-KeyPracticeTabPanel-inner-KeyPracticePanel")
		return
	}

	return
}

/*
	Show panel funcs.

	Call these from the controller, presenter and caller.
*/

// showKeyPracticePanel shows the panel you named KeyPracticePanel while hiding any other panels in this panel group.
// This panel will become visible only when this group of panels becomes visible.
/* Your note for this panel is:
A page to let the user Key without recording the results to the repo.
*/
func (group *panelGroup) showKeyPracticePanel() {
	tools.ShowPanelInTabGroup(group.keyPracticePanel)
}

