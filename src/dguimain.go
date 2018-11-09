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
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/aybabtme/rgbterm"
	"github.com/jroimartin/gocui"
)

var uiMainTransfersTicker dgticker

func uiMainStartTransfersTicker(ui *gocui.Gui) error {
	uiMainTransfersTicker.gears = time.NewTicker(500 * time.Millisecond)
	uiMainTransfersTicker.active = true
	go func() {
		for range uiMainTransfersTicker.gears.C {
			ui.Update(func(g *gocui.Gui) error {
				uiMainTransfersViewUpdate()
				return nil
			})
		}
	}()
	return nil
}

func uiMainStopTransfersTicker() error {
	uiMainTransfersTicker.gears.Stop()
	uiMainTransfersTicker.active = false
	return nil
}

func uiMainInit(ui *gocui.Gui) error {
	print("\033[H\033[2J")
	dgState.mainWindow.state.visible = true
	dgState.mainWindow.menuView.Frame = true
	dgState.mainWindow.menuView.Title = strings.Join(
		[]string{"DiskGem", dgState.mainWindow.state.version}, " ",
	)
	dgState.mainWindow.menuView.BgColor = gocui.ColorWhite
	dgState.mainWindow.menuView.FgColor = gocui.ColorBlack
	dgState.mainWindow.menuStatusView.Frame = false
	dgState.mainWindow.menuStatusView.FgColor = gocui.ColorBlack
	dgState.mainWindow.leftPaneView.Frame = true
	dgState.mainWindow.leftPaneView.Title = "Local"
	dgState.mainWindow.rightPaneView.Frame = true
	dgState.mainWindow.rightPaneView.Title = "Archive"
	dgState.mainWindow.statusView.Frame = true
	dgState.mainWindow.statusView.Title = "Status"
	dgState.mainWindow.statusView.Wrap = true
	dgState.mainWindow.statusView.Autoscroll = true
	dgState.mainWindow.transfersView.Frame = true
	dgState.mainWindow.transfersView.Title = "Transfers (0)"
	dgState.mainWindow.transfersView.Wrap = true
	ui.SetCurrentView("leftPane")
	uiMainMenuViewUpdate()
	uiMainMenuStatusViewUpdate()
	uiMainListLocalFolder()
	uiKeysBind(ui)
	coinFlip, _ := rand.Int(rand.Reader, big.NewInt(10))
	if coinFlip.Int64() == int64(5) {
		ui.Update(func(g *gocui.Gui) error {
			go dgUpdateCheck(func() {
				uiMainStatusViewMessage(0, "Could not check for DiskGem software updates.")
			}, func() {
				uiMainStatusViewMessage(
					1,
					"DiskGem software update is available! Download it from https://diskgem.info.",
				)
			})
			return nil
		})
	}
	return nil
}

func uiMainManagerLayout(ui *gocui.Gui) error {
	maxX, maxY := ui.Size()
	dgState.mainWindow.menuView, _ = ui.SetView("menu", 0, 0, maxX-1, 2)
	dgState.mainWindow.menuStatusView, _ = ui.SetView("menuStatus", maxX-9, 0, maxX-1, 2)
	dgState.mainWindow.leftPaneView, _ = ui.SetView("leftPane", 0, 3, maxX/2-1, maxY-10)
	dgState.mainWindow.rightPaneView, _ = ui.SetView("rightPane", maxX/2, 3, maxX-1, maxY-10)
	dgState.mainWindow.statusView, _ = ui.SetView("status", 0, maxY-9, (66*maxX/100)-1, maxY-1)
	dgState.mainWindow.transfersView, _ = ui.SetView("transfers", (66 * maxX / 100), maxY-9, maxX-1, maxY-1)
	if !dgState.mainWindow.state.visible {
		uiMainInit(ui)
	}
	return nil
}

