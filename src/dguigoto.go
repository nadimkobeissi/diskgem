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
	"path"

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
		uiMainStatusViewMessage(0, "Please enter a folder name.")
		return nil
	}
	if len(goToFolderPath) > 256 {
		uiMainStatusViewMessage(0, "Folder name is too long.")
		return nil
	}
	uiGoToToggle(ui, v)
	go func() {
		ui.Update(func(g *gocui.Gui) error {
			uiMainGoToFolder(ui, v, goToFolderPath)
			return nil
		})
	}()
	return nil
}

func uiGoToAutocomplete(ui *gocui.Gui, v *gocui.View) error {
	goToPath, _ := dgState.goToWindow.view.Line(0)
	if (path.Dir(goToPath) == dgState.goToWindow.state.lastPath) &&
		(path.Base(goToPath)[:1] == dgState.goToWindow.state.lastInitial) {
		dgState.goToWindow.state.index++
	} else {
		dgState.goToWindow.state.lastPath = path.Dir(goToPath)
		dgState.goToWindow.state.lastInitial = path.Base(goToPath)[:1]
	}

	ui.Update(func(g *gocui.Gui) error {
		newPath, newIndex := dgFileFolderPathAutocomplete(
			path.Join(
				dgState.goToWindow.state.lastPath,
				dgState.goToWindow.state.lastInitial,
			),
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
