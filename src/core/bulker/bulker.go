package bulker

import (
	model "challenge/src/model"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type bulkResponse struct {
	Index   string
	Records []*model.Mail
}

type IBulker interface {
	GetCommand() string
	GetData() string
	Bulk()
}

type Bulker struct {
	Mails []*model.Mail
	data  string
}

func (bulk Bulker) GetData() string {
	return bulk.data
}

type BulkerV1 struct {
	Bulker
}

func (bulk BulkerV1) GetCommand() string {
	return "_bulkv2"
}

func (bulk *BulkerV1) Bulk() {
	bulkResponse := bulkResponse{"mails", bulk.Mails}
	json, err := json.Marshal(bulkResponse)

	if err != nil {
		log.Println(err)
	}

	bulk.data = string(json)
}

// Bulker V2

type BulkerV2 struct {
	Bulker
}

func (bulk BulkerV2) GetCommand() string {
	return "_bulk"
}

func (bulk *BulkerV2) Bulk() {
	index := "{ \"index\" : { \"_index\" : \"mails\" }} "
	json := ""

	for i := 0; i < len(bulk.Mails); i++ {
		json += index + "\n"
		json += bulk.Mails[i].String() + "\n"
	}

	bulk.data = json
}

func RequestBulk(bulker IBulker) {

	url := "http://localhost:4080/api/" + bulker.GetCommand()
	bulker.Bulk()

	data := strings.NewReader(bulker.GetData())

	req, err := http.NewRequest("POST", url, data)
	if err != nil {
		log.Fatal(err)
	}

	req.SetBasicAuth("admin", "Complexpass#123")
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
