package parser

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"

	model "github.com/FranMT-S/Challenge-Go/src/model"
)

const (
	MESSAGE_ID                = "Message-ID:"
	DATE                      = "Date:"
	FROM                      = "From:"
	TO                        = "To:"
	SUBJECT                   = "Subject:"
	CC                        = "Cc:"
	MIME_VERSION              = "Mime-Version:"
	CONTENT_TYPE              = "Content-Type:"
	CONTENT_TRANSFER_ENCODING = "Content-Transfer-Encoding:"
	BCC                       = "Bcc:"
	X_FROM                    = "X-From:"
	X_TO                      = "X-To:"
	X_CC                      = "X-cc:"
	X_BCC                     = "X-bcc:"
	X_FOLDER                  = "X-Folder:"
	X_ORIGIN                  = "X-Origin:"
	X_FILENAME                = "X-FileName:"

	// Para usar un map con la misma estructura del model
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

	K_FATHER = "Father" // Utilizado en la interace lineMail
)

// Map con los campos del correo
func newMapMail() map[string]string {

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

type lineMail struct {
	lineFather    *lineMail
	field         string
	data          string
	lineToAnalize string
	numberLine    int
	lock          *sync.Mutex
}

func (lineMail *lineMail) getField() string {
	lineMail.lock.Lock()
	defer lineMail.lock.Unlock()

	if lineMail.field == K_FATHER {
		lineMail.field = lineMail.lineFather.getField()
	}

	return lineMail.field
}

func newLineMail(_fatherLine *lineMail, _lineToAnalize string, _numberLine int) *lineMail {
	return &lineMail{lineFather: _fatherLine, lineToAnalize: _lineToAnalize, numberLine: _numberLine, lock: &sync.Mutex{}}
}

// Line Reader String, lee una linea en string para formatearla
type ILineReader[T string | *lineMail] interface {
	Read(line T)
	getMapData() map[string]string
}

type lineByLineReader struct {
	mailMap       map[string]string
	beforeLecture string
}

func newLineByLineReader() *lineByLineReader {
	return &lineByLineReader{mailMap: newMapMail()}
}

func (lineReader lineByLineReader) getMapData() map[string]string {

	return lineReader.mailMap
}

func (lineReader *lineByLineReader) Read(line string) {

	if lineReader.mailMap[K_X_FILENAME] != "" {
		lineReader.mailMap[K_CONTENT] += line
	} else if strings.Contains(line, X_FROM) && lineReader.mailMap[K_X_FROM] == "" {
		lineReader.mailMap[K_X_FROM] = line[len(X_FROM):]
		lineReader.beforeLecture = K_X_FROM
	} else if strings.Contains(line, X_TO) && lineReader.mailMap[K_X_TO] == "" {
		lineReader.mailMap[K_X_TO] = line[len(X_TO):]
		lineReader.beforeLecture = K_X_TO
	} else if strings.Contains(line, X_CC) && lineReader.mailMap[K_X_CC] == "" {
		lineReader.mailMap[K_X_CC] = line[len(X_CC):]
		lineReader.beforeLecture = K_X_CC
	} else if strings.Contains(line, X_BCC) && lineReader.mailMap[K_X_BCC] == "" {
		lineReader.mailMap[K_X_BCC] = line[len(X_BCC):]
		lineReader.beforeLecture = K_X_BCC
	} else if strings.Contains(line, X_FOLDER) && lineReader.mailMap[K_X_FOLDER] == "" {
		lineReader.mailMap[K_X_FOLDER] = line[len(X_FOLDER):]
		lineReader.beforeLecture = K_X_FOLDER
	} else if strings.Contains(line, X_ORIGIN) && lineReader.mailMap[K_X_ORIGIN] == "" {
		lineReader.mailMap[K_X_ORIGIN] = line[len(X_ORIGIN):]
		lineReader.beforeLecture = K_X_ORIGIN
	} else if strings.Contains(line, X_FILENAME) && lineReader.mailMap[K_X_FILENAME] == "" {
		lineReader.mailMap[K_X_FILENAME] = line[len(X_FILENAME):]
		lineReader.beforeLecture = K_X_FILENAME
	} else if strings.Contains(line, MESSAGE_ID) && lineReader.mailMap[K_MESSAGE_ID] == "" {
		lineReader.mailMap[K_MESSAGE_ID] = line[len(MESSAGE_ID):]
		lineReader.beforeLecture = K_MESSAGE_ID
	} else if strings.Contains(line, DATE) && lineReader.mailMap[K_DATE] == "" {
		lineReader.mailMap[K_DATE] = line[len(DATE):]
		lineReader.beforeLecture = K_DATE
	} else if strings.Contains(line, FROM) && lineReader.mailMap[K_FROM] == "" {
		lineReader.mailMap[K_FROM] = line[len(FROM):]
		lineReader.beforeLecture = K_FROM
	} else if strings.Contains(line, TO) && lineReader.mailMap[K_TO] == "" {
		lineReader.mailMap[K_TO] = line[len(TO):]
		lineReader.beforeLecture = K_TO
	} else if strings.Contains(line, SUBJECT) && lineReader.mailMap[K_SUBJECT] == "" {
		lineReader.mailMap[K_SUBJECT] = line[len(SUBJECT):]
		lineReader.beforeLecture = K_SUBJECT
	} else if strings.Contains(line, CC) && lineReader.mailMap[K_CC] == "" {
		lineReader.mailMap[K_CC] = line[len(CC):]
		lineReader.beforeLecture = K_CC
	} else if strings.Contains(line, BCC) && lineReader.mailMap[K_BCC] == "" {
		lineReader.mailMap[K_BCC] = line[len(BCC):]
		lineReader.beforeLecture = K_BCC
	} else if strings.Contains(line, MIME_VERSION) && lineReader.mailMap[K_MIME_VERSION] == "" {
		lineReader.mailMap[K_MIME_VERSION] = line[len(MIME_VERSION):]
		lineReader.beforeLecture = K_MIME_VERSION
	} else if strings.Contains(line, CONTENT_TYPE) && lineReader.mailMap[K_CONTENT_TYPE] == "" {
		lineReader.mailMap[K_CONTENT_TYPE] = line[len(CONTENT_TYPE):]
		lineReader.beforeLecture = K_CONTENT_TYPE
	} else if strings.Contains(line, CONTENT_TRANSFER_ENCODING) && lineReader.mailMap[K_CONTENT_TRANSFER_ENCODING] == "" {
		lineReader.mailMap[K_CONTENT_TRANSFER_ENCODING] = line[len(CONTENT_TRANSFER_ENCODING):]
		lineReader.beforeLecture = K_CONTENT_TRANSFER_ENCODING
	} else if lineReader.beforeLecture != "" {
		lineReader.mailMap[lineReader.beforeLecture] += line
	}
}

//

func newLineByLineReaderAsync() *lineByLineReaderAsync {
	return &lineByLineReaderAsync{lock: &sync.Mutex{}, headLineFlag: -1}
}

type lineByLineReaderAsync struct {
	line         *lineMail
	lock         *sync.Mutex
	headLineFlag int
}

func (lineReader *lineByLineReaderAsync) Read(line *lineMail) {

	if lineReader.headLineFlag > 0 && lineReader.headLineFlag < line.numberLine {
		line.data = line.lineToAnalize
		line.field = K_CONTENT
	} else if strings.Contains(line.lineToAnalize, X_FROM) {
		line.data = line.lineToAnalize[len(X_FROM):]
		line.field = K_X_FROM
	} else if strings.Contains(line.lineToAnalize, X_TO) {
		line.data = line.lineToAnalize[len(X_TO):]
		line.field = K_X_TO
	} else if strings.Contains(line.lineToAnalize, X_CC) {
		line.data = line.lineToAnalize[len(X_CC):]
		line.field = K_X_CC
	} else if strings.Contains(line.lineToAnalize, X_BCC) {
		line.data = line.lineToAnalize[len(X_BCC):]
		line.field = K_X_BCC
	} else if strings.Contains(line.lineToAnalize, X_FOLDER) {
		line.data = line.lineToAnalize[len(X_FOLDER):]
		line.field = K_X_FOLDER
	} else if strings.Contains(line.lineToAnalize, X_ORIGIN) {
		line.data = line.lineToAnalize[len(X_ORIGIN):]
		line.field = K_X_ORIGIN
	} else if strings.Contains(line.lineToAnalize, X_FILENAME) {
		line.data = line.lineToAnalize[len(X_FILENAME):]
		line.field = K_X_FILENAME
		lineReader.headLineFlag = line.numberLine
	} else if strings.Contains(line.lineToAnalize, MESSAGE_ID) {
		line.data = line.lineToAnalize[len(MESSAGE_ID):]
		line.field = K_MESSAGE_ID
	} else if strings.Contains(line.lineToAnalize, DATE) {
		line.data = line.lineToAnalize[len(DATE):]
		line.field = K_DATE
	} else if strings.Contains(line.lineToAnalize, FROM) {
		line.data = line.lineToAnalize[len(FROM):]
		line.field = K_FROM
	} else if strings.Contains(line.lineToAnalize, TO) {
		line.data = line.lineToAnalize[len(TO):]
		line.field = K_TO
	} else if strings.Contains(line.lineToAnalize, SUBJECT) {
		line.data = line.lineToAnalize[len(SUBJECT):]
		line.field = K_SUBJECT
	} else if strings.Contains(line.lineToAnalize, CC) {
		line.data = line.lineToAnalize[len(CC):]
		line.field = K_CC
	} else if strings.Contains(line.lineToAnalize, BCC) {
		line.data = line.lineToAnalize[len(BCC):]
		line.field = K_BCC
	} else if strings.Contains(line.lineToAnalize, MIME_VERSION) {
		line.data = line.lineToAnalize[len(MIME_VERSION):]
		line.field = K_MIME_VERSION
	} else if strings.Contains(line.lineToAnalize, CONTENT_TYPE) {
		line.data = line.lineToAnalize[len(CONTENT_TYPE):]
		line.field = K_CONTENT_TYPE
	} else if strings.Contains(line.lineToAnalize, CONTENT_TRANSFER_ENCODING) {
		line.data = line.lineToAnalize[len(CONTENT_TRANSFER_ENCODING):]
		line.field = K_CONTENT_TRANSFER_ENCODING
	} else {
		line.data = line.lineToAnalize
		line.field = K_FATHER
	}

}

func (lineReader *lineByLineReaderAsync) getMapData() map[string]string {
	temp := lineReader.line
	mailMap := newMapMail()
	for {
		if temp == nil {
			break
		}

		if lineReader.headLineFlag < temp.numberLine {
			mailMap[K_CONTENT] = temp.lineToAnalize + mailMap[K_CONTENT]
		} else {
			mailMap[temp.getField()] = temp.data + mailMap[temp.getField()]
		}

		temp = temp.lineFather
	}
	return mailMap
}

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
	var mail model.Mail
	var mailMap map[string]string
	lineByLineReader := newLineByLineReader()
	reader := bufio.NewReader(file)
	// beforeLine := ""
	for {
		lineByte, err := reader.ReadBytes('\n')

		if err != nil {
			line := string(lineByte)
			if len(line) > 0 {
				lineByLineReader.Read(line)

			}
			if err != io.EOF {
				log.Println("Error al parserar el archivo: ", file.Name())

			}
			break
		}

		line := string(lineByte)
		lineByLineReader.Read(line)

	}
	mailMap = lineByLineReader.getMapData()

	jsonStr, err := json.Marshal(mailMap)
	if err != nil {
		fmt.Println(err)
	}

	mail = model.MailFromJson(jsonStr)

	return mail
}

