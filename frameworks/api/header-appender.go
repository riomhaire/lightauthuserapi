package api

import (
	"net/http"
	"os"
)

// AddWorkerHeader - adds header of which node actually processed request
func (r *RestAPI) AddWorkerHeader(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	host, err := os.Hostname()
	if err != nil {
		host = "Unknown"
	}
	rw.Header().Add("X-Worker", host)
	if next != nil {
		next(rw, req)
	}
}

// AddWorkerVersion - adds header of which version is installed
func (r *RestAPI) AddWorkerVersion(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	version := r.Registry.Configuration.Version
	if len(version) == 0 {
		version = "UNKNOWN"
	}
	rw.Header().Add("X-Worker-Version", version)
	if next != nil {
		next(rw, req)
	}
}

// AddWorkerHeader - adds coors header
func (r *RestAPI) AddCoorsHeader(rw http.ResponseWriter, request *http.Request, next http.HandlerFunc) {

	//rw.Header().Add("Access-Control-Allow-Origin", request.Header.Get("Origin"))
	rw.Header().Add("Access-Control-Allow-Credentials", "true")
	rw.Header().Add("Access-Control-Allow-Origin", "*")
	rw.Header().Add("Access-Control-Allow-Methods", "POST, PUT, GET, OPTIONS, DELETE")
	rw.Header().Add("Access-Control-Max-Age", "3600")
	rw.Header().Add("Access-Control-Allow-Headers", "Content-Type, Accept, X-Requested-With, remember-me, authorization, Authorization")

	if next != nil {
		next(rw, request)
	}
}
