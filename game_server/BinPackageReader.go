package game_server

import "sync"

type BinPackageReader struct {
	bytes []byte
	rwMU  sync.RWMutex
}

func (r *BinPackageReader) Read() []*BinPackage {
	r.rwMU.RLock()
	defer r.rwMU.RUnlock()
	list := []*BinPackage{}
	var pack *BinPackage
	var i uint32
	var err error

	for len(r.bytes) > 13 && err == nil {
		pack, i, err = BytesToBinPackage(r.bytes)
		if err == nil {
			list = append(list, pack)
			r.bytes = r.bytes[i:]
		}

	}

	return list
}
func (r *BinPackageReader) Append(bytes []byte) {
	r.rwMU.Lock()
	defer r.rwMU.Unlock()
	r.bytes = append(r.bytes, bytes...)
}

func NewBinPackageReader() *BinPackageReader {
	return &BinPackageReader{bytes: []byte{}}
}
