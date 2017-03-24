package game_server

import (
	"sync"
	"time"

	"golang.org/x/net/context"
)

type RequestSession struct {
	context.Context
	requestedAt       time.Time
	requestName       string
	connectionSession *ConnectionSession
	gameClienter      GameClienter
	sn                uint32
	actived           bool

	rwMU sync.RWMutex
}

func (rs *RequestSession) Name() string {
	return rs.requestName
}

func (rs *RequestSession) SN() uint32 {
	return rs.sn
}
func (rs *RequestSession) RequestName(uint32) string {
	return rs.requestName
}
func (rs *RequestSession) WithValue(key, val interface{}) {
	rs.rwMU.Lock()
	defer rs.rwMU.Unlock()
	if rs.Context == nil {
		rs.Context = context.Background()
	}
	rs.Context = context.WithValue(rs.Context, key, val)

}
func (rs *RequestSession) GameClienter() GameClienter {
	return rs.gameClienter
}

func (rs *RequestSession) ConnectionSession() *ConnectionSession {
	return rs.connectionSession
}

func (rs *RequestSession) Close() {
	rs.rwMU.Lock()
	defer rs.rwMU.Unlock()
	rs.actived = false
}

type ConnectionSession struct {
	context.Context
	rwMU         sync.RWMutex
	requests     map[uint32]*RequestSession // sequence number as a key
	sn           uint32
	gameClienter GameClienter
	connectedAt  time.Time
	responser    Responser
}

func (ms *ConnectionSession) GameClienter() GameClienter {
	return ms.gameClienter

}

func (ms *ConnectionSession) StartRequest(name string) (uint32, *RequestSession) {
	bs := &RequestSession{
		requestedAt:       time.Now().UTC(),
		requestName:       name,
		connectionSession: ms,
		actived:           true,
		gameClienter:      ms.gameClienter,
	}

	ms.rwMU.Lock()
	defer ms.rwMU.Unlock()
	ms.sn++
	bs.sn = ms.sn
	if ms.requests == nil {
		ms.requests = make(map[uint32]*RequestSession)
	}
	ms.requests[ms.sn] = bs
	return ms.sn, bs
}
func (ms *ConnectionSession) EndSession() {
	ms.rwMU.Lock()
	defer ms.rwMU.Unlock()

	for _, rs := range ms.requests {
		rs.Close()
	}
}
func (ms *ConnectionSession) EndRequest(sn uint32) {
	ms.rwMU.Lock()
	defer ms.rwMU.Unlock()

	if s, ok := ms.requests[sn]; ok {
		s.connectionSession = nil
		s.gameClienter = nil
	}
	delete(ms.requests, sn)

}

func (ms *ConnectionSession) Responser() Responser {
	return ms.responser
}

func (ms *ConnectionSession) ConnectedAt() time.Time {
	return ms.connectedAt
}
func (ms *ConnectionSession) RemoteIP() string {
	return ""
}

func (cs *ConnectionSession) WithValue(key, val interface{}) {
	cs.rwMU.Lock()
	defer cs.rwMU.Unlock()
	if cs.Context == nil {
		cs.Context = context.Background()
	}
	cs.Context = context.WithValue(cs.Context, key, val)

}

type FillFunSession struct {
	context.Context
}
