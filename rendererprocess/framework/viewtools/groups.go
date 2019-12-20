// +build js, wasm

package viewtools

import (
	"syscall/js"
)

// ShowPanelInButtonGroup shows a panel in a button pad button group and hides the other panels in the group.
func ShowPanelInButtonGroup(panel js.Value, force bool) {
	if force && hereIsVisible() {
		// show this panel.
		// the app is resized by HideShow.
		HideShow(here, panel)
	} else {
		// not forcing so just set the panel to be visible when the user makes it visible.
		_, isVisible := ShowInGroup(panel, ToBeSeenClassName, ToBeUnSeenClassName)
		if isVisible {
			// this panel is visible anyway so resize the app.
			SizeApp()
		}
	}
}

// ShowPanelInTabGroup shows a panel in a tab button group and hides the other panels in the group.
func ShowPanelInTabGroup(panel js.Value) {
	_, isVisible := ShowInGroup(panel, SeenClassName, UnSeenClassName)
	if isVisible {
		// this tab panel is visible anyway so resize the app.
		SizeApp()
	}
}

// ShowInGroup only works with panel groups decended from the slider collection.  ( panels shown with the back button on the left side. )
// It shows one panel while hiding the other panels in a panel group.
// It does so by adding and removing classes to panel classLists.
// Param target is the panel to be shown.
// Param showClass is the class name for showing target.
// Param hideClass is the class name for hiding the other panels in target's group.
// Returns 2 params
// 1. if param target has an ancestor which is the slider collections. ( panels shown with the back button on the left side. )
// 2. if the param target becomes visible.
func ShowInGroup(target js.Value, showClass, hideClass string) (isSliderSub, isVisible bool) {
	targetParent := target.Get("parentNode")
	isSliderSub = targetParent == mainMasterviewHomeSliderCollection
	// tab sibling panels are in sliders but they are special.
	classList := target.Get("classList")
	isTabSibling := classList.Call("contains", SliderPanelInnerSiblingClassName).Bool()
	if !(isSliderSub || isTabSibling) {
		// not in the slider collection
		isSliderSub = (isSliderSub || isTabSibling)
		return
	}
	targetInGroup := false
	var divs []js.Value
	for _, divs = range buttonPanelsMap {
		for _, div := range divs {
			if target == div {
				// target is in this group
				targetInGroup = true
				break
			}
		}
		if targetInGroup {
			// target is in this group
			break
		}
	}
	if !targetInGroup {
		// target not in group so not a slider sub.
		isSliderSub = false
		return
	}
	// yes target is a slider div
	setInGroup(divs, target, showClass, hideClass)
	// check for visibility
	for _, div := range divs {
		classList := div.Get("classList")
		if isVisible = !classList.Call("contains", UnSeenClassName).Bool(); isVisible {
			// only really visible if slider is visible
			classList = mainMasterviewHomeSlider.Get("classList")
			isVisible = !classList.Call("contains", UnSeenClassName).Bool()
			break
		}
	}
	if isVisible {
		// set here
		if isTabSibling {
			// here is never a tab panel sibling.
			// here must be a slider panel.
			// Find the correct ancestor.
			sliderPanel := targetParent
			for {
				classList = sliderPanel.Get("classList")
				if classList.Call("contains", SliderPanelClassName).Bool() {
					break
				}
				sliderPanel = sliderPanel.Get("parentNode")
			}
			here = sliderPanel
		} else {
			// target is the new here.
			here = target
		}
		// here is now this slider sub panel.
		// here is never a tab panel sibling.
		if showClass == ToBeSeenClassName {
			// set target for class seen
			// set non targets for class unseen
			setInGroup(divs, target, SeenClassName, UnSeenClassName)
		}
	}
	isSliderSub = isSliderSub || isTabSibling
	return
}

