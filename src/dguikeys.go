/* SPDX-License-Identifier: MIT
 * Copyright © 2018 Nadim Kobeissi <nadim@nadim.computer>. All Rights Reserved.
 */

package main

import (
	"github.com/jroimartin/gocui"
)

func uiKeysRuneBinder(key rune) func(*gocui.Gui, *gocui.View) error {
	return func(u *gocui.Gui, v *gocui.View) error {
		uiKeysRegularFilter(key)
		return nil
	}
}

func uiKeysBind(ui *gocui.Gui) error {
	var regularKeyRunes = []rune{
		'0', '1', '2', '3', '4', '5', '6', '7',
		'8', '9', 'a', 'b', 'c', 'd', 'e', 'f',
		'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n',
		'o', 'p', 'q', 'r', 's', 't', 'u', 'v',
		'w', 'x', 'y', 'z', 'A', 'B', 'C', 'D',
		'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L',
		'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T',
		'U', 'V', 'W', 'X', 'Y', 'Z',
	}
	for _, key := range regularKeyRunes {
		ui.SetKeybinding("leftPane", key, gocui.ModNone, uiKeysRuneBinder(key))
		ui.SetKeybinding("rightPane", key, gocui.ModNone, uiKeysRuneBinder(key))
	}
	ui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, uiKeysCtrlC)
	ui.SetKeybinding("", gocui.KeyCtrlD, gocui.ModNone, uiKeysCtrlD)
	ui.SetKeybinding("", gocui.KeyCtrlN, gocui.ModNone, uiKeysCtrlN)
	ui.SetKeybinding("", gocui.KeyCtrlG, gocui.ModNone, uiKeysCtrlG)
	ui.SetKeybinding("", gocui.KeyCtrlP, gocui.ModNone, uiKeysCtrlP)
	ui.SetKeybinding("", gocui.KeyCtrlR, gocui.ModNone, uiKeysCtrlR)
	ui.SetKeybinding("", gocui.KeyCtrlA, gocui.ModNone, uiKeysCtrlA)
	ui.SetKeybinding("", gocui.KeyCtrlQ, gocui.ModNone, uiKeysCtrlQ)
	ui.SetKeybinding("", gocui.KeyTab, gocui.ModNone, uiKeysTab)
	ui.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, uiKeysEnter)
	ui.SetKeybinding("", gocui.KeyDelete, gocui.ModNone, uiKeysDelete)
	ui.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, uiKeysArrowUp)
	ui.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, uiKeysArrowDown)
	ui.SetKeybinding("", gocui.KeyArrowLeft, gocui.ModNone, uiKeysArrowLeft)
	ui.SetKeybinding("", gocui.KeyArrowRight, gocui.ModNone, uiKeysArrowRight)
	ui.SetKeybinding("", gocui.KeyPgup, gocui.ModNone, uiKeysPgup)
	ui.SetKeybinding("", gocui.KeyPgdn, gocui.ModNone, uiKeysPgdn)
	return nil
}

func uiKeysRegularFilter(key rune) error {
	if dgState.connectWindow.state.visible {

	}
	if dgState.newFolderWindow.state.visible {
		return nil
	}
	if dgState.goToWindow.state.visible {
	}
	if dgState.propertiesWindow.state.visible {

	}
	if dgState.aboutWindow.state.visible {
		return nil
	}
	uiMainJumpToRune(key)
	return nil
}

func uiKeysCtrlC(ui *gocui.Gui, v *gocui.View) error {
	uiConnectToggle(ui, v)
	return nil
}

func uiKeysCtrlD(ui *gocui.Gui, v *gocui.View) error {
	uiMainDisconnect(ui, v)
	return nil
}

func uiKeysCtrlN(ui *gocui.Gui, v *gocui.View) error {
	uiNewFolderToggle(ui, v)
	return nil
}

func uiKeysCtrlG(ui *gocui.Gui, v *gocui.View) error {
	uiGoToToggle(ui, v)
	return nil
}

func uiKeysCtrlP(ui *gocui.Gui, v *gocui.View) error {
	uiPropertiesToggle(ui, v)
	return nil
}

func uiKeysCtrlR(ui *gocui.Gui, v *gocui.View) error {
	ui.Update(func(g *gocui.Gui) error {
		go uiMainRefresh(ui, v)
		return nil
	})
	return nil
}

func uiKeysCtrlA(ui *gocui.Gui, v *gocui.View) error {
	uiAboutToggle(ui, v)
	return nil
}

func uiKeysCtrlQ(ui *gocui.Gui, v *gocui.View) error {
	uiMainQuit(ui, v)
	return nil
}

func uiKeysTab(ui *gocui.Gui, v *gocui.View) error {
	if dgState.connectWindow.state.visible {
		return uiConnectHandleTab(ui, v)
	}
	if dgState.newFolderWindow.state.visible {
		return nil
	}
	if dgState.goToWindow.state.visible {
		go uiGoToAutocomplete(ui, v)
		return nil
	}
	if dgState.propertiesWindow.state.visible {
		uiPropertiesHandleTab(ui, v)
		return nil
	}
	if dgState.aboutWindow.state.visible {
		return nil
	}
	if dgState.mainWindow.state.leftPane.focused {
		if !dgState.mainWindow.state.connected {
			return nil
		}
		dgState.mainWindow.state.leftPane.focused = false
		dgState.mainWindow.state.rightPane.focused = true
		ui.SetCurrentView("leftPane")
	} else if dgState.mainWindow.state.rightPane.focused {
		dgState.mainWindow.state.leftPane.focused = true
		dgState.mainWindow.state.rightPane.focused = false
		ui.SetCurrentView("rightPane")
	}
	uiMainRenderPane(false)
	uiMainRenderPane(true)
	return nil
}

