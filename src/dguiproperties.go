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
	"fmt"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/jroimartin/gocui"
)

func uiProperties(ui *gocui.Gui, v *gocui.View) error {
	maxX, maxY := ui.Size()
	dgState.propertiesWindow.view, _ = ui.SetView("properties",
		maxX/2-25, maxY/2-10,
		maxX/2+25, maxY/2+11,
	)
	dgState.propertiesWindow.nameView, _ = ui.SetView("propertiesName",
		maxX/2-22, maxY/2-8,
		maxX/2+22, maxY/2-6,
	)
	dgState.propertiesWindow.sizeView, _ = ui.SetView("propertiesSize",
		maxX/2-22, maxY/2-3,
		maxX/2+22, maxY/2-1,
	)
	dgState.propertiesWindow.modDateView, _ = ui.SetView("propertiesModDate",
		maxX/2-22, maxY/2+2,
		maxX/2+22, maxY/2+4,
	)
	dgState.propertiesWindow.permissionsView, _ = ui.SetView("propertiesPermissions",
		maxX/2-22, maxY/2+7,
		maxX/2+22, maxY/2+9,
	)
	dgState.propertiesWindow.view.Title = "Properties"
	dgState.propertiesWindow.nameView.Title = "Name (edit to rename)"
	dgState.propertiesWindow.nameView.Editable = true
	dgState.propertiesWindow.sizeView.Title = "Size"
	dgState.propertiesWindow.modDateView.Title = "Last Modified"
	dgState.propertiesWindow.permissionsView.Title = "Permissions (edit to change)"
	dgState.propertiesWindow.permissionsView.Editable = true
	var selectedIndex *int
	var selectedFile os.FileInfo
	var selectedFileLstat os.FileInfo
	var err error
	if dgState.mainWindow.state.leftPane.focused {
		selectedIndex = &dgState.mainWindow.state.leftPane.selectedIndex
		selectedFile = dgState.mainWindow.state.leftPane.folderContents[*selectedIndex]
		selectedFileLstat, err = os.Lstat(path.Join(
			dgState.mainWindow.state.leftPane.cwd,
			selectedFile.Name(),
		))
	} else if dgState.mainWindow.state.rightPane.focused {
		selectedIndex = &dgState.mainWindow.state.rightPane.selectedIndex
		selectedFile = dgState.mainWindow.state.rightPane.folderContents[*selectedIndex]
		selectedFileLstat, err = sftpClient.Lstat(path.Join(
			dgState.mainWindow.state.rightPane.cwd,
			selectedFile.Name(),
		))
	}
	if err != nil {
		uiMainStatusViewMessage(0, strings.Join([]string{
			"Could not read permissions for ", selectedFile.Name(), ".",
		}, ""))
		return nil
	}
	formatPerm := strconv.FormatInt(int64(selectedFileLstat.Mode().Perm()), 8)
	dgState.propertiesWindow.state.fileName = selectedFileLstat.Name()
	dgState.propertiesWindow.state.permissions = formatPerm
	fmt.Fprintln(dgState.propertiesWindow.nameView, selectedFileLstat.Name())
	fmt.Fprintln(dgState.propertiesWindow.sizeView, dgFileSizeFormat(selectedFileLstat.Size()))
	fmt.Fprintln(dgState.propertiesWindow.modDateView, selectedFileLstat.ModTime())
	fmt.Fprintln(dgState.propertiesWindow.permissionsView, formatPerm)
	uiPropertiesRenderSelection(ui, v)
	return nil
}

func uiPropertiesToggle(ui *gocui.Gui, v *gocui.View) error {
	if dgState.propertiesWindow.state.visible {
		dgState.propertiesWindow.state.visible = false
		ui.Cursor = false
		ui.DeleteView("properties")
		ui.DeleteView("propertiesName")
		ui.DeleteView("propertiesSize")
		ui.DeleteView("propertiesModDate")
		ui.DeleteView("propertiesPermissions")
		if dgState.mainWindow.state.leftPane.focused {
			ui.SetCurrentView("leftPane")
		} else if dgState.mainWindow.state.leftPane.focused {
			ui.SetCurrentView("rightPane")
		}
	} else {
		uiMainHideAllWindowsExcept("properties", ui, v)
		dgState.propertiesWindow.state.visible = true
		ui.Cursor = true
		uiProperties(ui, v)
	}
	return nil
}

func uiPropertiesHandleTab(ui *gocui.Gui, v *gocui.View) error {
	if dgState.propertiesWindow.state.selected == "propertiesName" {
		dgState.propertiesWindow.state.selected = "propertiesPermissions"
	} else if dgState.propertiesWindow.state.selected == "propertiesPermissions" {
		dgState.propertiesWindow.state.selected = "propertiesName"
	}
	uiPropertiesRenderSelection(ui, v)
	return nil
}

func uiPropertiesRenderSelection(ui *gocui.Gui, v *gocui.View) error {
	ui.SetCurrentView(dgState.propertiesWindow.state.selected)
	line, _ := ui.CurrentView().Line(0)
	ui.CurrentView().SetCursor(len(line), 0)
	return nil
}

func uiPropertiesHandleEnter(ui *gocui.Gui, v *gocui.View) error {
	dgState.propertiesWindow.state.selected = "propertiesName"
	propertiesName, _ := dgState.propertiesWindow.nameView.Line(0)
	fileName := dgState.propertiesWindow.state.fileName
	propertiesPermissions, _ := dgState.propertiesWindow.permissionsView.Line(0)
	permissions := dgState.propertiesWindow.state.permissions
	validPermissions, err := regexp.MatchString("^[0-7]{3}$", propertiesPermissions)
	if err != nil || !validPermissions {
		uiMainStatusViewMessage(0, "Invalid file permissions.")
		return nil
	}
	if propertiesPermissions != permissions {
		go uiMainChmodFile(propertiesPermissions, fileName)
	}
	if propertiesName != fileName {
		go uiMainRenameFile(ui, v, propertiesName, fileName)
	}
	uiPropertiesToggle(ui, v)
	return nil
}
