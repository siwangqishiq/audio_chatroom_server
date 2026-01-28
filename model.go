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

type RemotePeer struct {
	AccountId int64  `json:"accountId"`
	Name      string `json:"name"`
}

type CreateRoom struct {
	RoomId   string       `json:"roomId"`
	ShowName string       `json:"showName"`
	Members  []RemotePeer `json:"members"`
}

type JoinRoom struct {
	RoomId   string       `json:"roomId"`
	ShowName string       `json:"showName"`
	Members  []RemotePeer `json:"members"`
}

type QuitRoom struct {
	RoomId string `json:"roomId"`
}

type FinishRoom struct {
	RoomId string `json:"roomId"`
}

func BuildFinishRoom(roomId string) Packet {
	data := FinishRoom{
		RoomId: roomId,
	}
	packet := Packet{
		Cmd:  CMD_FINISH_ROOM,
		Data: data,
		Cid:  0,
		Code: CODE_SUCCESS,
	}
	return packet
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
	room.Members = FindMembersByRoomId(roomId)

	packet := Packet{
		Cmd:  CMD_CREATE_ROOM_JOIN_RESP,
		Data: room,
		Cid:  cid,
		Code: CODE_SUCCESS,
	}
	return packet
}

func BuildJoinRoomSuccess(cid int, room *ChatRoom) Packet {
	joinRoom := JoinRoom{
		RoomId: room.roomId,
	}
	joinRoom.Members = FindMembersByRoomId(room.roomId)

	packet := Packet{
		Cmd:  CMD_JOIN_ROOM_RESP,
		Data: joinRoom,
		Cid:  cid,
		Code: CODE_SUCCESS,
	}
	return packet
}

func FindMembersByRoomId(roomId string) []RemotePeer {
	var members []RemotePeer = make([]RemotePeer, 0)
	room := roomManager.FindRoomById(roomId)
	if room == nil {
		return members
	}

	for k, v := range room.members {
		members = append(members, RemotePeer{
			AccountId: k,
			Name:      v,
		})
	} //end for
	return members
}

func BuildQuitRoomError(cid int, eCode int) Packet {
	packet := Packet{
		Cmd:  CMD_QUIT_ROOM_RESP,
		Data: nil,
		Cid:  cid,
		Code: eCode,
	}
	return packet
}

func BuildQuitRoomSuccess(cid int, room *ChatRoom) Packet {
	joinRoom := QuitRoom{
		RoomId: room.roomId,
	}

	packet := Packet{
		Cmd:  CMD_QUIT_ROOM_RESP,
		Data: joinRoom,
		Cid:  cid,
		Code: CODE_SUCCESS,
	}
	return packet
}