func uiMainMenuViewUpdate() error {
	dgState.mainWindow.menuView.Clear()
	fmt.Fprint(dgState.mainWindow.menuView, strings.Join([]string{
		rgbterm.BgString("Ctrl+", 78, 152, 184), "  ",
	}, ""))
	if dgState.mainWindow.state.connected {
		fmt.Fprint(dgState.mainWindow.menuView, strings.Join([]string{
			rgbterm.FgString("D", 225, 100, 100),
			"isconnect   ",
		}, ""))
	} else {
		fmt.Fprint(dgState.mainWindow.menuView, strings.Join([]string{
			rgbterm.FgString("C", 68, 132, 161),
			"onnect   ",
		}, ""))
	}
	fmt.Fprint(dgState.mainWindow.menuView, strings.Join([]string{
		rgbterm.FgString("N", 68, 132, 161),
		"ew Folder   ",
		rgbterm.FgString("G", 68, 132, 161),
		"o To   ",
		rgbterm.FgString("P", 68, 132, 161),
		"roperties   ",
		rgbterm.FgString("R", 68, 132, 161),
		"efresh   ",
		rgbterm.FgString("A", 68, 132, 161),
		"bout   ",
		rgbterm.FgString("Q", 225, 100, 100),
		"uit",
	}, ""))
	return nil
}

func uiMainMenuStatusViewUpdate() error {
	dgState.mainWindow.menuStatusView.Clear()
	if dgState.mainWindow.state.connected {
		fmt.Fprint(dgState.mainWindow.menuStatusView, rgbterm.BgString(" Online", 74, 197, 160))
	} else {
		fmt.Fprint(dgState.mainWindow.menuStatusView, rgbterm.BgString("Offline", 225, 100, 100))
	}
	return nil
}

func uiMainStatusViewMessage(kind int, message string) error {
	preface := rgbterm.FgString("Error", 225, 100, 100)
	if kind == 1 {
		preface = rgbterm.FgString("Note", 78, 152, 184)
	} else if kind == 2 {
		preface = rgbterm.FgString("Success", 74, 197, 160)
	}
	fmt.Fprintln(dgState.mainWindow.statusView, strings.Join([]string{
		preface, message,
	}, ": "))
	return nil
}

func uiMainTransfersViewUpdate() error {
	for i := len(dgState.mainWindow.state.transfers) - 1; i >= 0; i-- {
		if dgState.mainWindow.state.transfers[i].finished {
			dgState.mainWindow.state.transfers = append(
				dgState.mainWindow.state.transfers[:i],
				dgState.mainWindow.state.transfers[i+1:]...,
			)
		}
	}
	dgState.mainWindow.transfersView.Title = strings.Join([]string{
		"Transfers ", "(",
		strconv.Itoa(len(dgState.mainWindow.state.transfers)), ")",
	}, "")
	dgState.mainWindow.transfersView.Clear()
	for _, t := range dgState.mainWindow.state.transfers {
		tName := t.name
		if len(t.name) > 20 {
			tName = strings.Join([]string{t.name[0:20], "..."}, "")
		}
		fmt.Fprintln(dgState.mainWindow.transfersView, strings.Join([]string{
			t.direction, " ", tName, " (",
			rgbterm.FgString(dgFileSizeFormat(t.size), 150, 150, 150),
			") ", rgbterm.FgString(strings.Join([]string{
				strconv.Itoa(t.progress), "%",
			}, ""), 78, 152, 184),
		}, ""))
	}
	if len(dgState.mainWindow.state.transfers) == 0 {
		uiMainStopTransfersTicker()
	}
	return nil
}

func uiMainListLocalFolder() error {
	files, err := ioutil.ReadDir(dgState.mainWindow.state.leftPane.cwd)
	if err != nil {
		uiMainStatusViewMessage(0, "Could not list local folder.")
	}
	dgState.mainWindow.state.leftPane.folderContents = dgFileInfoSort(files)
	for i, file := range dgState.mainWindow.state.leftPane.folderContents {
		if file.Name() == dgState.mainWindow.state.leftPane.lastFolder {
			dgState.mainWindow.state.leftPane.selectedIndex = i
			break
		}
	}
	uiMainRenderPane(false)
	return nil
}

func uiMainListArchiveFolder(ui *gocui.Gui) error {
	if !dgState.mainWindow.state.connected {
		ui.Update(func(g *gocui.Gui) error {
			return uiMainRenderPane(true)
		})
		return nil
	}
	archiveFiles, err := dgSFTPClient.ReadDir(dgState.mainWindow.state.rightPane.cwd)
	if err != nil {
		ui.Update(func(g *gocui.Gui) error {
			return uiMainStatusViewMessage(0, err.Error())
		})
		return err
	}
	dgState.mainWindow.state.rightPane.folderContents = dgFileInfoSort(archiveFiles)
	for i, file := range dgState.mainWindow.state.rightPane.folderContents {
		if file.Name() == dgState.mainWindow.state.rightPane.lastFolder {
			dgState.mainWindow.state.leftPane.selectedIndex = i
			break
		}
	}
	ui.Update(func(g *gocui.Gui) error {
		return uiMainRenderPane(true)
	})
	return nil
}

