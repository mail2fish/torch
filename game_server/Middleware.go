package game_server

import "reflect"

type Middlewarer interface {
	Call(*MiddlewareIndexer, *RequestSession, Responser, *RequestData, interface{}) error
}

type MiddlewareIndexer struct {
	middlewares []Middlewarer
	index       int
	handlerFun  reflect.Value
	hasCalled   bool
}

func (m *MiddlewareIndexer) Call(s *RequestSession, rp Responser, pack *RequestData, param interface{}, handleFun reflect.Value) error {
	m.handlerFun = handleFun

	err := NextMiddleware(m, s, rp, pack, param)

	return err

}
func NextMiddleware(m *MiddlewareIndexer, s *RequestSession, rp Responser, pack *RequestData, param interface{}) error {
	var err error
	if m.index < len(m.middlewares) {
		fun := (m.middlewares)[m.index]
		m.index++
		err = fun.Call(m, s, rp, pack, param)
	}

	if err == nil && m.hasCalled == false {

		m.hasCalled = true

		in := make([]reflect.Value, 3)
		in[0] = reflect.ValueOf(s)
		in[1] = reflect.ValueOf(rp)
		in[2] = reflect.ValueOf(param)

		m.handlerFun.Call(in)

	}
	return err

}
