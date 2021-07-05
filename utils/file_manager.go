package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

const UP_PATH = "up"
const PROCS_PATH = "procs"

type FileManager struct {
	filePath string
}

func (f FileManager) EnsureDirectory() error {
	if _, err := os.Stat(f.getUpDirectory()); os.IsNotExist(err) {
		return err
	}

	if _, err := os.Stat(f.getProcDirectory()); os.IsNotExist(err) {
		return err
	}

	return nil
}

func (f FileManager) getUpDirectory() string {
	return f.getSubDirectory(UP_PATH)
}

func (f FileManager) getProcDirectory() string {
	return f.getSubDirectory(PROCS_PATH)
}

func (f FileManager) GetUpScripts() ([]string, error) {
	directory := f.getUpDirectory()

	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	paths := make([]string, len(files))
	for f, file := range files {
		paths[f] = filepath.Join(directory, file.Name())
	}

	return paths, nil
}

func (f FileManager) GetProcScripts() ([]string, error) {
	directory := f.getProcDirectory()

	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	paths := make([]string, len(files))
	for f, file := range files {
		paths[f] = filepath.Join(directory, file.Name())
	}

	return paths, nil
}

func (f FileManager) getSubDirectory(subdirectory string) string {
	path := fmt.Sprintf("%s/%s", f.filePath, subdirectory)
	p := filepath.FromSlash(path)
	return p
}
