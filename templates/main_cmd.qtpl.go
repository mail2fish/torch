// This file is automatically generated by qtc from "main_cmd.qtpl".
// See https://github.com/valyala/quicktemplate for details.

//line main_cmd.qtpl:1
package templates

//line main_cmd.qtpl:1
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line main_cmd.qtpl:1
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line main_cmd.qtpl:1
func StreamGenerateMain(qw422016 *qt422016.Writer, appPath string) {
	//line main_cmd.qtpl:1
	qw422016.N().S(`


package main

import "`)
	//line main_cmd.qtpl:6
	qw422016.E().S(appPath)
	//line main_cmd.qtpl:6
	qw422016.N().S(`/cmd/server/cmd"

func main() {
	cmd.Execute()

}



`)
//line main_cmd.qtpl:15
}

//line main_cmd.qtpl:15
func WriteGenerateMain(qq422016 qtio422016.Writer, appPath string) {
	//line main_cmd.qtpl:15
	qw422016 := qt422016.AcquireWriter(qq422016)
	//line main_cmd.qtpl:15
	StreamGenerateMain(qw422016, appPath)
	//line main_cmd.qtpl:15
	qt422016.ReleaseWriter(qw422016)
//line main_cmd.qtpl:15
}

//line main_cmd.qtpl:15
func GenerateMain(appPath string) string {
	//line main_cmd.qtpl:15
	qb422016 := qt422016.AcquireByteBuffer()
	//line main_cmd.qtpl:15
	WriteGenerateMain(qb422016, appPath)
	//line main_cmd.qtpl:15
	qs422016 := string(qb422016.B)
	//line main_cmd.qtpl:15
	qt422016.ReleaseByteBuffer(qb422016)
	//line main_cmd.qtpl:15
	return qs422016
//line main_cmd.qtpl:15
}
