package delivery

import (
	"TKPM/configs"
	"TKPM/internals/domain"
	"TKPM/internals/models"
	"TKPM/utils"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type RoomDelivery interface {
	UploadImage(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	GetRoomList(w http.ResponseWriter, r *http.Request)
	GetRoomById(w http.ResponseWriter, r *http.Request)
}

type roomDelivery struct {
	roomDomain    domain.Room
	storageConfig configs.Storage
}

func NewRoomDelivery(roomDomain domain.Room, storageConfig configs.Storage) RoomDelivery {
	return &roomDelivery{
		roomDomain:    roomDomain,
		storageConfig: storageConfig,
	}
}

func (d *roomDelivery) UploadImage(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	file, header, err := r.FormFile("image")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	name := strings.Split(header.Filename, ".")
	io.Copy(&buf, file)
	res, err := d.roomDomain.UploadImage(context.Background(), buf.Bytes(), name[0], header.Size, d.storageConfig)
	if err != nil {
		utils.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"message": err.Error()})
		return
	}

	utils.ResponseWithJson(w, http.StatusOK, map[string]string{"url": *res})
}

func (d *roomDelivery) Create(w http.ResponseWriter, r *http.Request) {

	var req models.Room

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"message": "Invalid body"})
		return
	}

	res, err := d.roomDomain.Create(context.Background(), &req)
	if err != nil {
		utils.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"message": err.Error()})
		return
	}
	utils.ResponseWithJson(w, http.StatusOK, res)

}

func (d *roomDelivery) Update(w http.ResponseWriter, r *http.Request) {

	var req models.Room

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"message": "Invalid body"})
		return
	}

	res, err := d.roomDomain.Update(context.Background(), &req)
	if err != nil {
		utils.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"message": err.Error()})
		return
	}
	utils.ResponseWithJson(w, http.StatusOK, res)
}

func (d *roomDelivery) GetRoomList(w http.ResponseWriter, r *http.Request) {
	offset := 0
	limit := 0
	offsetPra := r.URL.Query()["offset"]
	limitPra := r.URL.Query()["limit"]

	if len(offsetPra) > 0 {
		i, err := strconv.Atoi(offsetPra[0])
		if err != nil {
			utils.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"message": "Invalid body"})
			return
		}
		offset = i
	}

	if len(limitPra) > 0 {
		i, err := strconv.Atoi(limitPra[0])
		if err != nil {
			utils.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"message": "Invalid body"})
			return
		}
		limit = i
	}

	res, err := d.roomDomain.GetRoomList(context.Background(), int64(offset), int64(limit))
	if err != nil {
		utils.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"message": err.Error()})
		return
	}

	utils.ResponseWithJson(w, http.StatusOK, res)

}

func (d *roomDelivery) GetRoomById(w http.ResponseWriter, r *http.Request) {
	roomId := r.URL.Query()["roomId"]
	if len(roomId) == 0 {
		utils.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"message": "Invalid body"})
		return
	}

	res, err := d.roomDomain.GetRoomById(context.Background(), roomId[0])
	if err != nil {
		utils.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"message": err.Error()})
		return
	}
	utils.ResponseWithJson(w, http.StatusOK, res)
}
