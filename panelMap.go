package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/josephbudd/cwt/domain/data/filepaths"
	"github.com/josephbudd/cwtsitepack"
)

/*

	DO NOT EDIT THIS FILE.

	This file is generated by kickasm and regenerated by rekickasm.

*/

const (
	mainTemplate = "main.tmpl"
	headTemplate = "head.tmpl"
)

// serviceEmptyInsidePanelNamePathMap maps each markup panel template name to it's file path.
var serviceEmptyInsidePanelNamePathMap = map[string]map[string][]string{"Copy": map[string][]string{"CopyNotReadyPanel": []string{"CopyButton"}, "CopyPracticePanel": []string{"CopyButton", "CopyReadyPanel", "PracticeTab"}, "CopyTestPanel": []string{"CopyButton", "CopyReadyPanel", "TestTab"}, "CopyWPMPanel": []string{"CopyButton", "CopyReadyPanel", "WPMTab"}}, "Key": map[string][]string{"KeyNotReadyPanel": []string{"KeyButton"}, "KeyPracticePanel": []string{"KeyButton", "KeyReadyPanel", "PracticeTab"}, "KeyTestPanel": []string{"KeyButton", "KeyReadyPanel", "TestTab"}, "KeyWPMPanel": []string{"KeyButton", "KeyReadyPanel", "WPMTab"}}, "Reference": map[string][]string{"LettersPanel": []string{"ReferenceButton", "SelectCodesPanel", "LettersTab"}, "NumbersPanel": []string{"ReferenceButton", "SelectCodesPanel", "NumbersTab"}, "PunctuationPanel": []string{"ReferenceButton", "SelectCodesPanel", "PunctuationTab"}, "SpecialPanel": []string{"ReferenceButton", "SelectCodesPanel", "SpecialTab"}}}

// serveMainHTML only serves up main.tmpl with all of the templates for your markup panels.
func serveMainHTML(w http.ResponseWriter) {
	fmt.Println("serving serveMainHTML")
	var err error
	var masterT, tmpl *template.Template
	var tpath, s string
	var bb []byte
	var found bool
	var fname string
	var l int

	templateFolderPath := filepaths.GetShortTemplatePath()
	// main.tmpl
	tpath = filepath.Join(templateFolderPath, mainTemplate)
	if bb, found = cwtsitepack.Contents(tpath); !found {
		http.Error(w, fmt.Sprintf("Not found. (%s)", mainTemplate), 404)
		return
	}
	l += len(bb)
	masterT = template.New(mainTemplate)
	s = string(bb)
	if _, err = masterT.Parse(s); err != nil {
		http.Error(w, err.Error(), 300)
		return
	}
	// head.tmpl
	// the head template which contains
	//  * any css imports
	//  * any javascript imports
	tpath = filepath.Join(templateFolderPath, headTemplate)
	if bb, found = cwtsitepack.Contents(tpath); !found {
		// add a head.tmpl template
		// it's ok if the template is not there
		// but if it's there use it.
		bb = []byte(fmt.Sprintf("%[1]s%[1]s define %[3]q %[2]s%[2]s<!-- You do not have a %[3]s file to import your css files. Feel free to add one in the render/template folder. -->%[1]s%[1]s end %[2]s%[2]s", "{", "}", headTemplate))
	}
	tmpl = masterT.New(headTemplate)
	l += len(bb)
	s = string(bb)
	if _, err = tmpl.Parse(s); err != nil {
		http.Error(w, err.Error(), 300)
	}
	// panel template files
	for _, namePathMap := range serviceEmptyInsidePanelNamePathMap {
		for name, folders := range namePathMap {
			fname = name + ".tmpl"
			folderPath := strings.Join(folders, string(os.PathSeparator))
			tpath := filepath.Join(templateFolderPath, folderPath, fname)
			if bb, found = cwtsitepack.Contents(tpath); !found {
				http.Error(w, fmt.Sprintf("Not found. (%s)", fname), 404)
				return
			}
			l += len(bb)
			tmpl = masterT.New(fname)
			s = string(bb)
			if _, err = tmpl.Parse(s); err != nil {
				http.Error(w, err.Error(), 300)
			}
		}
	}
	// send the html
	if err = masterT.ExecuteTemplate(w, mainTemplate, nil); err != nil {
		if !strings.Contains(err.Error(), "reset") {
			http.Error(w, err.Error(), 300)
		}
	}
}
