package constants

import (
	"fmt"
	"os"
)

func InitializeVarEnviroment() {

	os.Setenv("INDEX", "mailsTest2")
	os.Setenv("URL", "http://localhost:4080/api/")
	os.Setenv("DB_USER", "admin")
	os.Setenv("DB_PASSWORD", "Complexpass")

	fmt.Println("Variables de entorno establecidas")
}