package main

type Packet struct {
	Cmd  int `json:"cmd"`
	Cid  int `json:"cid"`
	Data any `json:"data"`
}

type LoginData struct {
	Account int64 `json:"account"`
}

type CreateRoom struct {
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
	}
	return packet
}