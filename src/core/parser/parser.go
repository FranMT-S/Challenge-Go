package parser

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	model "github.com/FranMT-S/Challenge-Go/src/model"
)

func cleanField(s string) string {

	s = strings.ReplaceAll(s, "\r", "")
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "\t", "")
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
	return &lineByLineReader{mailMap: model.NewMapMail()}
}

func (lineReader lineByLineReader) getMapData() map[string]string {

	return lineReader.mailMap
}

func (lineReader *lineByLineReader) Read(line string) {

	if lineReader.mailMap[model.K_X_FILENAME] != "" {
		lineReader.mailMap[model.K_CONTENT] += line
	} else if strings.HasPrefix(line, X_FROM) && lineReader.mailMap[model.K_X_FROM] == "" {
		lineReader.mailMap[model.K_X_FROM] = line[len(X_FROM):]
		lineReader.beforeLecture = model.K_X_FROM
	} else if strings.HasPrefix(line, X_TO) && lineReader.mailMap[model.K_X_TO] == "" {
		lineReader.mailMap[model.K_X_TO] = line[len(X_TO):]
		lineReader.beforeLecture = model.K_X_TO
	} else if strings.HasPrefix(line, X_CC) && lineReader.mailMap[model.K_X_CC] == "" {
		lineReader.mailMap[model.K_X_CC] = line[len(X_CC):]
		lineReader.beforeLecture = model.K_X_CC
	} else if strings.HasPrefix(line, X_BCC) && lineReader.mailMap[model.K_X_BCC] == "" {
		lineReader.mailMap[model.K_X_BCC] = line[len(X_BCC):]
		lineReader.beforeLecture = model.K_X_BCC
	} else if strings.HasPrefix(line, X_FOLDER) && lineReader.mailMap[model.K_X_FOLDER] == "" {
		lineReader.mailMap[model.K_X_FOLDER] = line[len(X_FOLDER):]
		lineReader.beforeLecture = model.K_X_FOLDER
	} else if strings.HasPrefix(line, X_ORIGIN) && lineReader.mailMap[model.K_X_ORIGIN] == "" {
		lineReader.mailMap[model.K_X_ORIGIN] = line[len(X_ORIGIN):]
		lineReader.beforeLecture = model.K_X_ORIGIN
	} else if strings.HasPrefix(line, X_FILENAME) && lineReader.mailMap[model.K_X_FILENAME] == "" {
		lineReader.mailMap[model.K_X_FILENAME] = line[len(X_FILENAME):]
		lineReader.beforeLecture = model.K_X_FILENAME
	} else if strings.HasPrefix(line, MESSAGE_ID) && lineReader.mailMap[model.K_MESSAGE_ID] == "" {
		lineReader.mailMap[model.K_MESSAGE_ID] = line[len(MESSAGE_ID):]
		lineReader.beforeLecture = model.K_MESSAGE_ID
	} else if strings.HasPrefix(line, DATE) && lineReader.mailMap[model.K_DATE] == "" {
		lineReader.mailMap[model.K_DATE] = line[len(DATE):]
		lineReader.beforeLecture = model.K_DATE
	} else if strings.HasPrefix(line, FROM) && lineReader.mailMap[model.K_FROM] == "" {
		lineReader.mailMap[model.K_FROM] = line[len(FROM):]
		lineReader.beforeLecture = model.K_FROM
	} else if strings.HasPrefix(line, TO) && lineReader.mailMap[model.K_TO] == "" {
		lineReader.mailMap[model.K_TO] = line[len(TO):]
		lineReader.beforeLecture = model.K_TO
	} else if strings.HasPrefix(line, SUBJECT) && lineReader.mailMap[model.K_SUBJECT] == "" {
		lineReader.mailMap[model.K_SUBJECT] = line[len(SUBJECT):]
		lineReader.beforeLecture = model.K_SUBJECT
	} else if strings.HasPrefix(line, CC) && lineReader.mailMap[model.K_CC] == "" {
		lineReader.mailMap[model.K_CC] = line[len(CC):]
		lineReader.beforeLecture = model.K_CC
	} else if strings.HasPrefix(line, BCC) && lineReader.mailMap[model.K_BCC] == "" {
		lineReader.mailMap[model.K_BCC] = line[len(BCC):]
		lineReader.beforeLecture = model.K_BCC
	} else if strings.HasPrefix(line, MIME_VERSION) && lineReader.mailMap[model.K_MIME_VERSION] == "" {
		lineReader.mailMap[model.K_MIME_VERSION] = line[len(MIME_VERSION):]
		lineReader.beforeLecture = model.K_MIME_VERSION
	} else if strings.HasPrefix(line, CONTENT_TYPE) && lineReader.mailMap[model.K_CONTENT_TYPE] == "" {
		lineReader.mailMap[model.K_CONTENT_TYPE] = line[len(CONTENT_TYPE):]
		lineReader.beforeLecture = model.K_CONTENT_TYPE
	} else if strings.HasPrefix(line, CONTENT_TRANSFER_ENCODING) && lineReader.mailMap[model.K_CONTENT_TRANSFER_ENCODING] == "" {
		lineReader.mailMap[model.K_CONTENT_TRANSFER_ENCODING] = line[len(CONTENT_TRANSFER_ENCODING):]
		lineReader.beforeLecture = model.K_CONTENT_TRANSFER_ENCODING
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
		line.field = model.K_CONTENT
	} else if strings.HasPrefix(line.lineToAnalize, X_FROM) {
		line.data = line.lineToAnalize[len(X_FROM):]
		line.field = model.K_X_FROM
	} else if strings.HasPrefix(line.lineToAnalize, X_TO) {
		line.data = line.lineToAnalize[len(X_TO):]
		line.field = model.K_X_TO
	} else if strings.HasPrefix(line.lineToAnalize, X_CC) {
		line.data = line.lineToAnalize[len(X_CC):]
		line.field = model.K_X_CC
	} else if strings.HasPrefix(line.lineToAnalize, X_BCC) {
		line.data = line.lineToAnalize[len(X_BCC):]
		line.field = model.K_X_BCC
	} else if strings.HasPrefix(line.lineToAnalize, X_FOLDER) {
		line.data = line.lineToAnalize[len(X_FOLDER):]
		line.field = model.K_X_FOLDER
	} else if strings.HasPrefix(line.lineToAnalize, X_ORIGIN) {
		line.data = line.lineToAnalize[len(X_ORIGIN):]
		line.field = model.K_X_ORIGIN
	} else if strings.HasPrefix(line.lineToAnalize, X_FILENAME) {
		line.data = line.lineToAnalize[len(X_FILENAME):]
		line.field = model.K_X_FILENAME
		lineReader.headLineFlag = line.numberLine
	} else if strings.HasPrefix(line.lineToAnalize, MESSAGE_ID) {
		line.data = line.lineToAnalize[len(MESSAGE_ID):]
		line.field = model.K_MESSAGE_ID
	} else if strings.HasPrefix(line.lineToAnalize, DATE) {
		line.data = line.lineToAnalize[len(DATE):]
		line.field = model.K_DATE
	} else if strings.HasPrefix(line.lineToAnalize, FROM) {
		line.data = line.lineToAnalize[len(FROM):]
		line.field = model.K_FROM
	} else if strings.HasPrefix(line.lineToAnalize, TO) {
		line.data = line.lineToAnalize[len(TO):]
		line.field = model.K_TO
	} else if strings.HasPrefix(line.lineToAnalize, SUBJECT) {
		line.data = line.lineToAnalize[len(SUBJECT):]
		line.field = model.K_SUBJECT
	} else if strings.HasPrefix(line.lineToAnalize, CC) {
		line.data = line.lineToAnalize[len(CC):]
		line.field = model.K_CC
	} else if strings.HasPrefix(line.lineToAnalize, BCC) {
		line.data = line.lineToAnalize[len(BCC):]
		line.field = model.K_BCC
	} else if strings.HasPrefix(line.lineToAnalize, MIME_VERSION) {
		line.data = line.lineToAnalize[len(MIME_VERSION):]
		line.field = model.K_MIME_VERSION
	} else if strings.HasPrefix(line.lineToAnalize, CONTENT_TYPE) {
		line.data = line.lineToAnalize[len(CONTENT_TYPE):]
		line.field = model.K_CONTENT_TYPE
	} else if strings.HasPrefix(line.lineToAnalize, CONTENT_TRANSFER_ENCODING) {
		line.data = line.lineToAnalize[len(CONTENT_TRANSFER_ENCODING):]
		line.field = model.K_CONTENT_TRANSFER_ENCODING
	} else {
		line.data = line.lineToAnalize
		line.field = K_FATHER
	}

}

