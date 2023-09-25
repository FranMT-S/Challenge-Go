package Helpers

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"

	constants_err "github.com/FranMT-S/Challenge-Go/src/constants/errors"
)

// returns the paths of the files to be indexed in the channel specified in the "ch" parameter.
//
// It is recommended to use a buffer channel to read the paths.
//
//	If an error occurs, it returns files=nil and error
func ListAllFilesQuoteChannel(path string, ch chan string) (err error) {
	_, err = os.ReadDir(path)
	if err != nil {
		log.Println(constants_err.ERROR_DIRECTORY_NOT_FOUND + ": " + path)
		close(ch)
		return err
	}

	quoteBasic := NewQueueBasic()
	quoteBasic.Push(path)

	for {
		currentPath := quoteBasic.Poll()
		if currentPath == "" {
			break
		}

		directorys, _ := os.ReadDir(currentPath)
		for _, dir := range directorys {
			newPath := currentPath + "/" + dir.Name()
			if dir.IsDir() {
				quoteBasic.Push(newPath)
			} else {
				ch <- newPath
			}
		}

	}

	close(ch)
	return nil
}

// return a list with the path of the files to be indexed
//
// If an error occurs, it returns files=nil and error
func ListAllFilesQuoteBasic(path string) (files []string, err error) {
	_, err = os.ReadDir(path)
	if err != nil {
		log.Println(constants_err.ERROR_DIRECTORY_NOT_FOUND + ": " + path)
		return nil, err
	}

	quoteBasic := NewQueueBasic()
	quoteBasic.Push(path)

	for {
		currentPath := quoteBasic.Poll()
		if currentPath == "" {
			break
		}

		directorys, _ := os.ReadDir(currentPath)
		for _, dir := range directorys {
			newPath := currentPath + "/" + dir.Name()
			if dir.IsDir() {
				quoteBasic.Push(newPath)
			} else {
				files = append(files, newPath)
			}
		}

	}

	return files, nil
}

// return a list with the path of the files to be indexed
//
// If an error occurs, it returns files=nil and error
func ListAllFilesRecursive(path string) (files []string, err error) {

	dir, err := os.ReadDir(path)

	if err != nil {
		log.Println(constants_err.ERROR_DIRECTORY_NOT_FOUND + ": " + path)
		return nil, err
	}

	for i := 0; i < len(dir); i++ {
		newpath := path + "/" + dir[i].Name()

		if dir[i].IsDir() {
			subFiles, _ := ListAllFilesRecursive(newpath)
			files = append(files, subFiles...)
		} else {
			files = append(files, newpath)
		}

	}

	return files, nil
}

// return a list with the path of the files to be indexed
//
// If an error occurs, it returns files=nil and error
func ListAllFilesIterative(path string) (files []string, err error) {

	folders := []string{path}

	_, err = os.ReadDir(path)
	if err != nil {
		log.Println(constants_err.ERROR_DIRECTORY_NOT_FOUND + ": " + path)
		return nil, err
	}
	i := 0

	for {
		currentPath := folders[i]

		directorys, _ := os.ReadDir(currentPath)
		for _, dir := range directorys {
			newPath := currentPath + "/" + dir.Name()
			if dir.IsDir() {
				folders = append(folders, newPath)
			} else {
				files = append(files, newPath)
			}
		}

		i++

		if i >= len(folders) {
			break
		}
	}

	return files, nil
}

type FolderWalker struct {
	files []string
}

func (w FolderWalker) GetFiles() []string {
	return w.files
}

func (w *FolderWalker) walk(path string, d fs.DirEntry, err error) error {
	if err != nil {
		log.Println(constants_err.ERROR_DIRECTORY_NOT_FOUND + ": " + path)
		return err
	}

	if !d.IsDir() {
		w.files = append(w.files, path)
	}

	return nil
}

func ListAllFilesWalkDir(path string) (files []string) {
	_folderWalker := FolderWalker{}
	filepath.WalkDir(path, _folderWalker.walk)

	return _folderWalker.files
}