func uiMainRenderPane(rightPane bool) error {
	paneView := dgState.mainWindow.leftPaneView
	paneState := &dgState.mainWindow.state.leftPane
	paneTitle := "Local"
	if rightPane {
		paneView = dgState.mainWindow.rightPaneView
		paneState = &dgState.mainWindow.state.rightPane
		paneTitle = "Archive"
	}
	for len(paneState.folderContents)-1 < paneState.selectedIndex {
		paneState.selectedIndex--
	}
	if paneState.selectedIndex < 0 {
		paneState.selectedIndex = 0
	}
	_, ySize := paneView.Size()
	origin := paneState.selectedIndex - (ySize / 2)
	if origin < 0 {
		origin = 0
	}
	paneView.SetOrigin(0, origin)
	paneView.Clear()
	folderCount := 0
	fileCount := 0
	for i, file := range paneState.folderContents {
		isHidden, _ := regexp.MatchString("^\\.", file.Name())
		fileString := rgbterm.FgString(strings.Join([]string{
			"📄", file.Name(),
		}, "  "), 200, 200, 200)
		if file.IsDir() {
			folderCount++
		} else {
			fileCount++
		}
		if isHidden {
			if file.IsDir() || dgFileIsSymlink(file) {
				fileString = rgbterm.FgString(strings.Join([]string{
					"📁", file.Name(),
				}, "  "), 0, 100, 150)
			} else {
				fileString = rgbterm.FgString(strings.Join([]string{
					"📄", file.Name(),
				}, "  "), 100, 100, 100)
			}
		} else {
			if file.IsDir() || dgFileIsSymlink(file) {
				fileString = rgbterm.FgString(strings.Join([]string{
					"📁", file.Name(),
				}, "  "), 0, 200, 255)
			}
		}
		if i == paneState.selectedIndex && paneState.focused {
			fmt.Fprintln(paneView,
				rgbterm.BgString(fileString, 255, 255, 255),
			)
		} else {
			fmt.Fprintln(paneView, fileString)
		}
	}
	if dgState.mainWindow.state.connected || !rightPane {
		folderS := "s"
		fileS := "s"
		if folderCount == 1 {
			folderS = ""
		}
		if fileCount == 1 {
			fileS = ""
		}
		pathBase := path.Base(paneState.cwd)
		if len(pathBase) > 20 {
			pathBase = strings.Join([]string{pathBase[0:20], "..."}, "")
		}
		paneView.Title = strings.Join([]string{
			paneTitle, " (", pathBase, ") ",
			strconv.Itoa(folderCount), " folder", folderS, ", ",
			strconv.Itoa(fileCount), " file", fileS,
		}, "")
	} else {
		paneView.Title = paneTitle
	}
	return nil
}

func uiMainRenderFocusedPane() error {
	uiMainRenderPane(dgState.mainWindow.state.rightPane.focused)
	return nil
}

func uiMainHandleChdir(ui *gocui.Gui, v *gocui.View, folder string) error {
	if len(folder) > 0 {
		uiMainStatusViewMessage(1, strings.Join([]string{
			"Entering ", folder, "...",
		}, ""))
	}
	if dgState.mainWindow.state.leftPane.focused {
		uiMainListLocalFolder()
	} else if dgState.mainWindow.state.rightPane.focused {
		go uiMainListArchiveFolder(ui)
	}
	return nil
}

func uiMainHideAllWindowsExcept(except string, ui *gocui.Gui, v *gocui.View) error {
	if dgState.connectWindow.state.visible {
		if except != "connect" {
			uiConnectToggle(ui, v)
		}
	}
	if dgState.newFolderWindow.state.visible {
		if except != "newFolder" {
			uiNewFolderToggle(ui, v)
		}
	}
	if dgState.goToWindow.state.visible {
		if except != "goTo" {
			uiGoToToggle(ui, v)
		}
	}
	if dgState.propertiesWindow.state.visible {
		if except != "properties" {
			uiPropertiesToggle(ui, v)
		}
	}
	if dgState.aboutWindow.state.visible {
		if except != "about" {
			uiAboutToggle(ui, v)
		}
	}
	return nil
}

