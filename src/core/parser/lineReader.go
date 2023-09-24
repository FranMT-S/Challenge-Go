package parser

import (
	"strings"
	"sync"
)

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

func newLineByLineReaderAsync() *lineByLineReaderAsync {
	return &lineByLineReaderAsync{lock: &sync.Mutex{}, headLineFlag: -1}
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
	return &lineByLineReader{mailMap: make(map[string]string)}
}

func (lineReader lineByLineReader) getMapData() map[string]string {
	return lineReader.mailMap
}

func (lineReader *lineByLineReader) Read(line string) {
	if lineReader.mailMap[X_FILENAME] != "" {
		lineReader.mailMap[CONTENT] += line
	} else if strings.HasPrefix(line, X_FROM) && lineReader.mailMap[X_FROM] == "" {
		lineReader.mailMap[X_FROM] = line[len(X_FROM):]
		lineReader.beforeLecture = X_FROM
	} else if strings.HasPrefix(line, X_TO) && lineReader.mailMap[X_TO] == "" {
		lineReader.mailMap[X_TO] = line[len(X_TO):]
		lineReader.beforeLecture = X_TO
	} else if strings.HasPrefix(line, X_CC) && lineReader.mailMap[X_CC] == "" {
		lineReader.mailMap[X_CC] = line[len(X_CC):]
		lineReader.beforeLecture = X_CC
	} else if strings.HasPrefix(line, X_BCC) && lineReader.mailMap[X_BCC] == "" {
		lineReader.mailMap[X_BCC] = line[len(X_BCC):]
		lineReader.beforeLecture = X_BCC
	} else if strings.HasPrefix(line, X_FOLDER) && lineReader.mailMap[X_FOLDER] == "" {
		lineReader.mailMap[X_FOLDER] = line[len(X_FOLDER):]
		lineReader.beforeLecture = X_FOLDER
	} else if strings.HasPrefix(line, X_ORIGIN) && lineReader.mailMap[X_ORIGIN] == "" {
		lineReader.mailMap[X_ORIGIN] = line[len(X_ORIGIN):]
		lineReader.beforeLecture = X_ORIGIN
	} else if strings.HasPrefix(line, X_FILENAME) && lineReader.mailMap[X_FILENAME] == "" {
		lineReader.mailMap[X_FILENAME] = line[len(X_FILENAME):]
		lineReader.beforeLecture = X_FILENAME
	} else if strings.HasPrefix(line, MESSAGE_ID) && lineReader.mailMap[MESSAGE_ID] == "" {
		lineReader.mailMap[MESSAGE_ID] = line[len(MESSAGE_ID):]
		lineReader.beforeLecture = MESSAGE_ID
	} else if strings.HasPrefix(line, DATE) && lineReader.mailMap[DATE] == "" {
		lineReader.mailMap[DATE] = line[len(DATE):]
		lineReader.beforeLecture = DATE
	} else if strings.HasPrefix(line, FROM) && lineReader.mailMap[FROM] == "" {
		lineReader.mailMap[FROM] = line[len(FROM):]
		lineReader.beforeLecture = FROM
	} else if strings.HasPrefix(line, TO) && lineReader.mailMap[TO] == "" {
		lineReader.mailMap[TO] = line[len(TO):]
		lineReader.beforeLecture = TO
	} else if strings.HasPrefix(line, SUBJECT) && lineReader.mailMap[SUBJECT] == "" {
		lineReader.mailMap[SUBJECT] = line[len(SUBJECT):]
		lineReader.beforeLecture = SUBJECT
	} else if strings.HasPrefix(line, CC) && lineReader.mailMap[CC] == "" {
		lineReader.mailMap[CC] = line[len(CC):]
		lineReader.beforeLecture = CC
	} else if strings.HasPrefix(line, BCC) && lineReader.mailMap[BCC] == "" {
		lineReader.mailMap[BCC] = line[len(BCC):]
		lineReader.beforeLecture = BCC
	} else if strings.HasPrefix(line, MIME_VERSION) && lineReader.mailMap[MIME_VERSION] == "" {
		lineReader.mailMap[MIME_VERSION] = line[len(MIME_VERSION):]
		lineReader.beforeLecture = MIME_VERSION
	} else if strings.HasPrefix(line, CONTENT_TYPE) && lineReader.mailMap[CONTENT_TYPE] == "" {
		lineReader.mailMap[CONTENT_TYPE] = line[len(CONTENT_TYPE):]
		lineReader.beforeLecture = CONTENT_TYPE
	} else if strings.HasPrefix(line, CONTENT_TRANSFER_ENCODING) && lineReader.mailMap[CONTENT_TRANSFER_ENCODING] == "" {
		lineReader.mailMap[CONTENT_TRANSFER_ENCODING] = line[len(CONTENT_TRANSFER_ENCODING):]
		lineReader.beforeLecture = CONTENT_TRANSFER_ENCODING
	} else if lineReader.beforeLecture != "" {
		lineReader.mailMap[lineReader.beforeLecture] += line
	}
}

