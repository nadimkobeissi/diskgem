/* SPDX-License-Identifier: MIT
 * Copyright © 2018-2019 Nadim Kobeissi <nadim@nadim.computer>. All Rights Reserved.
 */

package main

import (
	"fmt"
	"path"
	"strings"

	"github.com/jroimartin/gocui"
)

func uiGoTo(ui *gocui.Gui, v *gocui.View) error {
	maxX, maxY := ui.Size()
	dgState.goToWindow.view, _ = ui.SetView("goTo",
		maxX/2-20, maxY/2-5,
		maxX/2+20, maxY/2-3,
	)
	dgState.goToWindow.view.Title = "Go To Folder..."
	dgState.goToWindow.view.Editable = true
	ui.SetCurrentView("goTo")
	return nil
}

func uiGoToToggle(ui *gocui.Gui, v *gocui.View) error {
	if dgState.goToWindow.state.visible {
		ui.Cursor = false
		dgState.goToWindow.state.visible = false
		ui.DeleteView("goTo")
		if dgState.mainWindow.state.leftPane.focused {
			ui.SetCurrentView("leftPane")
		} else if dgState.mainWindow.state.leftPane.focused {
			ui.SetCurrentView("rightPane")
		}
	} else {
		uiMainHideAllWindowsExcept("goTo", ui, v)
		ui.Cursor = true
		dgState.goToWindow.state.visible = true
		dgState.goToWindow.state.lastPath = ""
		dgState.goToWindow.state.lastInitial = ""
		dgState.goToWindow.state.index = 0
		uiGoTo(ui, v)
	}
	return nil
}

func uiGoToHandleEnter(ui *gocui.Gui, v *gocui.View) error {
	goToFolderPath, _ := dgState.goToWindow.view.Line(0)
	if len(goToFolderPath) < 1 {
		uiMainStatusViewMessage(ui, 0, "Please enter a folder name.")
		return nil
	}
	if len(goToFolderPath) > 256 {
		uiMainStatusViewMessage(ui, 0, "Folder name is too long.")
		return nil
	}
	uiGoToToggle(ui, v)
	go ui.Update(func(g *gocui.Gui) error {
		uiMainGoToFolder(ui, v, goToFolderPath)
		return nil
	})
	return nil
}

func uiGoToAutocomplete(ui *gocui.Gui, v *gocui.View) error {
	var autoCompleteInput = "/"
	goToPath, _ := dgState.goToWindow.view.Line(0)
	if (path.Dir(goToPath) == dgState.goToWindow.state.lastPath) &&
		(path.Base(goToPath)[:1] == dgState.goToWindow.state.lastInitial) {
		dgState.goToWindow.state.index++
		if !strings.HasSuffix(goToPath, "/") {
			autoCompleteInput = path.Base(goToPath)[:1]
		}
	} else {
		dgState.goToWindow.state.lastPath = path.Dir(goToPath)
		dgState.goToWindow.state.lastInitial = path.Base(goToPath)[:1]
		if !strings.HasSuffix(goToPath, "/") {
			autoCompleteInput = path.Base(goToPath)
		}
	}
	ui.Update(func(g *gocui.Gui) error {
		newPath, newIndex := dgFileFolderPathAutocomplete(
			path.Join(path.Dir(goToPath), autoCompleteInput),
			dgState.goToWindow.state.index,
		)
		dgState.goToWindow.state.index = newIndex
		dgState.goToWindow.view.Clear()
		fmt.Fprint(dgState.goToWindow.view, newPath)
		dgState.goToWindow.view.SetCursor(len(newPath), 0)
		return nil
	})
	return nil
}
