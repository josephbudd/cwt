// +build js, wasm

package keytestpanel

import (
	"fmt"
	"syscall/js"

	"github.com/josephbudd/cwt/rendererprocess/api/markup"
	"github.com/josephbudd/cwt/rendererprocess/framework/viewtools"
)

/*

	DO NOT EDIT THIS FILE.

	This file is generated by kickasm and regenerated by rekickasm.

*/

// panelGroup is a group of 1 panel.
// It also has a show panel func for each panel in this panel group.
type panelGroup struct {
	keyTestPanel js.Value
}

func (group *panelGroup) defineMembers() (err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("(group *panelGroup) defineMembers(): %w", err)
		}
	}()

    var panel *markup.Element
 if panel = document.ElementByID("mainMasterView_home_pad_KeyButton_KeyReadyPanel_tab_bar-KeyTestTabPanel-inner-KeyTestPanel"); panel == nil {
	err = fmt.Errorf("unable to find #mainMasterView_home_pad_KeyButton_KeyReadyPanel_tab_bar-KeyTestTabPanel-inner-KeyTestPanel")
		return
    }
    group.keyTestPanel = panel.JSValue()

	return
}

/*
	Show panel funcs.

	Call these from the controller, presenter and messenger.
*/

// showKeyTestPanel shows the panel you named KeyTestPanel while hiding any other panels in this panel group.
// This panel will become visible only when this group of panels becomes visible.
/* Your note for this panel is:
A page to let the user Key and record the results to the repo.
*/
func (group *panelGroup) showKeyTestPanel() {
	viewtools.ShowPanelInTabGroup(group.keyTestPanel)
}
