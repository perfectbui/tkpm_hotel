package delivery

import (
	"TKPM/common/enums"
	"TKPM/internals/domain"
	"TKPM/internals/models"
	"TKPM/utils"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
)

type ContractDelivery interface {
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	GetContractList(w http.ResponseWriter, r *http.Request)
	GetContractById(w http.ResponseWriter, r *http.Request)
}

type contractDelivery struct {
	contractDomain domain.Contract
	roomDomain     domain.Room
}

func NewContractDelivery(contractDomain domain.Contract, roomDomain domain.Room) ContractDelivery {
	return &contractDelivery{
		contractDomain: contractDomain,
		roomDomain:     roomDomain,
	}
}

func (d *contractDelivery) Create(w http.ResponseWriter, r *http.Request) {

	var req models.Contract

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"message": "Invalid body"})
		return
	}

	ctx := context.Background()
	// Need transaction //
	contract, err := d.contractDomain.Create(ctx, &req)
	if err != nil {
		utils.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"message": err.Error()})
		return
	}

	res, err := d.roomDomain.Update(ctx, &models.Room{
		RoomID: contract.RoomID,
		Status: enums.Booked.String(),
	})
	if err != nil {
		utils.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"message": err.Error()})
		return
	}
	// Need transaction //

	utils.ResponseWithJson(w, http.StatusOK, res)

}

func (d *contractDelivery) Update(w http.ResponseWriter, r *http.Request) {

	var req models.Contract

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"message": "Invalid body"})
		return
	}

	// Need transaction //
	ctx := context.Background()
	contract, err := d.contractDomain.Update(ctx, &req)
	if err != nil {
		utils.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"message": err.Error()})
		return
	}

	if contract.Status == enums.Completed.String() || contract.Status == enums.Cancel.String() {
		_, err := d.roomDomain.Update(ctx, &models.Room{
			RoomID: contract.RoomID,
			Status: enums.Available.String(),
		})
		if err != nil {
			utils.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"message": err.Error()})
			return
		}
	}
	// Need transaction //

	utils.ResponseWithJson(w, http.StatusOK, contract)
}

func (d *contractDelivery) GetContractList(w http.ResponseWriter, r *http.Request) {

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

	res, err := d.contractDomain.GetContractList(context.Background(), int64(offset), int64(limit))
	if err != nil {
		utils.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"message": err.Error()})
		return
	}
	utils.ResponseWithJson(w, http.StatusOK, res)

}

func (d *contractDelivery) GetContractById(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query()["userId"]
	if len(userId) == 0 {
		utils.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"message": "Invalid body"})
		return
	}

	res, err := d.contractDomain.GetContractByUserId(context.Background(), userId[0])
	if err != nil {
		utils.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"message": err.Error()})
		return
	}
	utils.ResponseWithJson(w, http.StatusOK, res)
}
