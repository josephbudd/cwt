

// You need to edit this file.
// This file compliments the javascript file at /renderer/js/copy.js
// Do not edit the file at /renderer/js/copy.js

/****************************************

	CopyControler

****************************************/

const CopyControler = {

	stateTest:0,
	statePractice:1,

    userIsCopying:false,
    keycodeRecords:{}, // map[id]record
    recordCount:0,
    DelaySeconds:5,
    wpmRecord:{}, // wpm record
    currentTestDitDah:"",
    currentTestDitDahLabels:"",
    currentTestDitDahCorrectAnswer:"",

	// CopyControler.Initialize is called on startup.
	Initialize: function() {
		CopyControler._Initialize()
        
        CopyControler.testCopy = document.getElementById("copy-test-copy")
        document.getElementById("copy-test-start").onclick = CopyControler.handleTestStart
        document.getElementById("copy-test-stop").onclick = CopyControler.handleTestStop

		CopyControler.practiceCopy = document.getElementById("copy-practice-copy")
        document.getElementById("copy-practice-start").onclick = CopyControler.handlePracticeStart
        document.getElementById("copy-practice-stop").onclick = CopyControler.handlePracticeStop

        document.getElementById("copy-wpm").onchange = function(e){CopyControler.handleWPMChange(e)}
		
	},

	// wpm

    handleWPMChange: function(e) {
        CopyControler.wpmRecord.wpm = window.parseInt(e.target.value)
        CopyLPC.upDateCopyWPM(CopyControler.wpmRecord)
    },

    updateWPM: function(record) {
        CopyControler.wpmRecord = record
        CopyPresenter.setWPM(record.wpm)
    },

	// keycode records

    updateKeycodeRecord: function(record) {
        if (record.id in CopyControler.keycodeRecords) {
            if (record.selected) {
                CopyControler.keycodeRecords[record.id] = record
            } else {
                delete(CopyControler.keycodeRecords[record.id])
                CopyControler.recordCount--
            }
        } else {
            if (record.selected) {
                CopyControler.keycodeRecords[record.id] = record
                CopyControler.recordCount++
            }
        }
        if (CopyControler.recordCount > 0) {
            CopyPresenter.showCopyTabReadyPanel()
        } else {
            CopyPresenter.showCopyTabNotReadyPanel()
        }
    },

    updateKeycodeRecords: function(records) {
        for (var i in records) {
            var record = records[i]
            if (record.id in CopyControler.keycodeRecords) {
                if (record.selected) {
                    CopyControler.keycodeRecords[record.id] = record
                } else {
                    CopyControler.keycodeRecords.remove(record.id)
                    CopyControler.recordCount--
                }
            } else {
                if (record.selected) {
                    CopyControler.keycodeRecords[record.id] = record
                    CopyControler.recordCount++
                }
            }
        }
        if (CopyControler.recordCount > 0) {
            CopyPresenter.showCopyTabReadyPanel()
        } else {
            CopyPresenter.showCopyTabNotReadyPanel()
        }
    },

	// test

    handleTestStart: function(e) {
        CopyPresenter.showTestStopButton()
        CopyPresenter.clearTestCopyDisplay()
        CopyControler.userIsCopying = true
        // generate some random stuff
        // input:
        //  * "." is used as a dit.
        //  * "-" is used as a dah.
        //  * " " is a character separator.
        //  * no word separator.
        // output:
        //  * "." is used as a dit.
        //  * "-" is used as a dah.
        //  * " " is a character separator.
        //  * "\t" is used as a word separator.
        goModalCB("The CW will begin 5 seconds after you click close. Enter you copy into the red square.", "Copy", function(){CopyControler.playTestCW()})
    },

    handleTestStop: function(e) {
        CopyControler.userIsCopying = false
        CopyPresenter.presentTestDitDahLabels(CopyControler.currentTestDitDahLabels)
        CopyPresenter.presentTestDitDah(CopyControler.currentTestDitDah)
        CopyPresenter.showTestStartButton()
        // check the answer
        var copy = CopyControler.testCopy.value.trim()
        if (copy.length == 0) {
            goModal("Lets try that again.", "Copy")
            return
        } 
        var solution = CopyControler.currentTestDitDahCorrectAnswer
        CopyLPC.checkKeyCodeCopy(copy, solution, CopyControler.stateTest)
    },

    playTestCW: function() {
        CopyPresenter.clearTestDitDahDisplay()
        CopyPresenter.focusTestCopyDisplay()
        var maxlen = CopyControler.keycodeRecords.length
        var aural = []
        var visual = []
        var answer = []
        for (var wordlen=1; wordlen<=5; wordlen++) {
            var auralchars = []
            var visualchars = []
            var answerchars = []
            for (var id in CopyControler.keycodeRecords) {
                var record = CopyControler.keycodeRecords[id]
                auralchars.push(record.ditdah)
                answerchars.push(record.character)
                switch (record.character) {
                case "<":
                    visualchars.push("&lt;")
                    break
                case ">":
                    visualchars.push("&gt;")
                    break
                case "&":
                    visualchars.push("&amp;")
                    break
                default:
                    visualchars.push(record.character)
                    break
                }
                if (visualchars.length == wordlen) {
                    aural.push(auralchars.join(" ")) // space is a char separator
                    auralchars = []
                    visual.push(visualchars.join(""))
                    visualchars = []
                    answer.push(answerchars.join(""))
                    answerchars = []
                }
            }
        }
        // save
        CopyControler.currentTestDitDahLabels = visual.join("&nbsp;&nbsp;&nbsp;&nbsp;")
        CopyControler.currentTestDitDah = aural.join("&nbsp;&nbsp;")
        CopyControler.currentTestDitDah.replace(" ", "&nbsp;")
        CopyControler.currentTestDitDahCorrectAnswer = answer.join(" ")
        // output aural
        CopyLPC.playDitDah(
			aural.join("\t"),
			CopyControler.wpmRecord.wpm,
			CopyControler.DelaySeconds,
			CopyControler.stateTest)
	},

	// practice

    handlePracticeStart: function(e) {
        CopyPresenter.showPracticeStopButton()
        CopyPresenter.clearPracticeCopyDisplay()
        CopyControler.userIsCopying = true
        // generate some random stuff
        // input:
        //  * "." is used as a dit.
        //  * "-" is used as a dah.
        //  * " " is a character separator.
        //  * no word separator.
        // output:
        //  * "." is used as a dit.
        //  * "-" is used as a dah.
        //  * " " is a character separator.
        //  * "\t" is used as a word separator.
        goModalCB("The CW will begin 5 seconds after you click close. Enter you copy into the red square.", "Copy", function(){CopyControler.playPracticeCW()})
    },

    handlePracticeStop: function(e) {
        CopyControler.userIsCopying = false
        CopyPresenter.presentPracticeDitDahLabels(CopyControler.currentPracticeDitDahLabels)
        CopyPresenter.presentPracticeDitDah(CopyControler.currentPracticeDitDah)
        CopyPresenter.showPracticeStartButton()
        // check the answer
        var copy = CopyControler.practiceCopy.value.trim()
        if (copy.length == 0) {
            goModal("Lets try that again.", "Copy")
            return
        } 
        var solution = CopyControler.currentPracticeDitDahCorrectAnswer
        CopyLPC.checkKeyCodeCopy(copy, solution, CopyControler.statePractice)
    },

    playPracticeCW: function() {
        CopyPresenter.clearPracticeDitDahDisplay()
        CopyPresenter.focusPracticeCopyDisplay()
        var maxlen = CopyControler.keycodeRecords.length
        var aural = []
        var visual = []
        var answer = []
        for (var wordlen=1; wordlen<=5; wordlen++) {
            var auralchars = []
            var visualchars = []
            var answerchars = []
            for (var id in CopyControler.keycodeRecords) {
                var record = CopyControler.keycodeRecords[id]
                auralchars.push(record.ditdah)
                answerchars.push(record.character)
                switch (record.character) {
                case "<":
                    visualchars.push("&lt;")
                    break
                case ">":
                    visualchars.push("&gt;")
                    break
                case "&":
                    visualchars.push("&amp;")
                    break
                default:
                    visualchars.push(record.character)
                    break
                }
                if (visualchars.length == wordlen) {
                    aural.push(auralchars.join(" ")) // space is a char separator
                    auralchars = []
                    visual.push(visualchars.join(""))
                    visualchars = []
                    answer.push(answerchars.join(""))
                    answerchars = []
                }
            }
        }
        // save
        CopyControler.currentPracticeDitDahLabels = visual.join("&nbsp;&nbsp;&nbsp;&nbsp;")
        CopyControler.currentPracticeDitDah = aural.join("&nbsp;&nbsp;")
        CopyControler.currentPracticeDitDah.replace(" ", "&nbsp;")
        CopyControler.currentPracticeDitDahCorrectAnswer = answer.join(" ")
        // output aural
        CopyLPC.playDitDah(
			aural.join("\t"),
			CopyControler.wpmRecord.wpm,
			CopyControler.DelaySeconds,
			CopyControler.statePractice)
	},

	// with state

    showResult: function(testResults, characterCount, correctCount, correctPercent, state) {
		if (state == CopyControler.stateTest) {
			CopyPresenter.presentTestDitDahResult("You got ".concat(correctCount.toString(), " out of ", characterCount.toString(), " correct. That's ", correctPercent.toString(), "% correct."))
			if (testResults.length > 0) {
				var problems = []
				for (var i in testResults) {
					problems.push("At position ".concat(testResults[i].position.toString(), " : <b>", testResults[i].copyChar, "</b> should have been <b>", testResults[i].solutionChar, "</b>"))
				}
				CopyPresenter.presentTestDitDahProblems(problems.join("<br/>"))
			}
		} else {
			CopyPresenter.presentPracticeDitDahResult("You got ".concat(correctCount.toString(), " out of ", characterCount.toString(), " correct. That's ", correctPercent.toString(), "% correct."))
			if (testResults.length > 0) {
				var problems = []
				for (var i in testResults) {
					problems.push("At position ".concat(testResults[i].position.toString(), " : <b>", testResults[i].copyChar, "</b> should have been <b>", testResults[i].solutionChar, "</b>"))
				}
				CopyPresenter.presentPracticeDitDahProblems(problems.join("<br/>"))
			}
		}
    },

}

