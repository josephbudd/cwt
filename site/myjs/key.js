

// You need to edit this file.
// This file compliments the javascript file at /renderer/js/key.js
// Do not edit the file at /renderer/js/key.js

/****************************************

	KeyControler

****************************************/

const KeyControler = {

	stateTest:0,
	statePractice:1,

    keycodeRecords:{}, // map[id]record
    recordCount:0,
    copy:undefined,

    wpmRecord:{}, // wpm record
    stack:[],
    downAt:0,
    upAt:0,

	// KeyControler.Initialize is called on startup.
	Initialize: function() {
		KeyControler._Initialize()

        document.getElementById("key-wpm").onchange = function(e){KeyControler.handleWPMChange(e)}

        KeyControler.testCopy = document.getElementById("key-test-copy")
        document.getElementById("key-test-start").onclick = KeyControler.handleTestStart
        document.getElementById("key-test-stop").onclick = KeyControler.handleTestStop

		KeyControler.practiceCopy = document.getElementById("key-practice-copy")
        document.getElementById("key-practice-start").onclick = KeyControler.handlePracticeStart
        document.getElementById("key-practice-stop").onclick = KeyControler.handlePracticeStop

		KeyPresenter.Initialize()
	},

	// wpm

    handleWPMChange: function(e) {
        KeyControler.wpmRecord.wpm = window.parseInt(e.target.value)
        KeyLPC.upDateKeyWPM(KeyControler.wpmRecord)
    },

    updateWPM: function(record) {
        KeyControler.wpmRecord = record
        KeyPresenter.setWPM(record.wpm)
	},
	
	// keycode records

    updateKeycodeRecord: function(record) {
        if (record.id in KeyControler.keycodeRecords) {
            if (record.selected) {
                KeyControler.keycodeRecords[record.id] = record
            } else {
                delete(KeyControler.keycodeRecords[record.id])
                KeyControler.recordCount--
            }
        } else {
            if (record.selected) {
                KeyControler.keycodeRecords[record.id] = record
                KeyControler.recordCount++
            }
        }
        if (KeyControler.recordCount > 0) {
            KeyPresenter.showKeyTabReadyPanel()
            KeyPresenter.clearTestCopy()
        } else {
            KeyPresenter.showKeyTabNotReadyPanel()
        }
    },

    updateKeycodeRecords: function(records) {
        for (var i in records) {
            var record = records[i]
            if (record.id in KeyControler.keycodeRecords) {
                if (record.selected) {
                    KeyControler.keycodeRecords[record.id] = record
                } else {
                    KeyControler.keycodeRecords.remove(record.id)
                    KeyControler.recordCount--
                }
            } else {
                if (record.selected) {
                    KeyControler.keycodeRecords[record.id] = record
                    KeyControler.recordCount++
                }
            }
        }
        if (KeyControler.recordCount > 0) {
            KeyPresenter.showKeyTabReadyPanel()
        } else {
            KeyPresenter.showKeyTabNotReadyPanel()
        }
    },

	// test

    handleTestStart: function(e) {
        KeyPresenter.showStopButton()
        KeyPresenter.clearTestCopy()
        KeyPresenter.clearTestText()
        KeyPresenter.showTestInstructions()
        // set mouse up and down handlers
        KeyControler.testCopy.onmousedown = KeyControler.handleOnMouseDown
        KeyControler.testCopy.onmouseup = KeyControler.handleOnMouseUp
        // clear the stack
        KeyControler.downAt = 0
        KeyControler.upAt = 0
        KeyControler.stack.length = 0
        KeyLPC.getTextToKey(KeyControler.stateTest)
    },

    handleTestStop: function(e) {
        KeyPresenter.showStartButton()
        // remove mouse up and down handlers
        KeyControler.testCopy.onmousedown = undefined
        KeyControler.testCopy.onmouseup = undefined
        KeyLPC.processKeyedMilliseconds(KeyControler.stack, KeyControler.stateTest)
    },

	// practice

    handlePracticeStart: function(e) {
        KeyPresenter.showStopButton()
        KeyPresenter.clearPracticeCopy()
        KeyPresenter.clearPracticeText()
        KeyPresenter.showPracticeInstructions()
        // set mouse up and down handlers
        KeyControler.practiceCopy.onmousedown = KeyControler.handleOnMouseDown
        KeyControler.practiceCopy.onmouseup = KeyControler.handleOnMouseUp
        // clear the stack
        KeyControler.downAt = 0
        KeyControler.upAt = 0
        KeyControler.stack.length = 0
        KeyLPC.getTextToKey(KeyControler.statePractice)
    },

    handlePracticeStop: function(e) {
        KeyPresenter.showStartButton()
        // remove mouse up and down handlers
        KeyControler.practiceCopy.onmousedown = undefined
        KeyControler.practiceCopy.onmouseup = undefined
        KeyLPC.processKeyedMilliseconds(KeyControler.stack, KeyControler.statePractice)
    },

	// mouse

    handleOnMouseDown: function(e) {
        e.stopPropagation()
        KeyControler.downAt = new Date().getTime()
         // push the pause
        KeyControler.stack.push (KeyControler.downAt - KeyControler.upAt)
    },

    handleOnMouseUp: function(e) {
        e.stopPropagation()
        KeyControler.upAt = new Date().getTime()
        // push key down time
        KeyControler.stack.push(KeyControler.upAt - KeyControler.downAt)
    }

}

