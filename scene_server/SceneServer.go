package scene_server

import (
	"errors"
	"sync"
	"torch/game_server"

	"github.com/go-playground/log"

	"gopkg.in/mgo.v2/bson"
)

type State uint16

const STARTED State = 1
const RUNNING State = 2
const ENDED State = 3

type SceneServer struct {
	rwMU   sync.RWMutex
	scenes map[bson.ObjectId]*scene
}

func (ss *SceneServer) StartScene(scener Scener, params ...interface{}) (bson.ObjectId, error) {

	s, err := startScene(scener, game_server.SResponser, params...)
	if err == nil {
		func(ss *SceneServer, s *scene) {
			ss.rwMU.Lock()
			defer ss.rwMU.Unlock()
			id := s.id
			if _, ok := ss.scenes[id]; !ok {
				ss.scenes[id] = s
			} else {
				log.Error("Scener ", id.Hex(), "has been created.")
			}
		}(ss, s)
		return s.id, err

	} else {
		log.Error("Failed to start scener", scener, params)
		return bson.NewObjectId(), err
	}

}

func (ss *SceneServer) EndScene(id bson.ObjectId) (bool, error) {
	ss.rwMU.RLock()
	defer ss.rwMU.RUnlock()
	if scene, ok := ss.scenes[id]; ok {
		return scene.end(game_server.SResponser)
	} else {
		log.Error("Failed to find scener ", id.Hex())
	}
	return false, errors.New("Failed to find scener to end")
}

func (ss *SceneServer) PushCmd(id bson.ObjectId, cmd *Cmd) (bool, error) {
	ss.rwMU.RLock()
	defer ss.rwMU.RUnlock()
	if scene, ok := ss.scenes[id]; ok {
		scene.pushCmd(cmd)
	} else {
		log.Error("Failed to find scener ", id.Hex())
	}
	return false, errors.New("Failed to find scener to push cmd")

}

func (ss *SceneServer) SceneState(id bson.ObjectId) (State, error) {
	ss.rwMU.RLock()
	defer ss.rwMU.RUnlock()
	if scene, ok := ss.scenes[id]; !ok {
		return scene.State(), nil
	} else {
		log.Error("Failed to find scene ", id.Hex())
	}
	return State(0), errors.New("Did not find a scene")
}

var LocalSceneServer *SceneServer

func init() {
	LocalSceneServer = &SceneServer{}
	LocalSceneServer.scenes = make(map[bson.ObjectId]*scene)
}
