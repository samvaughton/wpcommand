package server

import (
	"github.com/gorilla/mux"
	"github.com/samvaughton/wpcommand/v2/pkg/auth"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	"github.com/samvaughton/wpcommand/v2/pkg/util"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
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
			resp.WriteHeader(http.StatusForbidden)

			return
		}

		handler(resp, req)
	}
}

func SetupApi(router *mux.Router) {
	router.HandleFunc("/auth", authHandler).Methods("POST")
	router.Handle("/storage/{hash}", http.TimeoutHandler(util.HttpWrapHandlerFn(loadFileFromHashHandler), 120*time.Second, "timeout")).Methods("GET")

	api := NewAuthorizedApi("/api", router)

	api.HandleFunc("/user", AuthWrapper(types.AuthObjectUser, types.AuthActionWrite, createUserHandler)).Methods("POST")
	api.HandleFunc("/access", notImplemented).Methods("GET")
	api.HandleFunc("/access", hasAccessHandler).Methods("POST")

	api.HandleFunc("/account/{accUuid}", AuthWrapper(types.AuthObjectAccount, types.AuthActionReadSpecial, loadAccountHandler)).Methods("GET")
	api.HandleFunc("/account/{accUuid}", AuthWrapper(types.AuthObjectAccount, types.AuthActionWriteSpecial, updateAccountHandler)).Methods("PUT")
	api.HandleFunc("/account/{accUuid}/user/{userUuid}", AuthWrapper(types.AuthObjectUser, types.AuthActionReadSpecial, loadAccountUserApiItemHandler)).Methods("GET")
	api.HandleFunc("/account/{accUuid}/user/{userUuid}", AuthWrapper(types.AuthObjectUser, types.AuthActionWriteSpecial, updateAccountUserHandler)).Methods("PUT")
	api.HandleFunc("/account/{accUuid}/user", AuthWrapper(types.AuthObjectUser, types.AuthActionReadSpecial, loadAccountUserApiItemsHandler)).Methods("GET")
	api.HandleFunc("/account/{accUuid}/user", AuthWrapper(types.AuthObjectUser, types.AuthActionWriteSpecial, createAccountUserHandler)).Methods("POST")
	api.HandleFunc("/account", AuthWrapper(types.AuthObjectAccount, types.AuthActionReadSpecial, loadAccountsHandler)).Methods("GET")
	api.HandleFunc("/account", AuthWrapper(types.AuthObjectAccount, types.AuthActionWriteSpecial, createAccountHandler)).Methods("POST")

	api.HandleFunc("/site", AuthWrapper(types.AuthObjectSite, types.AuthActionWrite, createSiteHandler)).Methods("POST")
	api.HandleFunc("/site", AuthWrapper(types.AuthObjectSite, types.AuthActionRead, loadSitesHandler)).Methods("GET")
	api.HandleFunc("/site/{key}", AuthWrapper(types.AuthObjectSite, types.AuthActionRead, loadSiteHandler)).Methods("GET")
	api.HandleFunc("/site/{siteUuid}", AuthWrapper(types.AuthObjectSite, types.AuthActionWrite, updateSiteHandler)).Methods("PUT")

	api.HandleFunc("/site/{key}/command", AuthWrapper(types.AuthObjectCommand, types.AuthActionRun, loadSiteCommandsHandler)).Methods("GET")
	api.HandleFunc("/site/{key}/command", AuthWrapper(types.AuthObjectCommand, types.AuthActionWrite, createSiteCommandHandler)).Methods("POST")
	api.HandleFunc("/site/{siteUuid}/command/{cmdUuid}", AuthWrapper(types.AuthObjectCommand, types.AuthActionWrite, updateSiteCommandHandler)).Methods("PUT")

	api.HandleFunc("/site/{key}/blueprint", AuthWrapper(types.AuthObjectBlueprint, types.AuthActionRead, loadSiteBlueprintsHandler)).Methods("GET")

	api.HandleFunc("/command/job", AuthWrapper(types.AuthObjectCommandJob, types.AuthActionWrite, createCommandJobHandler)).Methods("POST")
	api.HandleFunc("/command/job", AuthWrapper(types.AuthObjectCommandJob, types.AuthActionRead, getCommandJobsHandler)).Methods("GET")
	api.HandleFunc("/command/job/{uuid}", AuthWrapper(types.AuthObjectCommandJob, types.AuthActionRead, getCommandJobHandler)).Methods("GET")

	api.HandleFunc("/command/job/{uuid}/event", AuthWrapper(types.AuthObjectCommandJobEvent, types.AuthActionRead, getCommandJobEventsHandler)).Methods("GET")

	api.HandleFunc("/blueprint/{bpUuid}/object/{objUuid}/revision/{revId}/file", AuthWrapper(types.AuthObjectBlueprintObject, types.AuthActionRead, loadBlueprintObjectFileHandler)).Methods("GET")
	api.HandleFunc("/blueprint/{bpUuid}/object/{objUuid}/revision/{revId}/version", AuthWrapper(types.AuthObjectBlueprintObject, types.AuthActionRead, createBlueprintObjectRevisionHandler)).Methods("POST")
	api.HandleFunc("/blueprint/{bpUuid}/object/{objUuid}/revision/{revId}", AuthWrapper(types.AuthObjectBlueprintObject, types.AuthActionRead, loadBlueprintObjectHandler)).Methods("GET")
	api.HandleFunc("/blueprint/{bpUuid}/object/{objUuid}/revision/{revId}", AuthWrapper(types.AuthObjectBlueprintObject, types.AuthActionDelete, deleteBlueprintObjectRevisionHandler)).Methods("DELETE")
	api.HandleFunc("/blueprint/{bpUuid}/object/{objUuid}/revision", AuthWrapper(types.AuthObjectBlueprintObject, types.AuthActionRead, loadBlueprintObjectRevisionsHandler)).Methods("GET")
	api.HandleFunc("/blueprint/{bpUuid}/object/{objUuid}", AuthWrapper(types.AuthObjectBlueprintObject, types.AuthActionDelete, deleteBlueprintObjectHandler)).Methods("DELETE")
	api.HandleFunc("/blueprint/{bpUuid}/object/{objUuid}", AuthWrapper(types.AuthObjectBlueprintObject, types.AuthActionWrite, updateBlueprintObjectHandler)).Methods("PUT")

	api.HandleFunc("/blueprint/{uuid}/object", AuthWrapper(types.AuthObjectBlueprintObject, types.AuthActionRead, loadBlueprintObjectsHandler)).Methods("GET")
	api.HandleFunc("/blueprint/{uuid}/object", AuthWrapper(types.AuthObjectBlueprintObject, types.AuthActionWrite, createBlueprintObjectHandler)).Methods("POST")

	api.HandleFunc("/blueprint/{uuid}", AuthWrapper(types.AuthObjectBlueprint, types.AuthActionRead, loadBlueprintHandler)).Methods("GET")
	api.HandleFunc("/blueprint/{uuid}", AuthWrapper(types.AuthObjectBlueprint, types.AuthActionDelete, deleteBlueprintHandler)).Methods("DELETE")

	api.HandleFunc("/blueprint", AuthWrapper(types.AuthObjectBlueprint, types.AuthActionRead, loadBlueprintsHandler)).Methods("GET")
	api.HandleFunc("/blueprint", AuthWrapper(types.AuthObjectBlueprint, types.AuthActionWrite, createBlueprintHandler)).Methods("POST")

	api.HandleFunc("/config", AuthWrapper(types.AuthObjectConfig, types.AuthActionRead, configHandler)).Methods("GET")
}

func NewAuthorizedApi(path string, router *mux.Router) *mux.Router {
	api := router.PathPrefix(path).Subrouter()
	api.Use(IsAuthorizedMiddleware)

	return api
}
