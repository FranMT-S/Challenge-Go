package bulker

import (
	"encoding/json"
	"fmt"
	"os"

	constants_log "github.com/FranMT-S/Challenge-Go/src/constants/logs"
	myDatabase "github.com/FranMT-S/Challenge-Go/src/db"
	_logs "github.com/FranMT-S/Challenge-Go/src/logs"
	model "github.com/FranMT-S/Challenge-Go/src/model"
)

// GetCommand - comando que se usara en el Request
// GetData - Obtener los datos usarse despues de un bulk
// SetMails - Establecer correos
// GetMails - Obtener correos
// Bulk - Transforma los correos en el formato que usare el Request al subir los datos.
type IBulker interface {
	GetCommand() string
	Bulk(mails []*model.Mail)
}

/*
-----------------------------------
Section Bulker V1
-----------------------------------
*/

type BulkerV1 struct {
}

func (bulk BulkerV1) GetCommand() string {
	return "_bulk"
}

func (bulk BulkerV1) Bulk(mails []*model.Mail) {
	index := fmt.Sprintf(`{ "index" : { "_index" : "%v" } }  `, os.Getenv("INDEX"))
	json := ""

	for i := 0; i < len(mails); i++ {
		json += index + "\n"
		json += mails[i].String() + "\n"
	}

	if err := myDatabase.BulkRequest(bulk.GetCommand(), json); err != nil {
		_logs.LogSVG(
			constants_log.FILE_NAME_ERROR_INDEXER,
			constants_log.OPERATION_BULKER,
			constants_log.ERROR_BULKER_FAILED,
			err,
		)
	}
}

/*
Tiene un funcion con pointer receiver por lo que es necesario
retornar un puntero para que la interface IBulker lo acepte
*/
func CreateBulkerV1() BulkerV1 {
	return BulkerV1{}
}

/*
-----------------------------------
End Bulker V1
-----------------------------------
*/

/*
-----------------------------------
Section Bulker V2
-----------------------------------
*/

// Formato para el Request Bulker V2
type bulkResponse struct {
	Index   string
	Records []*model.Mail
}

type BulkerV2 struct {
}

func (bulk BulkerV2) GetCommand() string {
	return "_bulkv2"
}

func (bulk BulkerV2) Bulk(mails []*model.Mail) {
	bulkResponse := bulkResponse{
		Index:   os.Getenv("INDEX"),
		Records: mails}

	json, err := json.Marshal(bulkResponse)
	if err != nil {
		_logs.LogSVG(
			constants_log.FILE_NAME_ERROR_INDEXER,
			constants_log.OPERATION_BULKER,
			constants_log.ERROR_JSON_PARSE,
			err,
		)
		return
	}

	if err := myDatabase.BulkRequest(bulk.GetCommand(), string(json)); err != nil {
		_logs.LogSVG(
			constants_log.FILE_NAME_ERROR_INDEXER,
			constants_log.OPERATION_BULKER,
			constants_log.ERROR_BULKER_FAILED,
			err,
		)
	}

}

func CreateBulkerV2() BulkerV2 {
	return BulkerV2{}
}
