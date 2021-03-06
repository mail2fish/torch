// This file is automatically generated by qtc from "binding_handers.qtpl".
// See https://github.com/valyala/quicktemplate for details.

//line binding_handers.qtpl:1
package templates

//line binding_handers.qtpl:1
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line binding_handers.qtpl:1
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line binding_handers.qtpl:1
func StreamGenerateBindingHandlers(qw422016 *qt422016.Writer, prefix string, appPath string) {
	//line binding_handers.qtpl:1
	qw422016.N().S(`
`)
	//line binding_handers.qtpl:2
	if len(prefix) > 0 {
		//line binding_handers.qtpl:2
		qw422016.N().S(`
package 	`)
		//line binding_handers.qtpl:3
		qw422016.E().S(prefix)
		//line binding_handers.qtpl:3
		qw422016.N().S(`_handlers

import (
	"`)
		//line binding_handers.qtpl:6
		qw422016.E().S(appPath)
		//line binding_handers.qtpl:6
		qw422016.N().S(`/game/`)
		//line binding_handers.qtpl:6
		qw422016.E().S(prefix)
		//line binding_handers.qtpl:6
		qw422016.N().S(`_global"
)   

type handler struct {
	id         uint16
	name       string
	fun        interface{}
	singleMode bool
}

var handlers = []*handler{}


func BindRequestHandlers() {

	for _, handler := range handlers {
		`)
		//line binding_handers.qtpl:22
		qw422016.E().S(prefix)
		//line binding_handers.qtpl:22
		qw422016.N().S(`_global.GS.BindRequestHandler(handler.id, handler.name, handler.fun, handler.singleMode)
	}
}

func addHandler(id uint16, name string, fun interface{}, singleMode bool) {
	handlers = append(handlers, &handler{id: id, name: name, fun: fun, singleMode: singleMode})
}

`)
		//line binding_handers.qtpl:30
	} else {
		//line binding_handers.qtpl:30
		qw422016.N().S(`
package 	handlers

import (
	"`)
		//line binding_handers.qtpl:34
		qw422016.E().S(appPath)
		//line binding_handers.qtpl:34
		qw422016.N().S(`/game/global"
)   

type handler struct {
	id         uint16
	name       string
	fun        interface{}
	singleMode bool
}

var handlers = []*handler{}


func BindRequestHandlers() {

	for _, handler := range handlers {
		global.GS.BindRequestHandler(handler.id, handler.name, handler.fun, handler.singleMode)
	}
}

func addHandler(id uint16, name string, fun interface{}, singleMode bool) {
	handlers = append(handlers, &handler{id: id, name: name, fun: fun, singleMode: singleMode})
}
`)
		//line binding_handers.qtpl:57
	}
	//line binding_handers.qtpl:57
	qw422016.N().S(`



`)
//line binding_handers.qtpl:61
}

//line binding_handers.qtpl:61
func WriteGenerateBindingHandlers(qq422016 qtio422016.Writer, prefix string, appPath string) {
	//line binding_handers.qtpl:61
	qw422016 := qt422016.AcquireWriter(qq422016)
	//line binding_handers.qtpl:61
	StreamGenerateBindingHandlers(qw422016, prefix, appPath)
	//line binding_handers.qtpl:61
	qt422016.ReleaseWriter(qw422016)
//line binding_handers.qtpl:61
}

//line binding_handers.qtpl:61
func GenerateBindingHandlers(prefix string, appPath string) string {
	//line binding_handers.qtpl:61
	qb422016 := qt422016.AcquireByteBuffer()
	//line binding_handers.qtpl:61
	WriteGenerateBindingHandlers(qb422016, prefix, appPath)
	//line binding_handers.qtpl:61
	qs422016 := string(qb422016.B)
	//line binding_handers.qtpl:61
	qt422016.ReleaseByteBuffer(qb422016)
	//line binding_handers.qtpl:61
	return qs422016
//line binding_handers.qtpl:61
}
