package ws

import (
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
)

//Hub alacen de clientes websocket
type Hub struct {
	clients   map[string]*user
	broadcast chan []byte
}

//NewHub devuelve el hub
func NewHub() *Hub {
	return &Hub{
		clients: make(map[string]*user),
	}
}

//IsConnect devuelve el estado de un usuario, conectado o desconectado.
func (hub Hub) IsConnect(user string) bool {
	usuario, ok := hub.clients[user]
	if ok && len(usuario.clients) > 0 {
		return true
	}
	return false
}

//Remove devuelve el estado de un usuario, conectado o desconectado.
func (hub Hub) Remove(user string, connection *connection) {
	if usuario, ok := hub.clients[user]; ok {
		delete(usuario.clients, connection)
	}
}

//Send enviar mensajes a los usuarios
func (hub Hub) Send(user string, message []byte) {
	log.Infof("Send %v", string(message))
	if usuario, ok := hub.clients[user]; ok {
		for client := range usuario.clients {
			client.WriteMessage(websocket.TextMessage, message)
		}
	}
}

//SendJSON enviar mensajes a los usuarios
func (hub Hub) SendJSON(user string, message interface{}) {
	if usuario, ok := hub.clients[user]; ok {
		for client := range usuario.clients {
			client.WriteJSON(message)
		}
	}
}
