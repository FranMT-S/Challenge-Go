package parser

import (
	"bufio"
	"io"
	"os"
	"regexp"
	"strings"
	"sync"

	model "github.com/FranMT-S/Challenge-Go/src/model"
)

const (
	MAX_CONCURRENT_LINES = 25
)

/*
IParser Mail
Proporciona el metodo para transformar un archivo a un formato de Correo.
*/
type IParserMail interface {
	Parse(file *os.File) (*model.Mail, error)
}

/*
--------------------
Parseador con Normal
--------------------

Lee linea por linea y asigna el contenido al correo
*/

type parserBasic struct{}

func NewParserBasic() *parserBasic {
	return &parserBasic{}
}

func (parser parserBasic) Parse(file *os.File) (*model.Mail, error) {
	// buf := make([]byte, 1024)
	var mail *model.Mail
	var mailMap map[string]string
	lineByLineReader := newLineByLineReader()
	reader := bufio.NewReader(file)
	// beforeLine := ""
	for {
		lineByte, err := reader.ReadBytes('\n')
		line := string(lineByte)
		if err != nil && len(line) <= 0 {

			if err != io.EOF {
				return nil, err

			}
			break
		}

		lineByLineReader.Read(line)

	}

	mailMap = lineByLineReader.getMapData()
	mail, err := mailFroMap(mailMap)
	if err != nil {
		return nil, err
	}

	return mail, nil
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

func (parser parserAsync) Parse(file *os.File) (*model.Mail, error) {
	// buf := make([]byte, 1024)
	var mail *model.Mail
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
				return nil, err
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
	mail, err := mailFroMap(mailMap)
	if err != nil {
		return nil, err
	}

	return mail, nil
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

func (parser ParserAsyncRegex) Parse(file *os.File) (*model.Mail, error) {

	var mail *model.Mail
	var wg sync.WaitGroup
	var semaphore chan struct{}
	var mutex = &sync.Mutex{}

	mailMap := make(map[string]string)
	indexMap := make(map[int]string)
	noMatchMap := make(map[int]string)
	i := -1 // counter for line index

	if parser.maxConcurrent > MAX_CONCURRENT_LINES {
		parser.maxConcurrent = MAX_CONCURRENT_LINES
	} else if parser.maxConcurrent <= 0 {
		parser.maxConcurrent = 1
	}

	semaphore = make(chan struct{}, parser.maxConcurrent)

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	content := string(bytes)

	re, _ := regexp.Compile(`(\r\n){2,}|\n{2,}`)
	reLine, _ := regexp.Compile(`^([\w-_]+:)(.+)`)

	match := re.Split(content, 2)

	header := match[0]
	body := match[1]

	mailMap[CONTENT] = body

	dataReader := strings.NewReader(header)
	reader := bufio.NewReader(dataReader)

	for {
		lineByte, err := reader.ReadBytes('\n')
		line := string(lineByte)
		i++
		indexLine := i
		if err != nil && len(line) <= 0 {
			if err != io.EOF {

				return nil, err

			}
			break
		}

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

	// Corregis los campos que no hicieron match
	for j := 0; j <= i; j++ {

		if indexMap[j] == "" {
			indexMap[j] = indexMap[j-1]
			mailMap[indexMap[j]] += noMatchMap[j]
		}
	}

	mail, err = mailFroMap(mailMap)
	if err != nil {
		return nil, err
	}

	return mail, nil
}
