package core

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/FranMT-S/Challenge-Go/src/core/bulker"
	"github.com/FranMT-S/Challenge-Go/src/core/parser"
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
	FilePaths  []string
	Parser     parser.IParserMail
	Bulker     bulker.IBulker
	Pagination int
}

func (indexer Indexer) Start() {

	if indexer.Pagination == 0 {
		indexer.Pagination = 1000
	}

	if indexer.Parser == nil {
		panic("Debe inicializar el campo Parse")
	}

	if indexer.Bulker == nil {
		panic("Debe inicializar el campo Bulker")
	}

	if indexer.FilePaths == nil || len(indexer.FilePaths) == 0 {
		panic("El arreglo no tiene datos")
	}

	FilePaths := indexer.FilePaths
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

			indexer.Bulker.SetMails(mails)

			bulkRequest(indexer.Bulker)
			mails = nil
			fmt.Println("---------------------------")
			fmt.Printf("---------Request %v Finalizada--------\n", i+1)
			fmt.Println("---------------------------")
		}
	}
}

func bulkRequest(bulker bulker.IBulker) {

	url := os.Getenv("URL") + bulker.GetCommand()
	bulker.Bulk()

	data := strings.NewReader(bulker.GetData())

	req, err := http.NewRequest("POST", url, data)
	if err != nil {
		log.Fatal(err)
	}

	req.SetBasicAuth(os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	log.Println(resp.StatusCode)
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(body))

}
