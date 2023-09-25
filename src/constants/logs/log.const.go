// contains variables with message errors for log files
package constants_log

const (
	FILE_NAME_ERROR_GENERAL  = "log_err_general"
	FILE_NAME_ERROR_PARSER   = "log_err_parser"
	FILE_NAME_ERROR_BULKER   = "log_err_bulker"
	FILE_NAME_ERROR_INDEXER  = "log_err_indexer"
	FILE_NAME_ERROR_DATABASE = "log_err_database"

	OPERATION_PARSER     = "parser"
	OPERATION_BULKER     = "bulker"
	OPERATION_DATABASE   = "database"
	OPERATION_OPEN_FILE  = "open file"
	OPERATION_LIST_FILES = "list files"

	ERROR_PARSER_FAILED    = "could not parse file"
	ERROR_BULKER_FAILED    = "data bulker failed"
	ERROR_DATA_BASE        = "there was an error querying the database"
	ERROR_CREATE_BASE      = "There was an error creating the database"
	ERROR_CREATE_LOG       = "log file could not be created"
	ERROR_JSON_PARSE       = "could not parse to json"
	ERROR_OPEN_FILE        = "could not open file"
	ERROR_NOT_IS_MIME_FILE = "the file is not a mime file"
)