/****************************************

	KeyPresenter

****************************************/

const KeyPresenter = {

    instructions: "Place the mouse over the red area and begin keying.",
	
	// KeyPresenter.Initialize is called on startup.
	Initialize: function() {
		KeyPresenter._Initialize()
        KeyPresenter.wpm = document.getElementById("key-wpm")

        KeyPresenter.testTextDisplay = document.getElementById("key-test-text")
        KeyPresenter.testCopyDisplay = document.getElementById("key-test-copy")
        KeyPresenter.testStartButton = document.getElementById("key-test-start")
        KeyPresenter.testStopButton = document.getElementById("key-test-stop")

        KeyPresenter.practiceTextDisplay = document.getElementById("key-practice-text")
        KeyPresenter.practiceCopyDisplay = document.getElementById("key-practice-copy")
        KeyPresenter.practiceStartButton = document.getElementById("key-practice-start")
        KeyPresenter.practiceStopButton = document.getElementById("key-practice-stop")

	},

	// wpm

    setWPM: function(wpm) {
        KeyPresenter.wpm.value = wpm.toString()
	},
	
	// test

    clearTestText: function() {
        while (KeyPresenter.testTextDisplay.firstChild != undefined) {
            KeyPresenter.testTextDisplay.removeChild(KeyPresenter.testTextDisplay.firstChild)
        }
    },

    displayTestText: function(text) {
        var p = document.createElement("p")
        for (var i in text) {
            if (i > 0) {
                p.appendChild(document.createElement("br"))
            }
            p.appendChild(document.createTextNode(text[i]))
        }
        KeyPresenter.testTextDisplay.appendChild(p)
    },

    displayTestCopy: function(ditdahs, texts) {
        var p = document.createElement("p")
        for (var i in ditdahs) {
            if (i > 0) {
                p.appendChild(document.createElement("br"))
            }
            var span = document.createElement("span")
            span.classList.add("ditdah")
            span.appendChild(document.createTextNode(ditdahs[i]))
            p.appendChild(span)
            p.appendChild(document.createTextNode(" : ".concat(texts[i])))
        }
        KeyPresenter.testCopyDisplay.appendChild(p)
    },

    showTestInstructions: function() {
        goModal(KeyPresenter.instructions, "Instructions")
    },

    clearTestCopy: function(html) {
        while (KeyPresenter.testCopyDisplay.firstChild != undefined) {
            KeyPresenter.testCopyDisplay.removeChild(KeyPresenter.testCopyDisplay.firstChild)
        }
    },

    showTestStartButton: function() {
        ElementHide(KeyPresenter.testStopButton)
        ElementShow(KeyPresenter.testStartButton)
    },

    showTestStopButton: function() {
        ElementHide(KeyPresenter.testStartButton)
        ElementShow(KeyPresenter.testStopButton)
    },
	
	// practice

    clearPracticeText: function() {
        while (KeyPresenter.practiceTextDisplay.firstChild != undefined) {
            KeyPresenter.practiceTextDisplay.removeChild(KeyPresenter.practiceTextDisplay.firstChild)
        }
    },

    displayPracticeText: function(text) {
        var p = document.createElement("p")
        for (var i in text) {
            if (i > 0) {
                p.appendChild(document.createElement("br"))
            }
            p.appendChild(document.createTextNode(text[i]))
        }
        KeyPresenter.practiceTextDisplay.appendChild(p)
    },

    displayPracticeCopy: function(ditdahs, texts) {
        var p = document.createElement("p")
        for (var i in ditdahs) {
            if (i > 0) {
                p.appendChild(document.createElement("br"))
            }
            var span = document.createElement("span")
            span.classList.add("ditdah")
            span.appendChild(document.createTextNode(ditdahs[i]))
            p.appendChild(span)
            p.appendChild(document.createTextNode(" : ".concat(texts[i])))
        }
        KeyPresenter.practiceCopyDisplay.appendChild(p)
    },

    showPracticeInstructions: function() {
        goModal(KeyPresenter.instructions, "Instructions")
    },

    clearPracticeCopy: function(html) {
        while (KeyPresenter.practiceCopyDisplay.firstChild != undefined) {
            KeyPresenter.practiceCopyDisplay.removeChild(KeyPresenter.practiceCopyDisplay.firstChild)
        }
    },

    showPracticeStartButton: function() {
        ElementHide(KeyPresenter.practiceStopButton)
        ElementShow(KeyPresenter.practiceStartButton)
    },

    showPracticeStopButton: function() {
        ElementHide(KeyPresenter.practiceStartButton)
        ElementShow(KeyPresenter.practiceStopButton)
    },

	// state

	displayText: function(text, state) {
		if (state == KeyControler.stateTest) {
			displayTestText(text)
		} else {
			displayPracticeText(text)
		}
	},

	displayCopy: function(text, state) {
		if (state == KeyControler.stateTest) {
			displayTestCopy(text)
		} else {
			displayPracticeCopy(text)
		}
	}
}