func uiKeysEnter(ui *gocui.Gui, v *gocui.View) error {
	if dgState.connectWindow.state.visible {
		if ui.CurrentView().Name() == "connectOKButton" {
			uiConnectHandleEnter(ui, v)
		}
		return nil
	}
	if dgState.newFolderWindow.state.visible {
		uiNewFolderHandleEnter(ui, v)
		return nil
	}
	if dgState.goToWindow.state.visible {
		uiGoToHandleEnter(ui, v)
		return nil
	}
	if dgState.propertiesWindow.state.visible {
		uiPropertiesHandleEnter(ui, v)
		return nil
	}
	if dgState.aboutWindow.state.visible {
		return nil
	}
	if !dgState.mainWindow.state.connected {
		uiMainStatusViewMessage(0, "Connection is not established.")
		return nil
	}
	if dgState.mainWindow.state.leftPane.focused &&
		len(dgState.mainWindow.state.leftPane.folderContents) > 0 {
		go uiMainFileUpload(ui, v)
		return nil
	} else if dgState.mainWindow.state.rightPane.focused &&
		len(dgState.mainWindow.state.rightPane.folderContents) > 0 {
		go uiMainFileDownload(ui, v)
		return nil
	}
	return nil
}

func uiKeysDelete(ui *gocui.Gui, v *gocui.View) error {
	if dgState.connectWindow.state.visible {
		return nil
	}
	if dgState.newFolderWindow.state.visible {
		return nil
	}
	if dgState.goToWindow.state.visible {
		return nil
	}
	if dgState.propertiesWindow.state.visible {
		return nil
	}
	if dgState.aboutWindow.state.visible {
		return nil
	}
	uiMainDeleteSelected(ui, v)
	return nil
}

func uiKeysArrowUp(ui *gocui.Gui, v *gocui.View) error {
	if dgState.connectWindow.state.visible {
		return nil
	}
	if dgState.newFolderWindow.state.visible {
		return nil
	}
	if dgState.goToWindow.state.visible {
		return nil
	}
	if dgState.propertiesWindow.state.visible {
		return nil
	}
	if dgState.aboutWindow.state.visible {
		return nil
	}
	uiMainNavigateUp(ui, v)
	return nil
}

func uiKeysArrowDown(ui *gocui.Gui, v *gocui.View) error {
	if dgState.connectWindow.state.visible {
		return nil
	}
	if dgState.newFolderWindow.state.visible {
		return nil
	}
	if dgState.goToWindow.state.visible {
		return nil
	}
	if dgState.propertiesWindow.state.visible {
		return nil
	}
	if dgState.aboutWindow.state.visible {
		return nil
	}
	uiMainNavigateDown(ui, v)
	return nil
}

func uiKeysArrowLeft(ui *gocui.Gui, v *gocui.View) error {
	if dgState.connectWindow.state.visible {
		ui.CurrentView().MoveCursor(-1, 0, true)
		return nil
	}
	if dgState.newFolderWindow.state.visible {
		ui.CurrentView().MoveCursor(-1, 0, true)
		return nil
	}
	if dgState.goToWindow.state.visible {
		ui.CurrentView().MoveCursor(-1, 0, true)
		return nil
	}
	if dgState.propertiesWindow.state.visible {
		ui.CurrentView().MoveCursor(-1, 0, true)
		return nil
	}
	if dgState.aboutWindow.state.visible {
		return nil
	}
	go func() {
		ui.Update(func(g *gocui.Gui) error {
			uiMainNavigateLeft(ui, v)
			return nil
		})
	}()
	return nil
}

func uiKeysArrowRight(ui *gocui.Gui, v *gocui.View) error {
	if dgState.connectWindow.state.visible {
		x, _ := ui.CurrentView().Cursor()
		line, _ := ui.CurrentView().Line(0)
		if len(line) > x {
			ui.CurrentView().MoveCursor(1, 0, true)
		}
		return nil
	}
	if dgState.newFolderWindow.state.visible {
		x, _ := ui.CurrentView().Cursor()
		line, _ := ui.CurrentView().Line(0)
		if len(line) > x {
			ui.CurrentView().MoveCursor(1, 0, true)
		}
		return nil
	}
	if dgState.goToWindow.state.visible {
		x, _ := ui.CurrentView().Cursor()
		line, _ := ui.CurrentView().Line(0)
		if len(line) > x {
			ui.CurrentView().MoveCursor(1, 0, true)
		}
		return nil
	}
	if dgState.propertiesWindow.state.visible {
		x, _ := ui.CurrentView().Cursor()
		line, _ := ui.CurrentView().Line(0)
		if len(line) > x {
			ui.CurrentView().MoveCursor(1, 0, true)
		}
		return nil
	}
	if dgState.aboutWindow.state.visible {
		return nil
	}
	go func() {
		ui.Update(func(g *gocui.Gui) error {
			uiMainNavigateRight(ui, v)
			return nil
		})
	}()
	return nil
}

func uiKeysPgup(ui *gocui.Gui, v *gocui.View) error {
	// This looks hacky, but seems to work
	// very well across the interface.
	for i := 0; i < 16; i++ {
		uiKeysArrowUp(ui, v)
	}
	return nil
}

func uiKeysPgdn(ui *gocui.Gui, v *gocui.View) error {
	// This looks hacky, but seems to work
	// very well across the interface.
	for i := 0; i < 16; i++ {
		uiKeysArrowDown(ui, v)
	}
	return nil
}
