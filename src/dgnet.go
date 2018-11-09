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
	sshClient, err := ssh.Dial("tcp", serverURI, sshConfig)
	if err != nil {
		return err
	}
	dgSFTPClient, err = sftp.NewClient(sshClient)
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
