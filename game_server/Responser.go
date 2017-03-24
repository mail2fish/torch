package game_server

import (
	"sync"

	"gopkg.in/mgo.v2/bson"
)

type ResponsKind int

const RK_WEBSOCKET ResponsKind = 1
const RK_SCENE ResponsKind = 2

type Responser interface {
	Write(ResponsKind, ...interface{})
}

type SuperResponser struct {
	rwMU sync.RWMutex

	gcResponsers                  map[bson.ObjectId]*GameClienterResponser
	closedGameClienterResponserCh chan bson.ObjectId
}

func (sr *SuperResponser) GameClienterResponser(gc GameClienter) Responser {
	gr := &GameClienterResponser{gc: gc, sr: sr}

	gr.chData = make(chan interface{}, 2)
	go gr.loopWrite()

	sr.rwMU.Lock()
	defer sr.rwMU.Unlock()
	sr.gcResponsers[gc.Id()] = gr
	return gr
}

func (sr *SuperResponser) Write(kind ResponsKind, data ...interface{}) {
	switch kind {
	case RK_WEBSOCKET:
		if len(data) >= 2 {
			switch data[0].(type) {
			case bson.ObjectId:
				if gr, ok := sr.getGameClienterResponser(data[0].(bson.ObjectId)); ok {
					for i, v := range data {
						if i == 0 {
							continue
						}
						gr.Write(RK_WEBSOCKET, v)
					}
				}
			}
		}
		// sr.gcResponser.Write(kind, data...)
	}

}

func (sr *SuperResponser) getGameClienterResponser(id bson.ObjectId) (*GameClienterResponser, bool) {
	sr.rwMU.RLock()
	defer sr.rwMU.RUnlock()
	gc, ok := sr.gcResponsers[id]
	return gc, ok
}

func (sr *SuperResponser) loopRemoveGameClienterResponser() {
	for id := range sr.closedGameClienterResponserCh {
		func(sr *SuperResponser) {
			sr.rwMU.Lock()
			defer sr.rwMU.Unlock()
			delete(sr.gcResponsers, id)
		}(sr)

	}
}
func newSueperResponser() *SuperResponser {

	sr := &SuperResponser{}
	sr.gcResponsers = make(map[bson.ObjectId]*GameClienterResponser)
	sr.closedGameClienterResponserCh = make(chan bson.ObjectId)

	go sr.loopRemoveGameClienterResponser()

	return sr
}

var SResponser *SuperResponser

func init() {
	SResponser = newSueperResponser()
}
