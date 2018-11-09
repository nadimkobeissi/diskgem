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
	"strings"

	"github.com/jroimartin/gocui"
)

func uiAbout(ui *gocui.Gui, v *gocui.View) error {
	maxX, maxY := ui.Size()
	dgState.aboutWindow.view, _ = ui.SetView("about",
		maxX/2-31, maxY/2-13,
		maxX/2+31, maxY/2+13,
	)
	view := dgState.aboutWindow.view
	view.Frame = true
	view.Title = "About"
	view.BgColor = gocui.ColorWhite
	view.FgColor = gocui.ColorBlack
	aboutText := `
	  DiskGem is software for secure file transfer over SFTP.
	  
	  DiskGem currently offers an easy to use, stable 
	  command-line user interface that supports parallel
	  file transfers and other useful features.

	  DiskGem will soon also support creating encrypted
	  archives on the server which offer encryption
	  of stored files as well as metadata obfuscation.

	  Using DiskGem:
		  - Arrow keys to navigate.
		  - Tab to switch between panes.
		  - Enter to upload/download files.
		  - Delete to delete files.
	`
	fmt.Fprintln(view, strings.Join([]string{
		"\n\n\n", aboutText, "\n  DiskGem ",
		dgState.mainWindow.state.version,
		"\n  https://diskgem.info",
	}, ""))
	return nil
}

func uiAboutToggle(ui *gocui.Gui, v *gocui.View) error {
	if dgState.aboutWindow.state.visible {
		dgState.aboutWindow.state.visible = false
		ui.DeleteView("about")
		if dgState.mainWindow.state.leftPane.focused {
			ui.SetCurrentView("leftPane")
		} else if dgState.mainWindow.state.leftPane.focused {
			ui.SetCurrentView("rightPane")
		}
	} else {
		uiMainHideAllWindowsExcept("about", ui, v)
		dgState.aboutWindow.state.visible = true
		uiAbout(ui, v)
	}
	return nil
}
