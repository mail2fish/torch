package game_server

import (
	"net"

	"gopkg.in/mgo.v2/bson"
)

type GameRole interface {
	Equal(GameRole) bool
}

type GameClienter interface {
	Close()
	Write([]byte) error // it should not be called directly
	Id() bson.ObjectId
	Equal(bson.ObjectId) bool

	Available() bool

	Role() GameRole
	EqualRole(GameRole) bool
	SetRole(GameRole)

	RemoteAddr() net.Addr
	ConnectionSession() *ConnectionSession
}
