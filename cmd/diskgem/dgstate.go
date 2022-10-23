/* SPDX-License-Identifier: MIT
 * Copyright © 2018-2019 Nadim Kobeissi <nadim@symbolic.software>. All Rights Reserved.
 */

package main

import (
	"errors"
	"io/ioutil"
	"os"
	"os/user"
	"path"

	"github.com/jroimartin/gocui"
)

type dgtransfer struct {
	name      string
	size      int64
	direction string
	progress  int
	remaining string
	fromPath  string
	toPath    string
	finished  bool
}

type panestate struct {
	focused        bool
	selectedIndex  int
	cwd            string
	folderContents []os.FileInfo
	lastFolder     string
	confirmPath    string
}

type mainstate struct {
	visible         bool
	connected       bool
	keyVerification bool
	version         string
	knownServers    []dgknownserver
	transfers       []*dgtransfer
	leftPane        panestate
	rightPane       panestate
}

type mainwindow struct {
	menuView       *gocui.View
	menuStatusView *gocui.View
	leftPaneView   *gocui.View
	rightPaneView  *gocui.View
	statusView     *gocui.View
	transfersView  *gocui.View
	state          mainstate
}

type connectstate struct {
	visible     bool
	serverURI   string
	username    string
	password    string
	fingerprint string
	selected    string
}

type connectwindow struct {
	view             *gocui.View
	serverURIView    *gocui.View
	usernameView     *gocui.View
	passwordView     *gocui.View
	passwordInfoView *gocui.View
	okButtonView     *gocui.View
	state            connectstate
}

type aboutstate struct {
	visible bool
}

type aboutwindow struct {
	view  *gocui.View
	state aboutstate
}

type gotostate struct {
	visible     bool
	lastPath    string
	lastInitial string
	index       int
}

type gotowindow struct {
	view  *gocui.View
	state gotostate
}

type newfolderstate struct {
	visible bool
}

type newfolderwindow struct {
	view  *gocui.View
	state newfolderstate
}

type propertiesstate struct {
	visible     bool
	fileName    string
	permissions string
	selected    string
}

type propertieswindow struct {
	view            *gocui.View
	nameView        *gocui.View
	sizeView        *gocui.View
	modDateView     *gocui.View
	permissionsView *gocui.View
	state           propertiesstate
}

type shellstate struct {
	visible bool
}

type shellwindow struct {
	view  *gocui.View
	state shellstate
}

type dgstate struct {
	mainWindow       mainwindow
	connectWindow    connectwindow
	goToWindow       gotowindow
	newFolderWindow  newfolderwindow
	propertiesWindow propertieswindow
	shellWindow      shellwindow
	aboutWindow      aboutwindow
}

var dgState dgstate

var dgStatePrototype = dgstate{
	mainWindow: mainwindow{
		menuView:       nil,
		menuStatusView: nil,
		leftPaneView:   nil,
		rightPaneView:  nil,
		statusView:     nil,
		state: mainstate{
			visible:         false,
			connected:       false,
			keyVerification: false,
			version:         dgVersion,
			knownServers:    []dgknownserver{},
			transfers:       []*dgtransfer{},
			leftPane: panestate{
				focused:        true,
				selectedIndex:  0,
				cwd:            "",
				folderContents: []os.FileInfo{},
				lastFolder:     "",
				confirmPath:    "",
			},
			rightPane: panestate{
				focused:        false,
				selectedIndex:  0,
				cwd:            "",
				folderContents: []os.FileInfo{},
				lastFolder:     "",
				confirmPath:    "",
			},
		},
	},
	connectWindow: connectwindow{
		view:             nil,
		serverURIView:    nil,
		usernameView:     nil,
		passwordView:     nil,
		passwordInfoView: nil,
		okButtonView:     nil,
		state: connectstate{
			visible:     false,
			serverURI:   "",
			username:    "",
			password:    "",
			fingerprint: "",
			selected:    "connectServerURI",
		},
	},
	goToWindow: gotowindow{
		view: nil,
		state: gotostate{
			visible:     false,
			lastPath:    "",
			lastInitial: "",
			index:       0,
		},
	},
	newFolderWindow: newfolderwindow{
		view: nil,
		state: newfolderstate{
			visible: false,
		},
	},
	propertiesWindow: propertieswindow{
		view:            nil,
		nameView:        nil,
		sizeView:        nil,
		modDateView:     nil,
		permissionsView: nil,
		state: propertiesstate{
			visible:     false,
			fileName:    "",
			permissions: "",
			selected:    "propertiesName",
		},
	},
	shellWindow: shellwindow{
		view: nil,
		state: shellstate{
			visible: false,
		},
	},
	aboutWindow: aboutwindow{
		view: nil,
		state: aboutstate{
			visible: false,
		},
	},
}

func dgStateReset(visible bool) error {
	lastFolderLocal := dgState.mainWindow.state.leftPane.cwd
	dgState.mainWindow.state = dgStatePrototype.mainWindow.state
	dgState.mainWindow.state.leftPane.cwd = dgStateInitLeftPaneCwd(lastFolderLocal)
	dgState.connectWindow.state = dgStatePrototype.connectWindow.state
	dgState.aboutWindow.state = dgStatePrototype.aboutWindow.state
	dgState.mainWindow.state.visible = visible
	dgConfigLoad()
	return nil
}

func dgStateInitLeftPaneCwd(lastFolderLocal string) string {
	wd, wdErr := os.Getwd()
	if len(lastFolderLocal) > 0 {
		return lastFolderLocal
	}
	if len(os.Args) > 1 && len(os.Args[1]) > 0 {
		startPath := os.Args[1]
		if path.IsAbs(startPath) {
			_, err := ioutil.ReadDir(path.Clean(startPath))
			if err == nil {
				return path.Clean(startPath)
			}
		} else if wdErr == nil {
			_, err := ioutil.ReadDir(path.Join(wd, startPath))
			if err == nil {
				return path.Join(wd, startPath)
			}
		}
	}
	if wdErr == nil {
		return wd
	}
	currentUser, err := user.Current()
	if err == nil {
		return path.Join(currentUser.HomeDir)
	}
	dgErrorCritical(errors.New("could not set any working directory"))
	return lastFolderLocal
}
