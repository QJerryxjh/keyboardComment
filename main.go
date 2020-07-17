package main

import (
	"keyboardComment/router"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	router.SetupRouter()
}
