// +build js, wasm

package viewtools

import (
	"syscall/js"

	"github.com/josephbudd/cwt/rendererprocess/framework/callback"
)

/*

WARNING:

DO NOT EDIT THIS FILE.

*/

// Visibility class names.
const (
	spawnIDReplacePattern  = "{{.SpawnID}}"
	TabClassName           = "tab"
	SelectedTabClassName   = "selected-tab"
	UnSelectedTabClassName = "unselected-tab"
	TabPanelClassName      = "panel-bound-to-tab"

	TabBarClassName      = "tab-bar"
	UnderTabBarClassName = "under-tab-bar"

	PanelClassName                   = "panel"
	PanelWithHeadingClassName        = "panel-with-heading"
	PanelWithTabBarClassName         = "panel-with-tab-bar"
	PanelHeadingClassName            = "heading-of-panel"
	PanelHeadingLevelPrefixClassName = "heading-of-panel-level-"
	TabPanelGroupClassName           = "inner-panel"
	UserContentClassName             = "user-content"
	ResizeMeWidthClassName           = "resize-me-width"

	SliderClassName                  = "slider"
	SliderPanelClassName             = "slider-panel"
	SliderPanelInnerClassName        = "slider-panel-pad"
	SliderPanelInnerSiblingClassName = "slider-panel-inner-sibling"
	SliderButtonPadClassName         = "slider-button-pad"

	SeenClassName       = "seen"
	UnSeenClassName     = "unseen"
	ToBeSeenClassName   = "tobe-seen"
	ToBeUnSeenClassName = "tobe-unseen"

	CookieCrumbClassName            = "cookie-crumb"
	CookieCrumbLevelPrefixClassName = "cookie-crumb-level-"

	VScrollClassName  = "vscroll"
	HVScrollClassName = "hvscroll"

	MasterID           = "mainMasterView"
	HomeID             = "mainMasterView-home"
	HomePadID          = "mainMasterView-home-pad"
	SliderID           = "mainMasterView-home-slider"
	SliderBackID       = "mainMasterView-home-slider-back"
	SliderCollectionID = "mainMasterView-home-slider-collection"

	BackIDAttribute         = "backid"
	BackColorLevelAttribute = "backColorLevel"
)

var (
	document  js.Value
	global    js.Value
	alert     js.Value
	undefined js.Value
	null      js.Value

	body                               js.Value
	mainMasterview                     js.Value
	mainMasterviewHome                 js.Value
	mainMasterviewHomeButtonPad        js.Value
	mainMasterviewHomeSlider           js.Value
	mainMasterviewHomeSliderBack       js.Value
	mainMasterviewHomeSliderCollection js.Value

	// modal
	modalMasterView        js.Value
	modalMasterViewCenter  js.Value
	modalMasterViewH1      js.Value
	modalMasterViewMessage js.Value
	modalMasterViewClose   js.Value
	modalQueue = make([]*modalViewData, 5, 5)
	modalQueueLastIndex    = -1
	beingModal             bool
	modalCallBack          func()
	// black
	blackMasterView js.Value
	// groups
	buttonPanelsMap = make(map[string][]js.Value, 100)
	// slider
	here =  js.Undefined()
	backStack []js.Value
	// tabber
	tabberLastPanelID     string
	tabberTabBarLastPanel map[string]string
	// button locking
	buttonsLocked             bool
	buttonsLockedMessageTitle string
	buttonsLockedMessageText  string

	// spawns

	spawnID               uint64
	SpawnIDReplacePattern string

	// user content

	panelNameHVScroll = map[string]bool{"CopyNotReadyPanel":false, "CopyPracticePanel":false, "CopyReadyPanel":false, "CopyTestPanel":false, "CopyWPMPanel":false, "KeyNotReadyPanel":false, "KeyPracticePanel":false, "KeyReadyPanel":false, "KeyTestPanel":false, "KeyWPMPanel":false, "LettersPanel":false, "NumbersPanel":false, "PunctuationPanel":false, "SelectCodesPanel":false, "SpecialPanel":false}

	// markup panels

	countMarkupPanels = 12
	countSpawnedMarkupPanels int
	countWidgetsWaiting      int

	// spawned widgets

	spawnedWidgets = make(map[uint64]spawnedWidgetInfo, 100)

	// printing

	extraHeight   float64
	documentTitle string
	printTitle    string
)

func getElementByID(document js.Value, id string) (e js.Value) {
	e = document.Call("getElementById", id)
	return
}

func init() {
	global = js.Global()
	document = global.Get("document")
	alert = global.Get("alert")
	undefined = js.Undefined()
	null = js.Null()

	documentTitle = document.Get("title").String()
	printTitle = documentTitle

	SpawnIDReplacePattern = spawnIDReplacePattern
	bodies := document.Call("getElementsByTagName", "BODY")
	body = bodies.Index(0)
	mainMasterview = getElementByID(document, MasterID)
	mainMasterviewHome = getElementByID(document, HomeID)
	mainMasterviewHomeButtonPad = getElementByID(document, HomePadID)
	mainMasterviewHomeSlider = getElementByID(document, SliderID)
	mainMasterviewHomeSliderBack = getElementByID(document, SliderBackID)
	mainMasterviewHomeSliderCollection = getElementByID(document, SliderCollectionID)
	// modal
	modalMasterView = getElementByID(document, "modalInformationMasterView")
	modalMasterViewCenter = getElementByID(document, "modalInformationMasterView-center")
	modalMasterViewH1 = getElementByID(document, "modalInformationMasterView-h1")
	modalMasterViewMessage = getElementByID(document, "modalInformationMasterView-message")
	modalMasterViewClose = getElementByID(document, "modalInformationMasterView-close")
	modalQueue = make([]*modalViewData, 5, 5)
	modalQueueLastIndex = -1
	callback.AddEventHandler(handleModalMasterViewClose, modalMasterViewClose, "click", false, 0)
	// black
	blackMasterView = getElementByID(document, "blackMasterView")
	// Setup printing markup panels.
	callback.AddEventHandler(beforePrint, global, "beforeprint", false, 0)
	callback.AddEventHandler(afterPrint, global, "afterprint", false, 0)

	initializeGroups()
	initializeSlider()
	initializeResize()
	initializeTabBar()
}