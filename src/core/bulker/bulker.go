package bulker

import (
	model "challenge/src/model"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// Formato para el Request Bulker V1
type bulkResponse struct {
	Index   string
	Records []model.Mail
}

// GetCommand - comando que se usara en el Request
// GetData - Obtener los datos usarse despues de un bulk
// SetMails - Establecer correos
// GetMails - Obtener correos
// Bulk - Transforma los correos en el formato que usare el Request al subir los datos.
type IBulker interface {
	GetCommand() string
	GetData() string
	SetMails(mails []model.Mail)
	GetMails() []model.Mail
	Bulk()
}

/*
Bulker estructura que heredaran los Bulker Concretos, simula una clase abstracta
no instanciar y usar por si sola.
*/
type Bulker struct {
	mails []model.Mail
	data  string
}

func (bulk Bulker) GetCommand() string {
	return ""
}

func (bulk Bulker) GetData() string {
	return bulk.data
}

func (bulk Bulker) GetMails() []model.Mail {
	return bulk.mails
}

func (bulk *Bulker) SetMails(mails []model.Mail) {
	bulk.mails = mails
}

func (bulk Bulker) Bulk() {}

/*
-----------------------------------
Section Bulker V1
-----------------------------------
*/

type BulkerV1 struct {
	Bulker
}

func (bulk BulkerV1) GetCommand() string {
	return "_bulkv2"
}

func (bulk *BulkerV1) Bulk() {
	bulkResponse := bulkResponse{os.Getenv("INDEX"), bulk.GetMails()}
	json, err := json.Marshal(bulkResponse)

	if err != nil {
		log.Println(err)
	}

	bulk.data = string(json)
}

/*
Tiene un funcion con pointer receiver por lo que es necesario
retornar un puntero para que la interface IBulker lo acepte
*/
func CreateBulkerV1() *BulkerV1 {
	var b Bulker

	return &BulkerV1{b}
}

/*
-----------------------------------
End Section Bulker V1
-----------------------------------
*/

/*
-----------------------------------
Section Bulker V2
-----------------------------------
*/

type BulkerV2 struct {
	Bulker
}

func (bulk BulkerV2) GetCommand() string {
	return "_bulk"
}

func (bulk *BulkerV2) Bulk() {
	index := fmt.Sprintf(`{ "index" : { "_index" : "%v" } }  `, os.Getenv("INDEX"))
	json := ""
	mails := bulk.GetMails()
	for i := 0; i < len(mails); i++ {
		json += index + "\n"
		json += mails[i].String() + "\n"
	}

	bulk.data = json
}

/*
Tiene un funcion con pointer receiver por lo que es necesario
retornar un puntero para que la interface IBulker lo acepte
*/
func CreateBulkerV2() *BulkerV2 {
	var b Bulker

	return &BulkerV2{b}
}

/*
-----------------------------------
End Bulker V2
-----------------------------------
*/
