package db

import "github.com/pkg/errors"

type UserUniqueEmailRule struct {
}

func (r UserUniqueEmailRule) Validate(value interface{}) error {
	email, _ := value.(string)

	if UserExists(email) {
		return errors.New("a user with that email already exists")
	}

	return nil
}
