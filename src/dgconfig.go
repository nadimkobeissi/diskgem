/* SPDX-License-Identifier: MIT
 * Copyright © 2018 Nadim Kobeissi <nadim@nadim.computer>. All Rights Reserved.
 */

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
	_, err = os.OpenFile(configFilePath, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		dgErrorCritical(errors.New("could not create or access config file"))
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
