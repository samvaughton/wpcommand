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
