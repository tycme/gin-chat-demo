package models

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func GetMessageListByRoomIdentity(roomIdentity string, limit, skip *int64) ([]*MsgBasic, error) {
	data := make([]*MsgBasic, 0)
	cursor, err := MongoDB.Collection(MsgBasic{}.CollectionName()).
		Find(context.Background(), bson.M{
			"room_identity": roomIdentity,
		}, &options.FindOptions{
			Limit: limit,
			Skip:  skip,
			Sort: bson.D{
				{"created_at", -1},
			},
		})
	if err != nil {
		return nil, err
	}
	for cursor.Next(context.Background()) {
		mb := new(MsgBasic)
		err = cursor.Decode(mb)
		if err != nil {
			return nil, err
		}
		data = append(data, mb)
	}
	return data, nil
}
