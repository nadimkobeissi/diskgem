/* SPDX-License-Identifier: MIT
 * Copyright © 2018 Nadim Kobeissi <nadim@nadim.computer>. All Rights Reserved.
 */

package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/machinebox/progress"
)

func dgFileFindSSHPrivateKeyFiles() [][]byte {
	currentUser, _ := user.Current()
	keyPattern := "^-----BEGIN \\w{1,16} PRIVATE KEY-----"
	sshFolder, err := ioutil.ReadDir(path.Join(currentUser.HomeDir, ".ssh"))
	if err != nil {
		return [][]byte{}
	}
	privateKeyFiles := [][]byte{}
	for _, file := range sshFolder {
		var filePartialContents = make([]byte, 44)
		filePath := path.Join(path.Join(currentUser.HomeDir, ".ssh"), file.Name())
		fileReader, err := os.Open(filePath)
		_, err = fileReader.Read(filePartialContents)
		if err != nil {
			continue
		}
		isPrivateKey, _ := regexp.MatchString(keyPattern, string(filePartialContents))
		if isPrivateKey {
			fileContents, err := ioutil.ReadFile(filePath)
			if err == nil {
				privateKeyFiles = append(privateKeyFiles, fileContents)
			}
		}
	}
	return privateKeyFiles
}

func dgFileInfoSort(files []os.FileInfo) []os.FileInfo {
	sortedFolders := []os.FileInfo{}
	sortedHiddenFolders := []os.FileInfo{}
	sortedFiles := []os.FileInfo{}
	sortedHiddenFiles := []os.FileInfo{}
	sort.Slice(files, func(i int, j int) bool {
		return files[i].Name() < files[j].Name()
	})
	for _, file := range files {
		isHidden, _ := regexp.MatchString("^\\.", file.Name())
		if isHidden {
			if file.IsDir() || dgFileIsSymlink(file) {
				sortedHiddenFolders = append(sortedHiddenFolders, file)
			} else {
				sortedHiddenFiles = append(sortedHiddenFiles, file)
			}
		} else {
			if file.IsDir() || dgFileIsSymlink(file) {
				sortedFolders = append(sortedFolders, file)
			} else {
				sortedFiles = append(sortedFiles, file)
			}
		}
	}
	return append(sortedFolders,
		append(sortedHiddenFolders,
			append(sortedFiles, sortedHiddenFiles...)...,
		)...,
	)
}

func dgFileUpload(
	selectedFile os.FileInfo, selectedFilePath string, archiveFilePath string,
	onStart func(), onProgress func(int, string), onFinish func(error),
) error {
	_, err := os.Lstat(selectedFilePath)
	if err != nil {
		onFinish(err)
		return nil
	}
	selectedFileDescriptor, err := os.Open(selectedFilePath)
	selectedFileReader := progress.NewReader(selectedFileDescriptor)
	if err != nil {
		onFinish(err)
		return nil
	}
	archiveFileWriter, err := dgSFTPClient.Create(archiveFilePath)
	if err != nil {
		onFinish(err)
		return nil
	}
	go func() {
		ctx := context.Background()
		progressChan := progress.NewTicker(
			ctx, selectedFileReader, selectedFile.Size(), 500*time.Millisecond,
		)
		for p := range progressChan {
			onProgress(
				int(p.Percent()),
				dgFileDurationFormat(p.Remaining()),
			)
		}
	}()
	onStart()
	archiveFileWriter.ReadFrom(selectedFileReader)
	onFinish(nil)
	return nil
}

func dgFileDownload(
	selectedFile os.FileInfo, selectedFilePath string, localFilePath string,
	onStart func(), onProgress func(int, string), onFinish func(error),
) error {
	_, err := dgSFTPClient.Lstat(selectedFilePath)
	if err != nil {
		onFinish(err)
		return nil
	}
	selectedFileReader, err := dgSFTPClient.Open(selectedFilePath)
	if err != nil {
		onFinish(err)
		return nil
	}
	localFileDescriptor, err := os.OpenFile(localFilePath, os.O_RDWR|os.O_CREATE, 0600)
	localFileWriter := progress.NewWriter(localFileDescriptor)
	if err != nil {
		onFinish(err)
		return nil
	}
	go func() {
		ctx := context.Background()
		progressChan := progress.NewTicker(
			ctx, localFileWriter, selectedFile.Size(), 500*time.Millisecond,
		)
		for p := range progressChan {
			onProgress(
				int(p.Percent()),
				dgFileDurationFormat(p.Remaining()),
			)
		}
	}()
	onStart()
	selectedFileReader.WriteTo(localFileWriter)
	onFinish(nil)
	return nil
}

func dgFileIsSymlink(file os.FileInfo) bool {
	return file.Mode()&os.ModeSymlink == os.ModeSymlink
}

func dgFileFolderPathAutocomplete(input string, index int) (string, int) {
	if input == "" {
		return input, index
	}
	var lastFolderContents []os.FileInfo
	var matches []string
	var err error
	lastFolder := path.Dir(input)
	danglingPath := path.Base(input)
	if path.IsAbs(input) {
		if dgState.mainWindow.state.leftPane.focused {
			lastFolderContents, err = ioutil.ReadDir(lastFolder)
		} else if dgState.mainWindow.state.rightPane.focused {
			lastFolderContents, err = dgSFTPClient.ReadDir(lastFolder)
		}
	} else {
		if dgState.mainWindow.state.leftPane.focused {
			lastFolderContents, err = ioutil.ReadDir(
				path.Join(dgState.mainWindow.state.leftPane.cwd, lastFolder),
			)
		} else if dgState.mainWindow.state.rightPane.focused {
			lastFolderContents, err = dgSFTPClient.ReadDir(
				path.Join(dgState.mainWindow.state.rightPane.cwd, lastFolder),
			)
		}
	}
	if err != nil {
		return input, index
	}
	for _, v := range lastFolderContents {
		if strings.HasPrefix(v.Name(), danglingPath) {
			matches = append(matches, v.Name())
		}
	}
	if len(matches) == 0 {
		return input, 0
	}
	if index >= len(matches) {
		index = 0
	}
	return path.Join(lastFolder, matches[index]), index
}

func dgFileSizeFormat(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f%cB", float64(b)/float64(div), "kMGTPE"[exp])
}

func dgFileDurationFormat(d time.Duration) string {
	return strings.Join([]string{
		strconv.Itoa(int(d.Minutes())), "min",
	}, "")
}
