package server

import (
	"github.com/gorilla/mux"
)

func SetupInternal(router *mux.Router) {
	public := router.PathPrefix("/internal").Subrouter()
	public.HandleFunc("/site/{siteUuid}/build-release", runSiteBuild).Methods("POST")
	public.HandleFunc("/site/{siteUuid}/build-preview", runSitePreview).Methods("POST")
}
