package server

import (
	"github.com/gorilla/mux"
	"github.com/samvaughton/wpcommand/v2/pkg/auth"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	log "github.com/sirupsen/logrus"
	"net/http"
)

var notImplemented = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
})

func AuthWrapper(obj string, act string, handler func(resp http.ResponseWriter, req *http.Request)) func(resp http.ResponseWriter, req *http.Request) {
	return func(resp http.ResponseWriter, req *http.Request) {
		ua := req.Context().Value("userAccount").(*types.UserAccount)

		allowed, err := auth.Enforcer.Enforce(ua.GetCasbinPolicyKey(), obj, act)

		if err != nil {
			log.Errorf("failed to enforce: %s", err)
			resp.WriteHeader(http.StatusInternalServerError)

			return
		}

		if allowed == false {
			resp.WriteHeader(http.StatusUnauthorized)

			return
		}

		handler(resp, req)
	}
}

func SetupApi(router *mux.Router) {
	router.HandleFunc("/auth", authHandler).Methods("POST")

	api := NewAuthorizedApi("/api", router)

	api.HandleFunc("/user", AuthWrapper(types.AuthObjectUser, types.AuthActionWrite, createUserHandler)).Methods("POST")
	api.HandleFunc("/access", notImplemented).Methods("GET")
	api.HandleFunc("/access", hasAccessHandler).Methods("POST")

	api.HandleFunc("/account", notImplemented)

	api.HandleFunc("/site", AuthWrapper(types.AuthObjectSite, types.AuthActionWrite, createSiteHandler)).Methods("POST")
	api.HandleFunc("/site", AuthWrapper(types.AuthObjectSite, types.AuthActionRead, loadSitesHandler)).Methods("GET")
	api.HandleFunc("/site/{key}", AuthWrapper(types.AuthObjectSite, types.AuthActionRead, loadSiteHandler)).Methods("GET")

	api.HandleFunc("/site/{key}/command", AuthWrapper(types.AuthObjectCommand, types.AuthActionRead, loadSiteCommandsHandler)).Methods("GET")
	api.HandleFunc("/site/{key}/blueprint", AuthWrapper(types.AuthObjectBlueprint, types.AuthActionRead, loadSiteBlueprintsHandler)).Methods("GET")

	api.HandleFunc("/command/job", AuthWrapper(types.AuthObjectCommandJob, types.AuthActionWrite, createCommandJobHandler)).Methods("POST")
	api.HandleFunc("/command/job", AuthWrapper(types.AuthObjectCommandJob, types.AuthActionRead, getCommandJobsHandler)).Methods("GET")
	api.HandleFunc("/command/job/{uuid}", AuthWrapper(types.AuthObjectCommandJob, types.AuthActionRead, getCommandJobHandler)).Methods("GET")

	api.HandleFunc("/command/job/{uuid}/events", AuthWrapper(types.AuthObjectCommandJobEvent, types.AuthActionRead, getCommandJobEventsHandler)).Methods("GET")

	api.HandleFunc("/blueprint", AuthWrapper(types.AuthObjectBlueprint, types.AuthActionRead, loadBlueprintsHandler)).Methods("GET")
	api.HandleFunc("/blueprint", AuthWrapper(types.AuthObjectBlueprint, types.AuthActionWrite, createBlueprintHandler)).Methods("POST")
	api.HandleFunc("/blueprint/{uuid}", AuthWrapper(types.AuthObjectBlueprint, types.AuthActionRead, loadBlueprintHandler)).Methods("GET")
	//api.HandleFunc("/blueprint/{uuid}", AuthWrapper(types.AuthObjectBlueprint, types.AuthActionWrite, updateBlueprintSetHandler)).Methods("POST")

	api.HandleFunc("/config", AuthWrapper(types.AuthObjectConfig, types.AuthActionRead, configHandler)).Methods("GET")
}

func NewAuthorizedApi(path string, router *mux.Router) *mux.Router {
	api := router.PathPrefix(path).Subrouter()
	api.Use(IsAuthorizedMiddleware)

	return api
}
