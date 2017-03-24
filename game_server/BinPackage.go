package game_server

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc32"

	"github.com/go-playground/log"
)

const (
	START_CHAR = uint8(253)
	END_CHAR   = uint8(254)
)

// START_CHAR(1 byte) - Handler Id (2 byte) -  Version(1 byte) - Data size(4 byte) - Data - CRC(4 byte) - END_CHAR

type BinPackage struct {
	HandlerId uint16
	Version   uint8
	Data      []byte
}

func (p *BinPackage) ToBytes() (data []byte) {
	var list []interface{}

	list = append(list, START_CHAR)
	list = append(list, p.HandlerId)
	list = append(list, p.Version)
	size := uint32(len(p.Data))
	list = append(list, size)
	for _, v := range p.Data {
		list = append(list, v)
	}

	list = append(list, crc32.ChecksumIEEE(p.Data))
	list = append(list, END_CHAR)

	buf := new(bytes.Buffer)
	for _, v := range list {
		err := binary.Write(buf, binary.BigEndian, v)
		if err != nil {
			log.Error("Met an error when parse some bytes. %v", err)

		}

	}
	result := buf.Bytes()
	return result
}

func BytesToBinPackage(byteSlice []byte) (pack *BinPackage, index uint32, err error) {

	minimalPackageSize := uint32(12)
	length := uint32(len(byteSlice))
	if length < minimalPackageSize {
		return nil, 0, errors.New(fmt.Sprintf("Data size  %d is less than mimial size ", length))
	}

	for i, v := range byteSlice {

		iSize := uint32(length - 1)

		if v == START_CHAR {
			suffix := END_CHAR

			handerIdStart := uint32(i + 1)
			handerIdEnd := handerIdStart + 2

			if handerIdStart > iSize || handerIdEnd > iSize {
				continue
			}
			handerId := binary.BigEndian.Uint16(byteSlice[handerIdStart:handerIdEnd])

			verIndex := handerIdEnd + 1

			if verIndex > iSize {
				continue
			}
			version := uint8(byteSlice[verIndex])
			sizeStart := verIndex
			sizeEnd := sizeStart + 4

			if length < uint32(i)+minimalPackageSize {
				continue
			}

			if sizeStart > iSize || sizeEnd > iSize {
				continue
			}

			dataSize := binary.BigEndian.Uint32(byteSlice[sizeStart:sizeEnd])
			dataStart := sizeEnd
			dataEnd := dataStart + dataSize

			crcStart := dataEnd
			crcEnd := crcStart + 4

			suffixIndex := crcEnd

			if dataStart > iSize || dataEnd > iSize || crcStart > iSize || crcEnd > iSize {
				continue
			}
			if uint8(byteSlice[suffixIndex]) != suffix {
				continue
			}

			data := byteSlice[dataStart:dataEnd]
			// crcCode := binary.BigEndian.Uint32(byteSlice[crcStart:crcEnd])
			// if crcCode != crc32.ChecksumIEEE(data) {

			// 	continue
			// }

			cloneData := make([]byte, len(data))
			copy(cloneData, data)
			pack = &BinPackage{
				HandlerId: handerId,
				Version:   version,
				Data:      cloneData}
			return pack, crcEnd + 1, nil
		}
	}
	return nil, 0, errors.New("No binpack exist. ")
}
