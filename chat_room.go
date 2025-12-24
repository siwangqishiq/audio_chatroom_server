package main

type ChatRoom struct {
	roomId   string
	roomName string
}

type ChatRoomManager struct {
	data map[string]*ChatRoom
}

func (c *ChatRoomManager) roomCount() int {
	return len(c.data)
}