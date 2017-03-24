package log_handler

import (
	"github.com/go-playground/log"
)

type PlayLogHandler struct {
}

func (h *PlayLogHandler) Run() chan<- *log.Entry {
	ch := make(chan *log.Entry)
	return ch
}
