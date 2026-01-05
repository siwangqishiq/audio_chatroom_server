package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Session struct {
	accountId int64
	conn *websocket.Conn
	sendPacketChan chan []byte 
	sendBinaryChan chan []byte
	close chan int

	wg sync.WaitGroup
}

func CreateNewSession(accountId int64, conn *websocket.Conn) *Session{
	var session *Session = new(Session)
	session.accountId = accountId
	session.conn = conn
	session.sendPacketChan = make(chan []byte, 4)
	session.sendBinaryChan = make(chan []byte, 8)
	session.close = make(chan int)
	return session
}

func (s *Session)RunLoop(){
	s.wg.Add(2)
	go s.ReadLoop()
	go s.WriteLoop()
	s.wg.Wait()
	s.conn.Close()
}

func (s *Session)ReadLoop(){
	defer func(){s.close <- 1}()
	defer accounts.RemoveAccount(s.accountId)
	defer s.wg.Done()

	Logi("will send login data to client",s.accountId)
	s.SendPacket(BuildLoginData(s.accountId))

	var isQuit bool = false
	for {
		// 读取消息
		msgType, msg, err := s.conn.ReadMessage()
		if err != nil {
			Loge("Read error:", err, "close this socket")
			break
		}
		
		switch msgType {
		case websocket.TextMessage: //文本消息
			textMsg := string(msg)
			Logi(fmt.Sprintf("msgType = %d recv: %s\n",msgType, textMsg))
			s.handlePacket(textMsg)
		case websocket.BinaryMessage: //二进制消息
		case websocket.CloseMessage:
			isQuit = true
		default:
			Logi("handle default")
		}//end switch

		if isQuit {
			break
		}
	}//end for each
	Logi(s.accountId, "websocket closed.")
}

func (s *Session)SendPacket(msg Packet){
	data,err := json.Marshal(msg)
	if(err != nil){
		return
	}
	s.sendPacketChan <- data
}
	

func (s *Session)WriteLoop(){
	defer s.wg.Done()

	for{
		select {
			case data := <-s.sendPacketChan:
			s.conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
			if err := s.conn.WriteMessage(websocket.TextMessage, data); err != nil {
				Loge("WriteLoop error on WriteMessage ", s.accountId)
				return
			}
			Logi("WriteLoop to send size",len(data))
			case data := <-s.sendBinaryChan:
			s.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := s.conn.WriteMessage(websocket.BinaryMessage, data); err != nil {
				Loge("ReadLoop error on write binary message ", s.accountId)
				return
			}
			Logi("WriteLoop to send size",len(data))
			case <-s.close:
				return
		}//end select
	}//end for each
}

func (s *Session)handlePacket(rawText string){
	packet := Packet{}
	err := json.Unmarshal([]byte(rawText), &packet)
	if(err != nil){
		Logi("handle packet error",err.Error())
		return
	}

	Logi("Packet cmd", packet.Cmd , s.accountId)
	switch packet.Cmd{
	case CMD_CREATE_ROOM_JOIN_REQ:
		s.handleCreateRoomAndJoin(&packet)
	case CMD_JOIN_ROOM_REQ:
		s.handleJoinRoom(&packet)
	case CMD_QUIT_ROOM_REQ:
		s.handleQuitRoom(&packet)
	default:
		Loge("Not support cmd",packet.Cmd)
	}//end switch
}

func (s *Session)handleQuitRoom(pkt *Packet){
	cid := pkt.Cid
	paramsMap := pkt.Data.(map[string]any)
	r,_ := paramsMap["roomId"]
	roomId := r.(string)

	Logi("handleQuitRoom cid",cid,"roomid",roomId)
	var room = roomManager.FindRoomById(roomId)
	if room == nil {
		Loge(roomId,"quit roomid has not exist!")
		s.SendPacket(BuildQuitRoomError(cid, CODE_ROOM_NOT_EXIST))
		return
	}

	result,msg := roomManager.QuitRoom(room.roomId, s.accountId)

	if result {
		Logi("QuitRoom success","roomid",room.roomId)
		s.SendPacket(BuildQuitRoomSuccess(cid, room))

		if(room.adminId == s.accountId){//主持人退出 需要结束会议
			s.finishRoom(room.roomId)
		}
	}else{
		Logi("QuitRoom failed","roomid",room.roomId,msg)
		s.SendPacket(BuildQuitRoomError(cid, CODE_QUIT_ROOM_ERROR))
	}
}

func (s *Session)finishRoom(roomId string){
	Logi("ChatRoom finished","roomid",roomId)
	if roomManager.FinishRoom(roomId) {
		s.SendPacket(BuildFinishRoom(roomId))
	}
}

func (s *Session)handleJoinRoom(pkt *Packet){
	cid := pkt.Cid

	paramsMap := pkt.Data.(map[string]any)
	r,_ := paramsMap["roomId"]
	roomId := r.(string)
	v,_ := paramsMap["showName"]
	showName := v.(string)

	Logi("handleJoinRoom cid",cid,"roomid",roomId, "showName", showName)
	var room = roomManager.FindRoomById(roomId)
	if room == nil {
		Loge(roomId,"join roomid has not exist!")
		s.SendPacket(BuildCreateRoomError(cid, CODE_ROOM_NOT_EXIST))
		return
	}
	
	room.AddMember(s.accountId)
	Logi("Join room",room.adminId,room.roomId)
	s.SendPacket(BuildJoinRoomSuccess(cid, room))
}

func (s *Session)handleCreateRoomAndJoin(pkt *Packet) {
	cid := pkt.Cid
	paramsMap := pkt.Data.(map[string]any)
	r,_ := paramsMap["roomId"]
	roomId := r.(string)
	v,_ := paramsMap["showName"]
	showName := v.(string)
	Logi("handleCreateRoomAndJoin cid",cid,"roomid",r , "showName", showName)

	if roomManager.CheckRoomExist(roomId) {
		Loge(roomId,"roomid has exist.")
		s.SendPacket(BuildCreateRoomError(cid, CODE_ERR_ROOMIDREPEAT))
		return
	}

	newRoom := roomManager.CreateNewRoom(roomId, s.accountId)
	Logi("create new Room",newRoom.adminId,newRoom.roomId)
	s.SendPacket(BuildCreateRoomSuccess(cid, newRoom.roomId))
}



