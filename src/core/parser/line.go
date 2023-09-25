package parser

import "sync"

/*
represents a line that will be analyzed by an ILineReader
parameters:
  - lineFather: previous line read
  - field: field of the email to which it belongs.
  - data: information that will be saved in the email.
  - lineToAnalize: field of the email to which it belongs.
  - lineNumber: line number in the file.
  - lock: ensures the correct assignment of the father when working with threads.
*/
type lineMail struct {
	lineFather    *lineMail
	field         string
	data          string
	lineToAnalize string
	lineNumber    int
	lock          *sync.Mutex
}

/*
Gets the field in which the information should be assigned
*/
func (lineMail *lineMail) getField() string {
	lineMail.lock.Lock()
	defer lineMail.lock.Unlock()

	if lineMail.field == K_FATHER {
		lineMail.field = lineMail.lineFather.getField()
	}

	return lineMail.field
}

func newLineMail(_fatherLine *lineMail, _lineToAnalize string, _numberLine int) *lineMail {
	return &lineMail{lineFather: _fatherLine, lineToAnalize: _lineToAnalize, lineNumber: _numberLine, lock: &sync.Mutex{}}
}
