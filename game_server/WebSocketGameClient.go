package game_server

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"github.com/go-playground/log"

	"gopkg.in/mgo.v2/bson"
)

// WebSocketGameClient  Implementation of GameClienter
type WebSocketGameClient struct {
	mu                sync.RWMutex
	id                bson.ObjectId
	role              GameRole
	requestDataReader *RequestDataReader
	conn              *websocket.Conn
	sev               *GameServer
	connSession       *ConnectionSession
	available         bool
}

func (ws *WebSocketGameClient) Close() {
	ws.mu.Lock()
	defer ws.mu.Unlock()
	ws.available = false

	ws.requestDataReader = nil
	ws.sev = nil

	if err := ws.conn.Close(); err != nil {
		log.Error("Failed to disconnect web socket connection", err)

	}

	ws.conn = nil

}

func (ws *WebSocketGameClient) Id() bson.ObjectId {
	return ws.id
}
func (ws *WebSocketGameClient) Equal(id bson.ObjectId) bool {
	ws.mu.Lock()
	defer ws.mu.Unlock()
	if ws.id == id {
		return true
	}

	return false

}
func (ws *WebSocketGameClient) ConnectionSession() *ConnectionSession {
	return ws.connSession
}

func (ws *WebSocketGameClient) Role() GameRole {
	return ws.role
}
func (ws *WebSocketGameClient) SetRole(role GameRole) {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	ws.role = role
}

func (ws *WebSocketGameClient) EqualRole(role GameRole) bool {

	if ws.role.Equal(role) {
		return true
	}

	return false

}

func (ws *WebSocketGameClient) Write(bytes []byte) error {
	// defer func() {
	// 	if r := recover(); r != nil {
	// 		fmt.Println("Recovered in f", r)
	// 	}
	// }()

	if ws.Available() {
		if err := ws.conn.WriteMessage(websocket.BinaryMessage, bytes); err != nil {
			ws.Close()
			log.Error("Failed to write to websocket. ", err)
			return err
		}
	}
	return nil
}

func (ws *WebSocketGameClient) RemoteAddr() net.Addr {
	return ws.conn.RemoteAddr()
}

func (ws *WebSocketGameClient) Available() bool {
	ws.mu.RLock()
	defer ws.mu.RUnlock()
	return ws.available
}

func NewWebSocketGameClient(c *websocket.Conn, sev *GameServer) *WebSocketGameClient {

	ws := &WebSocketGameClient{
		id:                bson.NewObjectId(),
		conn:              c,
		sev:               sev,
		requestDataReader: NewRequestDataReader(sev),
		available:         true,
	}
	ws.connSession = &ConnectionSession{gameClienter: ws, connectedAt: time.Now().UTC(), responser: SResponser.GameClienterResponser(ws)}

	go ws.read()

	return ws
}

func (ws *WebSocketGameClient) read() {
	if ws.Available() {
		for {
			ws.conn.SetReadDeadline(time.Now().Add(6 * time.Minute))
			if typ, bytes, err := ws.conn.ReadMessage(); err == nil {
				if typ == websocket.BinaryMessage {
					fmt.Println("TEST TTTT")
					ws.requestDataReader.Append(bytes)
					for _, pack := range ws.requestDataReader.Read() {
						cr := &ClientRequest{Package: pack, Client: ws}
						ws.sev.Post(cr)
					}

				}

			} else {
				ws.Close()
				log.Debug("conn closed for read error")
				break
			}

		}

	}
}
