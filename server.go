package main

import (
	"github.com/go-websocket/impl"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

var(
	upgrader = websocket.Upgrader{
		// 允许跨域
		CheckOrigin: func(r *http.Request) bool{
			return true
		},
	}
)
func wsHandler(w http.ResponseWriter, r *http.Request) {
	var(
		wsConn *websocket.Conn
		err error
		data []byte
		conn *impl.Connection
	)
	// upgrade:websocket, http 转换为websocket协议
	if wsConn, err = upgrader.Upgrade(w, r, nil); err != nil{
		return
	}
	defer conn.Close()
	if conn, err = impl.InitConnection(wsConn); err != nil{
		return
	}

	go func(){
		var(
			err error
		)
		for{
			if err = conn.WriteMessage([]byte("heartbeat")); err != nil{
				return
			}
			time.Sleep(1 * time.Second)
		}
	}()
	// websocket conn
	for{
		if data, err = conn.ReadMessage(); err != nil{
			return
		}
		if err = conn.WriteMessage(data); err != nil{
			return
		}
	}
}
func main()  {
	http.HandleFunc("/ws", wsHandler)
	http.ListenAndServe("127.0.0.1:7777", nil)
}
