package main

type Message struct {
	Cmd  int `json:"cmd"`
	Sid  int `json:"sid"`
	Data any `json:"data"`
}

type LoginData struct {
	Account int64 `json:"account"`
}

func BuildLoginData(account int64) Message {
	data := LoginData{
		Account: account,
	}
	msg := Message{
		Cmd:  CMD_LOGIN,
		Data: data,
		Sid:  0,
	}
	return msg
}