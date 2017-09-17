package game_server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/log"
	"github.com/gorilla/websocket"
)

// StartWebSocketServer 启动 WebScokcet 服务器
// webSocketURL 是 websocket client 建立连接时的那个 http 相对路径，类似: /ws_connet
// ipAndPort 是服务器的监听地址和端口，可以简单写成 :8080
func StartWebSocketServer(gameServer *GameServer, ipAndPort string, webSocketURL string, readBufferSize int, writeBufferSize int) {

	r := gin.Default()

	var wsupgrader = websocket.Upgrader{
		ReadBufferSize:  readBufferSize,
		WriteBufferSize: writeBufferSize,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	r.GET(webSocketURL, func(c *gin.Context) {
		func(w http.ResponseWriter, r *http.Request) {
			conn, err := wsupgrader.Upgrade(w, r, nil)

			if err != nil {
				log.Panic("Failed to set websocket upgrade:", err)
				return
			}

			wsClient := NewWebSocketGameClient(conn, gameServer)
			gameServer.AddGameClient(wsClient)

		}(c.Writer, c.Request)
	})
	r.Run(ipAndPort)

}