func uiMainFileUpload(ui *gocui.Gui, v *gocui.View) error {
	selectedIndex := &dgState.mainWindow.state.leftPane.selectedIndex
	selectedFile := dgState.mainWindow.state.leftPane.folderContents[*selectedIndex]
	selectedFilePath := path.Join(dgState.mainWindow.state.leftPane.cwd, selectedFile.Name())
	archiveFilePath := path.Join(dgState.mainWindow.state.rightPane.cwd, selectedFile.Name())
	if selectedFile.IsDir() || dgFileIsSymlink(selectedFile) {
		return nil
	}
	for _, v := range dgState.mainWindow.state.transfers {
		if v.toPath == archiveFilePath && v.direction == "→" && !v.finished {
			uiMainStatusViewMessage(1, strings.Join([]string{
				"Cannot upload ", selectedFile.Name(),
				", the destination is already being written to.",
			}, ""))
			return nil
		}
	}
	var thisTransfer dgtransfer
	dgFileUpload(selectedFile, selectedFilePath, archiveFilePath, func() {
		uiMainStatusViewMessage(1, strings.Join([]string{
			"Uploading ", selectedFile.Name(), "...",
		}, ""))
		thisTransfer = dgtransfer{
			name:      selectedFile.Name(),
			size:      selectedFile.Size(),
			direction: "→",
			progress:  0,
			fromPath:  selectedFilePath,
			toPath:    archiveFilePath,
			finished:  false,
		}
		dgState.mainWindow.state.transfers = append(
			dgState.mainWindow.state.transfers, &thisTransfer,
		)
		if !uiMainTransfersTicker.active {
			uiMainStartTransfersTicker(ui)
		}
		go uiMainListArchiveFolder(ui)
	}, func(p int) {
		thisTransfer.progress = p
	}, func(err error) {
		thisTransfer.finished = true
		if err != nil {
			uiMainStatusViewMessage(0, strings.Join([]string{
				"Could not upload ", selectedFile.Name(), ".",
			}, ""))
			return
		}
		thisTransfer.progress = 100
		uiMainStatusViewMessage(2, strings.Join([]string{
			"Uploaded ", selectedFile.Name(), ".",
		}, ""))
	})
	return nil
}

func uiMainFileDownload(ui *gocui.Gui, v *gocui.View) error {
	selectedIndex := &dgState.mainWindow.state.rightPane.selectedIndex
	selectedFile := dgState.mainWindow.state.rightPane.folderContents[*selectedIndex]
	selectedFilePath := path.Join(dgState.mainWindow.state.rightPane.cwd, selectedFile.Name())
	localFilePath := path.Join(dgState.mainWindow.state.leftPane.cwd, selectedFile.Name())
	if selectedFile.IsDir() || dgFileIsSymlink(selectedFile) {
		return nil
	}
	for _, v := range dgState.mainWindow.state.transfers {
		if v.toPath == localFilePath && v.direction == "←" && !v.finished {
			uiMainStatusViewMessage(1, strings.Join([]string{
				"Cannot download ", selectedFile.Name(),
				", the destination is already being written to.",
			}, ""))
			return nil
		}
	}
	var thisTransfer dgtransfer
	dgFileDownload(selectedFile, selectedFilePath, localFilePath, func() {
		uiMainStatusViewMessage(1, strings.Join([]string{
			"Downloading ", selectedFile.Name(), "...",
		}, ""))
		thisTransfer = dgtransfer{
			name:      selectedFile.Name(),
			size:      selectedFile.Size(),
			direction: "←",
			progress:  0,
			fromPath:  selectedFilePath,
			toPath:    localFilePath,
			finished:  false,
		}
		dgState.mainWindow.state.transfers = append(
			dgState.mainWindow.state.transfers, &thisTransfer,
		)
		if !uiMainTransfersTicker.active {
			uiMainStartTransfersTicker(ui)
		}
		ui.Update(func(g *gocui.Gui) error {
			uiMainListLocalFolder()
			return nil
		})
	}, func(p int) {
		thisTransfer.progress = p
	}, func(err error) {
		thisTransfer.finished = true
		if err != nil {
			uiMainStatusViewMessage(0, strings.Join([]string{
				"Could not download ", selectedFile.Name(), ".",
			}, ""))
			return
		}
		thisTransfer.progress = 100
		uiMainStatusViewMessage(2, strings.Join([]string{
			"Downloaded ", selectedFile.Name(), ".",
		}, ""))
	})
	return nil
}

