package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/samvaughton/wpcommand/v2/pkg/config"
	"github.com/samvaughton/wpcommand/v2/pkg/db"
	"github.com/samvaughton/wpcommand/v2/pkg/execution"
	"github.com/samvaughton/wpcommand/v2/pkg/registry"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

var notImplemented = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Not Implemented"))
})

func SetupApi(router *mux.Router) {
	router.HandleFunc("/auth", authHandler).Methods("POST")

	api := router.PathPrefix("/api").Subrouter()
	api.Use(IsAuthorizedMiddleware)

	api.HandleFunc("/user", notImplemented)
	api.HandleFunc("/account", notImplemented)

	api.HandleFunc("/site", createSiteHandler).Methods("POST")
	api.HandleFunc("/site", loadSitesHandler).Methods("GET")
	api.HandleFunc("/site/{key}", loadSiteHandler).Methods("GET")
	api.HandleFunc("/command/job", createCommandJobHandler).Methods("POST")

	api.HandleFunc("/config", configHandler).Methods("GET")
}

func loadSiteHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")

	account := req.Context().Value("account").(*types.Account)

	vars := mux.Vars(req)
	key := vars["key"]

	sites, err := db.SiteGetByKey(key, account.Id)

	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	json.NewEncoder(resp).Encode(sites)
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

func createCommandJobHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")

	jobReq, err := types.NewApiCreateCommandJobRequest(req)

	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	account := req.Context().Value("account").(*types.Account)

	// we now need to validate this job request check command & site selector
	if registry.CommandExists(jobReq.Command) == false {
		resp.WriteHeader(http.StatusNotFound)
		json.NewEncoder(resp).Encode(map[string]interface{}{
			"Status": "COMMAND_NOT_FOUND",
		})
		return
	}

	sites, err := db.SelectSites(jobReq.Selector, account.Id)

	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(resp).Encode(map[string]interface{}{
			"Status": err.Error(),
		})
		return
	}

	if len(sites) == 0 {
		resp.WriteHeader(http.StatusNotFound)
		json.NewEncoder(resp).Encode(map[string]interface{}{
			"Status": "SITE_NOT_FOUND",
		})
		return
	}

	// create command job
	jobs := db.CreateCommandJobs(jobReq.Command, sites)

	if len(jobs) == 0 {
		log.Error(fmt.Sprintf("something went wrong creating jobs. command=%s selector=%s", jobReq.Command, jobReq.Selector))
		resp.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(resp).Encode(map[string]interface{}{
			"Status": "ERROR_CREATING_JOBS",
		})
		return
	}

	json.NewEncoder(resp).Encode(types.ApiCreateCommandJobResponse{
		Request:   *jobReq,
		JobStatus: types.CommandJobStatusCreated,
		Sites:     sites,
	})
}

func authHandler(resp http.ResponseWriter, req *http.Request) {
	var authPayload types.Authentication

	bytes, err := ioutil.ReadAll(req.Body)

	resp.Header().Set("Content-Type", "application/json")

	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(bytes, &authPayload)

	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	user := db.UserGetByEmailAndAccountKey(authPayload.Email, authPayload.Account)

	if user == nil {
		resp.WriteHeader(http.StatusNotFound)
		return
	}

	if CheckPasswordHash(authPayload.Password, user.Password) == false {
		resp.WriteHeader(http.StatusUnauthorized)
		return
	}

	// get the actual account with the key
	var account types.Account
	for _, accItem := range user.Accounts {
		if accItem.Key == authPayload.Account {
			account = *accItem
			break
		}
	}

	tokenString, err := GenerateJWT(user.Email, account.Uuid)

	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusInternalServerError)

		return
	}

	resp.WriteHeader(http.StatusOK)

	json.NewEncoder(resp).Encode(map[string]interface{}{
		"Status": "AUTHENTICATED",
		"Token":  tokenString,
		"Email":  user.Email,
		"Account": map[string]string{
			"Name": account.Name,
			"Key":  account.Key,
		},
	})
}

func configHandler(resp http.ResponseWriter, req *http.Request) {
	bytes, err := json.Marshal(config.Config)

	if err != nil {
		log.WithFields(log.Fields{"endpoint": "/api/config"}).Error(err)

		resp.WriteHeader(500)

		return
	}

	_, err = resp.Write(bytes)

	if err != nil {
		log.WithFields(log.Fields{"endpoint": "/api/config"}).Error(err)
	}
}
