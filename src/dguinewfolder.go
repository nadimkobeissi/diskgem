/* SPDX-License-Identifier: MIT
 * Copyright © 2018-2019 Nadim Kobeissi <nadim@nadim.computer>. All Rights Reserved.
 */

package main

import (
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
	if len(newFolderName) < 1 {
		uiMainStatusViewMessage(ui, 0, "Please enter a folder name.")
		return nil
	}
	if len(newFolderName) > 128 {
		uiMainStatusViewMessage(ui, 0, "Folder name is too long.")
		return nil
	}
	uiNewFolderToggle(ui, v)
	ui.Update(func(g *gocui.Gui) error {
		go uiMainCreateFolder(ui, v, newFolderName)
		return nil
	})
	return nil
}
