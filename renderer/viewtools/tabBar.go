package viewtools

import (
	"strings"
	"syscall/js"
)

/*
	WARNING:

	DO NOT EDIT THIS FILE.

*/

// ForceTabButtonClick implements the behavior of a tab button being clicked by the user.
func (tools *Tools) ForceTabButtonClick(button js.Value) {
	tools.handleTabButtonOnClick(button)
}

func (tools *Tools) initializeTabBar() {
	notJS := tools.notJS
	tools.tabberLastPanelID = ""
	tools.tabberLastPanelLevels = make(map[string]string)

	tools.tabberLastPanelLevels["tabsMasterView_home_pad_CopyButton_CopyReadyPanel_tab_bar"] = "tabsMasterView_home_pad_CopyButton_CopyReadyPanel_tab_bar-WPMTabPanel"
	tools.tabberLastPanelLevels["tabsMasterView_home_pad_KeyButton_KeyReadyPanel_tab_bar"] = "tabsMasterView_home_pad_KeyButton_KeyReadyPanel_tab_bar-WPMTabPanel"
	tools.tabberLastPanelLevels["tabsMasterView_home_pad_ReferenceButton_SelectCodesPanel_tab_bar"] = "tabsMasterView_home_pad_ReferenceButton_SelectCodesPanel_tab_bar-LettersTabPanel"
	cb := tools.notJS.RegisterCallBack(
		func(args []js.Value) {
			target := notJS.GetEventTarget(args[0])
			tools.handleTabButtonOnClick(target)
		},
	)
	for id := range tools.tabberLastPanelLevels {
		tabbar := notJS.GetElementByID(id)
		tools.setTabBarOnClicks(tabbar, cb)
	}
}

func (tools *Tools) setTabBarOnClicks(tabbar js.Value, cb js.Callback) {
	notJS := tools.notJS
	children := notJS.ChildrenSlice(tabbar)
	for _, ch := range children {
		if notJS.TagName(ch) == "BUTTON" {
			ch.Set("onclick", cb)
		}
	}
}

func (tools *Tools) handleTabButtonOnClick(button js.Value) {
	if !tools.HandleButtonClick() {
		return
	}
	tools.setTabButtonFocus(button)
	nextpanelid := tools.notJS.ID(button) + "Panel"
	if nextpanelid != tools.tabberLastPanelID {
		// clear this level
		parts := strings.Split(nextpanelid, "-")
		nextpanellevel := parts[0]
		tools.IDHide(tools.tabberLastPanelLevels[nextpanellevel])
		// show the next panel
		tools.IDShow(nextpanelid)
		// remember next panel. it is now the last panel.
		tools.tabberLastPanelID = nextpanelid
		tools.tabberLastPanelLevels[nextpanellevel] = nextpanelid
	}
	tools.SizeApp()
}

func (tools *Tools) setTabButtonFocus(tabinfocus js.Value) {
	// focus the tab now in focus
	notJS := tools.notJS
	notJS.ClassListReplaceClass(tabinfocus, UnSelectedTabClassName, SelectedTabClassName)
	p := notJS.ParentNode(tabinfocus)
	children := notJS.ChildrenSlice(p)
	for _, ch := range children {
		if ch != tabinfocus && notJS.TagName(ch) == "BUTTON" && notJS.ClassListContains(ch, SelectedTabClassName) {
			notJS.ClassListReplaceClass(ch, SelectedTabClassName, UnSelectedTabClassName)
		}
	}
}
