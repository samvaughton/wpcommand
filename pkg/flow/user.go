package flow

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/samvaughton/wpcommand/v2/pkg/auth"
	"github.com/samvaughton/wpcommand/v2/pkg/db"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	log "github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

func CreateDefaultAccountAndUser() {
	accounts, err := db.AccountsGetAll()

	if err != nil {
		log.Error(err)

		return
	}

	if len(accounts) > 0 {
		log.Info("accounts exist")

		return
	}

	tx, err := db.Db.BeginTx(context.Background(), nil)

	account, err := db.AccountCreate(tx, "Default", "default")

	if err != nil {
		log.Error(err)
		tx.Rollback()

		return
	}

	roles := []string{types.RoleAdmin}

	user, err := UserCreate(tx, "admin@admin.com", "Admin", "Admin", "password", roles, account.Id)

	if err != nil {
		log.Error(err)
		tx.Rollback()

		return
	}

	err = tx.Commit()

	if err != nil {
		log.Error(err)

		return
	}

	log.Info(fmt.Sprintf("Default account \"%s\" has been created with user %s, password=\"%s\"", account.Key, user.Email, "password"))
}

func UserCreate(bdb bun.IDB, email string, firstName string, lastName string, passwordPlain string, roles []string, accountId int64) (*types.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passwordPlain), 10)

	if err != nil {
		log.Error("could not create user, an error occurred with password hashing", err)
	}

	user := &types.User{
		Uuid:      uuid.New().String(),
		Email:     strings.ToLower(strings.Trim(email, " ,-_=/\n\r\t.@")),
		Password:  string(hashedPassword),
		FirstName: strings.Title(strings.ToLower(firstName)),
		LastName:  strings.Title(strings.ToLower(lastName)),
		Enabled:   true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = bdb.NewInsert().Model(user).Returning("*").Exec(context.Background())

	if err != nil {
		log.Error("could not create a new user", err)

		return nil, err
	}

	// insert relation
	ua := &types.UserAccount{UserId: user.Id, AccountId: accountId, Roles: roles, Uuid: uuid.New().String()}

	_, err = bdb.NewInsert().Model(ua).Exec(context.Background())
	if err != nil {
		log.Error("could not create user <-> account relation", err)

		return nil, err
	}

	// add role to user
	for _, role := range ua.Roles {
		_, err = auth.Enforcer.AddRoleForUser(ua.GetCasbinPolicyKey(), role)

		if err != nil {
			log.Error("failed to create user_account role in casbin", err)

			return nil, err
		}
	}

	return user, nil
}

func UserUpdate(bdb bun.IDB, user *types.User, passwordPlain string) error {

	// check if we need to set the password
	if passwordPlain != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passwordPlain), 10)

		if err != nil {
			log.Error("could not update user, an error occurred with password hashing", err)
		}

		user.Password = string(hashedPassword)
	}

	_, err := bdb.NewUpdate().Model(user).Where("id = ?", user.Id).Returning("*").Exec(context.Background())

	if err != nil {
		log.Error("could not update account", err)

		return err
	}

	return nil
}
