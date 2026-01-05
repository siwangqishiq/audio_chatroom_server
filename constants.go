package main

const CODE_SUCCESS = 200
const CODE_FAILED = 500
const CODE_ERR_ROOMIDREPEAT = 501
const CODE_ROOM_NOT_EXIST = 502  //房间不存在
const CODE_QUIT_ROOM_ERROR = 503 //退出房间错误

const CMD_LOGIN int = 1 //登录

const CMD_CREATE_ROOM_JOIN_REQ = 10  // 创建房间并加入
const CMD_CREATE_ROOM_JOIN_RESP = 11 // 创建房间并加入 响应

const CMD_JOIN_ROOM_REQ = 12  //加入房间
const CMD_JOIN_ROOM_RESP = 13 //加入房间响应

const CMD_QUIT_ROOM_REQ = 14  // 退出房间
const CMD_QUIT_ROOM_RESP = 15 // 退出房间响应

const CMD_FINISH_ROOM = 22 //结束会议
