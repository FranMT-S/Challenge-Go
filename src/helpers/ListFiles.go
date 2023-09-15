package Helpers

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func ListAllFilesRecursive(path string) (files []string) {

	dir, err := os.ReadDir(path)

	if err != nil {
		fmt.Println("no se encontro el directorio: " + path)
		return
	}

	for i := 0; i < len(dir); i++ {
		newpath := path + "/" + dir[i].Name()

		if dir[i].IsDir() {
			subFiles := ListAllFilesRecursive(newpath)
			files = append(files, subFiles...)
		} else {
			files = append(files, newpath)
		}

	}

	return files
}

func ListAllFilesQuoteBasic(path string) (files []string) {

	_, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("no se encontro el directorio: " + path)
		return
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

	return files
}

func ListAllFilesIterative(path string) (files []string) {

	folders := []string{path}

	_, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("no se encontro el directorio: " + path)
		return
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

	return files
}

////////////

type FolderWalker struct {
	files []string
}

func (w FolderWalker) GetFiles() []string {
	return w.files
}

func (w *FolderWalker) walk(path string, d fs.DirEntry, err error) error {
	if err != nil {
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
