package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

/**
*
	1. 启动websocket 服务等待客户端连接
	2. 客户端连接后调用createRoom创建房间
	3. 客户端创建房间后加入房间 获取传输过来的opus音频包
	4. 在同一房间内的 转发音频包
	5. 退出房间
	6. 销毁房间
*/

var upgrader = websocket.Upgrader{
	// 允许跨域（开发时使用）
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func startWebsocketService(port int){
	portStr := strconv.Itoa(port)
	http.HandleFunc("/ws", wsHandler)
	fmt.Println("WebSocket server started: ws://0.0.0.0:"+ portStr +"/ws")
	fmt.Println("Current room count:",roomManager.roomCount())
	if err := http.ListenAndServe(portStr, nil); err != nil {
		fmt.Println(err)
		return
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	// 升级为 WebSocket 连接
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	clientIP := r.RemoteAddr
	fmt.Println("Client connected","remote addr", clientIP)
	for {
		// 读取消息
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Read error:", err, "close this socket")
			return
		}
		
		textMsg := string(msg)
		fmt.Printf("msgType = %d recv: %s\n",msgType, textMsg)
		if err := conn.WriteMessage(msgType, msg); err != nil {
			fmt.Println("Write error:", err)
			break
		}

		switch msgType {
		case websocket.TextMessage: //文本消息
			// handleTextMsg(textMsg, conn)
		case websocket.BinaryMessage: //二进制消息
			// handleBinaryMsg(textMsg, conn)
		case websocket.CloseMessage:
		default:
			fmt.Println("handle default")
		}
	}//end for each
}

var roomManager ChatRoomManager = ChatRoomManager{
	data : make(map[string]*ChatRoom),
}

func main() {
	fmt.Println("Start audio chatroom server.")
	startWebsocketService(8910)
}