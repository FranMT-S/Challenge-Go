package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/FranMT-S/Challenge-Go/src/constants"
	"github.com/FranMT-S/Challenge-Go/src/core"
	"github.com/FranMT-S/Challenge-Go/src/core/bulker"
	"github.com/FranMT-S/Challenge-Go/src/core/parser"
	myDatabase "github.com/FranMT-S/Challenge-Go/src/db"
	Helpers "github.com/FranMT-S/Challenge-Go/src/helpers"
	// "github.com/FranMT-S/Challenge-Go/src/core"
	// "github.com/FranMT-S/Challenge-Go/src/core/bulker"
	// "github.com/FranMT-S/Challenge-Go/src/core/parser"
)

func main() {
	constants.InitializeVarEnviroment()

	// Registra el tiempo de inicio
	startTime := time.Now()
	path := "db/test"
	// // path := "src/db/maildir/allen-p"

	listFiles := ListAllFilesRecursive(path)
	// // listFiles := ListAllFilesQuoteBasic(path)
	// // listFiles := ListAllFilesQueueSafe(path, 5)
	// listFiles := []string{"src/db/maildir/beck-s/aec/2"}
	// // listFiles := []string{"src/db/maildir/arora-h/sent_items/26"}
	// // listFiles := []string{"src/db/maildir/allen-p/straw/7"}

	// // for _, v := range listFiles {
	// // 	fmt.Println(v)
	// // }

	myDatabase.ZincDatabase().CreateIndex()

	// indexer := core.Indexer{listFiles, parser.ParserNormal{}, bulker.CreateBulkerV1(), 5000}
	indexer := core.Indexer{listFiles, parser.NewParserAsyn(50), bulker.CreateBulkerV1(), 5000}
	indexer.Start()

	endTime := time.Now()
	duration := endTime.Sub(startTime)
	seconds := duration.Seconds()

	fmt.Printf("El código se ejecutó en %.2f segundos\n", seconds)
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

func ListAllFilesRecursive(path string) (files []string) {

	dir, err := os.ReadDir(path)

	if err != nil {
		fmt.Println("no se encontro el directorio: " + path)
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

// Necesita arreglos
func ListAllFilesQueueSafe(path string, maxConcurrent int) (files []string) {
	// Crea un semáforo para limitar el número de hilos concurrentes
	semaphore := make(chan struct{}, maxConcurrent)
	var wg sync.WaitGroup
	semaphore <- struct{}{}
	<-semaphore

	_, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("no se encontro el directorio: " + path)
	}

	queueSafe := Helpers.NewQueueSafe()
	queueSafe.Push(path)

	currentPath := queueSafe.Poll()

	var visit func(quote *Helpers.QueueSafe, _path string)

	visit = func(quote *Helpers.QueueSafe, _path string) {
		defer wg.Done()

		if _path != "" {
			directorys, _ := os.ReadDir(_path)

			for _, dir := range directorys {

				newPath := _path + "/" + dir.Name()
				if dir.IsDir() {
					quote.Push(newPath)
				} else {
					files = append(files, newPath)
				}
			}
		}

		for {

			p := quote.Poll()
			if p == "" {
				break
			}
			wg.Add(1)
			go func() {
				visit(Helpers.NewQueueSafe(), p)
			}()

		}

	}

	wg.Add(1) // Agrega una goroutine para el directorio raíz
	visit(queueSafe, currentPath)
	wg.Wait() // Espera a que todas las goroutines hayan terminado

	return files
}

func ListAllFilesQuoteBasic(path string) (files []string) {

	_, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("no se encontro el directorio: " + path)
	}

	quoteBasic := Helpers.NewQueueBasic()
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
