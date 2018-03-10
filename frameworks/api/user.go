package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/riomhaire/lightauthuserapi/entities"
	"github.com/riomhaire/lightauthuserapi/usecases"
)

func (r *RestAPI) HandleGenericUser(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	code := http.StatusNotImplemented
	data := []byte("Not Implemented")
	var err usecases.LightAuthError

	valid, err := verifyAPIKey(request.Header.Get("Authorization"), r.Registry.Configuration.APIKey)

	if err.Code == usecases.NoError && valid {
		// Read
		switch request.Method {
		case http.MethodGet:
			// Extract search and number of results and page
			queryValues := request.URL.Query()
			page := -1     // All pages
			pageSize := -1 // One page
			search := queryValues.Get("search")
			val := queryValues.Get("page")
			i, err := strconv.Atoi(val)
			if err == nil {
				page = i
			}
			val = queryValues.Get("pageSize")
			i, err = strconv.Atoi(val)
			if err == nil {
				pageSize = i
			}

			users := r.Registry.Usecases.ListUsers(search, page, pageSize)
			data, _ = json.Marshal(users)
		case http.MethodPost:
			decoder := json.NewDecoder(request.Body)
			var u entities.User
			derr := decoder.Decode(&u)
			if derr == nil {
				var user entities.User
				user, err = r.Registry.Usecases.CreateUser(u)
				data, _ = json.Marshal(user)
			} else {
				err = usecases.NewError(usecases.Invalid, derr)
			}
			defer request.Body.Close()
		default:
			err = usecases.NewError(usecases.NotImplemented, errors.New("Not Implemented"))
		}
	}
	// Final encode
	code, returnData := applicationErrorToHttpStatus(err.Code)
	if err.Code == usecases.NoError {
		returnData = data
	}

	response.WriteHeader(code)
	response.Write(returnData)
	if code != http.StatusOK {
		msg := fmt.Sprintf("App Error %v : %v", code, string(data))
		r.Registry.Logger.Log("ERROR", msg)
	}
}

func (r *RestAPI) HandleSpecificUser(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	username := mux.Vars(request)["name"]
	code := http.StatusNotImplemented
	data := []byte("Not Implemented")
	var user entities.User
	var err usecases.LightAuthError

	valid, err := verifyAPIKey(request.Header.Get("Authorization"), r.Registry.Configuration.APIKey)

	if err.Code == usecases.NoError && valid {
		// Read
		switch request.Method {
		case http.MethodGet:
			user, err = r.Registry.Usecases.ReadUser(username)
		case http.MethodPut:
			decoder := json.NewDecoder(request.Body)
			var u entities.User
			derr := decoder.Decode(&u)
			if derr == nil {
				user, err = r.Registry.Usecases.UpdateUser(u)
			} else {
				err = usecases.NewError(usecases.Invalid, derr)
			}
			defer request.Body.Close()
		case http.MethodDelete:
			err = r.Registry.Usecases.DeleteUser(username)

		default:
			err = usecases.NewError(usecases.NotImplemented, errors.New("Not Implemented"))
		}
	}
	// Final encode
	code, data = applicationErrorToHttpStatus(err.Code)
	if err.Code == usecases.NoError {
		data, _ = json.Marshal(user)
	}

	response.WriteHeader(code)
	response.Write(data)
	if code != http.StatusOK {
		msg := fmt.Sprintf("App Error %v : %v", code, string(data))
		r.Registry.Logger.Log("ERROR", msg)
	}
}
