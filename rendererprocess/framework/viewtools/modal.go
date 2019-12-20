// +build js, wasm

package viewtools

import (
	"fmt"

	"github.com/josephbudd/cwt/rendererprocess/api/event"

)

/*
	WARNING:

	DO NOT EDIT THIS FILE.

*/

type modalViewData struct {
	title   string
	message string
	cb      func()
}

// GoModal adds a title and message and call back to the modalQueue.
func GoModal(message, title string, callback func()) {
	queueModal(
		&modalViewData{
			title:   title,
			message: fmt.Sprintf("<p>%s</p>", message),
			cb:      callback,
		},
	)
}

// GoModalHTML adds a title and html message and call back to the modalQueue.
// Param message is html.
// Param title is plain text.
func GoModalHTML(htmlMessage, title string, callback func()) {
	queueModal(
		&modalViewData{
			title:   title,
			message: htmlMessage,
			cb:      callback,
		},
	)
}

func beModal() {
	wasModal := beingModal
	m := unQueueModal()
	if beingModal = m != nil; !beingModal {
		return
	}
	modalMasterViewH1.Set("innerText", m.title)
	modalMasterViewMessage.Set("innerHTML", m.message)
	modalCallBack = m.cb
	ElementShow(modalMasterView)
	if !wasModal {
		SizeApp()
	}
}

func beNotModal() {
	if modalQueueLastIndex >= 0 {
		beModal()
		return
	}
	ElementHide(modalMasterView)
	SizeApp()
	beingModal = false
}

func queueModal(m *modalViewData) {
	if modalQueueLastIndex < 4 {
		modalQueueLastIndex++
		modalQueue[modalQueueLastIndex] = m
	}
	if !beingModal {
		beModal()
	}
}

func unQueueModal() *modalViewData {
	if modalQueueLastIndex < 0 {
		return nil
	}
	m := modalQueue[0]
	for i := 0; i < modalQueueLastIndex; i++ {
		modalQueue[i] = modalQueue[i+1]
	}
	modalQueue[modalQueueLastIndex] = nil
	modalQueueLastIndex--
	return m
}

func handleModalMasterViewClose(e event.Event) (nilReturn interface{}) {
	beNotModal()
	if modalCallBack != nil {
		cb := modalCallBack
		modalCallBack = nil
		cb()
	}
	return
}
