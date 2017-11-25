package ws

import (
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/juliotorresmoreno/pomodoro-server/models"
	"github.com/pingcap/tidb/terror"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	maxMessageSize = 8
)

type connection struct {
	*websocket.Conn
}

type user struct {
	clients []*connection
}

// ServeWs aca es donde establenemos la conexion websocket con el usuario
func (hub *Hub) ServeWs(w http.ResponseWriter, r *http.Request, session models.Session) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		terror.Log(err)
		return
	}
	username := session.Username
	if _, ok := hub.clients[username]; ok {
		hub.clients[username] = &user{
			clients: make([]*connection, 0),
		}
	}
	client := &connection{
		Conn: conn,
	}
	hub.clients[username].clients = append(hub.clients[username].clients, client)
	client.Listen()
}

//Listen
func (c *connection) Listen() {
	for {
		_, _, err := c.ReadMessage()
		if err != nil {
			break
		}
		log.Info("conn close")
	}
}
