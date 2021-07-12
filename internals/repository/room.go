package repository

import (
	"TKPM/common/enums"
	"TKPM/internals/models"
	"context"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

type Room interface {
	Create(ctx context.Context, room *models.Room) (*models.Room, error)
	Update(ctx context.Context, room *models.Room) (*models.Room, error)
	GetRoomList(ctx context.Context, offset, limit int64) ([]*models.Room, error)
	GetRoomById(ctx context.Context, roomId string) (*models.Room, error)
}

type roomRepository struct {
	coll *mongo.Collection
}

func NewRoomRepository(coll *mongo.Collection) Room {
	return &roomRepository{
		coll: coll,
	}
}

func (r *roomRepository) Create(ctx context.Context, room *models.Room) (*models.Room, error) {
	room.RoomID = uuid.New().String()
	room.Status = enums.Available.String()
	_, err := r.coll.InsertOne(ctx, room)
	return room, err
}

func (r *roomRepository) Update(ctx context.Context, room *models.Room) (*models.Room, error) {
	after := options.After
	res := r.coll.FindOneAndUpdate(ctx, bson.M{"room_id": room.RoomID}, bson.M{
		"$set": room,
	}, &options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	})
	if res.Err() != nil {
		return nil, res.Err()
	}

	var updatedRoom models.Room
	err := res.Decode(&updatedRoom)
	if err != nil {
		return nil, err
	}

	return &updatedRoom, nil
}

func (r *roomRepository) GetRoomList(ctx context.Context, offset, limit int64) ([]*models.Room, error) {
	result := []*models.Room{}

	opt := options.FindOptions{
		Skip:  &offset,
		Limit: &limit,
	}
	cur, err := r.coll.Find(ctx, &models.Room{}, &opt)

	if err != nil {
		return nil, err
	}

	if cur.Err() != nil {
		return nil, cur.Err()
	}

	if err := cur.All(ctx, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *roomRepository) GetRoomById(ctx context.Context, roomId string) (*models.Room, error) {
	var room models.Room
	err := r.coll.FindOne(ctx, bson.M{"room_id": roomId}).Decode(&room)
	if err != nil {
		return nil, err
	}

	return &room, nil
}