func uiMainRefresh(ui *gocui.Gui, v *gocui.View) error {
	if dgState.mainWindow.state.leftPane.focused {
		uiMainListLocalFolder()
		uiMainStatusViewMessage(1, "Refreshed local folder.")
	} else if dgState.mainWindow.state.rightPane.focused {
		uiMainListArchiveFolder(ui)
		uiMainStatusViewMessage(1, "Refreshed archive folder.")
	}
	return nil
}

func uiMainNavigateUp(ui *gocui.Gui, v *gocui.View) error {
	if dgState.mainWindow.state.leftPane.focused {
		if dgState.mainWindow.state.leftPane.selectedIndex > 0 {
			dgState.mainWindow.state.leftPane.selectedIndex--
		}
	} else if dgState.mainWindow.state.rightPane.focused {
		if dgState.mainWindow.state.rightPane.selectedIndex > 0 {
			dgState.mainWindow.state.rightPane.selectedIndex--
		}
	}
	uiMainRenderFocusedPane()
	return nil
}

func uiMainNavigateDown(ui *gocui.Gui, v *gocui.View) error {
	if dgState.mainWindow.state.leftPane.focused {
		leftMaxSel := len(dgState.mainWindow.state.leftPane.folderContents) - 1
		if dgState.mainWindow.state.leftPane.selectedIndex < leftMaxSel {
			dgState.mainWindow.state.leftPane.selectedIndex++
		}
	} else if dgState.mainWindow.state.rightPane.focused {
		rightMaxSel := len(dgState.mainWindow.state.rightPane.folderContents) - 1
		if dgState.mainWindow.state.rightPane.selectedIndex < rightMaxSel {
			dgState.mainWindow.state.rightPane.selectedIndex++
		}
	}
	uiMainRenderFocusedPane()
	return nil
}

func uiMainNavigateLeft(ui *gocui.Gui, v *gocui.View) error {
	var err error
	paneState := &dgState.mainWindow.state.leftPane
	if dgState.mainWindow.state.rightPane.focused {
		paneState = &dgState.mainWindow.state.rightPane
	}
	cwd := path.Join(paneState.cwd, "..")
	if len(cwd) == 0 {
		return nil
	}
	if dgState.mainWindow.state.leftPane.focused {
		_, err = ioutil.ReadDir(cwd)
	} else if dgState.mainWindow.state.rightPane.focused {
		_, err = dgSFTPClient.ReadDir(cwd)
	}
	if err != nil {
		uiMainStatusViewMessage(0, strings.Join([]string{
			"Cannot access ", cwd, ".",
		}, ""))
		return err
	}
	paneState.lastFolder = path.Base(paneState.cwd)
	paneState.cwd = cwd
	return uiMainHandleChdir(ui, v, "")
}

func uiMainNavigateRight(ui *gocui.Gui, v *gocui.View) error {
	var err error
	paneState := &dgState.mainWindow.state.leftPane
	if dgState.mainWindow.state.rightPane.focused {
		paneState = &dgState.mainWindow.state.rightPane
	}
	if len(paneState.folderContents) == 0 {
		return nil
	}
	selectedFile := paneState.folderContents[paneState.selectedIndex]
	if selectedFile.IsDir() {
		if dgState.mainWindow.state.leftPane.focused {
			_, err = ioutil.ReadDir(path.Join(paneState.cwd, selectedFile.Name()))
		} else if dgState.mainWindow.state.rightPane.focused {
			_, err = dgSFTPClient.ReadDir(path.Join(paneState.cwd, selectedFile.Name()))
		}
		if err != nil {
			uiMainStatusViewMessage(0, strings.Join([]string{
				"Cannot access ", path.Join(paneState.cwd, selectedFile.Name()), ".",
			}, ""))
			return err
		}
		paneState.cwd = path.Join(paneState.cwd, selectedFile.Name())
		return uiMainHandleChdir(ui, v, selectedFile.Name())
	} else if dgFileIsSymlink(selectedFile) {
		var link string
		if dgState.mainWindow.state.leftPane.focused {
			link, err = os.Readlink(path.Join(paneState.cwd, selectedFile.Name()))
		} else if dgState.mainWindow.state.rightPane.focused {
			link, err = dgSFTPClient.ReadLink(path.Join(paneState.cwd, selectedFile.Name()))
		}
		if err != nil {
			uiMainStatusViewMessage(0, "Could not read symbolic link.")
			return nil
		}
		if path.IsAbs(link) {
			paneState.cwd = path.Clean(link)
		} else {
			paneState.cwd = path.Join(paneState.cwd, link)
		}
		return uiMainHandleChdir(ui, v, selectedFile.Name())
	}
	return nil
}

