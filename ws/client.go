package ws

import (
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
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
func (hub *Hub) ServeWs(w http.ResponseWriter, r *http.Request, session string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		terror.Log(err)
		return
	}
	if _, ok := hub.clients[session]; ok {
		hub.clients[session] = &user{
			clients: make([]*connection, 0),
		}
	}
	client := &connection{
		Conn: conn,
	}
	hub.clients[session].clients = append(hub.clients[session].clients, client)
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
