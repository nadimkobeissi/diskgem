// Written by Nadim Kobeissi, <nadim@symbolic.software> November 2018
// Copyright (c) 2018 Nadim Kobeissi
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package main

import (
	"regexp"

	"github.com/jroimartin/gocui"
)

func uiNewFolder(ui *gocui.Gui, v *gocui.View) error {
	maxX, maxY := ui.Size()
	dgState.newFolderWindow.view, _ = ui.SetView("newFolder",
		maxX/2-20, maxY/2-5,
		maxX/2+20, maxY/2-3,
	)
	dgState.newFolderWindow.view.Title = "New Folder"
	dgState.newFolderWindow.view.Editable = true
	ui.SetCurrentView("newFolder")
	return nil
}

func uiNewFolderToggle(ui *gocui.Gui, v *gocui.View) error {
	if dgState.newFolderWindow.state.visible {
		ui.Cursor = false
		dgState.newFolderWindow.state.visible = false
		ui.DeleteView("newFolder")
		if dgState.mainWindow.state.leftPane.focused {
			ui.SetCurrentView("leftPane")
		} else if dgState.mainWindow.state.leftPane.focused {
			ui.SetCurrentView("rightPane")
		}
	} else {
		uiMainHideAllWindowsExcept("newFolder", ui, v)
		ui.Cursor = true
		dgState.newFolderWindow.state.visible = true
		uiNewFolder(ui, v)
	}
	return nil
}

func uiNewFolderHandleEnter(ui *gocui.Gui, v *gocui.View) error {
	newFolderName, _ := dgState.newFolderWindow.view.Line(0)
	validFolderName, _ := regexp.MatchString(
		"^([a-zA-Z0-9][^*/><?\"|:]*)$",
		newFolderName,
	)
	if len(newFolderName) < 1 {
		uiMainStatusViewMessage(0, "Please enter a folder name.")
		return nil
	}
	if len(newFolderName) > 128 {
		uiMainStatusViewMessage(0, "Folder name is too long.")
		return nil
	}
	if !validFolderName {
		uiMainStatusViewMessage(0, "Invalid folder name.")
		return nil
	}
	uiNewFolderToggle(ui, v)
	go uiMainCreateFolder(ui, v, newFolderName)
	return nil
}
