package models

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

type UserRoom struct {
	UserIdentity    string `bson:"user_identity"`
	RoomIdentity    string `bson:"room_identity"`
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