/****************************************

	KeyLPC

****************************************/

const KeyLPC = {

	// KeyLPC.Initialize is called on startup.
	Initialize: function() {
		lpc.callInGetKeyCodesCB(KeyLPC.getKeyCodesCB)
        lpc.callInUpdateKeyCodeCB(KeyLPC.updateKeyCodeCB)
		lpc.callInUpdateKeyCodesCB(KeyLPC.updateKeyCodesCB)
		lpc.UpdateKeyWPMCB = KeyLPC.updateKeyWPMCB
		lpc.GetKeyWPMCB = KeyLPC.getKeyWPMCB

        lpc.ProcessKeyedMillisecondsCB = KeyLPC.processKeyedMillisecondsCB
        lpc.GetTextWPMToKeyCB = KeyLPC.getTextToKeyCB
	},

	// KeyLPC.InitialCall is called on startup.
	InitialCall: function() {
		lpc.GetKeyWPM()
	},

    processKeyedMilliseconds: function(stack, state) {
        lpc.ProcessKeyedMilliseconds({stack:stack, state:state})
    },

    getTextToKey: function(state) {
        lpc.GetTextWPMToKey(state)
    },

    getKeyCodesCB: function(params) {
        if (params.error == true) {
            KeyPresenter.error("GetKeyCodesCB: ".concat(params.message))
            return undefined
        }
        KeyControler.updateKeycodeRecords(params.records)
    },

    updateKeyWPMCB: function(params) {
        if (params.error == true) {
            KeyPresenter.error("UpdateKeyWPMCB: ".concat(params.message))
            return
        }
    },

    getKeyWPMCB: function(params) {
        if (params.error == true) {
            KeyPresenter.error("GetKeyWPMCB: ".concat(params.message))
            return
        }
        KeyControler.updateWPM(params.record)
    },

    updateKeyCodeCB: function(params) {
        if (params.error == true) {
            KeyPresenter.error("UpdateKeyCodeCB: ".concat(params.message))
            return
        }
        KeyControler.updateKeycodeRecord(params.record)
    },

    updateKeyCodesCB: function(params) {
        if (params.error == true) {
            KeyPresenter.error("UpdateKeyCodesCB: ".concat(params.message))
            return
        }
        KeyControler.updateKeycodeRecords(params.records)
    },

    processKeyedMillisecondsCB: function(params) {
        if (params.error == true) {
            KeyPresenter.Error("ProcessKeyedMillisecondsCB: ".concat(params.message))
            return
        }
        KeyPresenter.displayCopy(params.ditdahs, params.texts, params.state)
    },

    getTextToKeyCB: function(params) {
        if (params.error == true) {
            KeyPresenter.Error("GetTextWPMToKeyCB: ".concat(params.message))
            return
        }
        KeyPresenter.displayText(params.lines, params.state)
    },

}
