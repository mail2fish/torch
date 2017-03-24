package game_server

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/go-playground/log"
)

type FillFunParams struct {
	Params interface{}
	ProtobufStructType
}

//id:   the client use id to identify the data type
//typ: Protobuf Struct Type
//fillFun: A function for filling protobuf struct data.
//fillFunParamType: The fill function's param type
// asap: Send it to the client as soon as posible
type ProtobufStructType struct {
	id               uint16
	typ              *reflect.Type
	fillFun          *reflect.Value
	fillFunParamType *reflect.Type
	asap             bool
}

func (r *ProtobufStructType) ASAP() bool {
	return r.asap
}

func (r *ProtobufStructType) Id() uint16 {
	return r.id
}

func (r *ProtobufStructType) Type() reflect.Type {
	return *r.typ
}

func (r *ProtobufStructType) FillFunParamType() *reflect.Type {
	return r.fillFunParamType
}

func (r *ProtobufStructType) GetFillFun() *reflect.Value {
	return r.fillFun
}

type ProtobufStructTypes struct {
	im map[uint16]*ProtobufStructType       // id map
	tm map[reflect.Type]*ProtobufStructType //type map
}

func (s *ProtobufStructTypes) All() []*ProtobufStructType {
	list := []*ProtobufStructType{}
	for _, st := range s.im {
		list = append(list, st)
	}
	return list
}

func (s *ProtobufStructTypes) FindById(id uint16) (typ *ProtobufStructType, ok bool) {
	typ, ok = s.im[id]
	return
}

func (s *ProtobufStructTypes) FindByType(t reflect.Type) (typ *ProtobufStructType, ok bool) {
	typ, ok = s.tm[t]
	return
}

func (s *ProtobufStructTypes) FindId(t reflect.Type) (id uint16, ok bool) {
	var tp reflect.Type
	if t.Kind() == reflect.Ptr {
		tp = t.Elem()
	}
	fmt.Println(tp)
	if typ, ok := s.tm[tp]; ok {
		return typ.id, ok
	} else {
		return 0, ok
	}

}

func (s *ProtobufStructTypes) FindType(id uint16) (typ *reflect.Type, ok bool) {
	if typ, ok := s.im[id]; ok {
		return typ.typ, ok
	}
	return nil, ok
}

func (s *ProtobufStructTypes) FindFillFunParamType(id uint16) (typ *reflect.Type, ok bool) {
	if typ, ok := s.im[id]; ok {
		if typ != nil && typ.fillFunParamType != nil {
			return typ.fillFunParamType, ok
		} else {
			return nil, false
		}

	}
	return nil, ok
}

func (s *ProtobufStructTypes) BindProtobufStructType(id uint16, typ interface{}, fillFun interface{}, asap bool) {
	t := reflect.TypeOf(typ)

	var rpStruct *ProtobufStructType
	if fillFun != nil {
		if reflect.TypeOf(fillFun).Kind() != reflect.Func {
			panic("The fillFun should be a function")
		}

		fun := reflect.ValueOf(fillFun)

		typ := reflect.TypeOf(fillFun)

		funAsParam := typ.In(0)

		if funAsParam == reflect.TypeOf(FillFunSession{}) {
			panic("The fillFun's first param should  be an asession.ASession")
		}

		funParam := typ.In(1)

		if len(funParam.Name()) == 0 {
			panic("The fillFun param should not be an interface{}")
		}

		if reflect.TypeOf(fillFun).Out(1).Name() != "error" {
			panic("The fillFun  should return an error param")
		}

		rpStruct = &ProtobufStructType{id: id, typ: &t, fillFun: &fun, fillFunParamType: &funParam, asap: asap}

	} else {
		rpStruct = &ProtobufStructType{id: id, typ: &t, fillFun: nil, fillFunParamType: nil, asap: asap}
	}

	if _, ok := s.im[id]; ok {
		panic(fmt.Sprintf("You can not set a same id for different proto struct %d", id))
	}
	if _, ok := s.tm[t]; ok {
		panic(fmt.Sprintf("You can not set a same type for different proto struct %v", t))
	}

	s.im[id] = rpStruct
	s.tm[t] = rpStruct

}

