package handler

import (
	"encoding/json"
	"fmt"
	"iMonitor/dao"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type msg struct {
	Id   string `json:"id"`
	Data string `json:"data"`
}

var ch = make(chan string)

// 开始监控
func Serve() {
	// 开启socket管理器
	go manager.start()
}

func ServerLog(c *gin.Context) {
	ws, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		http.NotFound(c.Writer, c.Request)
		return
	}
	client := &client{time.Now().String(), ws, make(chan msg, 1024)}
	manager.register <- client
	fmt.Println("Conn")
	fmt.Println(client.id)
	go client.read()
	go client.write()
}

//  websocket客户端
type client struct {
	id     string
	socket *websocket.Conn
	send   chan msg
}

// 客户端管理
type clientManager struct {
	clients    map[*client]bool
	broadcast  chan msg
	register   chan *client
	unregister chan *client
}

var manager = clientManager{
	broadcast:  make(chan msg),
	register:   make(chan *client),
	unregister: make(chan *client),
	clients:    make(map[*client]bool),
}

func (manager *clientManager) start() {
	defer func() {
		if err := recover(); err != nil {
			logrus.Error(err)
		}
	}()

	for {
		select {
		case conn := <-manager.register:
			manager.clients[conn] = true
		case conn := <-manager.unregister:
			if _, ok := manager.clients[conn]; ok {
				fmt.Println("Close")
				fmt.Println(conn.id)

				close(conn.send)
				conn.socket.Close()
				delete(manager.clients, conn)
			}
		case msg := <-manager.broadcast:
			for conn := range manager.clients {
				if conn.id == msg.Id {
					conn.send <- msg
				}
			}
		}
	}
}

func (c *client) write() {

	defer func() {
		c.socket.Close()
		if err := recover(); err != nil {
			logrus.Error(err)
		}
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			msgByte, _ := json.Marshal(message)
			c.socket.WriteMessage(websocket.BinaryMessage, msgByte)
		}
	}

}

func (c *client) read() {

	defer func() {
		// ch <- struct{}{}
		manager.unregister <- c
		c.socket.Close()
		if err := recover(); err != nil {
			logrus.Error(err)
		}
	}()

	for {
		msgType, msg, err := c.socket.ReadMessage()
		if err != nil {
			// fmt.Println(err)
			ch <- c.id
			manager.unregister <- c
			c.socket.Close()
			break
		}
		if msgType != websocket.CloseMessage {
			// type recv struct {
			// 	ServerId int `json:"serverId"`
			// }

			// var rcv = &recv{}
			// if err := json.Unmarshal(msg, &rcv); err != nil {
			// 	ch <- struct{}{}
			// 	manager.unregister <- c
			// 	logrus.Error(err)
			// 	break
			// }
			serverId, err := strconv.Atoi(string(msg))
			if err != nil {
				ch <- c.id
				manager.unregister <- c
				logrus.Error(err)
				break
			}
			go Connect(serverId, c)
		} else {
			ch <- c.id
			manager.unregister <- c
			c.socket.Close()
		}

	}
}

func Connect(serverId int, c *client) {
	defer func() {
		if err := recover(); err != nil {
			logrus.Error(err)
		}
	}()
	server := dao.Server()
	server.ServerId = serverId
	if err := server.Get(); err != nil {
		logrus.Info(err)
		return
	}
	u := url.URL{Scheme: "ws", Host: server.Host + ":" + server.Port, Path: "/ws"}
	var dialer *websocket.Dialer
	conn, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		log.Println(err)
	}
	go func() {
		select {
		case id := <-ch:
			if id == c.id {
				conn.Close()
			}
		}
	}()
	for {
		type Message struct {
			Data string `json:"data"`
		}
		m := Message{}
		err := conn.ReadJSON(&m)
		if err != nil {
			log.Println(err)
			break
		}
		fmt.Println(m)
		manager.broadcast <- msg{c.id, m.Data}
	}
}
