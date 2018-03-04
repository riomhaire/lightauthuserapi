package api

import (
	"net/http"
)

func (r *RestAPI) HandleHealth(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.Write([]byte("{ \"status\":\"up\"}"))

}
