# 集合表

## 用户表
```json
{
    "identity": "用户唯一标识",
	"account": "账号",
    "password": "密码",
    "nickname": "昵称",
    "sex": 1, // 0-未知 1-男 2-女
    "email": "邮箱",
    "avatar": "头像",
    "create_at": 1, //创建时间
    "update_at": 1, // 更新时间
}
```
## 消息集合
```json
{
    "user_identity": "用户的唯一标识",
	"room_identity": "房间的唯一标识",
	"data": "发送的数w[据",
	"create_at": 1, // 创建时间
	"update_at": 1, // 更新时间
}
```
## 房间集合
```json
{
    "number": "房间号",
	"name": "房间名称",
	"info": "房间简介",
	"user_identity": "房间创建者的唯一标识",
	"create_at": 1,
	"update_at": 1,
}
```
## 用户房间集合
```json
{
	"user_identity": "用户的唯一标识",
	"room_identity": "房间的唯一标识",
	"message_identity": "消息的唯一标识",
	"create_at": 1, // 创建时间
	"update_at": 1, // 更新时间
}
```