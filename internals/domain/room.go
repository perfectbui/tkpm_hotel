package domain

import (
	"TKPM/configs"
	"TKPM/internals/models"
	"TKPM/internals/repository"
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Room interface {
	UploadImage(ctx context.Context, data []byte, name string, size int64, storageConfig configs.Storage) (*string, error)
	Create(ctx context.Context, room *models.Room) (*models.Room, error)
	Update(ctx context.Context, room *models.Room) (*models.Room, error)
	GetRoomList(ctx context.Context, offset, limit int64) ([]*models.Room, error)
	GetRoomById(ctx context.Context, roomId string) (*models.Room, error)
}

type roomDomain struct {
	roomRepository repository.Room
}

func NewRoomDomain(roomRepository repository.Room) Room {
	return &roomDomain{
		roomRepository: roomRepository,
	}
}

func (d *roomDomain) UploadImage(ctx context.Context, data []byte, fileName string, size int64, storageConfig configs.Storage) (*string, error) {
	creds := credentials.NewStaticCredentials(storageConfig.AccessKey, storageConfig.SecretKey, "")
	_, err := creds.Get()
	if err != nil {
		fmt.Printf("Bad credentials: %s", err)
		return nil, err
	}

	cfg := aws.NewConfig().WithRegion(storageConfig.Region).WithCredentials(creds)
	ss, err := session.NewSession()
	if err != nil {
		return nil, err
	}

	svc := s3.New(ss, cfg)
	fileBytes := bytes.NewReader(data)
	fileType := http.DetectContentType(data)
	path := "/" + fileName
	params := &s3.PutObjectInput{
		Bucket:        aws.String(storageConfig.BucketName),
		Key:           aws.String(path),
		Body:          fileBytes,
		ContentLength: aws.Int64(size),
		ContentType:   aws.String(fileType),
	}

	_, err = svc.PutObject(params)
	if err != nil {
		fmt.Printf("can not upload image : %s", err)
		return nil, err
	}

	url := fmt.Sprintf("https://%v.s3.%v.amazonaws.com%v", storageConfig.BucketName, storageConfig.Region, path)

	return &url, nil
}

func (d *roomDomain) Create(ctx context.Context, room *models.Room) (*models.Room, error) {
	res, err := d.roomRepository.Create(ctx, room)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (d *roomDomain) Update(ctx context.Context, room *models.Room) (*models.Room, error) {
	res, err := d.roomRepository.Update(ctx, room)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (d *roomDomain) GetRoomList(ctx context.Context, offset, limit int64) ([]*models.Room, error) {
	res, err := d.roomRepository.GetRoomList(ctx, offset, limit)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (d *roomDomain) GetRoomById(ctx context.Context, roomId string) (*models.Room, error) {
	res, err := d.roomRepository.GetRoomById(ctx, roomId)
	if err != nil {
		return nil, err
	}

	return res, nil
}