/****************************************

	CopyPresenter

****************************************/

const CopyPresenter = {
	
	// CopyPresenter.Initialize is called on startup.
	Initialize: function() {
		CopyPresenter._Initialize()
		CopyPresenter.wpm = document.getElementById("copy-wpm")
	
        CopyPresenter.testCopyDisplay = document.getElementById("copy-test-copy")
        CopyPresenter.testDitDahDisplay = document.getElementById("copy-test-ditdah")
        CopyPresenter.testStartButton = document.getElementById("copy-test-start")
        CopyPresenter.testStopButton = document.getElementById("copy-test-stop")
	
        CopyPresenter.practiceCopyDisplay = document.getElementById("copy-practice-copy")
        CopyPresenter.practiceDitDahDisplay = document.getElementById("copy-practice-ditdah")
        CopyPresenter.practiceStartButton = document.getElementById("copy-practice-start")
        CopyPresenter.practiceStopButton = document.getElementById("copy-practice-stop")
	},

	// wpm

    setWPM: function(wpm) {
        CopyPresenter.wpm.value = wpm.toString()
    },

	// test

    showTestStartButton: function() {
        ElementHide(CopyPresenter.testStopButton)
        ElementShow(CopyPresenter.testStartButton)        
    },

    showTestStopButton: function() {
        ElementHide(CopyPresenter.testStartButton)
        ElementShow(CopyPresenter.testStopButton)        
    },

    clearTestDitDahDisplay: function() {
        while (CopyPresenter.testDitDahDisplay.firstChild != undefined) {
            CopyPresenter.testDitDahDisplay.removeChild(CopyPresenter.testDitDahDisplay.firstChild)
        }
    },

    focusTestCopyDisplay: function() {
        CopyPresenter.testCopyDisplay.focus = true
    },

    clearTestCopyDisplay: function() {
        CopyPresenter.testCopyDisplay.value = ""
    },

    presentTestDitDahLabels: function(text) {
        var p = document.createElement("p")
        p.innerHTML = text
        CopyPresenter.testDitDahDisplay.appendChild(p)
    },

    presentTestDitDah: function(ditdah) {
        var p = document.createElement("p")
        p.classList.add("ditdah")
        p.innerHTML = ditdah
        CopyPresenter.testDitDahDisplay.appendChild(p)
    },

    presentTestDitDahResult: function(message) {
        var p = document.createElement("p")
        p.classList.add("ditdah-answer")
        p.innerHTML = message
        CopyPresenter.testDitDahDisplay.appendChild(p)
    },

    presentTestDitDahProblems: function(message) {
        var p = document.createElement("p")
        p.classList.add("ditdah-problem")
        p.innerHTML = message
        CopyPresenter.testDitDahDisplay.appendChild(p)
    },

	// practice

    showPracticeStartButton: function() {
        ElementHide(CopyPresenter.practiceStopButton)
        ElementShow(CopyPresenter.practiceStartButton)        
    },

    showPracticeStopButton: function() {
        ElementHide(CopyPresenter.practiceStartButton)
        ElementShow(CopyPresenter.practiceStopButton)        
    },

    clearPracticeDitDahDisplay: function() {
        while (CopyPresenter.practiceDitDahDisplay.firstChild != undefined) {
            CopyPresenter.practiceDitDahDisplay.removeChild(CopyPresenter.practiceDitDahDisplay.firstChild)
        }
    },

    focusPracticeCopyDisplay: function() {
        CopyPresenter.practiceCopyDisplay.focus = true
    },

    clearPracticeCopyDisplay: function() {
        CopyPresenter.practiceCopyDisplay.value = ""
    },

    presentPracticeDitDahLabels: function(text) {
        var p = document.createElement("p")
        p.innerHTML = text
        CopyPresenter.practiceDitDahDisplay.appendChild(p)
    },

    presentPracticeDitDah: function(ditdah) {
        var p = document.createElement("p")
        p.classList.add("ditdah")
        p.innerHTML = ditdah
        CopyPresenter.practiceDitDahDisplay.appendChild(p)
    },

    presentPracticeDitDahResult: function(message) {
        var p = document.createElement("p")
        p.classList.add("ditdah-answer")
        p.innerHTML = message
        CopyPresenter.practiceDitDahDisplay.appendChild(p)
    },

    presentPracticeDitDahProblems: function(message) {
        var p = document.createElement("p")
        p.classList.add("ditdah-problem")
        p.innerHTML = message
        CopyPresenter.practiceDitDahDisplay.appendChild(p)
    },

}

