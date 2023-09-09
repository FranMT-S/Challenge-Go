package myDatabase

import (
	"fmt"
	"io"
	"net/http"
	"os"

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
