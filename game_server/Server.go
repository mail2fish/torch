package game_server

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"sync"

	goerrors "github.com/go-errors/errors"
	"github.com/go-playground/log"
	"github.com/golang/protobuf/proto"
)

type FunMap struct {
	ID         uint16
	Name       string
	Fun        reflect.Value
	ParamType  reflect.Type
	SingleMode bool
}

type GameServer struct {
	muRW    sync.RWMutex
	clients []GameClienter

	requests map[uint16]FunMap

	serialHandlers     []uint16
	concurrentHandlers []uint16
	middlewares        []Middlewarer
	requestRecevier    chan *ClientRequest
	wg                 *sync.WaitGroup
	runningStae        bool
}

func NewGameServer(wg *sync.WaitGroup) *GameServer {
	gs := &GameServer{wg: wg}

	gs.requests = make(map[uint16]FunMap)
	gs.middlewares = []Middlewarer{}
	gs.requestRecevier = make(chan *ClientRequest, 100)

	gs.wg.Add(1)

	gs.start()
	return gs
}

func (s *GameServer) start() {
	if !s.runningStae {
		s.runningStae = true
		for i := 0; i < runtime.NumCPU(); i++ {
			go s.Deliver()
		}

	}
}
func (s *GameServer) Shutdown() {
	s.runningStae = false
	s.wg.Done()
}

func (s *GameServer) HanderFunMap() map[uint16]FunMap {
	return s.requests
}

func (s *GameServer) Post(pack *ClientRequest) {
	s.requestRecevier <- pack
}
func rescue() {
	if r := recover(); r != nil {
		msg := goerrors.Wrap(r, 2).ErrorStack()
		log.Panic("Deliver request panic", msg)
	}
}
func (s *GameServer) Deliver() {
	defer rescue()
	s.wg.Add(1)
	defer s.wg.Done()
	log.Info("Game server start to deliver package.")
	defer log.Info("Game server has stoped to deliver package.")

	for rq := range s.requestRecevier {
		if !s.runningStae {
			return
		}
		go func(rq *ClientRequest, s *GameServer) {
			session := rq.Client.ConnectionSession()
			pack := rq.Package

			if param, err := s.getRequestParamater(pack.HandlerId, pack.Data); err == nil {
				handerId := pack.HandlerId

				if rH, ok := s.requests[handerId]; ok {
					handler := rH.Fun
					sn, rs := session.StartRequest(rH.Name)
					defer session.EndRequest(sn)

					log.Info("Dispatch a pack to handler: ", handerId, " name: ", s.requests[handerId].Name)

					middlewareIndexer := &MiddlewareIndexer{middlewares: s.middlewares, index: 0}
					middlewareIndexer.Call(rs, session.Responser(), rq.Package, param, handler)
					// middlewareIndexer.Call(nil, rq.Client, rq.Package, param, handler)
				} else {
					log.Notice("No Handler for the handle id: ", handerId)
				}

			}

		}(rq, s)

	}
}
func (s *GameServer) getRequestParamater(id uint16, bytes []byte) (interface{}, error) {
	rqParamsType := s.requests[id].ParamType
	if rqParamsType == nil {
		log.Error("Could not find a param to match the handler id: ", id)
		return nil, errors.New("Could not find a param to match the handler id")
	}
	var rqParams interface{}

	if rqParamsType.Kind() == reflect.Ptr {
		rqParams = reflect.New(rqParamsType.Elem()).Interface()
	}
	if err := proto.Unmarshal(bytes, rqParams.(proto.Message)); err != nil {
		log.Error("Unmarshal user data error. Handler Id: ", id, ", name: ", s.requests[id].Name, ", Error:", err)
		return nil, errors.New("Unmarshal request paramater error.")
	}

	return rqParams, nil

}

func (s *GameServer) AddGameClient(c GameClienter) {

	for _, c := range s.clients {
		if c.Equal(c.Id()) {
			return

		}
	}
	s.muRW.Lock()
	defer s.muRW.Unlock()
	s.clients = append(s.clients, c)
}

func (s *GameServer) RemoveGameClient(c GameClienter) {
	s.muRW.Lock()
	defer s.muRW.Unlock()

	for i, c := range s.clients {

		if c.Equal(c.Id()) {
			s.clients = append(s.clients[:i], s.clients[i+1:]...)
			return
		}

	}
}
func (s *GameServer) RequestHandlers(id uint16, name string, handler interface{}, singleMode bool) {

}
func (s *GameServer) BindRequestHandler(id uint16, name string, handler interface{}, singleMode bool) {
	if id > 5000 {
		log.Error("The request handler id should not greater than 5000, the id is %d %s", id, name)
		return
	}
	if _, ok := s.requests[id]; ok {
		log.Error("There were a handler for %d %s", id, name)
	} else {
		fmt.Printf("Binding request handler %d for %s\n", id, name)
		hParam, hFun := ParseRequestHandlerFun(handler)

		s.requests[id] = FunMap{ID: id,
			Name:       name,
			Fun:        hFun,
			ParamType:  hParam,
			SingleMode: singleMode}

		if singleMode {
			s.serialHandlers = append(s.serialHandlers, id)
		} else {
			s.concurrentHandlers = append(s.concurrentHandlers, id)
		}
	}
}

func ParseRequestHandlerFun(handler interface{}) (reflect.Type, reflect.Value) {
	hType := reflect.TypeOf(handler)
	if hType.Kind() != reflect.Func {

		panic("Request handler need be a func")
	}

	numArgs := hType.NumIn()

	var testFunc func(*RequestSession, Responser)
	testType := reflect.TypeOf(testFunc)
	numIn := testType.NumIn()
	if numArgs != numIn+1 {
		panic(fmt.Sprintf("Request handler func need  %d args", numIn+1))
	}

	for i := 0; i < numIn-1; i++ {
		if hType.In(i) != testType.In(i) {
			panic(fmt.Sprintf("Request handler func params  %v should be %v args", hType.In(i), testType.In(i)))
		}

	}

	return hType.In(numIn), reflect.ValueOf(handler)
}

// AppendMiddleware append a middleware to server.
func (s *GameServer) AppendMiddleware(m Middlewarer) {
	s.middlewares = append(s.middlewares, m)
}
