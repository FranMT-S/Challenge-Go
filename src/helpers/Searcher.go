package helpers

import "os"

func listAllFilesRecursive(path string) (files []string) {

	dir, err := os.ReadDir(path)

	if err != nil {
		fmt.Println("no se encontro el directorio: " + path)

	for i := 0; i < len(dir); i++ {
		newpath := path + "/" + dir[i].Name()

		if dir[i].IsDir() {
			subFiles := listAllFiles(newpath)
			files = append(files, subFiles...)
		} else {
			files = append(files, newpath)
		}

	}

	return files
}
