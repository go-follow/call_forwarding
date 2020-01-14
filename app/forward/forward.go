package forward

import (
	"fmt"
	"io"
	"net"
	"strings"

	"github.com/go-follow/call_forwarding/app/models"
	"github.com/go-follow/call_forwarding/logger"
)

//Forward структура для инициализации объекта переадресации
type Forward struct {
	s *models.Settings
	l net.Listener
}

//NewForward один объект переадресации
func NewForward(s *models.Settings) (*Forward, error) {
	if s == nil {
		return nil, fmt.Errorf("empty settings for forwarding")
	}
	if strings.Trim(s.ListnerIP, " ") == "" {
		return nil, fmt.Errorf("empty IP for listner")
	}
	if strings.Trim(s.ForwardIP, " ") == "" {
		return nil, fmt.Errorf("empty IP for forward")
	}
	if s.ListnerPort <= 0 {
		return nil, fmt.Errorf("empty port for lisner")
	}
	if s.ForwardPort <= 0 {
		return nil, fmt.Errorf("empty port for forward")
	}
	urlAccept := fmt.Sprintf("%s:%d", s.ListnerIP, s.ListnerPort)
	listner, err := net.Listen("tcp", urlAccept)
	if err != nil {
		return nil, fmt.Errorf("failed install listen in %s: %v", urlAccept, err)
	}
	return &Forward{s, listner}, nil
}

//StartListner - запуск для одной переадресации
func (f *Forward) StartListner() {
	for {
		connListner, err := f.l.Accept()
		if err != nil {
			logger.Errorf("failed to accept lisner in %s: %v",
				fmt.Sprintf("%s:%d", f.s.ListnerIP, f.s.ListnerPort), err)
			continue
		}
		defer connListner.Close()
		if err := f.forward(connListner); err != nil {
			logger.Errorf("failed to forward in %s: %v",
				fmt.Sprintf("%s:%d", f.s.ForwardIP, f.s.ForwardPort), err)
		}
	}
}

func (f *Forward) forward(connListner net.Conn) error {
	urlForward := fmt.Sprintf("%s:%d", f.s.ForwardIP, f.s.ForwardPort)
	connForward, err := net.Dial("tcp", urlForward)
	if err != nil {
		return err
	}
	go func() {
		defer connForward.Close()
		_, err = io.Copy(connForward, connListner)
		if err != nil && isLog(err.Error()) {
			logger.Error(err)
		}
	}()
	go func() {
		defer connForward.Close()
		_, err = io.Copy(connListner, connForward)
		if err != nil && isLog(err.Error()) {
			logger.Error(err)
		}
	}()

	return nil
}

func isLog(textError string) bool {
	if strings.Contains(textError, "use of closed network connection") {
		return false
	}
	return true
}