func uiMainJumpToRune(key rune) error {
	if dgState.mainWindow.state.leftPane.focused {
		for i, v := range dgState.mainWindow.state.leftPane.folderContents {
			if strings.HasPrefix(v.Name(), string(key)) {
				dgState.mainWindow.state.leftPane.selectedIndex = i
				break
			}
		}
	} else if dgState.mainWindow.state.rightPane.focused {
		for i, v := range dgState.mainWindow.state.rightPane.folderContents {
			if strings.HasPrefix(v.Name(), string(key)) {
				dgState.mainWindow.state.rightPane.selectedIndex = i
				break
			}
		}
	}
	uiMainRenderFocusedPane()
	return nil
}

func uiMainGoToFolder(ui *gocui.Gui, v *gocui.View, goToFolderPath string) error {
	if dgState.mainWindow.state.leftPane.focused {
		var cwd string
		if path.IsAbs(goToFolderPath) {
			cwd = path.Clean(goToFolderPath)
		} else {
			cwd = path.Join(dgState.mainWindow.state.leftPane.cwd, goToFolderPath)
		}
		_, err := ioutil.ReadDir(cwd)
		if err != nil {
			uiMainStatusViewMessage(0, "Could not access folder.")
			return err
		}
		dgState.mainWindow.state.leftPane.cwd = cwd
		uiMainListLocalFolder()
	} else if dgState.mainWindow.state.rightPane.focused {
		var cwd string
		if path.IsAbs(goToFolderPath) {
			cwd = path.Clean(goToFolderPath)
		} else {
			cwd = path.Join(dgState.mainWindow.state.rightPane.cwd, goToFolderPath)
		}
		_, err := dgSFTPClient.ReadDir(cwd)
		if err != nil {
			uiMainStatusViewMessage(0, "Could not access folder.")
			return err
		}
		dgState.mainWindow.state.rightPane.cwd = cwd
		go uiMainListArchiveFolder(ui)
	}
	return nil
}

func uiMainCreateFolder(ui *gocui.Gui, v *gocui.View, newFolderName string) error {
	if dgState.mainWindow.state.leftPane.focused {
		newFolder := path.Join(
			dgState.mainWindow.state.leftPane.cwd,
			newFolderName,
		)
		err := os.Mkdir(newFolder, 0700)
		if err != nil {
			uiMainStatusViewMessage(0, "Could not create folder.")
			return err
		}
		uiMainStatusViewMessage(2, strings.Join([]string{
			"Created new folder ", newFolderName, ".",
		}, ""))
		uiMainListLocalFolder()
	} else if dgState.mainWindow.state.rightPane.focused {
		newFolder := path.Join(
			dgState.mainWindow.state.rightPane.cwd,
			newFolderName,
		)
		err := dgSFTPClient.Mkdir(newFolder)
		if err != nil {
			uiMainStatusViewMessage(0, "Could not create folder.")
			return err
		}
		uiMainStatusViewMessage(2, strings.Join([]string{
			"Created new folder ", newFolderName, ".",
		}, ""))
		go uiMainListArchiveFolder(ui)
	}
	return nil
}

func uiMainChmodFile(permissions string, fileName string) error {
	var err error
	var filePath string
	permInt, _ := strconv.ParseInt(permissions, 8, 64)
	if dgState.mainWindow.state.leftPane.focused {
		filePath = path.Join(dgState.mainWindow.state.leftPane.cwd, fileName)
		err = os.Chmod(filePath, os.FileMode(permInt))
	} else if dgState.mainWindow.state.rightPane.focused {
		filePath = path.Join(dgState.mainWindow.state.rightPane.cwd, fileName)
		err = dgSFTPClient.Chmod(filePath, os.FileMode(permInt))
	}
	if err != nil {
		uiMainStatusViewMessage(0, strings.Join([]string{
			"Could not modify permissions for ", fileName, ".",
		}, ""))
		return err
	}
	uiMainStatusViewMessage(2, strings.Join([]string{
		"Modified permissions for ", fileName, ".",
	}, ""))
	return nil
}