/****************************************

	CopyLPC

****************************************/

const CopyLPC = {

	// CopyLPC.Initialize is called on startup.
	Initialize: function() {
        lpc.callInGetKeyCodesCB(CopyLPC.getKeyCodesCB)
        lpc.callInUpdateKeyCodeCB(CopyLPC.updateKeyCodeCB)
        lpc.callInUpdateKeyCodesCB(CopyLPC.updateKeyCodesCB)
		lpc.UpdateCopyWPMCB = CopyLPC.updateCopyWPMCB
		lpc.GetCopyWPMCB = CopyLPC.getCopyWPMCB
        lpc.PlayDitDahCB = CopyLPC.playDitDahCB
        lpc.CheckKeyCodeCopyCB = CopyLPC.checkKeyCodeCopyCB
	},

	// CopyLPC.InitialCall is called on startup.
	InitialCall: function() {
		lpc.GetCopyWPM()
    },
    
    upDateCopyWPM: function(record) {
        lpc.UpdateCopyWPM(record)
    },

    checkKeyCodeCopy: function(copy, solution, state) {
        lpc.CheckKeyCodeCopy({copy:copy, solution:solution, state:state})
    },

    playDitDah: function(ditdah, wpm, delay, state) {
        lpc.PlayDitDah({ditdah:ditdah, wpm:wpm, delay:delay, state:state})
    },

    getKeyCodesCB: function(params) {
        if (params.error == true) {
            CopyPresenter.error("GetKeyCodesCB: ".concat(params.message))
            return undefined
        }
        CopyControler.updateKeycodeRecords(params.records)
    },

    updateCopyWPMCB: function(params) {
        if (params.error == true) {
            CopyPresenter.error("UpdateCopyWPMCB: ".concat(params.message))
            return
        }
    },

    getCopyWPMCB: function(params) {
        if (params.error == true) {
            CopyPresenter.error("GetCopyWPMCB: ".concat(params.message))
            return
        }
        CopyControler.updateWPM(params.record)
    },

    updateKeyCodeCB: function(params) {
        if (params.error == true) {
            CopyPresenter.error("UpdateKeyCodeCB: ".concat(params.message))
            return
        }
        CopyControler.updateKeycodeRecord(params.record)
    },

    updateKeyCodesCB: function(params) {
        if (params.error == true) {
            CopyPresenter.error("UpdateKeyCodesCB: ".concat(params.message))
            return undefined
        }
        CopyControler.updateKeycodeRecords(params.records)
    },

    playDitDahCB: function(params) {
        if (params.error) {
            CopyPresenter.error("Audio play error: ".concat(params.message))
        }
    },

    checkKeyCodeCopyCB: function(params) {
        if (params.error == true) {
            CopyPresenter.error("CheckKeyCodeCopyCB: ".concat(params.message))
            return undefined
        }
        CopyControler.showResult(params.testResults, params.characterCount, params.correctCount, params.correctPercent, params.state)
    },

}
