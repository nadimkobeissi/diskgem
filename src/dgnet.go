/* SPDX-License-Identifier: MIT
 * Copyright © 2018 Nadim Kobeissi <nadim@symbolic.software>. All Rights Reserved.
 */

package main

import (
	"errors"
	"net"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

var dgSSHClient *ssh.Client
var dgSFTPClient *sftp.Client
var dgSFTPConfirmationChannel = make(chan bool)

func dgSFTPConnect(serverURI string, username string, password string) error {
	var authMethod []ssh.AuthMethod
	if len(password) > 0 {
		authMethod = append(authMethod, ssh.Password(password))
	} else {
		privateKeyFiles := dgFileFindSSHPrivateKeyFiles()
		for _, privateKeyFile := range privateKeyFiles {
			signer, err := ssh.ParsePrivateKey(privateKeyFile)
			if err != nil {
				continue
			}
			authMethod = append(authMethod, ssh.PublicKeys(signer))
		}
	}
	sshConfig := &ssh.ClientConfig{
		User:            username,
		Auth:            authMethod,
		HostKeyCallback: dgSFTPInitializeHostKeyVerification,
	}
	dgSSHClient, err := ssh.Dial("tcp", serverURI, sshConfig)
	if err != nil {
		return err
	}
	dgSFTPClient, err = sftp.NewClient(dgSSHClient)
	return err
}

func dgSFTPInitializeHostKeyVerification(hostname string, remote net.Addr, key ssh.PublicKey) error {
	dgState.mainWindow.state.keyVerification = true
	fp := ssh.FingerprintSHA256(key)
	for _, server := range dgState.mainWindow.state.knownServers {
		if server.Hostname == hostname && server.Fingerprint == fp {
			dgState.mainWindow.state.keyVerification = false
			return nil
		}
	}
	dgState.connectWindow.state.fingerprint = fp
	return errors.New("host key not recognized")
}

func dgSFTPConfirmHostKeyVerification(onConfirm func()) error {
	for i, server := range dgState.mainWindow.state.knownServers {
		if server.Hostname == dgState.connectWindow.state.serverURI {
			server.Username = dgState.connectWindow.state.username
			server.Fingerprint = dgState.connectWindow.state.fingerprint
			dgState.mainWindow.state.knownServers[i] = server
			onConfirm()
			return nil
		}
	}
	dgState.mainWindow.state.knownServers = append(
		dgState.mainWindow.state.knownServers, dgknownserver{
			Hostname:    dgState.connectWindow.state.serverURI,
			Username:    dgState.connectWindow.state.username,
			Fingerprint: dgState.connectWindow.state.fingerprint,
		},
	)
	onConfirm()
	return nil
}

func dgSFTPDisconnect() error {
	dgSFTPClient.Close()
	return nil
}
