package main

import (
	"sync"
)

type ChatRoom struct {
	roomId   string
	roomName string
	adminId  int64
	members  map[int64]string
	lock     sync.Mutex
}

func (c *ChatRoom) AddMember(accountId int64) {
	defer c.lock.Unlock()
	c.lock.Lock()
	c.members[accountId] = ""

	//set attach room
	memberSession := accounts.GetSessionById(accountId)
	if(memberSession != nil){
		memberSession.attachRoom = c
		Logi(accountId, "attach chatroom",c.roomId)
	}
}

func (c *ChatRoom) ForwardBytesData(bytesData []byte, senderAccount int64) {
	for accountId := range c.members {
		if accountId != senderAccount {
			session := accounts.GetSessionById(accountId)
			if session == nil {
				continue
			}

			// Logi("ForwardBytesData size ", len(bytesData))
			session.sendBinaryChan <- bytesData
		}
	} //end for each members
}

type ChatRoomManager struct {
	data  map[string]*ChatRoom
	mutex sync.Mutex
}

func (c *ChatRoomManager) RoomCount() int {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return len(c.data)
}

// 结束会议
func (c *ChatRoomManager) FinishRoom(roomId string) (delRoom *ChatRoom, result bool) {
	if c.CheckRoomExist(roomId) {
		c.mutex.Lock()
		defer c.mutex.Unlock()
		delRoom = c.data[roomId]
		delete(c.data, roomId)
		return delRoom, true
	}
	return nil, false
}

// 检查房间是否已经存在
func (c *ChatRoomManager) CheckRoomExist(roomId string) bool {
	// Logi("CheckRoomExist")
	c.mutex.Lock()
	defer c.mutex.Unlock()
	_, ok := c.data[roomId]

	// Logi("CheckRoomExist end")
	return ok
}

func (c *ChatRoomManager) FindRoomById(roomId string) *ChatRoom {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	room, ok := c.data[roomId]
	if ok {
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

	account := accounts.value[accountId]
	if account != nil {
		account.attachRoom = room
	}
	return room
}

func (c *ChatRoomManager) QuitRoom(roomId string, accountId int64) (ret bool, msg string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	room, ok := c.data[roomId]
	if !ok {
		return false, "room not exist"
	}

	_, ok = room.members[accountId]
	if !ok {
		return false, "account not in this room"
	}
	delete(room.members, accountId)

	account := accounts.value[accountId]
	if account != nil && account.attachRoom == room {
		account.attachRoom = nil
	}
	return true, ""
}
