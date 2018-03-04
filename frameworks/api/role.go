package api

import (
	"encoding/json"
	"net/http"

	"github.com/riomhaire/lightauthuserapi/usecases"
)

func (r *RestAPI) HandleReadRoles(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	valid, err := verifyAPIKey(req.Header.Get("Authorization"), r.Registry.Configuration.APIKey)

	if err.Code != usecases.NoError || !valid {
		http.Error(w, "Not Authorized", http.StatusUnauthorized)
		return
	}
	b, _ := json.Marshal(r.Registry.Usecases.ReadRoles())
	w.Write(b)

}
