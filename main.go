package main

import (
	"net/http"
	"strconv"
	"sync"

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
	defer wg.Done()

	portStr := strconv.Itoa(port)
	http.HandleFunc(WS_URL, wsHandler)
	Logi("WebSocket server started: ws://0.0.0.0:"+ portStr + WS_URL)
	Logi("Current room count:",roomManager.roomCount())
	if err := http.ListenAndServe(":" + portStr, nil); err != nil {
		Logi(err)
		return
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	// 升级为 WebSocket 连接
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		Loge("Upgrade error:", err)
		return
	}

	clientIP := r.RemoteAddr
	Logi("Client connected","remote addr", clientIP)

	accountId,session := accounts.AddNewAccount(conn)
	Logi("Client account id is",accountId)
	session.RunLoop()
}

var roomManager ChatRoomManager = ChatRoomManager{
	data : make(map[string] *ChatRoom),
}

var accounts ChatAccounts = ChatAccounts{
	value : make(map[int64] *Session),
}

const WS_URL string = "/chat"
var wg sync.WaitGroup

func main() {
	wg.Add(1)
	Logi("Start audio chatroom server.")
	go startWebsocketService(8910)
	wg.Wait()
}