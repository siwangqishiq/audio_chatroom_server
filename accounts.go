package main

import (
	"crypto/md5"
	"encoding/binary"

	"github.com/gorilla/websocket"
)

type ChatAccounts struct {
	value map[int64]*websocket.Conn
}

func (c *ChatAccounts) AddNewAccount(conn *websocket.Conn) int64{
	remoteUrl := conn.RemoteAddr().String()
	var account int64 = HashStringMD5(remoteUrl)
	c.value[account] = conn
	return account
}

func (c *ChatAccounts) RemoveAccount(account int64) bool {
	conn,ok := c.value[account]
	if(!ok){
		return false
	}
	defer conn.Close()
	return true
}

func HashStringMD5(s string) int64 {
	h := md5.New()
    h.Write([]byte(s))
    bytes := h.Sum(nil)
    value := int64(binary.LittleEndian.Uint64(bytes[:8]))
	if(value < 0){
		return -value
	}
	return value
}