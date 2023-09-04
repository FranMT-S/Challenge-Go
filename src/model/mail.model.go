package model

import (
	"encoding/json"
	"log"
)

/*
	Modelo de la estructura del correo
*/

type Mail struct {
	Message_ID                string
	Date                      string
	From                      string
	To                        string
	Subject                   string
	Cc                        string
	Mime_Version              string
	Content_Type              string
	Content_Transfer_Encoding string
	Bcc                       string
	X_From                    string
	X_To                      string
	X_cc                      string
	X_bcc                     string
	X_Folder                  string
	X_Origin                  string
	X_FileName                string
	Content                   string
}

// Devuelve un string del correo en formato Json
func (mail Mail) String() string {
	return mail.ToJson()
}

// Transforma el correo a JSON
func (mail Mail) ToJson() string {
	bytes, err := mail.ToJsonBytes()

	if err != nil {
		log.Println(err)
		return ""
	}

	return string(bytes)
}

// Transforma el correo a Json pero devuelve los datos en un arreglo de bytes
func (mail Mail) ToJsonBytes() ([]byte, error) {
	return json.Marshal(mail)
}