type lineByLineReaderAsync struct {
	line         *lineMail
	lock         *sync.Mutex
	headLineFlag int
}

func (lineReader *lineByLineReaderAsync) Read(line *lineMail) {

	if lineReader.headLineFlag > 0 && lineReader.headLineFlag < line.numberLine {
		line.data = line.lineToAnalize
		line.field = CONTENT
	} else if strings.HasPrefix(line.lineToAnalize, X_FROM) {
		line.data = line.lineToAnalize[len(X_FROM):]
		line.field = X_FROM
	} else if strings.HasPrefix(line.lineToAnalize, X_TO) {
		line.data = line.lineToAnalize[len(X_TO):]
		line.field = X_TO
	} else if strings.HasPrefix(line.lineToAnalize, X_CC) {
		line.data = line.lineToAnalize[len(X_CC):]
		line.field = X_CC
	} else if strings.HasPrefix(line.lineToAnalize, X_BCC) {
		line.data = line.lineToAnalize[len(X_BCC):]
		line.field = X_BCC
	} else if strings.HasPrefix(line.lineToAnalize, X_FOLDER) {
		line.data = line.lineToAnalize[len(X_FOLDER):]
		line.field = X_FOLDER
	} else if strings.HasPrefix(line.lineToAnalize, X_ORIGIN) {
		line.data = line.lineToAnalize[len(X_ORIGIN):]
		line.field = X_ORIGIN
	} else if strings.HasPrefix(line.lineToAnalize, X_FILENAME) {
		line.data = line.lineToAnalize[len(X_FILENAME):]
		line.field = X_FILENAME
		lineReader.headLineFlag = line.numberLine
	} else if strings.HasPrefix(line.lineToAnalize, MESSAGE_ID) {
		line.data = line.lineToAnalize[len(MESSAGE_ID):]
		line.field = MESSAGE_ID
	} else if strings.HasPrefix(line.lineToAnalize, DATE) {
		line.data = line.lineToAnalize[len(DATE):]
		line.field = DATE
	} else if strings.HasPrefix(line.lineToAnalize, FROM) {
		line.data = line.lineToAnalize[len(FROM):]
		line.field = FROM
	} else if strings.HasPrefix(line.lineToAnalize, TO) {
		line.data = line.lineToAnalize[len(TO):]
		line.field = TO
	} else if strings.HasPrefix(line.lineToAnalize, SUBJECT) {
		line.data = line.lineToAnalize[len(SUBJECT):]
		line.field = SUBJECT
	} else if strings.HasPrefix(line.lineToAnalize, CC) {
		line.data = line.lineToAnalize[len(CC):]
		line.field = CC
	} else if strings.HasPrefix(line.lineToAnalize, BCC) {
		line.data = line.lineToAnalize[len(BCC):]
		line.field = BCC
	} else if strings.HasPrefix(line.lineToAnalize, MIME_VERSION) {
		line.data = line.lineToAnalize[len(MIME_VERSION):]
		line.field = MIME_VERSION
	} else if strings.HasPrefix(line.lineToAnalize, CONTENT_TYPE) {
		line.data = line.lineToAnalize[len(CONTENT_TYPE):]
		line.field = CONTENT_TYPE
	} else if strings.HasPrefix(line.lineToAnalize, CONTENT_TRANSFER_ENCODING) {
		line.data = line.lineToAnalize[len(CONTENT_TRANSFER_ENCODING):]
		line.field = CONTENT_TRANSFER_ENCODING
	} else {
		line.data = line.lineToAnalize
		line.field = K_FATHER
	}

}

func (lineReader *lineByLineReaderAsync) getMapData() map[string]string {
	temp := lineReader.line
	mailMap := make(map[string]string)
	for {
		if temp == nil {
			break
		}

		if lineReader.headLineFlag < temp.numberLine {
			mailMap[CONTENT] = temp.lineToAnalize + mailMap[CONTENT]
		} else {
			mailMap[temp.getField()] = temp.data + mailMap[temp.getField()]
		}

		temp = temp.lineFather
	}
	return mailMap
}
