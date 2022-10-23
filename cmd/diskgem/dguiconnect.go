/* SPDX-License-Identifier: MIT
 * Copyright © 2018-2019 Nadim Kobeissi <nadim@symbolic.software>. All Rights Reserved.
 */

package main

import (
	"fmt"
	"os/user"
	"regexp"
	"strings"

	"github.com/jroimartin/gocui"
)

func uiConnect(ui *gocui.Gui, v *gocui.View) error {
	maxX, maxY := ui.Size()
	connectHeight := 22
	connectWidth := 50
	dgState.connectWindow.state.fingerprint = ""
	dgState.connectWindow.view, _ = ui.SetView("connect",
		maxX/2-(connectWidth/2), maxY/2-(connectHeight/2),
		maxX/2+(connectWidth/2), maxY/2+(connectHeight/2),
	)
	dgState.connectWindow.serverURIView, _ = ui.SetView("connectServerURI",
		maxX/2-(connectWidth/2)+2, maxY/2-(connectHeight/2)+2,
		maxX/2+(connectWidth/2)-2, maxY/2-(connectHeight/2)+4,
	)
	dgState.connectWindow.usernameView, _ = ui.SetView("connectUsername",
		maxX/2-(connectWidth/2)+2, maxY/2-(connectHeight/2)+6,
		maxX/2+(connectWidth/2)-2, maxY/2-(connectHeight/2)+8,
	)
	dgState.connectWindow.passwordView, _ = ui.SetView("connectPassword",
		maxX/2-(connectWidth/2)+2, maxY/2-(connectHeight/2)+10,
		maxX/2+(connectWidth/2)-2, maxY/2-(connectHeight/2)+12,
	)
	dgState.connectWindow.passwordInfoView, _ = ui.SetView("connectPasswordInfo",
		maxX/2-(connectWidth/2)+2, maxY/2-(connectHeight/2)+13,
		maxX/2+(connectWidth/2)-2, maxY/2-(connectHeight/2)+15,
	)
	dgState.connectWindow.okButtonView, _ = ui.SetView("connectOKButton",
		maxX/2-(connectWidth/2)+2, maxY/2-(connectHeight/2)+16,
		maxX/2-(connectWidth/2)+15, maxY/2-(connectHeight/2)+20,
	)
	dgState.connectWindow.view.Frame = true
	dgState.connectWindow.view.Title = "Connect"
	dgState.connectWindow.serverURIView.Frame = true
	dgState.connectWindow.serverURIView.Title = "Server URI"
	dgState.connectWindow.serverURIView.Editable = true
	dgState.connectWindow.okButtonView.Frame = true
	fmt.Fprint(dgState.connectWindow.okButtonView, "\n\n     OK")
	dgState.connectWindow.usernameView.Frame = true
	dgState.connectWindow.usernameView.Title = "Username"
	dgState.connectWindow.usernameView.Editable = true
	fmt.Fprint(dgState.connectWindow.serverURIView, dgState.connectWindow.state.serverURI)
	if len(dgState.connectWindow.state.username) == 0 {
		currentUser, _ := user.Current()
		fmt.Fprint(dgState.connectWindow.usernameView, currentUser.Username)
	} else {
		fmt.Fprint(dgState.connectWindow.usernameView, dgState.connectWindow.state.username)
	}
	fmt.Fprint(dgState.connectWindow.passwordView, dgState.connectWindow.state.password)
	dgState.connectWindow.passwordView.Frame = true
	dgState.connectWindow.passwordView.Title = "Password"
	dgState.connectWindow.passwordView.Editable = true
	dgState.connectWindow.passwordView.Mask = '*'
	dgState.connectWindow.passwordInfoView.Frame = false
	fmt.Fprint(
		dgState.connectWindow.passwordInfoView,
		"Leave password empty to use SSH keys.",
	)
	if len(dgState.connectWindow.state.serverURI) > 0 {
		dgState.connectWindow.state.selected = "connectOKButton"
	}
	uiConnectRenderSelection(ui, v)
	return nil
}

func uiConnectToggle(ui *gocui.Gui, v *gocui.View) error {
	if dgState.mainWindow.state.keyVerification {
		dgSFTPConfirmHostKeyVerification(func() {
			uiConnectSftpConnect(ui, v)
		})
	}
	if dgState.mainWindow.state.connected {
		return nil
	}
	if dgState.connectWindow.state.visible {
		dgState.connectWindow.state.visible = false
		ui.Cursor = false
		ui.DeleteView("connect")
		ui.DeleteView("connectServerURI")
		ui.DeleteView("connectUsername")
		ui.DeleteView("connectPassword")
		ui.DeleteView("connectPasswordInfo")
		ui.DeleteView("connectOKButton")
		if dgState.mainWindow.state.leftPane.focused {
			ui.SetCurrentView("leftPane")
		} else if dgState.mainWindow.state.leftPane.focused {
			ui.SetCurrentView("rightPane")
		}
	} else {
		uiMainHideAllWindowsExcept("connect", ui, v)
		dgState.connectWindow.state.visible = true
		uiConnect(ui, v)
	}
	return nil
}

