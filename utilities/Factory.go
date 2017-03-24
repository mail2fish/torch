package utilities

// import (
// 	"cocos/poem_combat/poem_server/game/proto_structs"
// 	"time"

// 	"github.com/jinzhu/gorm"
// )

// func GenerateRpHeader(sq uint32) *proto_structs.ResponseHeader {
// 	return &proto_structs.ResponseHeader{SentAt: time.Now().Unix(), Sq: sq}
// }

// func GenerateGormModel() gorm.Model {
// 	return gorm.Model{CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()}
// }

// func GenerateRpRc(sq uint32) *proto_structs.RpRc {
// 	return &proto_structs.RpRc{Header: GenerateRpHeader(sq)}
// }
