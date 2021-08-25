package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/samvaughton/wpcommand/v2/pkg/config"
	"github.com/samvaughton/wpcommand/v2/pkg/db"
	"github.com/samvaughton/wpcommand/v2/pkg/execution"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func loadSiteHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")

	account := req.Context().Value("account").(*types.Account)

	vars := mux.Vars(req)
	key := vars["key"]

	site, err := db.SiteGetByKey(key, account.Id)

	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	json.NewEncoder(resp).Encode(site)
}

func loadSiteCommandsHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")

	account := req.Context().Value("account").(*types.Account)

	vars := mux.Vars(req)
	key := vars["key"]

	site, err := db.SiteGetByKey(key, account.Id)

	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	commands, err := db.CommandsGetForSite(site.Id, account.Id)

	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	json.NewEncoder(resp).Encode(commands)
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
