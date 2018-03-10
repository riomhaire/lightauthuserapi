package api

import (
	"net/http"
)

func (r *RestAPI) HandleOptions(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	code := http.StatusOK
	response.WriteHeader(code)

}
