package main

import (
	"bufio"

	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/FranMT-S/Challenge-Go/src/constants"
	"github.com/FranMT-S/Challenge-Go/src/core"
	"github.com/FranMT-S/Challenge-Go/src/core/bulker"
	"github.com/FranMT-S/Challenge-Go/src/core/parser"
)

func main() {
	constants.InitializeVarEnviroment()

	createDirectoryIfNotExist()
	// path := "src/db/maildir/"
	// path := "src/db/maildir/arora-h"
	path := "src/db/maildir/arora-h/all_documents"

	FilePaths := listAllFiles(path)[2:5]

	_Parser := parser.ParserNormal{}
	// bulk := bulker.CreateBulkerV2()
	_Bulker := bulker.CreateBulkerV1()
	// pagination := 0

	indexer := core.Indexer{}
	// indexer := core.Indexer{FilePaths, myParse, bulk, pagination}

	indexer.FilePaths = FilePaths
	indexer.Parser = _Parser
	indexer.Bulker = _Bulker
	// indexer.Start()

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
