package model

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

// Constantes para el Map
const (
	K_MESSAGE_ID                = "Message_ID"
	K_DATE                      = "Date"
	K_FROM                      = "From"
	K_TO                        = "To"
	K_SUBJECT                   = "Subject"
	K_CC                        = "Cc"
	K_MIME_VERSION              = "Mime_Version"
	K_CONTENT_TYPE              = "Content_Type"
	K_CONTENT_TRANSFER_ENCODING = "Content_Transfer_Encoding"
	K_BCC                       = "Bcc"
	K_X_FROM                    = "X_From"
	K_X_TO                      = "X_To"
	K_X_CC                      = "X_cc"
	K_X_BCC                     = "X_bcc"
	K_X_FOLDER                  = "X_Folder"
	K_X_ORIGIN                  = "X_Origin"
	K_X_FILENAME                = "X_FileName"
	K_CONTENT                   = "Content"
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

func cleanField(s string) string {

	s = strings.ReplaceAll(s, "\r", "")
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.TrimSpace(s)
	return s
}

func parseDate(s string) string {

	// Parse the date and time string

	t, err := time.Parse("Mon, _2 Jan 2006 15:04:05 -0700 (MST)", s)
	if err != nil {
		log.Panic("Error al parsear la fecha:", err)
	}

	return t.Format("2006-01-02T15:04:05Z")
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

func MailFromJson(_json []byte) Mail {
	var mail Mail

	if err := json.Unmarshal(_json, &mail); err != nil {
		fmt.Println(err)
		return mail
	}

	return mail
}

func MailFromMap(_map map[string]string) Mail {
	var mail Mail

	mail.Message_ID = cleanField(_map[K_MESSAGE_ID])
	mail.Date = parseDate(cleanField(_map[K_DATE]))
	mail.From = cleanField(_map[K_FROM])
	mail.To = cleanField(_map[K_TO])
	mail.Subject = cleanField(_map[K_SUBJECT])
	mail.Cc = cleanField(_map[K_CC])
	mail.Mime_Version = cleanField(_map[K_MIME_VERSION])
	mail.Content_Type = cleanField(_map[K_CONTENT_TYPE])
	mail.Content_Transfer_Encoding = cleanField(_map[K_CONTENT_TRANSFER_ENCODING])
	mail.Bcc = cleanField(_map[K_BCC])
	mail.X_From = cleanField(_map[K_X_FROM])
	mail.X_To = cleanField(_map[K_X_TO])
	mail.X_cc = cleanField(_map[K_X_CC])
	mail.X_bcc = cleanField(_map[K_X_BCC])
	mail.X_Folder = cleanField(_map[K_X_FOLDER])
	mail.X_Origin = cleanField(_map[K_X_ORIGIN])
	mail.X_FileName = cleanField(_map[K_X_FILENAME])
	mail.Content = _map[K_CONTENT]

	return mail
}

// Map con los campos del correo
func NewMapMail() map[string]string {
	return map[string]string{
		K_MESSAGE_ID:                "",
		K_DATE:                      "",
		K_FROM:                      "",
		K_TO:                        "",
		K_SUBJECT:                   "",
		K_CC:                        "",
		K_MIME_VERSION:              "",
		K_CONTENT_TYPE:              "",
		K_CONTENT_TRANSFER_ENCODING: "",
		K_BCC:                       "",
		K_X_FROM:                    "",
		K_X_TO:                      "",
		K_X_CC:                      "",
		K_X_BCC:                     "",
		K_X_FOLDER:                  "",
		K_X_ORIGIN:                  "",
		K_X_FILENAME:                "",
		K_CONTENT:                   "",
	}
}

// Para Obtener el numero total de mensajes
type Stats struct {
	DocNum int `json:"doc_num"`
}

type ResponseIndexData struct {
	Stats Stats `json:"stats"`
}