func (lineReader *lineByLineReaderAsync) getMapData() map[string]string {
	temp := lineReader.line
	mailMap := model.NewMapMail()
	for {
		if temp == nil {
			break
		}

		if lineReader.headLineFlag < temp.numberLine {
			mailMap[model.K_CONTENT] = temp.lineToAnalize + mailMap[model.K_CONTENT]
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
		line := string(lineByte)
		if err != nil && len(line) <= 0 {

			if err != io.EOF {
				log.Println("Error al parserar el archivo: ", file.Name())

			}
			break
		}

		lineByLineReader.Read(line)

	}

	mailMap = lineByLineReader.getMapData()
	mail = model.MailFromMap(mailMap)

	return mail
}

// Maxima cantidad de hilos es 25
type parserAsync struct {
	maxConcurrent int
}

/*
Parseador Asincrono acepta un valor que especifica el limite de lineas que leera al mismo tiempo

Maximo 50 hilos. Minimo 1.

-1 Para usarlo sin limite de hilos pero deberia evitarse.
*/
func NewParserAsync(_maxConcurrent int) *parserAsync {
	return &parserAsync{maxConcurrent: _maxConcurrent}
}

func (parser parserAsync) Parse(file *os.File) model.Mail {
	// buf := make([]byte, 1024)
	var mail model.Mail
	var mailMap map[string]string
	var wg sync.WaitGroup
	var semaphore chan struct{}

	if parser.maxConcurrent > 50 {
		parser.maxConcurrent = 50
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

		if err != nil && len(line) <= 0 {

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
	mail = model.MailFromMap(mailMap)

	return mail
}

/*
--------------------
Parseador con Split
--------------------

Usa Expresiones Regulares para parsear el contenido
*/

type ParserAsyncSplit struct {
	maxConcurrent int
}

func NewParserAsyncSpliter(_maxConcurrent int) *ParserAsyncSplit {
	return &ParserAsyncSplit{maxConcurrent: _maxConcurrent}
}

func (parser ParserAsyncSplit) Parse(file *os.File) model.Mail {

	var mail model.Mail
	var mailMap map[string]string
	var wg sync.WaitGroup
	var semaphore chan struct{}
	lineByLineReaderAsync := newLineByLineReaderAsync()

	if parser.maxConcurrent > 50 {
		parser.maxConcurrent = 50
	} else if parser.maxConcurrent <= 0 {
		parser.maxConcurrent = 1
	}

	semaphore = make(chan struct{}, parser.maxConcurrent)

	bytes, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	re, _ := regexp.Compile(`X-FileName:.+\n{1,}`)

	content := string(bytes)

	index := re.FindStringIndex(content)
	endIndex := index[1]

	dataReader := strings.NewReader(content[:endIndex])
	reader := bufio.NewReader(dataReader)

	for {
		lineByte, err := reader.ReadBytes('\n')
		line := string(lineByte)

		if err != nil && len(line) <= 0 {

			if err != io.EOF {
				log.Println("Error al parserar el archivo: ", file.Name())
			}
			break
		}

		var _newLineMail *lineMail

		if lineByLineReaderAsync.line == nil {
			lineByLineReaderAsync.line = newLineMail(nil, line, 0)
			_newLineMail = lineByLineReaderAsync.line
		} else {
			_newLineMail = newLineMail(lineByLineReaderAsync.line, line, lineByLineReaderAsync.line.numberLine+1)
			lineByLineReaderAsync.line = _newLineMail
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

	mailMap[model.K_CONTENT] = content[:endIndex]
	mail = model.MailFromMap(mailMap)

	return mail
}

/*
--------------------
Parseador con Split
--------------------

Usa Expresiones Regulares para parsear el contenido
*/

type ParserAsyncRegex struct {
	maxConcurrent int
}

func NewParserAsyncRegex(_maxConcurrent int) *ParserAsyncRegex {
	return &ParserAsyncRegex{maxConcurrent: _maxConcurrent}
}

func (parser ParserAsyncRegex) Parse(file *os.File) model.Mail {

	var mail model.Mail
	var wg sync.WaitGroup
	var semaphore chan struct{}
	var mutex = &sync.Mutex{}

	mailMap := map[string]string{}
	indexMap := map[int]string{}
	noMatchMap := map[int]string{}
	i := -1

	if parser.maxConcurrent > 50 {
		parser.maxConcurrent = 50
	} else if parser.maxConcurrent <= 0 {
		parser.maxConcurrent = 1
	}

	semaphore = make(chan struct{}, parser.maxConcurrent)

	bytes, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	content := string(bytes)

	// re, _ := regexp.Compile(`X-FileName:.+\n{1,}`)

	// index := re.FindStringIndex(content)
	// endIndex := index[1]
	// header := strings.TrimSpace(content[:endIndex])
	// body := strings.TrimSpace(content[endIndex:])

	re, _ := regexp.Compile(`(\r\n){2,}|\n{2,}`)
	reLine, _ := regexp.Compile(`^([\w-_]+:)(.+)`)

	match := re.Split(content, 2)

	header := match[0]
	body := match[1]

	dataReader := strings.NewReader(header)
	reader := bufio.NewReader(dataReader)

	for {
		lineByte, err := reader.ReadBytes('\n')
		line := string(lineByte)
		i++
		indexLine := i

		if err != nil && len(line) <= 0 {

			if err != io.EOF {
				log.Println("Error al parserar el archivo: ", file.Name())
			}

			break
		}

		fmt.Println("Linea:", line)

		wg.Add(1)
		semaphore <- struct{}{}
		// fmt.Println("Entrando: ", line)
		go func() {
			defer wg.Done()
			match := reLine.FindStringSubmatch(line)
			if len(match) > 0 {
				// match[1] el campo, match[2] el contenido del campo
				mutex.Lock()
				mailMap[match[1]] = match[2]
				indexMap[indexLine] = match[1]
				mutex.Unlock()
			} else {
				// El campo sera el de la linea anterior
				mutex.Lock()
				indexMap[indexLine] = ""
				noMatchMap[indexLine] = line
				mutex.Unlock()
			}

			<-semaphore

		}()
	}

	wg.Wait()
	close(semaphore)

	for j := 0; j <= i; j++ {

		if indexMap[j] == "" {
			indexMap[j] = indexMap[j-1]
			mailMap[indexMap[j]] += noMatchMap[j]
		}
	}

	// Corregis los campos que no hicieron match

	// mailMap = lineByLineReaderAsync.getMapData()

	mail.Message_ID = cleanField(mailMap[MESSAGE_ID])
	mail.Date = parseDate(cleanField(mailMap[DATE]))
	mail.From = cleanField(mailMap[FROM])
	mail.To = cleanField(mailMap[TO])
	mail.Subject = cleanField(mailMap[SUBJECT])
	mail.Cc = cleanField(mailMap[CC])
	mail.Mime_Version = cleanField(mailMap[MIME_VERSION])
	mail.Content_Type = cleanField(mailMap[CONTENT_TYPE])
	mail.Content_Transfer_Encoding = cleanField(mailMap[CONTENT_TRANSFER_ENCODING])
	mail.Bcc = cleanField(mailMap[BCC])
	mail.X_From = cleanField(mailMap[X_FROM])
	mail.X_To = cleanField(mailMap[X_TO])
	mail.X_cc = cleanField(mailMap[X_CC])
	mail.X_bcc = cleanField(mailMap[X_BCC])
	mail.X_Folder = cleanField(mailMap[X_FOLDER])
	mail.X_Origin = cleanField(mailMap[X_ORIGIN])
	mail.X_FileName = cleanField(mailMap[X_FILENAME])
	mail.Content = body

	return mail
}
