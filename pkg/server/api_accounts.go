package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/samvaughton/wpcommand/v2/pkg/db"
	"github.com/samvaughton/wpcommand/v2/pkg/flow"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	"github.com/samvaughton/wpcommand/v2/pkg/util"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func loadAccountsHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")

	items, err := db.AccountsGetAll()

	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	json.NewEncoder(resp).Encode(items)
}

func loadAccountHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(req)
	accUuid := vars["accUuid"]

	item, err := db.AccountGetByUuid(accUuid)

	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	json.NewEncoder(resp).Encode(item)
}

func loadAccountUserApiItemsHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(req)

	accUuid := vars["accUuid"]

	account, err := db.AccountGetByUuid(accUuid)

	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	items, err := db.UsersGetByAccountIdSafe(account.Id)

	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	var apiItems = make([]*types.UserApiItem, 0)

	for _, user := range items {
		apiItems = append(apiItems, user.ConvertToApiItem())
	}

	json.NewEncoder(resp).Encode(apiItems)
}

func loadAccountUserApiItemHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(req)

	accUuid := vars["accUuid"]
	userUuid := vars["userUuid"]

	account, err := db.AccountGetByUuid(accUuid)

	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	item, err := db.UserGetByUuidSafe(userUuid, account.Id)

	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	json.NewEncoder(resp).Encode(item.ConvertToApiItem())
}

func createAccountHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")

	payload, err := types.NewCreateAccountPayloadFromHttpRequest(req)

	if err != nil {
		log.Error(err)
		util.HttpErrorEncode(resp, util.HttpStatusContentMalformed, "Could not decode payload.", util.HttpEmptyErrors())
		return
	}

	validationErrors := payload.Validate()

	if validationErrors != nil {
		resp.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(resp).Encode(map[string]interface{}{"Status": "VALIDATION_ERRORS", "Errors": validationErrors})

		return
	}

	account, err := db.AccountCreate(db.Db, payload.Name, util.Slug(payload.Name))

	if err != nil {
		log.Error(err)
		util.HttpErrorEncode(resp, util.HttpStatusInternalServerError, "Something went wrong when creating.", util.HttpEmptyErrors())
		return
	}

	json.NewEncoder(resp).Encode(account)
}

func updateAccountHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(req)
	accUuid := vars["accUuid"]

	account, err := db.AccountGetByUuid(accUuid)

	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	payload, err := types.NewUpdateAccountPayloadFromHttpRequest(req)

	if err != nil {
		log.Error(err)
		util.HttpErrorEncode(resp, util.HttpStatusContentMalformed, "Could not decode payload.", util.HttpEmptyErrors())
		return
	}

	validationErrors := payload.Validate()

	if validationErrors != nil {
		resp.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(resp).Encode(map[string]interface{}{"Status": "VALIDATION_ERRORS", "Errors": validationErrors})

		return
	}

	account.Name = payload.Name
	account.Key = util.Slug(payload.Name)

	err = db.AccountUpdate(db.Db, account)

	if err != nil {
		log.Error(err)
		util.HttpErrorEncode(resp, util.HttpStatusInternalServerError, "Something went wrong when updating.", util.HttpEmptyErrors())
		return
	}

	json.NewEncoder(resp).Encode(account)
}

func createAccountUserHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(req)

	accUuid := vars["accUuid"]

	account, err := db.AccountGetByUuid(accUuid)

	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	payload, err := types.NewAccountUserCreatePayloadFromHttpRequest(req)

	if err != nil {
		log.Error(err)
		util.HttpErrorEncode(resp, util.HttpStatusContentMalformed, "Could not decode payload.", util.HttpEmptyErrors())
		return
	}

	validationErrors := payload.Validate(db.UserUniqueEmailRule{})

	if validationErrors != nil {
		resp.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(resp).Encode(map[string]interface{}{"Status": "VALIDATION_ERRORS", "Errors": validationErrors})

		return
	}

	user, err := flow.UserCreate(db.Db, payload.Email, payload.FirstName, payload.LastName, payload.Password, []string{types.RoleMember}, account.Id)

	if err != nil {
		log.Error(err)
		util.HttpErrorEncode(resp, util.HttpStatusInternalServerError, "Something went wrong when creating.", util.HttpEmptyErrors())
		return
	}

	json.NewEncoder(resp).Encode(user.ConvertToApiItem())
}

func updateAccountUserHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(req)

	accUuid := vars["accUuid"]
	userUuid := vars["userUuid"]

	account, err := db.AccountGetByUuid(accUuid)

	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := db.UserGetByUuidSafe(userUuid, account.Id)

	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	payload, err := types.NewAccountUserUpdatePayloadFromHttpRequest(req)

	if err != nil {
		log.Error(err)
		util.HttpErrorEncode(resp, util.HttpStatusContentMalformed, "Could not decode payload.", util.HttpEmptyErrors())
		return
	}

	validationErrors := payload.Validate(db.UserUniqueEmailRule{IgnoreUserId: user.Id})

	if validationErrors != nil {
		resp.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(resp).Encode(map[string]interface{}{"Status": "VALIDATION_ERRORS", "Errors": validationErrors})

		return
	}

	user.Email = payload.Email
	user.FirstName = payload.FirstName
	user.LastName = payload.LastName

	err = flow.UserUpdate(db.Db, user, payload.Password)

	if err != nil {
		log.Error(err)
		util.HttpErrorEncode(resp, util.HttpStatusInternalServerError, "Something went wrong when updating.", util.HttpEmptyErrors())
		return
	}

	json.NewEncoder(resp).Encode(user.ConvertToApiItem())
}
