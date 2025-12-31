package main

import (
	"crypto/md5"
	"encoding/binary"
	"sync"

	"github.com/gorilla/websocket"
)

type ChatAccounts struct {
	value map[int64]*Session
	mutex sync.Mutex
}

func (c *ChatAccounts) AddNewAccount(conn *websocket.Conn) (int64, *Session){
	c.mutex.Lock()
	defer c.mutex.Unlock()

	remoteUrl := conn.RemoteAddr().String()
	var account int64 = HashStringMD5(remoteUrl)

	session := CreateNewSession(account, conn)
	c.value[account] = session

	return account, session
}


func (c *ChatAccounts) RemoveAccount(account int64) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	_,ok := c.value[account]
	if(!ok){
		return false
	}
	// defer conn.Close()
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