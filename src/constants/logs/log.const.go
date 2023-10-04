package constants_log

//  log error messages and api responses
const (
	FILE_NAME_ERROR_GENERAL   = "log_err_general"
	FILE_NAME_ERROR_PARSER    = "log_err_parser"
	FILE_NAME_ERROR_BULKER    = "log_err_bulker"
	FILE_NAME_ERROR_INDEXER   = "log_err_indexer"
	FILE_NAME_ERROR_DATABASE  = "log_err_database"
	FILE_NAME_ERROR_PROFILING = "log_err_profiling"

	OPERATION_PARSER     = "parser"
	OPERATION_BULKER     = "bulker"
	OPERATION_DATABASE   = "database"
	OPERATION_OPEN_FILE  = "open file"
	OPERATION_LIST_FILES = "list files"
	OPERATION_PROFILING  = "profiling"

	ERROR_PARSER_FAILED        = "could not parse file"
	ERROR_BULKER_FAILED        = "data bulker failed"
	ERROR_DATA_BASE            = "there was an error querying the database"
	ERROR_CREATE_BASE          = "There was an error creating the database"
	ERROR_CREATE_LOG           = "log file could not be created"
	ERROR_MEM_PROFILING_CREATE = "mem profiling file could not be created"
	ERROR_MEM_PROFILING_WRITE  = "could not write memory profile"
	ERROR_CPU_PROFILING_CREATE = "CPU profiling file could not be created"
	ERROR_CPU_PROFILING_START  = "could not start CPU profile"
	ERROR_JSON_PARSE           = "could not parse to json"
	ERROR_OPEN_FILE            = "could not open file"
	ERROR_NOT_IS_MIME_FILE     = "the file is not a mime file"
)
