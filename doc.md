# Audio ChatRoom

## Android Client API

```
// 连接服务
void connect(wsServer:String,callback:Callback) 
```

```
// 关闭服务
void close()
```

```
//创建房间
void createRoom(roomId:String, account : String, callback:Callback) 
```

```
//加入房间
void joinRoom(roomId:String, account: String, callback:Callback)
```

```
//退出房间
void quitRoom(roomId:String, callback:Callback)
```

```
//关闭房间
void closeRoom(roomId:String, callback:Callback)
```
