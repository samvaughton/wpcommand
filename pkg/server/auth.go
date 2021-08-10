package server

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"github.com/samvaughton/wpcommand/v2/pkg/db"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
	"net/http"
	"time"
)

const SecretKey = "ai9ufgh94873yh8t0924hgt0[84wghneo8ridvoiah93-"

func GenerateJWT(email string, accountUuid string) (string, error) {
	var signingKey = []byte(SecretKey)
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["userEmail"] = email
	claims["accountUuid"] = accountUuid
	claims["exp"] = time.Now().Add(time.Hour * 7 * 24).Unix()

	tokenString, err := token.SignedString(signingKey)

	if err != nil {
		fmt.Errorf("something went wrong: %s", err)
		return "", err
	}

	return tokenString, nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func IsAuthorizedMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] == nil {
			w.WriteHeader(http.StatusUnauthorized)

			err := errors.New("token not found")
			log.Error(err)
			json.NewEncoder(w).Encode(err)

			return
		}

		var signingKey = []byte(SecretKey)

		token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("there was an error in parsing")
			}

			return signingKey, nil
		})

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Error(err)

			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			account := db.AccountGetByUuid(fmt.Sprintf("%s", claims["accountUuid"]))

			if account != nil {
				ctx := context.WithValue(r.Context(), "account", account)
				r = r.WithContext(ctx)
			}

			next.ServeHTTP(w, r)

			return
		}

		w.WriteHeader(http.StatusUnauthorized)
		resErr := errors.New("unauthorized")
		log.Error(resErr)
		json.NewEncoder(w).Encode(resErr)
	})
}
