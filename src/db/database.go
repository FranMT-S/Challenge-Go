package myDatabase

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	constants_log "github.com/FranMT-S/Challenge-Go/src/constants/logs"
	_logs "github.com/FranMT-S/Challenge-Go/src/logs"
	myMiddleware "github.com/FranMT-S/Challenge-Go/src/middleware"
)

var z_database *zincDatabase

type zincDatabase struct {
	client *http.Client
}

// ZincDatabase returns a single instance of the database
func ZincDatabase() *zincDatabase {
	if z_database == nil {
		z_database = &zincDatabase{client: &http.Client{}}
	}

	return z_database
}

// CreateIndex Try to create the database index in case it does not exist.
func (db zincDatabase) CreateIndex() error {
	index := fmt.Sprintf(`{
		"name": "%v",
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
			"highlightable": true
		  },
		  "Bcc": {
			"type": "text",
			"index": true,
			"store": false,
			"sortable": true,
			"aggregatable": true,
			"highlightable": true
		  },
		  "Cc": {
			"type": "text",
			"index": true,
			"store": false,
			"sortable": true,
			"aggregatable": true,
			"highlightable": true
		  },
		  "Content": {
			"type": "text",
			"index": true,
			"store": false,
			"sortable": true,
			"aggregatable": true,
			"highlightable": true
		  },
		  "Content_Transfer_Encoding": {
			"type": "text",
			"index": true,
			"store": false,
			"sortable": true,
			"aggregatable": true,
			"highlightable": true
		  },
		  "Content_Type": {
			"type": "text",
			"index": true,
			"store": false,
			"sortable": true,
			"aggregatable": true,
			"highlightable": true
		  },
		  "From": {
			"type": "text",
			"index": true,
			"store": false,
			"sortable": true,
			"aggregatable": true,
			"highlightable": true
		  },
		  "Message_ID": {
			"type": "text",
			"index": true,
			"store": false,
			"sortable": true,
			"aggregatable": true,
			"highlightable": true
		  },
		  "Mime_Version": {
			"type": "text",
			"index": true,
			"store": false,
			"sortable": true,
			"aggregatable": true,
			"highlightable": true
		  },
		  "Subject": {
			"type": "text",
			"index": true,
			"store": false,
			"sortable": true,
			"aggregatable": true,
			"highlightable": true
		  },
		  "To": {
			"type": "text",
			"index": true,
			"store": false,
			"sortable": true,
			"aggregatable": true,
			"highlightable": true
		  },
		  "X_FileName": {
			"type": "text",
			"index": true,
			"store": false,
			"sortable": true,
			"aggregatable": true,
			"highlightable": true
		  },
		  "X_Folder": {
			"type": "text",
			"index": true,
			"store": false,
			"sortable": true,
			"aggregatable": true,
			"highlightable": true
		  },
		  "X_From": {
			"type": "text",
			"index": true,
			"store": false,
			"sortable": true,
			"aggregatable": true,
			"highlightable": true
		  },
		  "X_Origin": {
			"type": "text",
			"index": true,
			"store": false,
			"sortable": true,
			"aggregatable": true,
			"highlightable": true
		  },
		  "X_To": {
			"type": "text",
			"index": true,
			"store": false,
			"sortable": true,
			"aggregatable": true,
			"highlightable": true
		  },
		  "X_bcc": {
			"type": "text",
			"index": true,
			"store": false,
			"sortable": true,
			"aggregatable": true,
			"highlightable": true
		  },
		  "X_cc": {
			"type": "text",
			"index": true,
			"store": false,
			"sortable": true,
			"aggregatable": true,
			"highlightable": true
		  },
		  "_id": {
			"type": "keyword",
			"index": true,
			"store": false,
			"sortable": true,
			"aggregatable": true,
			"highlightable": true
		  }
		}
		}
	}`, os.Getenv("INDEX"))

	url := os.Getenv("URL") + "index"

	data := strings.NewReader(index)
	req, err := http.NewRequest("POST", url, data)
	if err != nil {
		_logs.Println(constants_log.ERROR_CREATE_BASE + ":" + err.Error())
		_logs.LogSVG(
			constants_log.FILE_NAME_ERROR_DATABASE,
			constants_log.OPERATION_DATABASE,
			constants_log.ERROR_CREATE_BASE,
			err,
		)

		return err
	}

	myMiddleware.ZincHeader(req)

	resp, err := db.client.Do(req)
	if err != nil {

		_logs.Println(constants_log.ERROR_DATA_BASE + ":" + err.Error())
		_logs.LogSVG(
			constants_log.FILE_NAME_ERROR_DATABASE,
			constants_log.OPERATION_DATABASE,
			constants_log.ERROR_DATA_BASE,
			err,
		)

		return err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		_logs.Println(constants_log.ERROR_DATA_BASE + ":" + err.Error())
		_logs.LogSVG(
			constants_log.FILE_NAME_ERROR_DATABASE,
			constants_log.OPERATION_DATABASE,
			constants_log.ERROR_DATA_BASE,
			err,
		)

		return err
	}

	if string(body) != fmt.Sprintf(`{"error":"index [%v] already exists"}`, os.Getenv("INDEX")) {
		_logs.ColorGreen()
		fmt.Println("index created")
		fmt.Println(string(body))
		_logs.ColorWhite()
	}

	return nil
}

// BulkRequest makes a request to the database to store the files.
func BulkRequest(command, mailsData string) error {

	url := os.Getenv("URL") + command

	data := strings.NewReader(mailsData)
	req, err := http.NewRequest("POST", url, data)
	if err != nil {
		return err
	}

	myMiddleware.ZincHeader(req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println("\n" + string(body))
	return nil
}
