package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/samvaughton/wpcommand/v2/pkg/db"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func loadBlueprintsHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")

	userAccount := req.Context().Value("userAccount").(*types.UserAccount)
	items, err := db.BlueprintsGetForAccountId(userAccount.AccountId)

	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	json.NewEncoder(resp).Encode(items)
}

func loadBlueprintHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")

	userAccount := req.Context().Value("userAccount").(*types.UserAccount)

	vars := mux.Vars(req)
	key := vars["uuid"]

	item, err := db.BlueprintGetByUuidSafe(key, userAccount.AccountId)

	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	json.NewEncoder(resp).Encode(item)
}

func deleteBlueprintHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")

	userAccount := req.Context().Value("userAccount").(*types.UserAccount)

	vars := mux.Vars(req)
	key := vars["uuid"]

	item, err := db.BlueprintGetByUuidSafe(key, userAccount.AccountId)

	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = db.Db.Query("DELETE FROM blueprint_sets WHERE id = ?", item.Id)

	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	json.NewEncoder(resp).Encode(item)
}

func loadBlueprintObjectsHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")

	userAccount := req.Context().Value("userAccount").(*types.UserAccount)

	vars := mux.Vars(req)
	key := vars["uuid"]

	item, err := db.BlueprintGetByUuidSafe(key, userAccount.AccountId)

	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	items, err := db.GetLatestBlueprintObjectsForBlueprintSetId(item.Id)

	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	json.NewEncoder(resp).Encode(items)
}

func createBlueprintHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")

	payload, err := types.NewCreateBlueprintSetPayloadFromHttpRequest(req)

	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	userAccount := req.Context().Value("userAccount").(*types.UserAccount)
	blueprint, err := db.BlueprintCreateFromPayload(payload, userAccount.AccountId)

	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(resp).Encode(map[string]string{
			"Status":  "COULD_NOT_CREATE",
			"Message": "Something went wrong when creating.",
		})
		return
	}

	json.NewEncoder(resp).Encode(blueprint)
}

func createBlueprintObjectHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")

	payload, err := types.NewCreateObjectBlueprintPayloadFromHttpRequest(req)

	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	validationErrors := payload.Validate()

	if validationErrors != nil {
		resp.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(resp).Encode(map[string]interface{}{"Status": "VALIDATION_ERRORS", "Errors": validationErrors})

		return
	}

	account := req.Context().Value("account").(*types.Account)
	object, err := db.BlueprintObjectCreateFromPayload(payload, account.Id)

	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(resp).Encode(map[string]string{
			"Status":  "COULD_NOT_CREATE",
			"Message": "Something went wrong creating the object.",
		})
		return
	}

	json.NewEncoder(resp).Encode(object)
}