func (s ProtobufStructTypes) MarshalProtostructTypeHash(rh ProtostructTypeHash) ([]byte, error) {
	var b = new(bytes.Buffer)
	binary.Write(b, binary.BigEndian, uint16(len(rh)))
	for k, v := range rh {
		binary.Write(b, binary.BigEndian, k)

		var typ *reflect.Type
		var ok bool
		checkTyp := false
		typName := GetTypeName(v)

		if typ, ok = s.FindFillFunParamType(k); ok {
			if (*typ).Name() == typName {
				checkTyp = true
			}
		}

		if typ, ok = s.FindType(k); ok {
			if (*typ).Name() == typName {
				checkTyp = true
			}
		}

		if checkTyp != true {
			panic("The RPHash value didn't bind")
		}

		if bytes, e := json.Marshal(v); e != nil {
			return nil, errors.New(fmt.Sprintf("Faild to marshal ProtostructTypeHash in RpStructs. The error is %v", e))
		} else {
			nameBytes := []byte(typName)
			binary.Write(b, binary.BigEndian, uint32(len(nameBytes)))
			binary.Write(b, binary.BigEndian, nameBytes)

			binary.Write(b, binary.BigEndian, uint32(len(bytes)))
			binary.Write(b, binary.BigEndian, bytes)
		}
	}
	return b.Bytes(), nil
}

func (s ProtobufStructTypes) UnMarshalProtostructTypeHash(bytes []byte) (ProtostructTypeHash, error) {
	rh := ProtostructTypeHash{}

	if len(bytes) <= 0 {
		return nil, errors.New("No data exists.")
	}
	c := binary.BigEndian.Uint16(bytes[0:2])

	index := uint32(2)
	for i := uint16(0); i < c; i++ {

		idBegin := index
		idEnd := idBegin + 2
		id := binary.BigEndian.Uint16(bytes[idBegin:idEnd])

		nameSizeBegin := idEnd
		nameSizeEnd := nameSizeBegin + 4

		nameSize := binary.BigEndian.Uint32(bytes[nameSizeBegin:nameSizeEnd])

		nameBegin := nameSizeEnd
		nameEnd := nameBegin + nameSize

		name := string(bytes[nameBegin:nameEnd])

		sizeBegin := nameEnd
		sizeEnd := sizeBegin + 4
		size := binary.BigEndian.Uint32(bytes[sizeBegin:sizeEnd])

		dataBegin := sizeEnd
		dataEnd := uint32(dataBegin) + size
		data := bytes[dataBegin:dataEnd]
		index = dataEnd

		if typ, ok := s.FindFillFunParamType(id); ok && (*typ).Name() == name {
			log.Debug("Find RpStruct param type by ", id, (*typ).Name())
			v := reflect.New(*typ).Interface()
			if e := json.Unmarshal(data, v); e != nil {
				log.Error("Faild to unmarshal type ", GetTypeName(typ), "  in UnMarshalProtostructTypeHash. The error is %v", e)
			}
			rh[id] = reflect.ValueOf(v).Elem().Interface()

		} else if typ, ok := s.FindType(id); ok {
			log.Debug("Find RpStruct type by Id:%d. %s", id, (*typ).Name())
			v := reflect.New(*typ).Interface()
			if e := json.Unmarshal(data, v); e != nil {
				log.Error("Faild to unmarshal type  ", GetTypeName(typ), "  in UnMarshalProtostructTypeHash. The error is %v", e)
			}
			rh[id] = v

		} else {
			panic("Faild to find ProtostructTypeHash Value")
		}

	}

	return rh, nil
}

func (s ProtobufStructTypes) IsFillFunParam(id uint16, i interface{}) bool {

	if s, ok := s.FindById(id); ok {
		typ := s.FillFunParamType()
		if typ != nil && (*typ).Name() == GetTypeName(i) {
			return true
		}
	}
	return false
}

func NewProtobufStructTypes() *ProtobufStructTypes {
	rpStructs := &ProtobufStructTypes{}
	rpStructs.im = make(map[uint16]*ProtobufStructType)
	rpStructs.tm = make(map[reflect.Type]*ProtobufStructType)
	return rpStructs
}

func GetTypeName(args interface{}) string {
	var argsTyp reflect.Type
	switch args.(type) {
	case reflect.Type:
		argsTyp = args.(reflect.Type)
	default:
		argsTyp = reflect.TypeOf(args)
	}

	var argsName string
	if argsTyp.Kind() == reflect.Ptr {
		argsName = argsTyp.Elem().Name()
	} else {
		argsName = argsTyp.Name()
	}
	return argsName
}

type ProtostructTypeHash map[uint16]interface{}

var GameProtobufStructTypes *ProtobufStructTypes

func init() {
	GameProtobufStructTypes = NewProtobufStructTypes()
}
