package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/samvaughton/wpcommand/v2/pkg/auth"
	"github.com/samvaughton/wpcommand/v2/pkg/db"
	"github.com/samvaughton/wpcommand/v2/pkg/flow"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	"github.com/samvaughton/wpcommand/v2/pkg/util"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func loadWpUsersHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(req)

	userAccount := req.Context().Value("userAccount").(*types.UserAccount)

	site, err := db.SiteGetByUuidSafe(vars["siteUuid"], userAccount.AccountId)

	if err != nil {
		log.Error(err)
		util.HttpError(util.HttpStatusNotFound, "could not find site", util.HttpEmptyErrors())
		return
	}

	cachedData, err := site.GetWpCachedData()

	if err != nil {
		log.Error(err)
		util.HttpError(util.HttpStatusInternalServerError, "failed to decode wp cached data", err)
		return
	}

	json.NewEncoder(resp).Encode(auth.FilterWpUserList(userAccount, cachedData.UserList))
}

func createWpUserHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(req)
	userAccount := req.Context().Value("userAccount").(*types.UserAccount)

	site, err := db.SiteGetByUuidSafe(vars["siteUuid"], userAccount.AccountId)

	if err != nil {
		log.Error(err)
		util.HttpError(util.HttpStatusNotFound, "could not find site", util.HttpEmptyErrors())
		return
	}

	newUserPayload, err := types.NewCreateWpUserPayloadFromHttpRequest(req)

	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	flow.RunWpUserCreate(newUserPayload, site, types.FlowOptions{LogSource: "API_CREATE_WP_USER"})

	json.NewEncoder(resp).Encode(map[string]string{"Status": "OK"})
}

func updateWpUserHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(req)
	userAccount := req.Context().Value("userAccount").(*types.UserAccount)

	site, err := db.SiteGetByUuidSafe(vars["siteUuid"], userAccount.AccountId)

	if err != nil {
		log.Error(err)
		util.HttpError(util.HttpStatusNotFound, "could not find site", util.HttpEmptyErrors())
		return
	}

	cachedData, err := site.GetWpCachedData()

	if err != nil {
		log.Error(err)
		util.HttpError(util.HttpStatusInternalServerError, "failed to decode wp cached data", err)
		return
	}

	updateUserPayload, err := types.NewUpdateWpUserPayloadFromHttpRequest(req)

	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	user := auth.FindInWpUserList(userAccount, cachedData.UserList, updateUserPayload.Id)

	if user == nil {
		log.Error(err)
		util.HttpError(util.HttpStatusNotFound, "could not find user", util.HttpEmptyErrors())
		return
	}

	err = flow.RunWpUserUpdate(updateUserPayload, site, types.FlowOptions{LogSource: "API_UPDATE_WP_USER"})

	if err != nil {
		log.Error(err)
		util.HttpError(util.HttpStatusInternalServerError, "failed to update user on website", err)
		return
	}

	json.NewEncoder(resp).Encode(map[string]string{"Status": "OK"})
}
