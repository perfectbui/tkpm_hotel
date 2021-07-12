package models

const RoomCollection = "room"

type Room struct {
	RoomID      string  `json:"roomId,omitempty" bson:"room_id,omitempty"`
	Name        string  `json:"name,omitempty" bson:"name,omitempty"`
	Price       int     `json:"price,omitempty" bson:"price,omitempty"`
	Status      string  `json:"status,omitempty" bson:"status,omitempty"`
	Description string  `json:"description,omitempty" bson:"description,omitempty"`
	Rate        float32 `json:"rate,omitempty" bson:"rate,omitempty"`
	Image       string  `json:"image,omitempty" bson:"image,omitempty"`
}
