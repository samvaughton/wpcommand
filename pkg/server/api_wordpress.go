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
	"strconv"
)

func loadWpUsersHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(req)

	userAccount := req.Context().Value("userAccount").(*types.UserAccount)

	site, err := db.SiteGetByUuidSafe(vars["siteUuid"], userAccount.AccountId)

	if err != nil {
		log.Error(err)
		util.HttpErrorEncode(resp, util.HttpStatusNotFound, "could not find site", util.HttpEmptyErrors())
		return
	}

	cachedData, err := site.GetWpCachedData()

	if err != nil {
		log.Error(err)
		util.HttpErrorEncode(resp, util.HttpStatusInternalServerError, "failed to decode wp cached data", err)
		return
	}

	json.NewEncoder(resp).Encode(types.NewApiWpUserListFromWpUserList(auth.FilterWpUserList(userAccount, cachedData.UserList)))
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

	payload, err := types.NewCreateWpUserPayloadFromHttpRequest(req)

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

	err = flow.RunWpUserCreate(payload, site, types.FlowOptions{LogSource: "API_CREATE_WP_USER"})

	if err != nil {
		util.HttpErrorEncode(resp, util.HttpStatusInternalServerError, "failed to create user on website", err)
		return
	}

	err = flow.RunWpUserSync(site, types.FlowOptions{LogSource: "API_CREATE_WP_USER"})

	if err != nil {
		util.HttpErrorEncode(resp, util.HttpStatusInternalServerError, "failed to create user on website", err)
		return
	}

	json.NewEncoder(resp).Encode(map[string]string{"Status": "OK"})
}

func updateWpUserHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(req)
	userAccount := req.Context().Value("userAccount").(*types.UserAccount)

	site, err := db.SiteGetByUuidSafe(vars["siteUuid"], userAccount.AccountId)

	if err != nil {
		util.HttpErrorEncode(resp, util.HttpStatusNotFound, "could not find site", util.HttpEmptyErrors())
		return
	}

	cachedData, err := site.GetWpCachedData()

	if err != nil {
		util.HttpErrorEncode(resp, util.HttpStatusInternalServerError, "failed to decode wp cached data", err)
		return
	}

	updateUserPayload, err := types.NewUpdateWpUserPayloadFromHttpRequest(req)

	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	wpId, err := strconv.Atoi(vars["userId"])

	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	user := auth.FindInWpUserList(userAccount, cachedData.UserList, wpId)

	if user == nil {
		util.HttpErrorEncode(resp, util.HttpStatusNotFound, "could not find user", util.HttpEmptyErrors())
		return
	}

	err = flow.RunWpUserUpdate(wpId, updateUserPayload, site, types.FlowOptions{LogSource: "API_UPDATE_WP_USER"})

	if err != nil {
		util.HttpErrorEncode(resp, util.HttpStatusInternalServerError, "failed to update user on website", err)
		return
	}

	err = flow.RunWpUserSync(site, types.FlowOptions{LogSource: "API_UPDATE_WP_USER"})

	if err != nil {
		util.HttpErrorEncode(resp, util.HttpStatusInternalServerError, "failed to update user on website", err)
		return
	}

	json.NewEncoder(resp).Encode(map[string]string{"Status": "OK"})
}

func deleteWpUserHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(req)
	userAccount := req.Context().Value("userAccount").(*types.UserAccount)

	site, err := db.SiteGetByUuidSafe(vars["siteUuid"], userAccount.AccountId)

	if err != nil {
		log.Error(err)
		util.HttpErrorEncode(resp, util.HttpStatusNotFound, "could not find site", util.HttpEmptyErrors())
		return
	}

	cachedData, err := site.GetWpCachedData()

	if err != nil {
		log.Error(err)
		util.HttpErrorEncode(resp, util.HttpStatusInternalServerError, "failed to decode wp cached data", err)
		return
	}

	wpId, err := strconv.Atoi(vars["userId"])

	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	user := auth.FindInWpUserList(userAccount, cachedData.UserList, wpId)

	if user == nil {
		log.Error(err)
		util.HttpErrorEncode(resp, util.HttpStatusNotFound, "could not find user", util.HttpEmptyErrors())
		return
	}

	err = flow.RunWpUserDelete(wpId, site, types.FlowOptions{LogSource: "API_DELETE_WP_USER"})

	if err != nil {
		log.Error(err)
		util.HttpErrorEncode(resp, util.HttpStatusInternalServerError, "failed to delete user on website", err)
		return
	}

	err = flow.RunWpUserSync(site, types.FlowOptions{LogSource: "API_DELETE_WP_USER"})

	if err != nil {
		util.HttpErrorEncode(resp, util.HttpStatusInternalServerError, "failed to update user on website", err)
		return
	}

	json.NewEncoder(resp).Encode(map[string]string{"Status": "OK"})
}
