package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/samvaughton/wpcommand/v2/pkg/db"
	"net/http"
)

func SetupPublic(router *mux.Router) {
	public := router.PathPrefix("/public").Subrouter()
	public.HandleFunc("/healthz", getHealthCheckHandler).Methods("GET")

	public.HandleFunc("/site/{uuid}/config", getSiteConfigHandler).Methods("GET")
	public.HandleFunc("/command/job/{jobUuid}/event/{eventUuid}", getCommandJobEventHandler).Methods("GET")
}

func getHealthCheckHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")

	json.NewEncoder(resp).Encode(map[string]string{
		"Status": "OK",
	})
}

func getSiteConfigHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(req)

	siteUuid, exists := vars["uuid"]

	if exists == false {
		resp.WriteHeader(http.StatusNotFound)
		return
	}

	site, err := db.SiteGetByUuid(siteUuid)

	if err != nil {
		resp.WriteHeader(http.StatusNotFound)
		return
	}

	resp.Write([]byte(site.SiteConfig))
}
