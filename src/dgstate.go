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

type newfolderstate struct {
	visible bool
}

type newfolderwindow struct {
	view  *gocui.View
	state newfolderstate
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

type dgstate struct {
	mainWindow       mainwindow
	connectWindow    connectwindow
	newFolderWindow  newfolderwindow
	goToWindow       gotowindow
	propertiesWindow propertieswindow
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
				cwd:            ".",
				folderContents: []os.FileInfo{},
				lastFolder:     "",
			},
			rightPane: panestate{
				focused:        false,
				selectedIndex:  0,
				cwd:            ".",
				folderContents: []os.FileInfo{},
				lastFolder:     "",
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
	newFolderWindow: newfolderwindow{
		view: nil,
		state: newfolderstate{
			visible: false,
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
	aboutWindow: aboutwindow{
		view: nil,
		state: aboutstate{
			visible: false,
		},
	},
}

func dgStateReset(visible bool) error {
	dgState.mainWindow.state = dgStatePrototype.mainWindow.state
	dgState.connectWindow.state = dgStatePrototype.connectWindow.state
	dgState.aboutWindow.state = dgStatePrototype.aboutWindow.state
	wd, err := os.Getwd()
	cwdSet := false
	if len(os.Args) > 1 && len(os.Args[1]) > 0 {
		startPath := os.Args[1]
		if path.IsAbs(startPath) {
			_, err := ioutil.ReadDir(path.Clean(startPath))
			if err == nil {
				dgState.mainWindow.state.leftPane.cwd = path.Clean(startPath)
				cwdSet = true
			}
		} else if err == nil {
			_, err := ioutil.ReadDir(path.Join(wd, startPath))
			if err == nil {
				dgState.mainWindow.state.leftPane.cwd = path.Join(wd, startPath)
				cwdSet = true
			}
		}
	}
	if !cwdSet {
		_, err = os.Getwd()
		if err != nil {
			currentUser, err := user.Current()
			if err != nil {
				dgErrorCritical(errors.New("could not read local folder"))
			}
			dgState.mainWindow.state.leftPane.cwd = path.Join(currentUser.HomeDir)
		}
	}
	dgState.mainWindow.state.visible = visible
	dgConfigLoad()
	return nil
}