// Maxima cantidad de hilos es 25
type parserAsync struct {
	maxConcurrent int
}

/*
Parseador Asincrono acepta un valor que especifica el limite de lineas que leera al mismo tiempo

Maximo 25 hilos. Minimo 1.

-1 Para usarlo sin limite de hilos pero deberia evitarse.
*/
func NewParserAsyn(_maxConcurrent int) *parserAsync {
	return &parserAsync{maxConcurrent: _maxConcurrent}
}

func (parser parserAsync) Parse(file *os.File) model.Mail {
	// buf := make([]byte, 1024)
	var mail model.Mail
	var mailMap map[string]string
	var wg sync.WaitGroup
	var semaphore chan struct{}

	if parser.maxConcurrent > 35 {
		parser.maxConcurrent = 35
	} else if parser.maxConcurrent <= 0 {
		parser.maxConcurrent = 1
	}

	semaphore = make(chan struct{}, parser.maxConcurrent)

	lineByLineReaderAsync := newLineByLineReaderAsync()
	reader := bufio.NewReader(file)

	for {
		lineByte, err := reader.ReadBytes('\n')
		line := string(lineByte)

		var _newLineMail *lineMail

		if lineByLineReaderAsync.line == nil {
			lineByLineReaderAsync.line = newLineMail(nil, line, 0)
			_newLineMail = lineByLineReaderAsync.line
		} else {
			_newLineMail = newLineMail(lineByLineReaderAsync.line, line, lineByLineReaderAsync.line.numberLine+1)
			lineByLineReaderAsync.line = _newLineMail
		}

		if err != nil {

			// En caso de que quede una ultima linea sin salto de linea
			if len(line) > 0 {
				wg.Add(1)

				semaphore <- struct{}{}
				go func() {
					defer wg.Done()
					lineByLineReaderAsync.Read(_newLineMail)
					<-semaphore
				}()

			}
			if err != io.EOF {
				log.Println("Error al parserar el archivo: ", file.Name())

			}
			break
		}

		wg.Add(1)
		semaphore <- struct{}{}
		go func() {
			defer wg.Done()
			lineByLineReaderAsync.Read(_newLineMail)
			<-semaphore

		}()
	}

	wg.Wait()
	close(semaphore)

	mailMap = lineByLineReaderAsync.getMapData()
	jsonStr, err := json.Marshal(mailMap)
	if err != nil {
		fmt.Println(err)
	}

	mail = model.MailFromJson(jsonStr)

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
