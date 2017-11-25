package controllers

import (
	"fmt"
	"net/http"

	"github.com/juliotorresmoreno/pomodoro-server/models"
	"github.com/juliotorresmoreno/pomodoro-server/ws"
)

//Timer
type Timer struct {
	hub *ws.Hub
}

//NewTimer
func NewTimer(hub *ws.Hub) Timer {
	return Timer{
		hub: hub,
	}
}

func (timer Timer) NewPomodoro(w http.ResponseWriter, r *http.Request, session models.Session) {
	fmt.Fprintf(w, "Hola mundo")
}
