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
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"os/user"
	"path"
)

type dgknownserver struct {
	Hostname    string
	Username    string
	Fingerprint string
	LastFolder  string
}

type dgconfig struct {
	ServerURI    string
	Username     string
	KnownServers []dgknownserver
}

func dgConfigLoad() error {
	parsedConfig := dgconfig{
		ServerURI:    "",
		Username:     "",
		KnownServers: []dgknownserver{},
	}
	currentUser, _ := user.Current()
	configFolderInfo, err := os.Stat(
		path.Join(currentUser.HomeDir, ".config"),
	)
	configFilePath := path.Join(path.Join(path.Join(
		currentUser.HomeDir, ".config"), "diskgem"), "diskgem.cfg",
	)
	if err != nil || !configFolderInfo.IsDir() {
		err = os.Mkdir(path.Join(currentUser.HomeDir, ".config"), 0700)
		if err != nil {
			dgErrorCritical(errors.New("could not create config folder"))
		}
	}
	configFolderInfo, err = os.Stat(
		path.Join(currentUser.HomeDir, path.Join(".config", "diskgem")),
	)
	if err != nil || !configFolderInfo.IsDir() {
		err = os.Mkdir(path.Join(currentUser.HomeDir, path.Join(".config", "diskgem")), 0700)
		if err != nil {
			dgErrorCritical(errors.New("could not create config folder"))
		}
	}
	_, err = os.Stat(configFilePath)
	if err != nil {
		configFile, err := os.Create(configFilePath)
		if err != nil {
			dgErrorCritical(errors.New("could not create config file"))
		}
		_ = configFile.Chmod(0600)
	}
	configFileContents, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		dgErrorCritical(errors.New("could not read config file"))
	}
	if len(configFileContents) == 0 {
		configFileContents, err = json.Marshal(parsedConfig)
		err = ioutil.WriteFile(configFilePath, configFileContents, 0600)
		if err != nil {
			dgErrorCritical(errors.New("could not write to config file"))
		}
	} else {
		err = json.Unmarshal(configFileContents, &parsedConfig)
		if err != nil {
			dgErrorCritical(errors.New("could not read config file"))
		}
	}
	dgState.connectWindow.state.serverURI = parsedConfig.ServerURI
	dgState.connectWindow.state.username = parsedConfig.Username
	dgState.mainWindow.state.knownServers = parsedConfig.KnownServers
	return nil
}

func dgConfigSave() error {
	parsedConfig := dgconfig{
		ServerURI:    dgState.connectWindow.state.serverURI,
		Username:     dgState.connectWindow.state.username,
		KnownServers: dgState.mainWindow.state.knownServers,
	}
	currentUser, _ := user.Current()
	configFileContents, _ := json.MarshalIndent(parsedConfig, "", "")
	configFilePath := path.Join(path.Join(path.Join(
		currentUser.HomeDir, ".config"), "diskgem"), "diskgem.cfg",
	)
	err := ioutil.WriteFile(configFilePath, configFileContents, 0600)
	return err
}

func dgConfigSetLastFolder(serverURI string, lastFolder string) error {
	for i, server := range dgState.mainWindow.state.knownServers {
		if server.Hostname == serverURI {
			server.LastFolder = lastFolder
			dgState.mainWindow.state.knownServers[i] = server
		}
	}
	return nil
}
