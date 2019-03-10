package impl

import (
	"errors"
	"github.com/gorilla/websocket"
	"sync"
)

/*
websocket 的ReadMessage和WriteMessage 不是线程安全的，封装线程安全的ReadMessage和WriteMessage
 */
type Connection struct{
	wsConn * websocket.Conn
	inChan chan []byte
	outChan chan []byte
	closeChan chan byte
	mutex sync.Mutex
	isClosed bool
}

func InitConnection(wsConn *websocket.Conn)(conn *Connection, err error){
	conn = &Connection{
		wsConn:wsConn,
		inChan:make(chan []byte, 1000),
		outChan:make(chan []byte, 1000),
		closeChan:make(chan byte, 1),
	}
	go conn.readLoop()
	go conn.writeLoop()
	return
}

// API
func (conn *Connection)ReadMessage()(data []byte, err error){
	select {
	case data = <- conn.inChan:
	case <- conn.closeChan:
		err = errors.New("conneciton is closed!!")
	}
	return
}

func (conn *Connection)WriteMessage(data []byte)(err error)  {
	select {
	case conn.outChan <- data:
	case <- conn.closeChan:
		err = errors.New("connection is closed")
	}

	return
}

func (conn *Connection)Close()  {
	conn.wsConn.Close()
	conn.mutex.Lock()
	if ! conn.isClosed{
		close(conn.closeChan)
		conn.isClosed = true
	}
	conn.mutex.Unlock()
}

func (conn *Connection)readLoop(){
	var(
		data []byte
		err error
	)
	defer conn.Close()
	for {
		if _, data, err = conn.wsConn.ReadMessage(); err != nil{
			return
		}
		select {
		case conn.inChan <- data:
		case <- conn.closeChan:
			//closeChan 被关闭的时候
			return
		}
		
	}
}

func (conn *Connection)writeLoop(){
	var (
		data []byte
		err error
	)
	defer conn.Close()
	for {
		select {
		case data = <- conn.outChan:
		case <- conn.closeChan:
			return
		}
		if err = conn.wsConn.WriteMessage(websocket.TextMessage, data); err != nil{
			return
		}

	}
}
