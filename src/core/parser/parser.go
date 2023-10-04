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
	MAX_CONCURRENT_LINES = 15
)

// Provides a method to transform files to emails.
//   - Parse: transform a file in mail
type IParserMail interface {
	Parse(file *os.File) (*model.Mail, error)
}

/*
Parser that implement IParserMail, internally it uses regular expressions to optimize performance, .

  - Maxconcurrent: defines the number of lines of the file that will be read at the same
    time from it with a maximum of 15 and a minimum of 1
*/
type parserAsyncRegex struct {
	maxConcurrent int
}

/*
Return a parserAsyncRegex that implement IParserMail.

Internally it uses regular expressions to optimize performance.

  - Maxconcurrent: defines the number of lines of the file that will be read
    at the same time from it with a maximum of 15 and a minimum of 1
*/
func NewParserAsyncRegex(_maxConcurrent int) *parserAsyncRegex {
	return &parserAsyncRegex{maxConcurrent: _maxConcurrent}
}

// Splits the file into header and body content to streamline the process.
// return a mail if successful.
//
//	if failure return  mail=nil and  error
func (parser parserAsyncRegex) Parse(file *os.File) (*model.Mail, error) {

	var mail *model.Mail
	var wg sync.WaitGroup
	var semaphore chan struct{}
	var mutex = &sync.Mutex{}

	mailMap := make(map[string]string)
	indexMap := make(map[int]string)   // index line
	noMatchMap := make(map[int]string) // for lines that matches fields that contain more than one line
	i := -1                            // counter for index line

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

	re, _ := regexp.Compile(`(\r\n){2,}|\n{2,}`)   // used to split the file
	reLine, _ := regexp.Compile(`^([\w-_]+:)(.+)`) // used to find the fields

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
		// fmt.Println("reading: ", line)
		go func() {
			defer wg.Done()
			match := reLine.FindStringSubmatch(line)
			if len(match) > 0 {
				// match[1] = field
				// match[2] = content line
				mutex.Lock()
				mailMap[match[1]] = match[2]
				indexMap[indexLine] = match[1]
				mutex.Unlock()
			} else {
				// if not match mean that line containt mutilines
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

	// correct the lines that did not match since they are multiline
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

// Parser that implement IParserMail, read the entire file line by line to create the email
type parserBasic struct{}

// Create a parserBasic that Implement IParserMail
//
// read the entire file line by line to create the email
func NewParserBasic() *parserBasic {
	return &parserBasic{}
}

// Read line by line and return a mail if successful.
//
//	if failure return  mail=nil and  error
func (parser parserBasic) Parse(file *os.File) (*model.Mail, error) {

	var mail *model.Mail
	var mailMap map[string]string
	lineByLineReader := newLineByLineReader()

	reader := bufio.NewReader(file)

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

/*
Parser that implement IParserMail , read the entire file and analyze it by reading several lines at the same time.

  - Maxconcurrent: defines the number of lines of the file that will be read at the
    same time from it with a maximum of 15 and a minimum of 1
*/
type parserAsync struct {
	maxConcurrent int
}

/*
return a parserAsync  that implement IParserMail, read the entire file and analyze it by reading several lines at the same time.

  - Maxconcurrent: defines the number of lines of the file that will be read at the
    sametime from it with a maximum of 15 and a minimum of 1
*/
func NewParserAsync(_maxConcurrent int) *parserAsync {
	return &parserAsync{maxConcurrent: _maxConcurrent}
}

// read the entire file and analyze it by reading several lines at the same time
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
			_newLineMail = newLineMail(lineByLineReaderAsync.line, line, lineByLineReaderAsync.line.lineNumber+1)
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
