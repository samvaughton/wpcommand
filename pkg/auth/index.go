package auth

import (
	sqladapter "github.com/Blank-Xu/sql-adapter"
	"github.com/casbin/casbin/v2"
	"github.com/samvaughton/wpcommand/v2/pkg/db"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	log "github.com/sirupsen/logrus"
)

var Enforcer *casbin.Enforcer

func InitAuth() {
	db := db.Db.DB

	adapter, err := sqladapter.NewAdapter(db, "postgres", "casbin")

	if err != nil {
		log.Fatalf("could not instantiate enforcer: %s", err)
	}

	e, err := casbin.NewEnforcer("casbin/model.conf", adapter)

	if err != nil {
		log.Fatalf("could not instantiate enforcer: %s", err)
	}

	var count int64

	row := db.QueryRow("select count(*) as count from casbin")
	err = row.Scan(&count)

	if err != nil {
		log.Fatalf("failed to load total policy count: %s", err)
	}

	log.Infof("total casbin policy count: %v", count)

	if count == 0 {
		// load default policy in
		e.AddPolicies([][]string{
			{types.RoleAdmin, types.AuthObjectSite, "*"},
			{types.RoleAdmin, types.AuthObjectCommand, "*"},
			{types.RoleAdmin, types.AuthObjectCommandJob, "*"},
			{types.RoleAdmin, types.AuthObjectCommandJobEvent, "*"},
			{types.RoleAdmin, types.AuthObjectBlueprint, "*"},
			{types.RoleAdmin, types.AuthObjectBlueprintObject, "*"},
			{types.RoleAdmin, types.AuthObjectUser, "*"},
			{types.RoleAdmin, types.AuthObjectAccount, "*"},
			{types.RoleAdmin, types.AuthObjectConfig, "read"},

			{types.RoleMember, types.AuthObjectSite, "read"},
			{types.RoleMember, types.AuthObjectCommand, "run"},
			{types.RoleMember, types.AuthObjectCommandJob, "read"},
			{types.RoleMember, types.AuthObjectCommandJobEvent, "read"},
			{types.RoleMember, types.AuthObjectUser, "read"},
			{types.RoleMember, types.AuthObjectAccount, "read"},
		})

		log.Info("default casbin policies added")
	}

	Enforcer = e
}
