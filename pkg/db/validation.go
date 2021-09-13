package db

import "github.com/pkg/errors"

type UserUniqueEmailRule struct {
	IgnoreUserId int64
}

func (r UserUniqueEmailRule) Validate(value interface{}) error {
	email, _ := value.(string)

	user := UserGetByEmail(email)

	if user == nil {
		return nil
	}

	if r.IgnoreUserId > 0 && user.Id == r.IgnoreUserId {
		return nil
	}

	return errors.New("user already exists")
}
