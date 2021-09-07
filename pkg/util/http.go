package util

import (
	"encoding/json"
	"net/http"
)

type HttpStatus string

const HttpStatusNotFound = "NOT_FOUND"
const HttpStatusContentMalformed = "CONTENT_MALFORMED"
const HttpStatusInvalidPayload = "INVALID_PAYLOAD"
const HttpStatusValidationErrors = "VALIDATION_ERRORS"
const HttpStatusInternalServerError = "INTERNAL_SERVER_ERROR"

func HttpEmptyErrors() map[string]string {
	return map[string]string{}
}

func HttpErrorEncode(resp http.ResponseWriter, status HttpStatus, message string, errors interface{}) {
	switch status {
	case HttpStatusNotFound:
		resp.WriteHeader(http.StatusNotFound)
	case HttpStatusContentMalformed:
	case HttpStatusValidationErrors:
	case HttpStatusInvalidPayload:
		resp.WriteHeader(http.StatusBadRequest)
		break
	case HttpStatusInternalServerError:
		resp.WriteHeader(http.StatusInternalServerError)
		break
	}

	json.NewEncoder(resp).Encode(HttpError(status, message, errors))
}

func HttpError(status HttpStatus, message string, errors interface{}) map[string]interface{} {
	return map[string]interface{}{
		"Status":  string(status),
		"Message": message,
		"Errors":  errors,
	}
}
