package quoter

import (
	"bytes"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

type Message struct {
	Typ     int
	Content []byte
}

type Socket struct {
	conn    *websocket.Conn
	errChan chan error
	msgChan chan Message
}

func NewSocket(wsURL string, header http.Header) (*Socket, error) {
	var conn, res, err = websocket.DefaultDialer.Dial(wsURL, header)
	if err != nil {
		return nil, fmt.Errorf("dial.error:%+v", err)
	}
	defer func() {
		_ = res.Body.Close()
	}()
	var s = &Socket{
		conn:    conn,
		errChan: make(chan error, 1),
		msgChan: make(chan Message, 10),
	}
	go s.loop()
	return s, nil
}

// Close connection
func (s *Socket) Close() error {
	return s.conn.Close()
}

// SendStr Send string message
func (s *Socket) SendStr(msg string) error {
	var content = &bytes.Buffer{}
	content.WriteString(msg)
	return s.conn.WriteMessage(websocket.TextMessage, content.Bytes())
}

func (s *Socket) Wait() <-chan error {
	return s.errChan
}

func (s *Socket) Message() <-chan Message {
	return s.msgChan
}

func (s *Socket) loop() {
	for {
		var t, msg, err = s.conn.ReadMessage()
		if err != nil {
			s.errChan <- err
			break
		}
		s.msgChan <- Message{
			Typ:     t,
			Content: msg,
		}
	}
}