func uiConnectHandleTab(ui *gocui.Gui, v *gocui.View) error {
	if dgState.connectWindow.state.selected == "connectServerURI" {
		dgState.connectWindow.state.selected = "connectUsername"
	} else if dgState.connectWindow.state.selected == "connectUsername" {
		dgState.connectWindow.state.selected = "connectPassword"
	} else if dgState.connectWindow.state.selected == "connectPassword" {
		dgState.connectWindow.state.selected = "connectOKButton"
	} else if dgState.connectWindow.state.selected == "connectOKButton" {
		dgState.connectWindow.state.selected = "connectServerURI"
	}
	uiConnectRenderSelection(ui, v)
	return nil
}

func uiConnectRenderSelection(ui *gocui.Gui, v *gocui.View) error {
	ui.Cursor = (dgState.connectWindow.state.selected != "connectOKButton")
	ui.SetCurrentView(dgState.connectWindow.state.selected)
	line, _ := ui.CurrentView().Line(0)
	ui.CurrentView().SetCursor(len(line), 0)
	if dgState.connectWindow.state.selected == "connectOKButton" {
		ui.Cursor = false
		dgState.connectWindow.okButtonView.BgColor = gocui.ColorWhite
		dgState.connectWindow.okButtonView.FgColor = gocui.ColorBlack
	} else {
		ui.Cursor = true
		dgState.connectWindow.okButtonView.BgColor = gocui.ColorBlack
		dgState.connectWindow.okButtonView.FgColor = gocui.ColorWhite
	}
	return nil
}

func uiConnectHandleEnter(ui *gocui.Gui, v *gocui.View) error {
	dgState.connectWindow.state.selected = "connectServerURI"
	serverURI, _ := dgState.connectWindow.serverURIView.Line(0)
	username, _ := dgState.connectWindow.usernameView.Line(0)
	password, _ := dgState.connectWindow.passwordView.Line(0)
	serverPattern := "^(\\w|\\.)+(:\\d{1,5})?$"
	usernamePattern := "^[a-z_]([a-z0-9_-]{0,31}|[a-z0-9_-]{0,30}\\$)$"
	validServer, _ := regexp.MatchString(serverPattern, serverURI)
	validUsername, _ := regexp.MatchString(usernamePattern, username)
	if !validServer {
		uiMainStatusViewMessage(ui, 0, "Invalid server.")
		return nil
	}
	if !validUsername {
		uiMainStatusViewMessage(ui, 0, "Invalid username.")
		return nil
	}
	uiConnectToggle(ui, v)
	serverHasPort, _ := regexp.MatchString(":\\d{1,5}$", serverURI)
	if !serverHasPort {
		serverURI = strings.Join([]string{serverURI, "22"}, ":")
		uiMainStatusViewMessage(ui, 1, "No port specified, assuming 22.")
	}
	dgState.connectWindow.state.serverURI = serverURI
	dgState.connectWindow.state.username = username
	dgState.connectWindow.state.password = password
	dgState.connectWindow.state.fingerprint = ""
	go uiConnectSftpConnect(ui, v)
	return nil
}

func uiConnectSftpConnect(ui *gocui.Gui, v *gocui.View) error {
	uiMainStatusViewMessage(ui, 1, strings.Join([]string{
		"Connecting to ", dgState.connectWindow.state.serverURI, "...",
	}, ""))
	err := dgSFTPConnect(
		dgState.connectWindow.state.serverURI,
		dgState.connectWindow.state.username,
		dgState.connectWindow.state.password,
	)
	if err != nil {
		if dgState.mainWindow.state.keyVerification {
			uiMainStatusViewMessage(ui, 1, strings.Join([]string{
				dgState.connectWindow.state.serverURI,
				" is offering a public key with fingerprint\n      ",
				dgState.connectWindow.state.fingerprint, "\n      ",
				"Press Ctrl+C to accept key or Ctrl+D to disconnect.",
			}, ""))
		} else {
			uiMainStatusViewMessage(ui, 0, strings.Join([]string{
				"Cannot connect to ", dgState.connectWindow.state.serverURI, ".",
			}, ""))
		}
		return err
	} else {
		uiMainStatusViewMessage(ui, 1, strings.Join([]string{
			"Server fingerprint accepted:", "\n      ",
			dgState.connectWindow.state.fingerprint,
		}, ""))
	}
	dgState.mainWindow.state.connected = true
	err = dgConfigSave()
	if err != nil {
		uiMainStatusViewMessage(ui, 0, "Could not write to config file.")
	}
	uiMainStatusViewMessage(ui, 2, strings.Join([]string{
		"Connected to ", dgState.connectWindow.state.serverURI, ".",
	}, ""))
	uiMainStatusViewMessage(ui, 1, "Listing archive...")
	var lastFolder string
	for _, server := range dgState.mainWindow.state.knownServers {
		if server.Hostname == dgState.connectWindow.state.serverURI {
			if len(server.LastFolder) > 0 {
				lastFolder = server.LastFolder
			}
		}
	}
	if len(lastFolder) > 0 {
		_, err = dgSFTPClient.ReadDir(lastFolder)
		if err == nil {
			dgState.mainWindow.state.rightPane.cwd = lastFolder
		} else {
			dgState.mainWindow.state.rightPane.cwd, err = dgSFTPClient.Getwd()
		}
	} else {
		dgState.mainWindow.state.rightPane.cwd, err = dgSFTPClient.Getwd()
	}
	if err != nil {
		uiMainStatusViewMessage(ui, 0, "Could not list remote starting folder. Disconnecting.")
		uiMainDisconnect(ui, v)
		return err
	}
	go uiMainListArchiveFolder(ui)
	uiMainMenuViewUpdate()
	uiMainMenuStatusViewUpdate()
	return nil
}
