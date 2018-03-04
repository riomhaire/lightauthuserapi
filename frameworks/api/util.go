package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/riomhaire/lightauthuserapi/usecases"
)

func extractAuthorization(header string) (string, error) {
	if len(header) < len(bearerPrefix) {
		return "", errors.New("Not Authorized")
	}

	prefix := strings.ToLower(header[:len(bearerPrefix)])

	if !strings.HasPrefix(prefix, bearerPrefix) {
		return "", errors.New("Not Authorized")
	}
	token := header[len(bearerPrefix):]
	return token, nil
}

func verifyAPIKey(header, key string) (bool, usecases.LightAuthError) {
	token, err := extractAuthorization(header)

	if err != nil {
		return false, usecases.NewError(usecases.NotAuthorized, errors.New("Not Authorized"))
	}

	if strings.Compare(key, token) != 0 {
		return false, usecases.NewError(usecases.NotAuthorized, errors.New("Not Authorized"))
	}
	return true, usecases.NewError(usecases.NoError, nil)
}

func applicationErrorToHttpStatus(appCode int) (int, []byte) {
	switch appCode {
	case usecases.NoError:
		return http.StatusOK, []byte("")
	case usecases.AlreadyExists:
		return http.StatusConflict, []byte("Already Exists")
	case usecases.NotImplemented:
		return http.StatusNotImplemented, []byte("Not Implemented")
	case usecases.Unknown:
		return http.StatusNotFound, []byte("Not Found")
	case usecases.Invalid:
		return http.StatusNotAcceptable, []byte("Invalid Request")
	case usecases.NotAuthorized:
		return http.StatusUnauthorized, []byte("Not Authorized")
	case usecases.InternalError:
		return http.StatusInternalServerError, []byte("Internal Error")
	}

	return http.StatusBadRequest, []byte("Bad Request")
}
