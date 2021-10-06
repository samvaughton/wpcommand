package server

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/samvaughton/wpcommand/v2/pkg/auth"
	"github.com/samvaughton/wpcommand/v2/pkg/config"
	"github.com/samvaughton/wpcommand/v2/pkg/db"
	"github.com/samvaughton/wpcommand/v2/pkg/execution"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	"github.com/samvaughton/wpcommand/v2/pkg/util"
	log "github.com/sirupsen/logrus"
	"gopkg.in/guregu/null.v3"
	"net/http"
)

func loadSiteHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")

	userAccount := req.Context().Value("userAccount").(*types.UserAccount)

	vars := mux.Vars(req)
	key := vars["key"]

	site, err := db.SiteGetByKey(key, userAccount.AccountId)

	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	json.NewEncoder(resp).Encode(site)
}

func updateSiteHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")

	userAccount := req.Context().Value("userAccount").(*types.UserAccount)

	vars := mux.Vars(req)
	key := vars["siteUuid"]

	site, err := db.SiteGetByUuidSafe(key, userAccount.AccountId)

	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusNotFound)
		return
	}

	payload, err := types.NewUpdateSitePayloadFromHttpRequest(req)

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

	payload.HydrateSite(site)

	// we want to update all objects with the same uuid for this
	res, err := db.Db.NewUpdate().Model(site).WherePK().Exec(context.Background())

	if err != nil {
		log.Error(err)
		util.HttpErrorEncode(resp, util.HttpStatusInternalServerError, "Something went wrong.", util.HttpEmptyErrors())
		return
	}

	ra, err := res.RowsAffected()

	if err != nil || ra == 0 {
		log.Error(err)
		util.HttpErrorEncode(resp, util.HttpStatusInternalServerError, "Something went wrong.", util.HttpEmptyErrors())
		return
	}

	json.NewEncoder(resp).Encode(site)
}

func createSiteCommandHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")

	userAccount := req.Context().Value("userAccount").(*types.UserAccount)

	vars := mux.Vars(req)
	siteUuid := vars["key"]

	site, err := db.SiteGetByUuidSafe(siteUuid, userAccount.AccountId)

	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusNotFound)
		return
	}

	payload, err := types.NewCreateCommandPayloadFromHttpRequest(req)

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

	cmd := &types.Command{
		SiteId:    null.NewInt(site.Id, true),
		AccountId: null.NewInt(userAccount.AccountId, true),
	}

	payload.HydrateCommand(cmd)
	cmd.Key = util.Slug(payload.Description)

	// we want to update all objects with the same uuid for this
	res, err := db.Db.NewInsert().Model(cmd).Exec(context.Background())

	if err != nil {
		log.Error(err)
		util.HttpErrorEncode(resp, util.HttpStatusInternalServerError, "Something went wrong.", util.HttpEmptyErrors())
		return
	}

	ra, err := res.RowsAffected()

	if err != nil || ra == 0 {
		log.Error(err)
		util.HttpErrorEncode(resp, util.HttpStatusInternalServerError, "Something went wrong.", util.HttpEmptyErrors())
		return
	}

	json.NewEncoder(resp).Encode(cmd)
}

func updateSiteCommandHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")

	userAccount := req.Context().Value("userAccount").(*types.UserAccount)

	vars := mux.Vars(req)
	siteUuid := vars["siteUuid"]
	cmdUuid := vars["cmdUuid"]

	site, err := db.SiteGetByUuidSafe(siteUuid, userAccount.AccountId)

	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusNotFound)
		return
	}

	cmd, err := db.CommandGetByUuidAccountSafe(cmdUuid, site.AccountId)

	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusNotFound)
		return
	}

	payload, err := types.NewUpdateCommandPayloadFromHttpRequest(req)

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

	payload.HydrateCommand(cmd)
	cmd.Key = util.Slug(payload.Description)

	// we want to update all objects with the same uuid for this
	res, err := db.Db.NewUpdate().Model(cmd).WherePK().Exec(context.Background())

	if err != nil {
		log.Error(err)
		util.HttpErrorEncode(resp, util.HttpStatusInternalServerError, "Something went wrong.", util.HttpEmptyErrors())
		return
	}

	ra, err := res.RowsAffected()

	if err != nil || ra == 0 {
		log.Error(err)
		util.HttpErrorEncode(resp, util.HttpStatusInternalServerError, "Something went wrong.", util.HttpEmptyErrors())
		return
	}

	json.NewEncoder(resp).Encode(cmd)
}

func loadSiteCommandsHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")

	userAccount := req.Context().Value("userAccount").(*types.UserAccount)

	vars := mux.Vars(req)
	key := vars["key"]

	site, err := db.SiteGetByKey(key, userAccount.AccountId)

	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	var commands = make([]*types.Command, 0)

	if req.URL.Query().Get("type") == "runnable" {
		commands, err = db.CommandsGetRunnableForSiteSafe(site.Id, userAccount.AccountId)
	} else if req.URL.Query().Get("type") == "attached" {
		commands, err = db.CommandsGetAttachedForSiteSafe(site.Id)
	}

	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	json.NewEncoder(resp).Encode(auth.FilterCommandList(userAccount, commands))
}

func loadSiteBlueprintsHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")

	userAccount := req.Context().Value("userAccount").(*types.UserAccount)

	vars := mux.Vars(req)
	key := vars["key"]

	site, err := db.SiteGetByKey(key, userAccount.AccountId)

	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	items, err := db.BlueprintsGetForSiteSafe(site.Id, userAccount.AccountId)

	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	json.NewEncoder(resp).Encode(items)
}

func loadSitesHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")

	selector := "*"
	query := req.URL.Query()
	if query.Get("selector") != "" {
		selector = query.Get("selector")
	}

	account := req.Context().Value("account").(*types.Account)
	sites, err := db.SelectSites(selector, account.Id)

	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	json.NewEncoder(resp).Encode(sites)
}

func createSiteHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")

	site, err := types.NewSiteFromHttpRequest(req)

	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	// validate that the provided label selector + namespace matches a pod
	_, err = execution.GetPodBySite(site.LabelSelector, site.Namespace, config.Config.K8.LabelSelector, config.Config.K8RestConfig)

	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(resp).Encode(map[string]string{
			"Status":  "INVALID_CONFIGURATION",
			"Message": "Could not find a wordpress instance with the given configuration.",
		})
		return
	}

	account := req.Context().Value("account").(*types.Account)
	err = db.SiteCreateFromStruct(site, account.Id)

	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(resp).Encode(map[string]string{
			"Status":  "COULD_NOT_CREATE",
			"Message": "Something went wrong creating the site.",
		})
		return
	}

	json.NewEncoder(resp).Encode(site)
}
