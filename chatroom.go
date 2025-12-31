package main

import (
	"sync"
)

type ChatRoom struct {
	roomId   string
	roomName string
	adminId  int64
	members  map[int64]string
	lock sync.Mutex
}

func (c *ChatRoom)AddMember(accountId int64){
	defer c.lock.Unlock()
	c.lock.Lock()
	c.members[accountId] = ""
}

type ChatRoomManager struct {
	data map[string]*ChatRoom
	mutex sync.Mutex
}

func (c *ChatRoomManager) roomCount() int {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return len(c.data)
}

// 检查房间是否已经存在
func (c *ChatRoomManager) CheckRoomExist(roomId string) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	_, ok := c.data[roomId]
	return ok
}

func (c *ChatRoomManager) FindRoomById(roomId string) *ChatRoom {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	room, ok := c.data[roomId]
	if(ok){
		return room
	}
	return nil
}

func (c *ChatRoomManager) CreateNewRoom(roomId string, accountId int64) *ChatRoom {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	room := &ChatRoom{
		roomId:   roomId,
		roomName: roomId,
		adminId:  accountId,
		members:  make(map[int64]string),
	}

	room.AddMember(accountId)

	c.data[roomId] = room
	return room
}

func (c *ChatRoomManager) QuitRoom(roomId string, accountId int64) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	room, ok := c.data[roomId]
	if !ok {
		return false
	}

	_, ok = room.members[accountId]
	if !ok {
		return false
	}
	delete(room.members, accountId)
	return true
}
