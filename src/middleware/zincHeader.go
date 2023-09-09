package myMiddleware

import (
	"net/http"
	"os"
)

func ZincHeader(req *http.Request) {

	req.SetBasicAuth(os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36")

}
