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
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"os/user"
	"path"
	"regexp"
	"sort"
	"strings"
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
		filePath := path.Join(path.Join(currentUser.HomeDir, ".ssh"), file.Name())
		fileContents, err := ioutil.ReadFile(filePath)
		if err != nil {
			continue
		}
		isPrivateKey, _ := regexp.MatchString(keyPattern, string(fileContents))
		if isPrivateKey {
			privateKeyFiles = append(privateKeyFiles, fileContents)
		}
	}
	return privateKeyFiles
}

// dgFileInfoSort sorts a a slice of FileInfo objects such that
// they are organized alphabetically and with folders preceding files.
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

// dgFileUpload throws a file unto the unknowable ether.
func dgFileUpload(
	selectedFile os.FileInfo, selectedFilePath string, archiveFilePath string,
	onStart func(), onProgress func(int), onFinish func(error),
) error {
	selectedFileLstat, err := os.Lstat(selectedFilePath)
	if err != nil {
		onFinish(err)
		return nil
	}
	selectedFileSize := selectedFileLstat.Size()
	selectedFileChunk := int64(256000)
	selectedFileReader, err := os.Open(selectedFilePath)
	if err != nil {
		onFinish(err)
		return nil
	}
	archiveFileWriter, err := sftpClient.Create(archiveFilePath)
	if err != nil {
		onFinish(err)
		return nil
	}
	onStart()
	for c := int64(0); c <= selectedFileSize; c += selectedFileChunk {
		selectedFileReader.Seek(int64(c), 0)
		var buffer []byte
		if c+selectedFileChunk > selectedFileSize {
			buffer = make([]byte, selectedFileSize-c)
		} else {
			buffer = make([]byte, selectedFileChunk)
		}
		selectedFileReader.Read(buffer)
		archiveFileWriter.Write(buffer)
		onProgress(int(math.Ceil(float64(c * 100 / selectedFileSize))))
	}
	// archiveFileWriter.Chmod(0600)
	onFinish(nil)
	return nil
}

func dgFileDownload(
	selectedFile os.FileInfo, selectedFilePath string, localFilePath string,
	onStart func(), onProgress func(int), onFinish func(error),
) error {
	selectedFileLstat, err := sftpClient.Lstat(selectedFilePath)
	if err != nil {
		onFinish(err)
		return nil
	}
	selectedFileSize := selectedFileLstat.Size()
	selectedFileChunk := int64(256000)
	selectedFileReader, err := sftpClient.Open(selectedFilePath)
	if err != nil {
		onFinish(err)
		return nil
	}
	localFileWriter, err := os.Create(localFilePath)
	if err != nil {
		onFinish(err)
		return nil
	}
	onStart()
	for c := int64(0); c <= selectedFileSize; c += selectedFileChunk {
		selectedFileReader.Seek(int64(c), 0)
		var buffer []byte
		if c+selectedFileChunk > selectedFileSize {
			buffer = make([]byte, selectedFileSize-c)
		} else {
			buffer = make([]byte, selectedFileChunk)
		}
		selectedFileReader.Read(buffer)
		localFileWriter.Write(buffer)
		onProgress(int(math.Ceil(float64(c * 100 / selectedFileSize))))
	}
	localFileWriter.Chmod(0600)
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
			lastFolderContents, err = sftpClient.ReadDir(lastFolder)
		}
	} else {
		if dgState.mainWindow.state.leftPane.focused {
			lastFolderContents, err = ioutil.ReadDir(path.Join(dgState.mainWindow.state.leftPane.cwd, lastFolder))
		} else if dgState.mainWindow.state.rightPane.focused {
			lastFolderContents, err = sftpClient.ReadDir(path.Join(dgState.mainWindow.state.rightPane.cwd, lastFolder))
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
