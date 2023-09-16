package core

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/FranMT-S/Challenge-Go/src/core/bulker"
	"github.com/FranMT-S/Challenge-Go/src/core/parser"
	Helpers "github.com/FranMT-S/Challenge-Go/src/helpers"
	"github.com/FranMT-S/Challenge-Go/src/model"
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
		panic("Debe inicializar el campo Parse")
	}

	if indexer.Bulker == nil {
		panic("Debe inicializar el campo Bulker")
	}

	if len(FilePaths) == 0 {
		panic("El arreglo no tiene datos")
	}
	indexer.work(FilePaths)
}

func (indexer Indexer) Start(path string) {

	if indexer.Pagination == 0 {
		indexer.Pagination = 1000
	}

	if indexer.Parser == nil {
		panic("Debe inicializar el campo Parse")
	}

	if indexer.Bulker == nil {
		panic("Debe inicializar el campo Bulker")
	}

	FilePaths := Helpers.ListAllFilesQuoteBasic(path)
	indexer.work(FilePaths)
}

func (indexer Indexer) work(FilePaths []string) {

	part := make([]string, len(FilePaths))
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
			part = FilePaths[start:]
		} else if start < len(FilePaths) {
			part = FilePaths[start:end]
		}

		if len(part) > 0 {

			var mails []model.Mail

			for j := 0; j < len(part); j++ {

				file, err := os.Open(part[j])
				if err != nil {
					log.Fatal(err)
				}

				fmt.Println("Parseando: " + part[j])

				mails = append(mails, indexer.Parser.Parse(file))
				file.Close()
			}

			indexer.Bulker.Bulk(mails)

			mails = nil
			fmt.Println("---------------------------")
			fmt.Printf("---------Request %v Finalizada--------\n", i+1)
			fmt.Println("---------------------------")
		}
	}
}

func (indexer Indexer) StartAsync(path string, maxConcurrent int) {

	if maxConcurrent < 0 {
		maxConcurrent = 1
	} else if maxConcurrent > 10 {
		maxConcurrent = 10
	}

	pathCh := make(chan string)
	mutex := new(sync.Mutex)
	wg := new(sync.WaitGroup)

	for i := 0; i < maxConcurrent; i++ {
		wg.Add(1)
		go indexer.workAsync(pathCh, mutex, wg, i+1)
	}

	Helpers.ListAllFilesQuoteChannel(path, pathCh)
	wg.Wait()
}

func (indexer Indexer) workAsync(pathCh chan string, mutex *sync.Mutex, wg *sync.WaitGroup, id int) {
	defer wg.Done()

	NumRequest := 0
	// part := make([]string, indexer.Pagination)
	var mailList []model.Mail
	for path := range pathCh {
		file, err := os.Open(path)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Worker %v parseando: %v\n", id, path)

		mailList = append(mailList, indexer.Parser.Parse(file))
		file.Close()

		if len(mailList) == indexer.Pagination {
			mutex.Lock()
			indexer.Bulker.Bulk(mailList)
			fmt.Println("---------------------------")
			fmt.Printf("--Worker %v, Request %v Finalizada--------\n", id, NumRequest)
			fmt.Println("---------------------------")
			mailList = nil
			NumRequest++
			mutex.Unlock()
		}
	}

	// Si quedaron pendientes
	if len(mailList)%indexer.Pagination != 0 {
		mutex.Lock()
		indexer.Bulker.Bulk(mailList)
		mutex.Unlock()
		mailList = nil
	}
}
