// This file is automatically generated by qtc from "request_protobuf.qtpl".
// See https://github.com/valyala/quicktemplate for details.

//line request_protobuf.qtpl:1
package templates

//line request_protobuf.qtpl:1
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line request_protobuf.qtpl:1
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line request_protobuf.qtpl:1
func StreamGenerateExampleRequestProtobuf(qw422016 *qt422016.Writer, prefix string, appPath string) {
	//line request_protobuf.qtpl:1
	qw422016.N().S(`
syntax = "proto3";
`)
	//line request_protobuf.qtpl:3
	if len(prefix) > 0 {
		//line request_protobuf.qtpl:3
		qw422016.N().S(`
package `)
		//line request_protobuf.qtpl:4
		qw422016.E().S(prefix)
		//line request_protobuf.qtpl:4
		qw422016.N().S(`_pstructs;
import "common.proto";

`)
		//line request_protobuf.qtpl:7
	} else {
		//line request_protobuf.qtpl:7
		qw422016.N().S(`

package pstructs;
import "common.proto";

`)
		//line request_protobuf.qtpl:12
	}
	//line request_protobuf.qtpl:12
	qw422016.N().S(`

message RequestHeader{
    int64 sentAt = 1;    
    uint32 requestSqNumber = 2;
}

// request id 1
message RqExample{
    RequestHeader header=1;
    string name=2;
}

`)
//line request_protobuf.qtpl:25
}

//line request_protobuf.qtpl:25
func WriteGenerateExampleRequestProtobuf(qq422016 qtio422016.Writer, prefix string, appPath string) {
	//line request_protobuf.qtpl:25
	qw422016 := qt422016.AcquireWriter(qq422016)
	//line request_protobuf.qtpl:25
	StreamGenerateExampleRequestProtobuf(qw422016, prefix, appPath)
	//line request_protobuf.qtpl:25
	qt422016.ReleaseWriter(qw422016)
//line request_protobuf.qtpl:25
}

//line request_protobuf.qtpl:25
func GenerateExampleRequestProtobuf(prefix string, appPath string) string {
	//line request_protobuf.qtpl:25
	qb422016 := qt422016.AcquireByteBuffer()
	//line request_protobuf.qtpl:25
	WriteGenerateExampleRequestProtobuf(qb422016, prefix, appPath)
	//line request_protobuf.qtpl:25
	qs422016 := string(qb422016.B)
	//line request_protobuf.qtpl:25
	qt422016.ReleaseByteBuffer(qb422016)
	//line request_protobuf.qtpl:25
	return qs422016
//line request_protobuf.qtpl:25
}
