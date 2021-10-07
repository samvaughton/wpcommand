package server

import (
	"encoding/json"
	"github.com/samvaughton/wpcommand/v2/pkg/auth"
	"net/http"
)

func reloadCasbinHandler(resp http.ResponseWriter, req *http.Request) {
	auth.Enforcer.LoadPolicy()

	json.NewEncoder(resp).Encode(map[string]string{"Status": "OK"})
}