//HideShow hides one panel and shows another panel.
// both panels must have the parentNode == SliderPresenter.sliderCollection.
func HideShow(hideDiv, showDiv js.Value) {
	// hide the hide div
	isSliderH := hideInGroup(hideDiv, SeenClassName, UnSeenClassName)
	// show the show div
	isSliderS, isVisibleS := ShowInGroup(showDiv, SeenClassName, UnSeenClassName)
	if isSliderS {
		// reset the back button's color class.
		backColorLevel := showDiv.Call("getAttribute", "backColorLevel").String()
		classList := mainMasterviewHomeSliderBack.Get("classList")
		firstClass := classList.Call("item", 0).String()
		classList.Call("replace", firstClass, backColorLevel)
	}
	if isSliderH && isSliderS {
		// the slider was visible for the hideDiv and so it still is for the showDiv
		SizeApp()
		return
	}
	// hideDiv and showDiv are not both sliders
	if !isVisibleS {
		// showDiv, the div to show is not visible
		if isSliderH {
			// hideDiv is in the slider collection
			hideSlider()
		} else {
			// hideDiv is not in the slider collection, its a master div or home or some little thing in a panel
			ElementHide(hideDiv)
		}
		if isSliderS {
			// showDiv is in the slider collection
			showSlider()
		} else {
			// showDiv is not in the slider collection
			ElementShow(showDiv)
		}
	}
	SizeApp()
}

// setInGroup works with a group of panels from buttonPanelsMap.
// It sets target's to setClass and removes unSetClass.
// It sets the other panel's to unSetClass and removes setClass.
func setInGroup(group []js.Value, target js.Value, setClass, unSetClass string) {
	var classList js.Value
	for _, panel := range group {
		if panel != target {
			classList = panel.Get("classList")
			if !classList.Call("replace", setClass, unSetClass).Bool() {
				classList.Call("add", unSetClass)
			}
		}
	}
	classList = target.Get("classList")
	if !classList.Call("replace", unSetClass, setClass).Bool() {
		classList.Call("add", setClass)
	}
}

// hideInGroup hides target in a group.
// Returns is the target is a slider sub panel, a child of the slider collection div.
func hideInGroup(target js.Value, showClass, hideClass string) (isSliderSub bool) {
	parentNode := target.Get("parentNode")
	isSliderSub = parentNode == mainMasterviewHomeSliderCollection
	if !isSliderSub {
		// not in the slider collection.
		return
	}
	classList := target.Get("classList")
	classList.Call("replace", showClass, hideClass)
	return
}

