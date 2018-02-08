package scp

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// CreateTempFileFromBytes creates temporary file from byte slice and returns path to the file.
func CreateTempFileFromBytes(file []byte, directory, fileName string) (string, error) {
	dirPath, err := ioutil.TempDir("", directory)
	if err != nil {
		return "", err
	}

	localPath, err := createDirStructure(dirPath, fileName)
	if err != nil {
		return "", err
	}

	if _, err := os.Stat(localPath); os.IsNotExist(err) {
		empty := []byte("")
		err := ioutil.WriteFile(localPath, empty, 0755)
		if err != nil {
			return "", err
		}
	}

	f, err := os.OpenFile(localPath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return "", err
	}

	_, err = f.WriteString(string(file))
	if err != nil {
		return "", err
	}
	defer f.Close()

	return localPath, nil
}

// createDirStructure creates directory on the filesystem.
func createDirStructure(dir, file string) (string, error) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.Mkdir(dir, 0777); err != nil {
			return "", err
		}
	}
	return filepath.Join(dir, file), nil
}
