package main

type ChatRoom struct {
	roomId   string
	roomName string
	adminId  int64
	members  map[int64]string
}

type ChatRoomManager struct {
	data map[string]*ChatRoom
}

func (c *ChatRoomManager) roomCount() int {
	return len(c.data)
}

// 检查房间是否已经存在
func (c *ChatRoomManager) CheckRoomExist(roomId string) bool {
	_, ok := c.data[roomId]
	return ok
}

func (c *ChatRoomManager) CreateNewRoom(roomId string, accountId int64) *ChatRoom {
	room := &ChatRoom{
		roomId:   roomId,
		roomName: roomId,
		adminId:  accountId,
		members:  make(map[int64]string),
	}
	room.members[int64(accountId)] = ""

	c.data[roomId] = room
	return room
}

func (c *ChatRoomManager) QuitRoom(roomId string, accountId int64) bool {
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
