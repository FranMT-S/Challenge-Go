package main

import (
	"bufio"
	"challenge/src/core/bulker"
	"challenge/src/core/parser"
	"challenge/src/model"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

type TestData struct {
	Name string
	Data []string
	Arr  []TestData
}

func main() {
	createDirectoryIfNotExist()
	pagination := 5000

	path := "src/db/maildir"

	// directorys, _ := listAll(path)

	// path := "src/db/maildir/arora-h"
	directorys := listAllFiles(path)

	// for _, v := range directorys {
	// 	fmt.Println(v)
	// }

	// var mails []*model.Mail
	// i := 0

	// for i := 0; i < len(directorys); i++ {

	// 	fmt.Println(directorys[i])

	// 	// file, err := os.Open(directorys[i])
	// 	// // file, err := os.Open("src/db/maildir/arora-h/inbox/1")

	// 	// if err != nil {
	// 	// 	log.Fatal(err)
	// 	// }
	// 	// myParse := new(parser.ParserNormal)
	// 	// mails = append(mails, myParse.Parse(file))
	// 	// fmt.Println("ruta parseada: " + directorys[i])
	// 	// file.Close()

	// }

	part := make([]string, len(directorys))
	myParse := new(parser.ParserNormal)
	var b bulker.Bulker

	bulk := &bulker.BulkerV1{b}

	// 	// bulk := &bulker.BulkerV2{b}

	// 	// bulk.Bulk()
	// 	// fmt.Println(bulk.GetData())
	//    bulker.RequestBulk(bulk)
	// 	// fmt.Println("llegue aqui")

	// Ciclo con paginacion
	for i := 0; i <= (len(directorys) / pagination); i++ {
		start := i * pagination
		end := (i + 1) * pagination

		if end > len(directorys) && len(directorys)%pagination != 0 {
			part = directorys[start:]
		} else if start < len(directorys) {
			part = directorys[start:end]
		}

		if len(part) > 0 {

			var mails []*model.Mail

			for j := 0; j < len(part); j++ {

				file, err := os.Open(part[j])
				if err != nil {
					log.Fatal(err)
				}

				mails = append(mails, myParse.Parse(file))
				fmt.Println("Parseando: " + part[j])
				file.Close()
			}

			bulk.Mails = mails
			bulker.RequestBulk(bulk)
			mails = nil
			fmt.Println("---------------------------")
			fmt.Printf("---------Request %v--------", i+1)
			fmt.Println("---------------------------")
		}

	}

	// 	var b bulker.Bulker
	// 	b.Mails = mails

	// 	bulk := &bulker.BulkerV1{b}
	// 	// bulk := &bulker.BulkerV2{b}

	// 	// bulk.Bulk()
	// 	// fmt.Println(bulk.GetData())
	// 	bulker.RequestBulk(bulk)
	// 	// fmt.Println("llegue aqui")
}

func ReadLineForLineBufio(file *os.File) {

	defer file.Close()
	reader := bufio.NewReader(file)

	json, _ := myWriter(true, false)
	defer json.Close()

	for {
		lines, err := reader.ReadBytes('\n')

		if err != nil {

			if err == io.EOF {
				log.Println(err)
			}

			break
		}

		json.Write(lines)
		fmt.Println(string(lines))

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
			log.Println(err)
			continue
		}
		if n > 0 {
			json.Write(buf[:n])
			fmt.Println(string(buf[:n]))
		}
	}
}

func ReadLineForLineScanner(file *os.File) {

	defer file.Close()

	json, _ := myWriter(true, false)
	defer json.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		json.Write(scanner.Bytes())
		fmt.Println(scanner.Text())
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

func listAllFiles(path string) (files []string) {

	dir, err := os.ReadDir(path)

	if err != nil {
		fmt.Println("no se encontro el directorio: " + path)
	}

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

func listAll(path string) ([]string, []string) {
	var directorys, files []string

	dir, err := os.ReadDir(path)

	if err != nil {
		fmt.Println("no se encontro el directorio:" + path)
	}

	for i := 0; i < len(dir); i++ {
		newpath := path + "/" + dir[i].Name()

		if dir[i].IsDir() {
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
