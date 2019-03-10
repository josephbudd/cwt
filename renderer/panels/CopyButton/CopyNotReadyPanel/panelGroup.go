package copynotreadypanel

import (
	"syscall/js"

	"github.com/pkg/errors"

	"github.com/josephbudd/cwt/renderer/notjs"
	"github.com/josephbudd/cwt/renderer/viewtools"
)

// PanelGroup is a group of 2 panels.
// It also has show panel funcs for each panel in this panel group.
type PanelGroup struct {
	tools *viewtools.Tools
	notJS *notjs.NotJS

	copyNotReadyPanel js.Value
	copyReadyPanel js.Value
}

func (panelGroup *PanelGroup) defineMembers() (err error) {

	defer func() {
		if err != nil {
			err = errors.WithMessage(err, "(panelGroup *PanelGroup) defineMembers()")
		}
	}()

	notJS := panelGroup.notJS
	null := js.Null()

	if panelGroup.copyNotReadyPanel = notJS.GetElementByID("tabsMasterView-home-pad-CopyButton-CopyNotReadyPanel"); panelGroup.copyNotReadyPanel == null {
		err = errors.New("unable to find #tabsMasterView-home-pad-CopyButton-CopyNotReadyPanel")
		return
	}
	if panelGroup.copyReadyPanel = notJS.GetElementByID("tabsMasterView-home-pad-CopyButton-CopyReadyPanel"); panelGroup.copyReadyPanel == null {
		err = errors.New("unable to find #tabsMasterView-home-pad-CopyButton-CopyReadyPanel")
		return
	}


	return
}

/*
	Show panel funcs.

	Call these from the controler, presenter and caller.
*/

// showCopyNotReadyPanel shows the panel you named CopyNotReadyPanel while hiding any other panels in this panel group.
// This panel's id is tabsMasterView-home-pad-CopyButton-CopyNotReadyPanel.
// This panel either becomes visible immediately or whenever this group of panels is made visible.  Whenever could be immediately if this panel group is currently visible.
// Param force boolean effects when this panel becomes visible.
//  * if force is true then
//    immediately if the home button pad is not currently displayed;
//    whenever if the home button pad is currently displayed.
//  * if force is false then whenever.
/* Your note for this panel is:
Shown when no codes are selected in references.
Display a message telling the user that no codes are selected.

*/
func (panelGroup *PanelGroup) showCopyNotReadyPanel(force bool) {
	panelGroup.tools.ShowPanelInButtonGroup(panelGroup.copyNotReadyPanel, force)
}

// showCopyReadyPanel shows the panel you named CopyReadyPanel while hiding any other panels in this panel group.
// That panel's id is tabsMasterView-home-pad-CopyButton-CopyReadyPanel.
// That panel either becomes visible immediately or whenever this group of panels is made visible.  Whenever could be immediately if this panel group is currently visible.
// Param force boolean effects when that panel becomes visible.
//  * if force is true then
//    immediately if the home button pad is not currently displayed;
//    whenever if the home button pad is currently displayed.
//  * if force is false then whenever.
/* Your note for that panel is:

*/
func (panelGroup *PanelGroup) showCopyReadyPanel(force bool) {
	panelGroup.tools.ShowPanelInButtonGroup(panelGroup.copyReadyPanel, force)
}

