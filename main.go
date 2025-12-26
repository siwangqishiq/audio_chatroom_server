package main

import (
	"encoding/json"
	"fmt"
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

const WS_URL string = "/chat"

var wg sync.WaitGroup

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
	fmt.Println("WebSocket server started: ws://0.0.0.0:"+ portStr + WS_URL)
	fmt.Println("Current room count:",roomManager.roomCount())
	if err := http.ListenAndServe(":" + portStr, nil); err != nil {
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
	accountId := accounts.AddNewAccount(conn)
	fmt.Println("Client account id is",accountId)

	SendPacket(conn, BuildLoginData(accountId))
	for {
		// 读取消息
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Read error:", err, "close this socket")
			break
		}
		
		textMsg := string(msg)
		fmt.Printf("msgType = %d recv: %s\n",msgType, textMsg)

		switch msgType {
		case websocket.TextMessage: //文本消息
			handlePacket(accountId, textMsg, conn)
		case websocket.BinaryMessage: //二进制消息
			// handleBinaryMsg(textMsg, conn)
		case websocket.CloseMessage:
		default:
			fmt.Println("handle default")
		}
	}//end for each
	accounts.RemoveAccount(accountId)
}

func SendPacket(conn *websocket.Conn, msg Packet) {
	data,err := json.Marshal(msg)
	if(err != nil){
		return
	}
	conn.WriteMessage(websocket.TextMessage, data)
}

func SendBinary(conn *websocket.Conn){

}

var roomManager ChatRoomManager = ChatRoomManager{
	data : make(map[string] *ChatRoom),
}

var accounts ChatAccounts = ChatAccounts{
	value : make(map[int64] *websocket.Conn),
}

func handlePacket(accountId int64, rawText string , conn *websocket.Conn){
	packet := Packet{}
	err := json.Unmarshal([]byte(rawText), &packet)
	if(err != nil){
		fmt.Println("handle packet error",err.Error())
		return
	}

	fmt.Println("Packet cmd", packet.Cmd)
	switch packet.Cmd{
	case CMD_CREATE_ROOM_JOIN_REQ:
		handleCreateRoomAndJoin(accountId, &packet, conn)
	default:
		fmt.Println("Not support cmd",packet.Cmd)
	}//end switch
}

func handleCreateRoomAndJoin(accountId int64, packt *Packet, conn *websocket.Conn) {
	cid := packt.Cid
	paramsMap := packt.Data.(map[string]any)
	roomId,_ := paramsMap["roomId"]
	r := roomId.(string)
	fmt.Println("handleCreateRoomAndJoin cid",cid,"roomid",r)
	// roomManager.CheckRoomExist()
}

func main() {
	wg.Add(1)
	fmt.Println("Start audio chatroom server.")
	go startWebsocketService(8910)
	wg.Wait()
}