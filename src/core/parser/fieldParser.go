package parser

// contains methods that help clean and format email data.

import (
	"fmt"
	"strings"
	"time"

	constants_log "github.com/FranMT-S/Challenge-Go/src/constants/logs"
	model "github.com/FranMT-S/Challenge-Go/src/model"
)

// cleanField cleans line breaks, tabs, backspaces at the beginning and end of each line.
//
// only must used in field header of mails.
func cleanField(s string) string {
	s = strings.ReplaceAll(s, "\r", "")
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "\t", "")
	s = strings.TrimSpace(s)
	return s
}

// parseDate if could not convert string to date, return empty string and error
func parseDate(s string) (date string, err error) {
	t, err := time.Parse("Mon, _2 Jan 2006 15:04:05 -0700 (MST)", s)
	if err != nil {
		return "", err
	}

	return t.Format("2006-01-02T15:04:05Z"), nil
}

// mailFroMap return a mail from a map structure containing keys of the mail fields
//
// if failed return a mail=nil and the error
func mailFroMap(_map map[string]string) (*model.Mail, error) {

	mail := new(model.Mail)

	mail.Mime_Version = cleanField(_map[MIME_VERSION])
	if mail.Mime_Version == "" {
		return nil, fmt.Errorf("%v", constants_log.ERROR_NOT_IS_MIME_FILE)
	}

	date, err := parseDate(cleanField(_map[DATE]))
	if err != nil {
		return nil, err
	}

	mail.Date = date
	mail.Message_ID = cleanField(_map[MESSAGE_ID])
	mail.From = cleanField(_map[FROM])
	mail.To = cleanField(_map[TO])
	mail.Subject = cleanField(_map[SUBJECT])
	mail.Cc = cleanField(_map[CC])
	mail.Content_Type = cleanField(_map[CONTENT_TYPE])
	mail.Content_Transfer_Encoding = cleanField(_map[CONTENT_TRANSFER_ENCODING])
	mail.Bcc = cleanField(_map[BCC])
	mail.X_From = cleanField(_map[X_FROM])
	mail.X_To = cleanField(_map[X_TO])
	mail.X_cc = cleanField(_map[X_CC])
	mail.X_bcc = cleanField(_map[X_BCC])
	mail.X_Folder = cleanField(_map[X_FOLDER])
	mail.X_Origin = cleanField(_map[X_ORIGIN])
	mail.X_FileName = cleanField(_map[X_FILENAME])
	mail.Content = _map[CONTENT]

	return mail, nil
}
