package models

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

type UserRoom struct {
	UserIdentity    string `bson:"user_identity"`
	RoomIdentity    string `bson:"room_identity"`
	RoomType        int    `bson:"room_type"` // 房间类型 1-独聊房间, 2-群聊房间
	MessageIdentity string `bson:"message_identity"`
	CreatedAt       int64  `bson:"created_at"`
	UpdateAt        int64  `bson:"updated_at"`
}

func (UserRoom) CollectionName() string {
	return "user_room"
}

func GetUserRoomByUserIdentityRoomIdentity(userIdentity, roomIdentity string) (*UserRoom, error) {
	ur := &UserRoom{}
	err := MongoDB.Collection(UserRoom{}.CollectionName()).
		FindOne(context.Background(), bson.D{
			{"user_identity", userIdentity},
			{"room_identity", roomIdentity},
		}).Decode(ur)
	return ur, err
}

func GetUserRoomByRoomIdentity(roomIdentity string) ([]*UserRoom, error) {
	cursor, err := MongoDB.Collection(UserRoom{}.CollectionName()).
		Find(context.Background(), bson.D{
			{
				"room_identity", roomIdentity,
			},
		})
	if err != nil {
		return nil, err
	}
	urs := make([]*UserRoom, 0)
	for cursor.Next(context.Background()) {
		ur := &UserRoom{}
		err := cursor.Decode(ur)
		if err != nil {
			return nil, err
		}
		urs = append(urs, ur)
	}
	return urs, nil
}

func JudgeUserIsFriend(u1, u2 string) bool {
	cursor, err := MongoDB.Collection(UserRoom{}.CollectionName()).
		Find(context.Background(), bson.D{
			{"user_identity", u1},
			{"room_type", 1},
		})
	roomIdentity := make([]string, 0)
	if err != nil {
		log.Printf("[DB ERROR: %V\n", err)
		return false
	}
	for cursor.Next(context.Background()) {
		ur := &UserRoom{}
		err := cursor.Decode(ur)
		if err != nil {
			return false
		}
		roomIdentity = append(roomIdentity, ur.RoomIdentity)
	}
	// 获取关联u2单聊房间个数
	cnt, err := MongoDB.Collection(UserRoom{}.CollectionName()).
		CountDocuments(context.Background(), bson.D{
			{"user_identity", u2},
			{"room_type", 1},
			{"room_identity", bson.M{"$in": roomIdentity}},
		})
	if err != nil {
		log.Printf("[DB ERROR : %v]\n", err)
		return false
	}
	if cnt > 0 {
		return true
	}
	return false
}

func InsertOneUserRoom(ur *UserRoom) error {
	_, err := MongoDB.Collection(UserRoom{}.CollectionName()).
		InsertOne(context.Background(), ur)
	return err
}
