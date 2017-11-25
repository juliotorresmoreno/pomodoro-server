package app

import (
	"net/http"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/juliotorresmoreno/pomodoro-server/db"
	"github.com/juliotorresmoreno/pomodoro-server/router"
)

type Server struct {
	router http.Handler
}

func (server Server) Start() {
	t := sync.WaitGroup{}
	t.Add(2)
	db.Sync()
	go func() {
		if err := http.ListenAndServe(":8080", server.router); err != nil {
			log.Fatal(err)
		}
	}()
	log.Info("Server listen on 0.0.0.0:8080")
	t.Wait()
}

func NewServer() Server {
	handler := router.NewRouter()
	return Server{
		router: handler,
	}
}
