package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/samvaughton/wpcommand/v2/pkg/db"
	"github.com/samvaughton/wpcommand/v2/pkg/registry"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func getCommandJobsHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")

	account := req.Context().Value("account").(*types.Account)
	items, err := db.CommandJobsGetForAccount(account.Id)

	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	json.NewEncoder(resp).Encode(items)
}

func getCommandJobHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")

	account := req.Context().Value("account").(*types.Account)

	vars := mux.Vars(req)
	uuid := vars["uuid"]

	item, err := db.CommandJobGetByUuidSafe(uuid, account.Id)

	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	json.NewEncoder(resp).Encode(item)
}

func getCommandJobEventsHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")

	account := req.Context().Value("account").(*types.Account)

	vars := mux.Vars(req)
	uuid := vars["uuid"]

	item, err := db.CommandJobGetByUuidSafe(uuid, account.Id)

	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	events, err := db.CommandJobEventLogGetByJob(item.Id)

	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	json.NewEncoder(resp).Encode(events)
}

func createCommandJobHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")

	jobReq, err := types.NewApiCreateCommandJobRequest(req)

	if err != nil || jobReq.CommandId == 0 || jobReq.Selector == "" {
		log.Error(err)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	account := req.Context().Value("account").(*types.Account)
	user := req.Context().Value("user").(*types.User)

	// check database
	command, err := db.CommandGetByIdAccountSafe(jobReq.CommandId, account.Id)

	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusNotFound)
		json.NewEncoder(resp).Encode(map[string]interface{}{
			"Status": "COMMAND_NOT_FOUND",
		})
		return
	}

	// we now need to validate this job request check command & site selector
	if command.Type == types.CommandTypeWpBuiltIn && registry.CommandExists(command.Key) == false {
		resp.WriteHeader(http.StatusNotFound)
		json.NewEncoder(resp).Encode(map[string]interface{}{
			"Status": "BUILT_IN_COMMAND_NOT_FOUND",
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
	jobs := db.CreateCommandJobs(command, sites, user.Id)

	if len(jobs) == 0 {
		log.Error(fmt.Sprintf("something went wrong creating jobs. command=%s selector=%s", command.Key, jobReq.Selector))
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
		Jobs:      jobs,
	})
}