func uiMainRenameFile(ui *gocui.Gui, v *gocui.View, givenName string, fileName string) error {
	var oldPath string
	var newPath string
	var err error
	isWholePath, _ := regexp.MatchString(`/`, givenName)
	if dgState.mainWindow.state.leftPane.focused {
		oldPath = path.Join(dgState.mainWindow.state.leftPane.cwd, fileName)
		if isWholePath {
			newPath = givenName
		} else {
			newPath = path.Join(dgState.mainWindow.state.leftPane.cwd, givenName)
		}
		err = os.Rename(oldPath, newPath)
		uiMainListLocalFolder()
	} else if dgState.mainWindow.state.rightPane.focused {
		oldPath = path.Join(dgState.mainWindow.state.rightPane.cwd, fileName)
		if isWholePath {
			newPath = givenName
		} else {
			newPath = path.Join(dgState.mainWindow.state.rightPane.cwd, givenName)
		}
		err = dgSFTPClient.Rename(oldPath, newPath)
		go uiMainListArchiveFolder(ui)
	}
	var messageFileName string
	if isWholePath {
		messageFileName = newPath
	} else {
		messageFileName = path.Base(newPath)
	}
	if err != nil {
		uiMainStatusViewMessage(0, strings.Join([]string{
			"Could not rename ", fileName, " to ", messageFileName, ".",
		}, ""))
		return err
	}
	uiMainStatusViewMessage(2, strings.Join([]string{
		"Renamed ", fileName, " to ", messageFileName, ".",
	}, ""))
	return nil
}

func uiMainDeleteSelected(ui *gocui.Gui, v *gocui.View) error {
	var err error
	var selectedIndex *int
	var selectedFile os.FileInfo
	if dgState.mainWindow.state.leftPane.focused &&
		len(dgState.mainWindow.state.leftPane.folderContents) > 0 {
		selectedIndex = &dgState.mainWindow.state.leftPane.selectedIndex
		selectedFile = dgState.mainWindow.state.leftPane.folderContents[*selectedIndex]
		err = os.Remove(path.Join(
			dgState.mainWindow.state.leftPane.cwd, selectedFile.Name(),
		))
		uiMainListLocalFolder()
	} else if dgState.mainWindow.state.rightPane.focused &&
		len(dgState.mainWindow.state.rightPane.folderContents) > 0 {
		selectedIndex = &dgState.mainWindow.state.rightPane.selectedIndex
		selectedFile = dgState.mainWindow.state.rightPane.folderContents[*selectedIndex]
		err = dgSFTPClient.Remove(path.Join(
			dgState.mainWindow.state.rightPane.cwd, selectedFile.Name(),
		))
		go uiMainListArchiveFolder(ui)
	}
	if err != nil {
		if selectedFile.IsDir() {
			uiMainStatusViewMessage(0, strings.Join([]string{
				"Could not delete ", selectedFile.Name(), ". ",
				"Make sure the folder is empty and try again.",
			}, ""))
		} else {
			uiMainStatusViewMessage(0, strings.Join([]string{
				"Could not delete ", selectedFile.Name(), ".",
			}, ""))
		}
		return err
	}
	uiMainStatusViewMessage(2, strings.Join([]string{
		"Deleted ", selectedFile.Name(), ".",
	}, ""))
	return nil
}

func uiMainDisconnect(ui *gocui.Gui, v *gocui.View) error {
	if !dgState.mainWindow.state.connected && !dgState.mainWindow.state.keyVerification {
		return nil
	}
	if dgState.mainWindow.state.connected {
		dgConfigSetLastFolder(
			dgState.connectWindow.state.serverURI,
			dgState.mainWindow.state.rightPane.cwd,
		)
		dgSFTPDisconnect()
		err := dgConfigSave()
		if err != nil {
			uiMainStatusViewMessage(0, "Could not write to config file.")
		}
	}
	dgStateReset(true)
	uiMainListLocalFolder()
	go uiMainListArchiveFolder(ui)
	uiMainStatusViewMessage(1, "Disconnected.")
	uiMainMenuViewUpdate()
	uiMainMenuStatusViewUpdate()
	return nil
}

func uiMainQuit(ui *gocui.Gui, v *gocui.View) error {
	uiMainDisconnect(ui, v)
	print("\033[H\033[2J")
	os.Exit(0)
	return nil
}
