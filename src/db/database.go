package myDatabase

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/FranMT-S/Challenge-Go/src/constants"
	myMiddleware "github.com/FranMT-S/Challenge-Go/src/middleware"
	"github.com/FranMT-S/Challenge-Go/src/model"
)

var z_database *zincDatabase

type zincDatabase struct {
	client *http.Client
}

func ZincDatabase() *zincDatabase {
	if z_database == nil {
		z_database = &zincDatabase{client: &http.Client{}}
	}

	return z_database
}

func (db zincDatabase) GetIndexData() (io.ReadCloser, *model.ResponseError) {
	url := os.Getenv("URL") + "index/mailsTest3"

	dbReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error al crear la solicitud:", err)
		return nil, model.NewResponseError(http.StatusBadRequest, constants.STATUS_ERROR, constants.ERROR_CREATE_REQUEST)
	}

	myMiddleware.ZincHeader(dbReq)

	// Realizar la solicitud
	dbResp, err := db.client.Do(dbReq)
	if err != nil {
		fmt.Println("Error al realizar la solicitud:", err)
		return nil, model.NewResponseError(http.StatusBadRequest, constants.STATUS_ERROR, constants.ERROR_REQUEST)
	}

	// Verificar el código de estado de la respuesta
	if dbResp.StatusCode != http.StatusOK {
		fmt.Println("Respuesta no exitosa. Código de estado:", dbResp.Status)
		return nil, model.NewResponseError(dbResp.StatusCode, constants.STATUS_ERROR, constants.ERROR_REQUEST)
	}

	return dbResp.Body, nil
}

func (db zincDatabase) CreateIndex() {
	index := `{
		"name": "Test",
		"storage_type": "disk",
		"mappings": {
		"properties": {
		  "Date": {
			"type": "date",
			"format":"2006-01-02T15:04:05Z",
			"index": true,
			"store": false,
			"sortable": true,
			"aggregatable": true,
			"highlightable": false
		  },
		  "Bcc": {
			"type": "text",
			"index": true,
			"store": false,
			"sortable": true,
			"aggregatable": true,
			"highlightable": false
		  },
		  "Cc": {
			"type": "text",
			"index": true,
			"store": false,
			"sortable": true,
			"aggregatable": true,
			"highlightable": false
		  },
		  "Content": {
			"type": "text",
			"index": true,
			"store": false,
			"sortable": true,
			"aggregatable": true,
			"highlightable": false
		  },
		  "Content_Transfer_Encoding": {
			"type": "text",
			"index": true,
			"store": false,
			"sortable": true,
			"aggregatable": true,
			"highlightable": false
		  },
		  "Content_Type": {
			"type": "text",
			"index": true,
			"store": false,
			"sortable": true,
			"aggregatable": true,
			"highlightable": false
		  },
		  "From": {
			"type": "text",
			"index": true,
			"store": false,
			"sortable": true,
			"aggregatable": true,
			"highlightable": false
		  },
		  "Message_ID": {
			"type": "text",
			"index": true,
			"store": false,
			"sortable": true,
			"aggregatable": true,
			"highlightable": false
		  },
		  "Mime_Version": {
			"type": "text",
			"index": true,
			"store": false,
			"sortable": true,
			"aggregatable": true,
			"highlightable": false
		  },
		  "Subject": {
			"type": "text",
			"index": true,
			"store": false,
			"sortable": true,
			"aggregatable": true,
			"highlightable": false
		  },
		  "To": {
			"type": "text",
			"index": true,
			"store": false,
			"sortable": true,
			"aggregatable": true,
			"highlightable": false
		  },
		  "X_FileName": {
			"type": "text",
			"index": true,
			"store": false,
			"sortable": true,
			"aggregatable": true,
			"highlightable": false
		  },
		  "X_Folder": {
			"type": "text",
			"index": true,
			"store": false,
			"sortable": true,
			"aggregatable": true,
			"highlightable": false
		  },
		  "X_From": {
			"type": "text",
			"index": true,
			"store": false,
			"sortable": true,
			"aggregatable": true,
			"highlightable": false
		  },
		  "X_Origin": {
			"type": "text",
			"index": true,
			"store": false,
			"sortable": true,
			"aggregatable": true,
			"highlightable": false
		  },
		  "X_To": {
			"type": "text",
			"index": true,
			"store": false,
			"sortable": true,
			"aggregatable": true,
			"highlightable": false
		  },
		  "X_bcc": {
			"type": "text",
			"index": true,
			"store": false,
			"sortable": true,
			"aggregatable": true,
			"highlightable": false
		  },
		  "X_cc": {
			"type": "text",
			"index": true,
			"store": false,
			"sortable": true,
			"aggregatable": true,
			"highlightable": false
		  },
		  "_id": {
			"type": "keyword",
			"index": true,
			"store": false,
			"sortable": true,
			"aggregatable": true,
			"highlightable": false
		  }
		}
		}
	}`

	url := os.Getenv("URL") + "index"

	data := strings.NewReader(index)

	req, err := http.NewRequest("POST", url, data)
	if err != nil {
		fmt.Println(`entro en error 1`)
		log.Fatal(err)
	}

	myMiddleware.ZincHeader(req)

	resp, err := db.client.Do(req)

	if err != nil {

		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {

		log.Fatal(err)
	}

	if string(body) != `{"error":"index [Test] already exists"}` {
		fmt.Println(string(body))
	}
}

func BulkRequest(command, mailsData string) {

	url := os.Getenv("URL") + command

	data := strings.NewReader(mailsData)

	req, err := http.NewRequest("POST", url, data)
	if err != nil {
		log.Fatal(err)
	}

	myMiddleware.ZincHeader(req)

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
