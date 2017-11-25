package ws

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

//Send enviar mensajes a los usuarios
func (hub Hub) Send(user string, mensaje []byte) {

}
