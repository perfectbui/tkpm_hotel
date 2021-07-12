package delivery

import (
	"TKPM/internals/domain"
	"TKPM/internals/models"
	"TKPM/utils"
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
)

type AccountDelivery interface {
	SignUp(w http.ResponseWriter, r *http.Request)
	SignIn(w http.ResponseWriter, r *http.Request)
	GetAccountById(w http.ResponseWriter, r *http.Request)

	CheckAuth(accountId string) (*models.Account, error)
}

type accountDelivery struct {
	accountDomain domain.Account
}

func NewAccountDelivery(accountDomain domain.Account) AccountDelivery {
	return &accountDelivery{
		accountDomain: accountDomain,
	}
}

func (d *accountDelivery) SignUp(w http.ResponseWriter, r *http.Request) {

	var req models.Account

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"message": "Invalid body"})
		return
	}

	res, err := d.accountDomain.SignUp(context.Background(), &req)
	if err != nil {
		utils.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"message": err.Error()})
		return
	}

	utils.ResponseWithJson(w, http.StatusOK, res)
}

func (d *accountDelivery) SignIn(w http.ResponseWriter, r *http.Request) {
	var req models.Account

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"message": "Invalid body"})
		return
	}

	res, err := d.accountDomain.SignIn(context.Background(), &req)
	if err != nil {
		utils.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"message": err.Error()})
		return
	}

	utils.ResponseWithJson(w, http.StatusOK, &models.Account{
		Token: genToken(res),
	})
}

func (d *accountDelivery) GetAccountById(w http.ResponseWriter, r *http.Request) {

	accountId := r.URL.Query()["accountId"]
	if len(accountId) == 0 {
		utils.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"message": "Invalid body"})
		return
	}

	res, err := d.accountDomain.CheckAuth(context.Background(), accountId[0])
	if err != nil {
		utils.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"message": err.Error()})
		return
	}

	utils.ResponseWithJson(w, http.StatusOK, res)
}

func (d *accountDelivery) CheckAuth(accountId string) (*models.Account, error) {
	res, err := d.accountDomain.CheckAuth(context.Background(), accountId)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func genToken(account *models.Account) string {

	bytes, err := json.Marshal(account)
	if err != nil {
		return ""
	}

	token := base64.StdEncoding.EncodeToString(bytes)
	return string(token)
}
