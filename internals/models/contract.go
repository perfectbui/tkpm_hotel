package models

import "time"

const ContractCollection = "contract"

type Contract struct {
	ContractID string     `json:"contractId,omitempty" bson:"contract_id,omitempty"`
	AdminID    string     `json:"adminId,omitempty" bson:"admin_id,omitempty"`
	UserID     string     `json:"userId,omitempty" bson:"user_id,omitempty"`
	RoomID     string     `json:"roomId,omitempty" bson:"room_id,omitempty"`
	Price      int        `json:"price,omitempty" bson:"price,omitempty"`
	CreatedAt  *time.Time `json:"createdAt,omitempty" bson:"created_at,omitempty"`
	StartTime  *time.Time `json:"startTime,omitempty" bson:"start_time,omitempty"`
	EndTime    *time.Time `json:"endTime,omitempty" bson:"end_time,omitempty"`
	Status     string     `json:"status,omitempty" bson:"status,omitempty"`
}
