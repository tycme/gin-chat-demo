package models

import (
	"context"
)

type MsgBasic struct {
	UserIdentity string `bson:"user_identity"`
	RoomIdentity string `bson:"room_identity"`
	Data         string `bson:"data"`
	CreatedAt    int64  `bson:"created_at"`
	UpdateAt     int64  `bson:"updated_at"`
}

func (MsgBasic) CollectionName() string {
	return "message_basic"
}

func InsertOneMsg(mb *MsgBasic) error {
	_, err := MongoDB.Collection(MsgBasic{}.CollectionName()).
		InsertOne(context.Background(), mb)
	return err
}
