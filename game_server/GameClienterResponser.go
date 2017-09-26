package game_server

import (
	"errors"
	"reflect"
	"sync"

	"github.com/go-playground/log"

	"github.com/gogo/protobuf/proto"
)

type GameClienterResponser struct {
	sr          *SuperResponser
	gc          GameClienter
	chData      chan interface{}
	chAvailable bool
	mu          sync.RWMutex
}

func (gr *GameClienterResponser) Write(kind ResponsKind, data ...interface{}) {
	if RK_WEBSOCKET == kind {
		if len(data) == 1 {

			gr.chData <- data[0]
		}

	} else {
		log.Error("Send wrong kind data to game  clienter responser")
	}

}

func (gr *GameClienterResponser) loopWrite() {

	for data := range gr.chData {

		if gr.gc == nil || !gr.gc.Available() {
			log.Error("Loop write game clienter is not available")
			gr.sr.closedGameClienterResponserCh <- gr.gc.Id()
			close(gr.chData)
			return
		}

		var rpStruct *ProtobufStructType
		var ok bool
		var result interface{}
		var err error

		switch data.(type) {
		case *FillFunParams:
			params := data.(*FillFunParams)
			rpStruct = &params.ProtobufStructType
			log.Debug("Writing *FillFunParams to nc. The type is %v. ", rpStruct.Type().String())

			if result, err = _callFillFun(*params); err != nil {
				result = nil
				log.Error("Call Fill fun failed. Error: %v", err)
			}

		case FillFunParams:
			params := data.(FillFunParams)
			rpStruct = &params.ProtobufStructType

			log.Debug("Writing FillFunParams to nc. The type is %v. ", rpStruct.Type().String())
			if result, err = _callFillFun(params); err != nil {
				result = nil
				log.Error("Call Fill fun failed. Error: %v", err)
			}

		default:

			t := reflect.TypeOf(data)
			rpStruct, ok = GameProtobufStructTypes.FindByType(t)

			if !ok {
				if reflect.ValueOf(data).Kind() == reflect.Ptr {
					t = reflect.ValueOf(data).Type().Elem()
				}

				rpStruct, ok = GameProtobufStructTypes.FindByType(t)
				if !ok {
					log.Debug("Writing error!!!!!")
					continue
				}
			}
			log.Debug("Writing the type ", t.Name(), " to client. ")

			result = data

		}

		if result == nil {
			continue
		}

		var pack ResponseData

		dataSending, err := proto.Marshal(result.(proto.Message))
		if err != nil {
			log.Error("Marshal proto error  ", err)
			continue
		}
		pack = ResponseData{HandlerId: rpStruct.Id(), Version: 1, Data: dataSending}

		bytesData := pack.ToBytes()
		gr.gc.Write(bytesData)
		// fmt.Println(bytesData)
	}
}

func _callFillFun(params FillFunParams) (interface{}, error) {

	fs := &FillFunSession{}

	fillFun := params.GetFillFun()
	if fillFun == nil {
		return nil, errors.New("Fill fun is nil")
	}

	results := fillFun.Call([]reflect.Value{reflect.ValueOf(fs), reflect.ValueOf(params.Params)})

	if results[1].Interface() != nil {
		err := results[1].Interface().(error)
		log.Error("Load FillFunParams failed. The error is %v ", err)
		return nil, err
	}
	return results[0].Interface(), nil

}
