package main

type Packet struct {
	Cmd  int `json:"cmd"`
	Cid  int `json:"cid"`
	Code int `json:"code"`
	Data any `json:"data"`
}

type LoginData struct {
	Account int64 `json:"account"`
}

type CreateRoom struct {
	RoomId   string `json:"roomId"`
	ShowName string `json:"showName"`
}

type JoinRoom struct {
	RoomId   string `json:"roomId"`
	ShowName string `json:"showName"`
}

func BuildLoginData(account int64) Packet {
	data := LoginData{
		Account: account,
	}
	packet := Packet{
		Cmd:  CMD_LOGIN,
		Data: data,
		Cid:  0,
		Code: CODE_SUCCESS,
	}
	return packet
}

func BuildCreateRoomError(cid int, eCode int) Packet {
	packet := Packet{
		Cmd:  CMD_CREATE_ROOM_JOIN_RESP,
		Data: nil,
		Cid:  cid,
		Code: eCode,
	}
	return packet
}

func BuildCreateRoomSuccess(cid int, roomId string) Packet {
	room := CreateRoom{
		RoomId: roomId,
	}
	packet := Packet{
		Cmd:  CMD_CREATE_ROOM_JOIN_RESP,
		Data: room,
		Cid:  cid,
		Code: CODE_SUCCESS,
	}
	return packet
}

func BuildJoinRoomSuccess(cid int, room *ChatRoom) Packet {
	// memberList := make([]int64, 0, len(room.members))
	// memberMap := room.members
	// for k, _ := range memberMap {
	// 	memberList = append(memberList, k)
	// }

	joinRoom := JoinRoom{
		RoomId: room.roomId,
	}

	packet := Packet{
		Cmd:  CMD_JOIN_ROOM_RESP,
		Data: joinRoom,
		Cid:  cid,
		Code: CODE_SUCCESS,
	}
	return packet
}