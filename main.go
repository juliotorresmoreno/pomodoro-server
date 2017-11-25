package main

import "github.com/juliotorresmoreno/pomodoro-server/app"
import _ "github.com/mattn/go-sqlite3"

func main() {
	app := app.NewServer()
	app.Start()
}
