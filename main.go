package main

import (
	"github.com/juliotorresmoreno/pomodoro-server/app"

	//_ "github.com/mattn/go-sqlite3"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	app := app.NewServer()
	app.Start()
}
