package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/samvaughton/wpcommand/v2/pkg/auth"
	"github.com/samvaughton/wpcommand/v2/pkg/db"
	"github.com/samvaughton/wpcommand/v2/pkg/flow"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strings"
)

func hasAccessHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")

	ua := req.Context().Value("userAccount").(*types.UserAccount)

	uack := ua.GetCasbinPolicyKey()

	var ep []interface{}
	ep = append(ep, uack)

	for _, p := range strings.Split(req.URL.Query().Get("params"), ",") {
		ep = append(ep, p)
	}

	allowed, err := auth.Enforcer.Enforce(ep...)

	if err != nil {
		log.Errorf("failed to enforce: %s", err)
		resp.WriteHeader(http.StatusInternalServerError)

		return
	}

	if allowed == false {
		resp.WriteHeader(http.StatusUnauthorized)

		return
	}

	resp.WriteHeader(http.StatusOK)
}

func createUserHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")

	up, err := types.NewUserCreatePayloadFromHttpRequest(req)

	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	validationErrors := up.Validate(db.UserUniqueEmailRule{})

	if validationErrors != nil {
		resp.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(resp).Encode(map[string]interface{}{"Status": "VALIDATION_ERRORS", "Errors": validationErrors})

		return
	}

	// if its validation only skip now
	if req.URL.Query().Get("validate") != "" {
		resp.WriteHeader(http.StatusOK)
		json.NewEncoder(resp).Encode(map[string]string{"Status": "VALIDATION_PASSED"})
		return
	}

	account := req.Context().Value("account").(*types.Account)

	user, err := flow.UserCreate(db.Db, up.Email, up.FirstName, up.LastName, up.Password, []string{types.RoleMember}, account.Id)

	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(resp).Encode(map[string]string{
			"Status":  "COULD_NOT_CREATE",
			"Message": "Something went wrong creating the user.",
		})
		return
	}

	json.NewEncoder(resp).Encode(user)
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
	var userAccount *types.UserAccount
	for _, accItem := range user.UserAccounts {
		if accItem.Account.Key == authPayload.Account {
			userAccount = accItem
			break
		}
	}

	if userAccount == nil {
		log.Error(errors.New("could not find user account for given credentials"))
		resp.WriteHeader(http.StatusInternalServerError)

		return
	}

	tokenString, err := GenerateJWT(userAccount)

	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusInternalServerError)

		return
	}

	account := userAccount.Account

	// we now want to load the casbin filtered policy doc for this user

	resp.WriteHeader(http.StatusOK)

	userGroups, err := auth.Enforcer.GetRolesForUser(userAccount.GetCasbinPolicyKey())

	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusInternalServerError)
	}

	conf := auth.Enforcer.GetModel().ToText()
	policy := auth.Enforcer.GetPolicy()

	rcpd := ""

	for _, group := range userGroups {
		for _, line := range policy {
			// if the policy line isn't in one of the user groups, we can skip it as they shouldn't see it and its not required
			if group != line[0] {
				continue
			}

			rcpd = rcpd + fmt.Sprintf("p,%s\n", strings.Join(line, ","))
		}
	}

	for _, group := range userGroups {
		rcpd = rcpd + fmt.Sprintf("g,%s,%s\n", group, userAccount.GetCasbinPolicyKey())
	}

	json.NewEncoder(resp).Encode(map[string]interface{}{
		"Status":          "AUTHENTICATED",
		"Token":           tokenString,
		"Email":           user.Email,
		"Roles":           userAccount.Roles,
		"AccountUserUuid": userAccount.Uuid,
		"Account": map[string]string{
			"Name": account.Name,
			"Key":  account.Key,
		},
		"Casbin": map[string]string{
			"Conf":   conf,
			"Policy": rcpd,
		},
	})
}
