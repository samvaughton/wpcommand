package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/samvaughton/wpcommand/v2/pkg/auth"
	"github.com/samvaughton/wpcommand/v2/pkg/flow"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	"github.com/samvaughton/wpcommand/v2/pkg/util"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func loadWpUsersHandler(resp http.ResponseWriter, req *http.Request) {
	site, userAccount := initHandlerWithSiteByUuid(resp, req)

	if site == nil {
		return // error handled by func
	}

	cachedData, err := site.GetWpCachedData()

	if err != nil {
		log.Error(err)
		util.HttpErrorEncode(resp, util.HttpStatusInternalServerError, "failed to decode wp cached data", err)
		return
	}

	json.NewEncoder(resp).Encode(types.NewApiWpUserListFromWpUserList(auth.FilterWpUserList(userAccount, cachedData.UserList)))
}

func createWpUserLoginUrlHandler(resp http.ResponseWriter, req *http.Request) {
	site, ua := initHandlerWithSiteByKey(resp, req)

	if site == nil {
		return // error handled by func
	}

	vars := mux.Vars(req)

	// verify user exists
	userId, err := strconv.Atoi(vars["userId"])

	if err != nil {
		util.HttpErrorEncode(resp, util.HttpStatusInvalidPayload, "invalid userId", err)
		return
	}

	var found *types.WpUser
	cachedData, err := site.GetWpCachedData()
	for _, item := range cachedData.UserList {
		if item.ID == userId {
			found = &item
			break
		}
	}

	if found == nil {
		log.Error(err)
		util.HttpErrorEncode(resp, util.HttpStatusNotFound, "could not find user", util.HttpEmptyErrors())
	}

	// does not have access
	if auth.WpUserHasImpersonateAccess(ua, found) == false {
		log.Error(err)
		util.HttpErrorEncode(resp, util.HttpStatusNotFound, "could not find user", util.HttpEmptyErrors())
	}

	loginUrl, err := flow.RunWpCreateUserLogin(site, vars["userId"], types.FlowOptions{LogSource: "API_CREATE_WP_USER_LOGIN"})

	if err != nil {
		log.Error(err)
		util.HttpErrorEncode(resp, util.HttpStatusInternalServerError, "failed to create user login", err)
		return
	}

	json.NewEncoder(resp).Encode(map[string]string{"LoginUrl": loginUrl})
}

func createWpUserHandler(resp http.ResponseWriter, req *http.Request) {
	site, _ := initHandlerWithSiteByUuid(resp, req)

	if site == nil {
		return // error handled by func
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
	vars := mux.Vars(req)
	site, userAccount := initHandlerWithSiteByUuid(resp, req)

	if site == nil {
		return // error handled by func
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
	vars := mux.Vars(req)
	site, userAccount := initHandlerWithSiteByUuid(resp, req)

	if site == nil {
		return // error handled by func
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
