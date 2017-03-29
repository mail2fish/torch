package scene_server

import (
	"sync"
	"time"
	"torch/game_server"

	"gopkg.in/mgo.v2/bson"
)

type CmdType uint16

type Scener interface {
	Id() bson.ObjectId

	// Start is called to start  running GameObject in the Scener
	Start(game_server.Responser, ...interface{}) (bson.ObjectId, error)

	// End is called to stop the Scener
	End(game_server.Responser) (bool, error)

	// 检查是否应该结束游戏场景
	WhetherEndScene() bool

	ProcessCmd(cmd *Cmd)
	ProcessTick(time.Time)
	ProcessSignal(int)
}

type Cmd struct {
	SentAt      int64
	ReceivedAt  int64
	PushCmdAt   int64
	ProcessedAt int64
	ClientId    bson.ObjectId
	Paramater   interface{}
}

type scene struct {
	id        bson.ObjectId
	rwMu      sync.RWMutex
	state     State
	startedAt time.Time
	endedAt   time.Time
	scener    Scener
	chCmd     chan *Cmd
	chSignal  chan int
	chTick    chan time.Time
}

func (s *scene) pushCmd(cmd *Cmd) {
	cmd.PushCmdAt = time.Now().Unix()
	s.chCmd <- cmd
}

func (s *scene) end(rp game_server.Responser) (bool, error) {
	s.rwMu.Lock()
	defer s.rwMu.Unlock()
	s.state = ENDED
	s.endedAt = time.Now().UTC()
	return s.scener.End(rp)

}
func (s *scene) State() State {
	s.rwMu.RLock()
	defer s.rwMu.Unlock()
	return s.state
}
func (s *scene) running() bool {
	s.rwMu.RLock()
	defer s.rwMu.RUnlock()
	if s.state == RUNNING {
		return true
	} else {
		return false
	}

}
func (s *scene) tickLoop() {
	for s.running() && !s.scener.WhetherEndScene() {
		select {
		case <-time.After(100 * time.Millisecond):
			s.chTick <- time.Now()

		}
	}
}
func (s *scene) loop() {
	if s.state == STARTED {
		s.rwMu.Lock()
		s.state = RUNNING
		s.rwMu.Unlock()

		for s.running() && !s.scener.WhetherEndScene() {
			select {
			case tick := <-s.chTick:
				s.scener.ProcessTick(tick)

			case cmd := <-s.chCmd:
				s.scener.ProcessCmd(cmd)

			case signal := <-s.chSignal:
				s.scener.ProcessSignal(signal)

			}

		}

	}
}
func startScene(sr Scener, rp game_server.Responser, params ...interface{}) (*scene, error) {
	id, err := sr.Start(rp, params...)
	if err == nil {
		scene := &scene{scener: sr,
			startedAt: time.Now().UTC(),
			state:     STARTED,
			id:        id,
			chCmd:     make(chan *Cmd),
			chSignal:  make(chan int),
			chTick:    make(chan time.Time),
		}

		go scene.loop()
		go scene.tickLoop()
		return scene, nil

	}
	return nil, err
}
