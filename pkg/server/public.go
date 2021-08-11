package server

import (
	"github.com/gorilla/mux"
	"github.com/samvaughton/wpcommand/v2/pkg/db"
	"net/http"
)

func SetupPublic(router *mux.Router) {
	public := router.PathPrefix("/public").Subrouter()
	public.HandleFunc("/site/{uuid}/config", getSiteConfigHandler).Methods("GET")
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
