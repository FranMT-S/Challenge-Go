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

/*
Provides an interface responsible for uploading data to the database.

  - getCommand: returns the string with the address of the end point that will be consumed from [Zincsearch].

  - Bulk: perform a [Bulk] operation to upload the data to [Zincsearch].

[Zincsearch]: https://zincsearch-docs.zinc.dev/
[Bulk]: https://zincsearch-docs.zinc.dev/api/document/bulk/
*/
type IBulker interface {
	getCommand() string
	Bulk(mails []*model.Mail)
}

/*
Provides an methods responsible for uploading data to the database.

This object is based on the [Bulk] endpoint provided by [Zincsearch].

[Bulk]: https://zincsearch-docs.zinc.dev/api/document/bulk/
[Zincsearch]: https://zincsearch-docs.zinc.dev/
*/
type BulkerV1 struct {
}

func (bulk BulkerV1) getCommand() string {
	return "_bulk"
}

// Upload the information to the database
func (bulk BulkerV1) Bulk(mails []*model.Mail) {
	index := fmt.Sprintf(`{ "index" : { "_index" : "%v" } }  `, os.Getenv("INDEX"))
	json := ""

	for i := 0; i < len(mails); i++ {
		json += index + "\n"
		json += mails[i].String() + "\n"
	}

	if err := myDatabase.BulkRequest(bulk.getCommand(), json); err != nil {
		_logs.LogSVG(
			constants_log.FILE_NAME_ERROR_INDEXER,
			constants_log.OPERATION_BULKER,
			constants_log.ERROR_BULKER_FAILED,
			err,
		)
	}
}

// Create a BulkerV1 object that implements IBulker
func CreateBulkerV1() BulkerV1 {
	return BulkerV1{}
}

/*
provides the structure accepted by the endpoint [BulkV2] provded by [Zincsearch].

[BulkV2]: https://zincsearch-docs.zinc.dev/api/document/bulkv2/
[Zincsearch]: https://zincsearch-docs.zinc.dev/
*/
type bulkResponse struct {
	Index   string        `json:"index"`
	Records []*model.Mail `json:"records"`
}

/*
Provides an methods responsible for uploading data to the database.

This object is based on the [BulkV2] endpoint provided by [Zincsearch]

[BulkV2]: https://zincsearch-docs.zinc.dev/api/document/bulkv2/
[Zincsearch]: https://zincsearch-docs.zinc.dev/
*/
type BulkerV2 struct {
}

func (bulk BulkerV2) getCommand() string {
	return "_bulkv2"
}

// Upload the information to the database
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

	if err := myDatabase.BulkRequest(bulk.getCommand(), string(json)); err != nil {
		_logs.LogSVG(
			constants_log.FILE_NAME_ERROR_INDEXER,
			constants_log.OPERATION_BULKER,
			constants_log.ERROR_BULKER_FAILED,
			err,
		)
	}

}

// Create a BulkerV2 object that implements IBulker
func CreateBulkerV2() BulkerV2 {
	return BulkerV2{}
}
