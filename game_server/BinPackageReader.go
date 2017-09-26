package game_server

// 这里 bytes 的内存使用可以优化，应该复用一个 Array
type RequestDataReader struct {
	gameServer *GameServer
	bytes      []byte
}

func (r *RequestDataReader) Read() []*RequestData {

	list := []*RequestData{}
	var pack *RequestData
	var i uint32
	var err error

	for len(r.bytes) > 13 && err == nil {
		pack, i, err = BytesToRequestData(r.gameServer, r.bytes)
		if err == nil {
			list = append(list, pack)
			r.bytes = r.bytes[i:]
		}

	}

	return list
}
func (r *RequestDataReader) Append(bytes []byte) {

	r.bytes = append(r.bytes, bytes...)
}

func NewRequestDataReader(s *GameServer) *RequestDataReader {
	return &RequestDataReader{bytes: []byte{}}
}
