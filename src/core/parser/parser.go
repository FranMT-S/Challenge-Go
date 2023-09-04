package parser

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"

	model "github.com/FranMT-S/Challenge-Go/src/model"
)

const MESSAGE_ID = "Message-ID:"
const DATE = "Date:"
const FROM = "From:"
const TO = "To:"
const SUBJECT = "Subject:"
const CC = "Cc:"
const MIME_VERSION = "Mime-Version:"
const CONTENT_TYPE = "Content-Type:"
const CONTENT_TRANSFER_ENCODING = "Content-Transfer-Encoding:"
const BCC = "Bcc:"
const X_FROM = "X-From:"
const X_TO = "X-To:"
const X_CC = "X-cc:"
const X_BCC = "X-bcc:"
const X_FOLDER = "X-Folder:"
const X_ORIGIN = "X-Origin:"
const X_FILENAME = "X-FileName:"

/*

IParser Mail

Proporciona el metodo para transformar un archivo a un formato de Correo.

*/

type IParserMail interface {
	Parse(file *os.File) model.Mail
}

/*
--------------------
Parseador con Normal
--------------------

Lee linea por linea y asigna el contenido al correo
*/

type ParserNormal struct{}

func (parser ParserNormal) Parse(file *os.File) model.Mail {
	// buf := make([]byte, 1024)
	mail := model.Mail{}

	reader := bufio.NewReader(file)

	for {
		lineByte, err := reader.ReadBytes('\n')

		if err != nil {
			if err != io.EOF {
				log.Println("Error al parserar el archivo: ", file.Name())
				log.Println(err)
			}
			break
		}

		line := string(lineByte)
		// fmt.Println(line)
		switch {
		case strings.Contains(line, X_FROM):
			mail.X_From = line[len(X_FROM):]
			break
		case strings.Contains(line, X_TO):
			mail.X_To = line[len(X_TO):]
			break
		case strings.Contains(line, X_CC):
			mail.X_cc = line[len(X_CC):]
			break
		case strings.Contains(line, X_BCC):
			mail.X_bcc = line[len(X_BCC):]
			break
		case strings.Contains(line, X_FOLDER):
			mail.X_Folder = line[len(X_FOLDER):]
			break
		case strings.Contains(line, X_ORIGIN):
			mail.X_Origin = line[len(X_ORIGIN):]
			break
		case strings.Contains(line, X_FILENAME):
			mail.X_FileName = line[len(X_FILENAME):]
			break
		case strings.Contains(line, MESSAGE_ID):
			mail.Message_ID = line[len(MESSAGE_ID):]
			break
		case strings.Contains(line, DATE):
			mail.Date = line[len(DATE):]
			break
		case strings.Contains(line, FROM):
			mail.From = line[len(FROM):]
			break
		case strings.Contains(line, TO):
			mail.To = line[len(TO):]
			break
		case strings.Contains(line, SUBJECT):
			mail.Subject = line[len(SUBJECT):]
			break
		case strings.Contains(line, CC):
			mail.Cc = line[len(CC):]
			break
		case strings.Contains(line, BCC):
			mail.Bcc = line[len(BCC):]
			break
		case strings.Contains(line, MIME_VERSION):
			mail.Mime_Version = line[len(MIME_VERSION):]
			break
		case strings.Contains(line, CONTENT_TYPE):
			mail.Content_Type = line[len(CONTENT_TYPE):]
			break
		case strings.Contains(line, CONTENT_TRANSFER_ENCODING):
			mail.Content_Transfer_Encoding = line[len(CONTENT_TRANSFER_ENCODING):]
			break

		default:
			mail.Content += line
		}

	}
	return mail
}

/*
--------------------
Parseador con Regex
--------------------

Usa Expresiones Regulares para parsear el contenido
*/

type ParserWithRegex struct{}

func (parser ParserWithRegex) Parse(file *os.File) model.Mail {
	// buf := make([]byte, 1024)
	mail := model.Mail{}

	bytes, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	data := string(bytes)
	// fmt.Println(data)
	var valid = regexp.MustCompile("^(?P<From>\n?From:.+\n)")

	fmt.Printf("%#v\n", data)
	fmt.Printf("%#v\n", valid.FindStringSubmatch(data))
	fmt.Printf("%#v\n", valid.SubexpNames())

	return mail
}
