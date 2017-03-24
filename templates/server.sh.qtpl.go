// This file is automatically generated by qtc from "server.sh.qtpl".
// See https://github.com/valyala/quicktemplate for details.

//line server.sh.qtpl:1
package templates

//line server.sh.qtpl:1
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line server.sh.qtpl:1
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line server.sh.qtpl:1
func StreamGenerateServerScript(qw422016 *qt422016.Writer, cmd string) {
	//line server.sh.qtpl:1
	qw422016.N().S(`#!/bin/bash
echo GOPATH: ${GOPATH}
`)
	//line server.sh.qtpl:3
	qw422016.E().S(cmd)
	//line server.sh.qtpl:3
	qw422016.N().S(`
`)
//line server.sh.qtpl:4
}

//line server.sh.qtpl:4
func WriteGenerateServerScript(qq422016 qtio422016.Writer, cmd string) {
	//line server.sh.qtpl:4
	qw422016 := qt422016.AcquireWriter(qq422016)
	//line server.sh.qtpl:4
	StreamGenerateServerScript(qw422016, cmd)
	//line server.sh.qtpl:4
	qt422016.ReleaseWriter(qw422016)
//line server.sh.qtpl:4
}

//line server.sh.qtpl:4
func GenerateServerScript(cmd string) string {
	//line server.sh.qtpl:4
	qb422016 := qt422016.AcquireByteBuffer()
	//line server.sh.qtpl:4
	WriteGenerateServerScript(qb422016, cmd)
	//line server.sh.qtpl:4
	qs422016 := string(qb422016.B)
	//line server.sh.qtpl:4
	qt422016.ReleaseByteBuffer(qb422016)
	//line server.sh.qtpl:4
	return qs422016
//line server.sh.qtpl:4
}
