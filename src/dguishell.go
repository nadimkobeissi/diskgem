/* SPDX-License-Identifier: MIT
 * Copyright © 2018 Nadim Kobeissi <nadim@nadim.computer>. All Rights Reserved.
 */

package main

import (
	"os/exec"
	"strings"

	"github.com/jroimartin/gocui"
)

func uiShell(ui *gocui.Gui, v *gocui.View) error {
	maxX, maxY := ui.Size()
	dgState.shellWindow.view, _ = ui.SetView("shell",
		maxX/2-20, maxY/2-5,
		maxX/2+20, maxY/2-3,
	)
	dgState.shellWindow.view.Title = "Shell"
	dgState.shellWindow.view.Editable = true
	ui.SetCurrentView("shell")
	return nil
}

func uiShellToggle(ui *gocui.Gui, v *gocui.View) error {
	if dgState.shellWindow.state.visible {
		ui.Cursor = false
		dgState.shellWindow.state.visible = false
		ui.DeleteView("shell")
		if dgState.mainWindow.state.leftPane.focused {
			ui.SetCurrentView("leftPane")
		} else if dgState.mainWindow.state.leftPane.focused {
			ui.SetCurrentView("rightPane")
		}
	} else {
		uiMainHideAllWindowsExcept("shell", ui, v)
		ui.Cursor = true
		dgState.shellWindow.state.visible = true
		uiShell(ui, v)
	}
	return nil
}

func uiShellHandleEnter(ui *gocui.Gui, v *gocui.View) error {
	shellCommand, _ := dgState.shellWindow.view.Line(0)
	if len(shellCommand) > 0 {
		go uiShellCommandRun(ui, v, shellCommand)
	}
	uiShellToggle(ui, v)
	return nil
}

func uiShellCommandRun(ui *gocui.Gui, v *gocui.View, shellCommand string) error {
	if dgState.mainWindow.state.leftPane.focused {
		shellCommandArgs := strings.Fields(shellCommand)
		cmd := exec.Command(shellCommandArgs[0])
		cmd.Args = shellCommandArgs
		cmd.Dir = dgState.mainWindow.state.leftPane.cwd
		ui.Update(func(g *gocui.Gui) error {
			uiMainStatusViewMessage(1, strings.Join([]string{
				"Dispatching shell command `", shellCommand, "` locally.",
			}, ""))
			return nil
		})
		cmd.Run()
		ui.Update(func(g *gocui.Gui) error {
			uiMainStatusViewMessage(1, strings.Join([]string{
				"Shell command `", shellCommand,
				"` has finished executing locally.\n",
				"      DiskGem cannot guarantee that execution was successful.",
			}, ""))
			return nil
		})
	} else if dgState.mainWindow.state.rightPane.focused {
		ui.Update(func(g *gocui.Gui) error {
			uiMainStatusViewMessage(1, strings.Join([]string{
				"Dispatching shell command `", shellCommand, "` remotely.",
			}, ""))
			return nil
		})
		sshSession, err := dgSSHClient.NewSession()
		if err != nil {
			ui.Update(func(g *gocui.Gui) error {
				uiMainStatusViewMessage(0, "Could not run shell command on archive.")
				return nil
			})
			return nil
		}
		sshSession.Run(strings.Join([]string{
			"cd ", dgState.mainWindow.state.rightPane.cwd,
			"; ", shellCommand,
		}, ""))
		sshSession.Close()
		ui.Update(func(g *gocui.Gui) error {
			uiMainStatusViewMessage(1, strings.Join([]string{
				"Shell command `", shellCommand,
				"` has finished executing remotely.\n",
				"      DiskGem cannot guarantee that execution was successful.",
			}, ""))
			return nil
		})
	}
	return nil
}
