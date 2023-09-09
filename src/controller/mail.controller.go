package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/FranMT-S/Challenge-Go/src/constants"
	myDatabase "github.com/FranMT-S/Challenge-Go/src/db"
	"github.com/FranMT-S/Challenge-Go/src/model"
	"github.com/go-chi/chi/v5"
)

const ()

func GetTotalMessage(res http.ResponseWriter, req *http.Request) {

	code := 0

	dbRespBody, responseError := myDatabase.ZincDatabase().GetIndexData()

	if responseError != nil {
		code = http.StatusInternalServerError
		res.WriteHeader(responseError.Status)
		res.Write([]byte(responseError.Error()))
		fmt.Println("Error al decodificar la respuesta JSON:", responseError.Err)
	}

	defer dbRespBody.Close()

	var ResponseIndexData model.ResponseIndexData
	err := json.NewDecoder(dbRespBody).Decode(&ResponseIndexData)

	if err != nil {
		code = http.StatusInternalServerError
		responseError = model.NewResponseError(code, constants.STATUS_ERROR, constants.ERROR_SERVER)

		res.WriteHeader(code)
		res.Write([]byte(responseError.Error()))
		fmt.Println("Error al decodificar la respuesta JSON:", err)
		return
	}

	code = http.StatusOK
	res.WriteHeader(code)
	res.Write([]byte(fmt.Sprintf(`
	{
		"status":%v,
		"msg":"%v",
		"total":%v
	}`, code, constants.STATUS_OK, ResponseIndexData.Stats.DocNum)))
}

func GetMessage(res http.ResponseWriter, req *http.Request) {

	code := 0

	id := chi.URLParam(req, "id")
	fmt.Println(id, code)

	// var ResponseIndexData model.ResponseIndexData
	// err = json.NewDecoder(dbRespBody).Decode(&ResponseIndexData)

	// if err != nil {
	// 	code = http.StatusInternalServerError
	// 	res.WriteHeader(code)
	// 	res.Write([]byte(fmt.Sprintf(`
	// 	{
	// 		"status":%v,
	// 		"msg":"%v",
	// 		"error":"%v"
	// 	}`, code, "ERROR", "Hubo un error en el servidor")))
	// 	fmt.Println("Error al decodificar la respuesta JSON:", err)
	// 	return
	// }

	// code = http.StatusOK
	// res.WriteHeader(code)
	// res.Write([]byte(fmt.Sprintf(`
	// {
	// 	"status":%v,
	// 	"msg":"%v",
	// 	"total":%v
	// }`, code, "OK", ResponseIndexData.Stats.DocNum)))
}
