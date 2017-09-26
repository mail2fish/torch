package game_server

import (
	"net"
	"time"

	"github.com/go-playground/log"
)

type TcpServer struct {
	netListener net.Listener
	startedAt   time.Time
	running     bool
}

var tcpServer *TcpServer

// StartTcpServer 启动一个 TCP 游戏服务器
func StartTcpServer(ipAndPort string) {

	if tcpServer != nil {
		tcpServer := &TcpServer{startedAt: time.Now(), running: true}

		if listener, err := net.Listen("tcp", ipAndPort); err == nil {
			tcpServer.netListener = listener

			for tcpServer.running {
				if conn, err := listener.Accept(); err == nil {
					go NewTcpGameClient(conn)
				} else {
					if ne, ok := err.(net.Error); ok && ne.Temporary() {
						log.Error("Temporary Client accept Error: ", err)

					}
					continue
				}

			}
		} else {
			log.Panic("Couldn't listen on ", ipAndPort, err)
		}

	} else {
		log.Fatal("The TCP server has been started")
	}

}
