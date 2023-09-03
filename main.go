package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

func main() {
	createDirectoryIfNotExist()

	// cpu, err := os.Create("cpu.prof")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// pprof.StartCPUProfile(cpu)
	// defer pprof.StopCPUProfile()

	// path := "src/db/maildir"

	// directorys, _ := listAll(path)

	// for _, v := range directorys {
	// 	fmt.Println(v)
	// }
	// mail := core.Mail{"", "3"}
	// mail2 := new(core.Mail)
	// mail2.Name = "Marta"
	// pBytes, _ := json.Marshal(mail)
	// pBytes2, _ := json.Marshal(mail2)

	// fmt.Println(string(pBytes))
	// fmt.Println(string(pBytes2))

	// f, err := myWriter(true)

	// if err != nil {
	// 	log.Println(err.Error())
	// }
	// defer f.Close()

	file, err := os.Open("src/db/maildir/arora-h/deleted_items/34")

	if err != nil {
		log.Fatal(err)
	}

	ReadLineForLineScanner(file)

}

func ReadLineForLineScanner(file *os.File) {

	defer file.Close()
	scanner := bufio.NewScanner(file)

	// optionally, resize scanner's capacity for lines over 64K, see next example
	const maxCapacity int = 1024 // your required line length
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	json, _ := myWriter(true, false)
	defer json.Close()

	for scanner.Scan() {
		json.Write(scanner.Bytes())
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func ReadLineForLine(file *os.File) {

	defer file.Close()
	buf := make([]byte, 1024)

	json, _ := myWriter(true, false)
	defer json.Close()

	for {
		n, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
			continue
		}
		if n > 0 {
			json.Write(buf[:n])
			fmt.Println(string(buf[:n]))
		}
	}
}

func myWriter(modeAppend bool, readWrite bool) (*os.File, error) {

	name := "json/data" + time.Now().Format("02_150405") + ".json"

	modes := os.O_CREATE
	if readWrite {
		modes |= os.O_RDWR
	} else {
		modes |= os.O_WRONLY
	}

	if modeAppend {
		modes |= os.O_APPEND
	} else {
		modes |= os.O_TRUNC
	}

	return os.OpenFile(name, modes, 0644)

}

func createDirectoryIfNotExist() {
	if _, err := os.Stat("json"); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir("json", os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}
}

func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), err
}

func listAll(path string) ([]string, []string) {
	var directorys, files []string

	dir, err := os.ReadDir(path)

	if err != nil {
		fmt.Println("no se encontro el directorio")

	}

	for i := 0; i < len(dir); i++ {
		newpath := path + "/" + dir[i].Name()
		dir, err := isDirectory(newpath)
		if err != nil {
			fmt.Println("No se encontro el directorio:" + newpath)
		}
		if dir {
			subDir, subFiles := listAll(newpath)
			directorys = append(directorys, newpath)
			directorys = append(directorys, subDir...)
			files = append(files, subFiles...)
		} else {
			files = append(files, newpath)
		}
	}

	return directorys, files
}
