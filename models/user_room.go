package models

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
