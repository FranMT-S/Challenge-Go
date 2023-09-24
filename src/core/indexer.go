package core

import (
	"fmt"
	"log"
	"os"
	"sync"

	constants_err "github.com/FranMT-S/Challenge-Go/src/constants/errors"
	constants_log "github.com/FranMT-S/Challenge-Go/src/constants/logs"
	"github.com/FranMT-S/Challenge-Go/src/core/bulker"
	"github.com/FranMT-S/Challenge-Go/src/core/parser"
	Helpers "github.com/FranMT-S/Challenge-Go/src/helpers"
	_logs "github.com/FranMT-S/Challenge-Go/src/logs"
	"github.com/FranMT-S/Challenge-Go/src/model"
)

const (
	maxConcurrentAllowed = 20
)

/*
Indexer - Indexa ua lista de archivos en la base de datos
FilePaths rutas de archivos a indexar.
Parse necesario para transformar los archivos de correo a JSON
Bulker hace la peticion para subir el contenido a la base de datos
Pagination ayudara a divir la cantidad de peticiones para reducir la carga.

	Si no se asigna una pagination (pagination = 0), se establece por defecto 1000

El objetivo es pueda intercambiar algoritmos para hacer el bulk a la base de datos
*/
type Indexer struct {
	Parser     parser.IParserMail
	Bulker     bulker.IBulker
	Pagination int
}

func (indexer Indexer) StartFromArray(FilePaths []string) {

	if indexer.Pagination == 0 {
		indexer.Pagination = 1000
	}

	if indexer.Parser == nil {
		panic(constants_err.ERROR_PARSER_UNINITIALIZED)
	}

	if indexer.Bulker == nil {
		panic(constants_err.ERROR_BULKER_UNINITIALIZED)
	}

	if len(FilePaths) == 0 {
		panic(constants_err.ERROR_ARRAY_EMPTY)
	}

	indexer.work(FilePaths)
}

func (indexer Indexer) Start(path string) {

	if indexer.Pagination == 0 {
		indexer.Pagination = 1000
	}

	if indexer.Parser == nil {
		panic(constants_err.ERROR_PARSER_UNINITIALIZED)
	}

	if indexer.Bulker == nil {
		panic(constants_err.ERROR_BULKER_UNINITIALIZED)
	}

	FilePaths, err := Helpers.ListAllFilesQuoteBasic(path)
	if err != nil {
		log.Println("Error:", err)

	} else {
		indexer.work(FilePaths)
	}
}

func (indexer Indexer) work(FilePaths []string) {

	paths := make([]string, len(FilePaths))
	count := (len(FilePaths) / indexer.Pagination)

	// Si hay residuos aumentamos en uno la cuenta para paginar
	if (len(FilePaths) % indexer.Pagination) != 0 {
		count++
	}

	// Ciclo con paginacion
	// Hecho de esta manera porque no deseo mutar el array.
	for i := 0; i < count; i++ {
		start := i * indexer.Pagination
		end := (i + 1) * indexer.Pagination

		// end debe ser menor tamaño del arreglo
		// start debe ser menor a la longitud del arreglo
		// el residuo al dividir entre la paginacion no debe ser 0
		if end > len(FilePaths) && len(FilePaths)%indexer.Pagination != 0 {
			paths = FilePaths[start:]
		} else if start < len(FilePaths) {
			paths = FilePaths[start:end]
		}

		if len(paths) > 0 {

			var mails []*model.Mail

			for j := 0; j < len(paths); j++ {

				file, err := os.Open(paths[j])
				if err != nil {
					_logs.LogSVG(
						constants_log.FILE_NAME_ERROR_INDEXER,
						constants_log.OPERATION_PARSER,
						constants_log.ERROR_OPEN_FILE+": "+paths[j],
						err,
					)
					continue
				}

				fmt.Print("\rParsing: " + paths[j])
				parsedMail, err := indexer.Parser.Parse(file)
				if err != nil {
					_logs.LogSVG(
						constants_log.FILE_NAME_ERROR_INDEXER,
						constants_log.OPERATION_PARSER,
						constants_log.ERROR_PARSER_FAILED+": "+paths[j],
						err,
					)
					continue
				} else {
					mails = append(mails, parsedMail)
				}

				file.Close()
			}

			indexer.Bulker.Bulk(mails)

			mails = nil
			fmt.Println("---------------------------")
			fmt.Printf("---------Request %v Completed--------\n", i+1)
			fmt.Println("---------------------------")
		}
	}
}

func (indexer Indexer) StartAsync(path string, maxConcurrent int) {

	if maxConcurrent < 0 {
		maxConcurrent = 1
	} else if maxConcurrent > maxConcurrentAllowed {
		maxConcurrent = maxConcurrentAllowed
	}

	pathCh := make(chan string)
	mutex := new(sync.Mutex)
	wg := new(sync.WaitGroup)

	for i := 0; i < maxConcurrent; i++ {
		wg.Add(1)
		go indexer.workerAsync(pathCh, mutex, wg, i+1)
	}

	if err := Helpers.ListAllFilesQuoteChannel(path, pathCh); err != nil {
		log.Println(err)
	}

	wg.Wait()
}

func (indexer Indexer) workerAsync(pathCh chan string, mutex *sync.Mutex, wg *sync.WaitGroup, id int) {
	defer wg.Done()

	NumRequest := 0

	var mails []*model.Mail
	for path := range pathCh {
		file, err := os.Open(path)
		if err != nil {
			_logs.LogSVG(
				constants_log.FILE_NAME_ERROR_INDEXER,
				constants_log.OPERATION_PARSER,
				constants_log.ERROR_OPEN_FILE+": "+path,
				err,
			)
			continue
		}

		fmt.Printf("\rWorker %v parsing: %v", id, path)

		parsedMail, err := indexer.Parser.Parse(file)
		if err != nil {
			_logs.LogSVG(
				constants_log.FILE_NAME_ERROR_INDEXER,
				constants_log.OPERATION_PARSER,
				constants_log.ERROR_PARSER_FAILED+": "+path,
				err,
			)
			continue
		} else {
			mails = append(mails, parsedMail)
		}

		file.Close()

		if len(mails) == indexer.Pagination {
			indexer.safeRequest(mails, mutex, id, NumRequest)
			mails = nil
			NumRequest++
		}
	}

	// Si quedaron pendientes
	if len(mails)%indexer.Pagination != 0 {
		indexer.safeRequest(mails, mutex, id, NumRequest)
		mails = nil

	}
}

func (indexer Indexer) safeRequest(mails []*model.Mail, mutex *sync.Mutex, id int, NumRequest int) {
	mutex.Lock()
	indexer.Bulker.Bulk(mails)
	fmt.Println("---------------------------")
	fmt.Printf("--Worker %v, Request %v Completed--------\n", id, NumRequest)
	fmt.Println("---------------------------")
	mutex.Unlock()
}