func initializeGroups() {
	// build the buttonPanelsMap
	var buttonid string
	var panel js.Value
	// CopyButton CopyButton button.
	buttonid = "mainMasterView-home-pad-CopyButton"
	buttonPanelsMap[buttonid] = make([]js.Value, 0, 5)
	panel = getElementByID(document, "mainMasterView-home-pad-CopyButton-CopyNotReadyPanel")
	if panel == js.Undefined() {
		message := "viewtools.initializeGroups: Cant find #mainMasterView-home-pad-CopyButton-CopyNotReadyPanel"
		alert.Invoke(message)
		return
	}
	buttonPanelsMap[buttonid] = append(buttonPanelsMap[buttonid], panel)
	panel = getElementByID(document, "mainMasterView-home-pad-CopyButton-CopyReadyPanel")
	if panel == js.Undefined() {
		message := "viewtools.initializeGroups: Cant find #mainMasterView-home-pad-CopyButton-CopyReadyPanel"
		alert.Invoke(message)
		return
	}
	buttonPanelsMap[buttonid] = append(buttonPanelsMap[buttonid], panel)
	// CopyButton CopyWPMTab button.
	buttonid = "mainMasterView_home_pad_CopyButton_CopyReadyPanel_tab_bar-CopyWPMTab"
	buttonPanelsMap[buttonid] = make([]js.Value, 0, 5)
	panel = getElementByID(document, "mainMasterView_home_pad_CopyButton_CopyReadyPanel_tab_bar-CopyWPMTabPanel-inner-CopyWPMPanel")
	if panel == js.Undefined() {
		message := "viewtools.initializeGroups: Cant find #mainMasterView_home_pad_CopyButton_CopyReadyPanel_tab_bar-CopyWPMTabPanel-inner-CopyWPMPanel"
		alert.Invoke(message)
		return
	}
	buttonPanelsMap[buttonid] = append(buttonPanelsMap[buttonid], panel)
	// CopyButton CopyPracticeTab button.
	buttonid = "mainMasterView_home_pad_CopyButton_CopyReadyPanel_tab_bar-CopyPracticeTab"
	buttonPanelsMap[buttonid] = make([]js.Value, 0, 5)
	panel = getElementByID(document, "mainMasterView_home_pad_CopyButton_CopyReadyPanel_tab_bar-CopyPracticeTabPanel-inner-CopyPracticePanel")
	if panel == js.Undefined() {
		message := "viewtools.initializeGroups: Cant find #mainMasterView_home_pad_CopyButton_CopyReadyPanel_tab_bar-CopyPracticeTabPanel-inner-CopyPracticePanel"
		alert.Invoke(message)
		return
	}
	buttonPanelsMap[buttonid] = append(buttonPanelsMap[buttonid], panel)
	// CopyButton CopyTestTab button.
	buttonid = "mainMasterView_home_pad_CopyButton_CopyReadyPanel_tab_bar-CopyTestTab"
	buttonPanelsMap[buttonid] = make([]js.Value, 0, 5)
	panel = getElementByID(document, "mainMasterView_home_pad_CopyButton_CopyReadyPanel_tab_bar-CopyTestTabPanel-inner-CopyTestPanel")
	if panel == js.Undefined() {
		message := "viewtools.initializeGroups: Cant find #mainMasterView_home_pad_CopyButton_CopyReadyPanel_tab_bar-CopyTestTabPanel-inner-CopyTestPanel"
		alert.Invoke(message)
		return
	}
	buttonPanelsMap[buttonid] = append(buttonPanelsMap[buttonid], panel)
	// KeyButton KeyButton button.
	buttonid = "mainMasterView-home-pad-KeyButton"
	buttonPanelsMap[buttonid] = make([]js.Value, 0, 5)
	panel = getElementByID(document, "mainMasterView-home-pad-KeyButton-KeyNotReadyPanel")
	if panel == js.Undefined() {
		message := "viewtools.initializeGroups: Cant find #mainMasterView-home-pad-KeyButton-KeyNotReadyPanel"
		alert.Invoke(message)
		return
	}
	buttonPanelsMap[buttonid] = append(buttonPanelsMap[buttonid], panel)
	panel = getElementByID(document, "mainMasterView-home-pad-KeyButton-KeyReadyPanel")
	if panel == js.Undefined() {
		message := "viewtools.initializeGroups: Cant find #mainMasterView-home-pad-KeyButton-KeyReadyPanel"
		alert.Invoke(message)
		return
	}
	buttonPanelsMap[buttonid] = append(buttonPanelsMap[buttonid], panel)
	// KeyButton KeyWPMTab button.
	buttonid = "mainMasterView_home_pad_KeyButton_KeyReadyPanel_tab_bar-KeyWPMTab"
	buttonPanelsMap[buttonid] = make([]js.Value, 0, 5)
	panel = getElementByID(document, "mainMasterView_home_pad_KeyButton_KeyReadyPanel_tab_bar-KeyWPMTabPanel-inner-KeyWPMPanel")
	if panel == js.Undefined() {
		message := "viewtools.initializeGroups: Cant find #mainMasterView_home_pad_KeyButton_KeyReadyPanel_tab_bar-KeyWPMTabPanel-inner-KeyWPMPanel"
		alert.Invoke(message)
		return
	}
	buttonPanelsMap[buttonid] = append(buttonPanelsMap[buttonid], panel)
	// KeyButton KeyPracticeTab button.
	buttonid = "mainMasterView_home_pad_KeyButton_KeyReadyPanel_tab_bar-KeyPracticeTab"
	buttonPanelsMap[buttonid] = make([]js.Value, 0, 5)
	panel = getElementByID(document, "mainMasterView_home_pad_KeyButton_KeyReadyPanel_tab_bar-KeyPracticeTabPanel-inner-KeyPracticePanel")
	if panel == js.Undefined() {
		message := "viewtools.initializeGroups: Cant find #mainMasterView_home_pad_KeyButton_KeyReadyPanel_tab_bar-KeyPracticeTabPanel-inner-KeyPracticePanel"
		alert.Invoke(message)
		return
	}
	buttonPanelsMap[buttonid] = append(buttonPanelsMap[buttonid], panel)
	// KeyButton KeyTestTab button.
	buttonid = "mainMasterView_home_pad_KeyButton_KeyReadyPanel_tab_bar-KeyTestTab"
	buttonPanelsMap[buttonid] = make([]js.Value, 0, 5)
	panel = getElementByID(document, "mainMasterView_home_pad_KeyButton_KeyReadyPanel_tab_bar-KeyTestTabPanel-inner-KeyTestPanel")
	if panel == js.Undefined() {
		message := "viewtools.initializeGroups: Cant find #mainMasterView_home_pad_KeyButton_KeyReadyPanel_tab_bar-KeyTestTabPanel-inner-KeyTestPanel"
		alert.Invoke(message)
		return
	}
	buttonPanelsMap[buttonid] = append(buttonPanelsMap[buttonid], panel)
	// ReferenceButton ReferenceButton button.
	buttonid = "mainMasterView-home-pad-ReferenceButton"
	buttonPanelsMap[buttonid] = make([]js.Value, 0, 5)
	panel = getElementByID(document, "mainMasterView-home-pad-ReferenceButton-SelectCodesPanel")
	if panel == js.Undefined() {
		message := "viewtools.initializeGroups: Cant find #mainMasterView-home-pad-ReferenceButton-SelectCodesPanel"
		alert.Invoke(message)
		return
	}
	buttonPanelsMap[buttonid] = append(buttonPanelsMap[buttonid], panel)
	// ReferenceButton LettersTab button.
	buttonid = "mainMasterView_home_pad_ReferenceButton_SelectCodesPanel_tab_bar-LettersTab"
	buttonPanelsMap[buttonid] = make([]js.Value, 0, 5)
	panel = getElementByID(document, "mainMasterView_home_pad_ReferenceButton_SelectCodesPanel_tab_bar-LettersTabPanel-inner-LettersPanel")
	if panel == js.Undefined() {
		message := "viewtools.initializeGroups: Cant find #mainMasterView_home_pad_ReferenceButton_SelectCodesPanel_tab_bar-LettersTabPanel-inner-LettersPanel"
		alert.Invoke(message)
		return
	}
	buttonPanelsMap[buttonid] = append(buttonPanelsMap[buttonid], panel)
	// ReferenceButton NumbersTab button.
	buttonid = "mainMasterView_home_pad_ReferenceButton_SelectCodesPanel_tab_bar-NumbersTab"
	buttonPanelsMap[buttonid] = make([]js.Value, 0, 5)
	panel = getElementByID(document, "mainMasterView_home_pad_ReferenceButton_SelectCodesPanel_tab_bar-NumbersTabPanel-inner-NumbersPanel")
	if panel == js.Undefined() {
		message := "viewtools.initializeGroups: Cant find #mainMasterView_home_pad_ReferenceButton_SelectCodesPanel_tab_bar-NumbersTabPanel-inner-NumbersPanel"
		alert.Invoke(message)
		return
	}
	buttonPanelsMap[buttonid] = append(buttonPanelsMap[buttonid], panel)
	// ReferenceButton PunctuationTab button.
	buttonid = "mainMasterView_home_pad_ReferenceButton_SelectCodesPanel_tab_bar-PunctuationTab"
	buttonPanelsMap[buttonid] = make([]js.Value, 0, 5)
	panel = getElementByID(document, "mainMasterView_home_pad_ReferenceButton_SelectCodesPanel_tab_bar-PunctuationTabPanel-inner-PunctuationPanel")
	if panel == js.Undefined() {
		message := "viewtools.initializeGroups: Cant find #mainMasterView_home_pad_ReferenceButton_SelectCodesPanel_tab_bar-PunctuationTabPanel-inner-PunctuationPanel"
		alert.Invoke(message)
		return
	}
	buttonPanelsMap[buttonid] = append(buttonPanelsMap[buttonid], panel)
	// ReferenceButton SpecialTab button.
	buttonid = "mainMasterView_home_pad_ReferenceButton_SelectCodesPanel_tab_bar-SpecialTab"
	buttonPanelsMap[buttonid] = make([]js.Value, 0, 5)
	panel = getElementByID(document, "mainMasterView_home_pad_ReferenceButton_SelectCodesPanel_tab_bar-SpecialTabPanel-inner-SpecialPanel")
	if panel == js.Undefined() {
		message := "viewtools.initializeGroups: Cant find #mainMasterView_home_pad_ReferenceButton_SelectCodesPanel_tab_bar-SpecialTabPanel-inner-SpecialPanel"
		alert.Invoke(message)
		return
	}
	buttonPanelsMap[buttonid] = append(buttonPanelsMap[buttonid], panel)
}
