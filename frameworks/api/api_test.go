package api

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/riomhaire/lightauthuserapi/entities"
	"github.com/riomhaire/lightauthuserapi/test"
	"github.com/riomhaire/lightauthuserapi/usecases"
)

// These are the API level tests ... This builds basic configuration
func createTestRegistry() usecases.Registry {
	logger := test.NewStringLogger()
	configuration := usecases.Configuration{}
	registry := usecases.Registry{}

	configuration.APIKey = "secret"
	configuration.RoleStore = "NONE"
	configuration.UserStore = "NONE"

	registry.Logger = logger
	registry.Configuration = configuration

	// Starting state for the DB
	userDb := make(map[string]entities.User)
	roleDb := make([]entities.Role, 0)
	roleDb = append(roleDb, entities.Role{"TEST"})

	registry.StorageInteractor = test.NewInMemoryDBInteractor(logger, userDb, roleDb)
	registry.Usecases = usecases.Usecases{&registry}

	return registry
}

// Tests API Key check works for create
func TestAPIKeySuccess(t *testing.T) {
	registry := createTestRegistry()
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	body := []byte("{\"username\":\"test\",\"password\":\"secret\"}")
	req, err := http.NewRequest("POST", "/api/v1/user/account", bytes.NewReader(body))

	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("bearer %v", "secret"))

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	restAPI := NewRestAPI(&registry)
	handler := http.HandlerFunc(restAPI.HandleGenericUser)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	status := rr.Code
	if status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestAPIKeyFailure(t *testing.T) {
	registry := createTestRegistry()
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	body := []byte("{\"username\":\"test\",\"password\":\"secret\"}")
	req, err := http.NewRequest("POST", "/api/v1/user/account", bytes.NewReader(body))

	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("bearer %v", "wrongsecret"))

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	restAPI := NewRestAPI(&registry)
	handler := http.HandlerFunc(restAPI.HandleGenericUser)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	status := rr.Code
	if status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}
}

func TestListRoles(t *testing.T) {
	registry := createTestRegistry()
	req, err := http.NewRequest("GET", "/api/v1/user/roles", nil)

	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("bearer %v", "secret"))

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	restAPI := NewRestAPI(&registry)
	handler := http.HandlerFunc(restAPI.HandleReadRoles)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	status := rr.Code
	if status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	// Check body is what we want
	body := rr.Body.String()
	expected := "[\"TEST\"]"
	if body != expected {
		t.Errorf("handler returned wrong content: got [%v] wanted [%v]", body, expected)
	}

}

func TestAddUser(t *testing.T) {
	registry := createTestRegistry()
	userName := "addTestUser"

	// Get Current User
	_, err := registry.Usecases.ReadUser(userName)
	if err.Code == usecases.NoError {
		t.Errorf("Unexpected - User exists before we test!!!")
	}

	// Create - Then Read
	user := entities.User{}
	user.Username = userName

	_, err = registry.Usecases.CreateUser(user)
	if err.Code != usecases.NoError {
		t.Errorf("Unexpected creation error - " + err.Error.Error())
	}

	// Should read ok
	user, err = registry.Usecases.ReadUser(userName)
	if err.Code != usecases.NoError {
		t.Errorf("Unexpected - User not exists after create !!")
	}

	// Check Names Match
	if user.Username != userName {
		t.Errorf("Unexpected - Read user does not match written ")
	}

}

func TestUpdateUser(t *testing.T) {
	registry := createTestRegistry()
	userName := "updateTestUser"

	// Get Current User
	_, err := registry.Usecases.ReadUser(userName)
	if err.Code == usecases.NoError {
		t.Errorf("Unexpected - User exists before we test!!!")
	}

	// Create - Then Read
	user := entities.User{}
	user.Username = userName
	user.Enabled = false

	_, err = registry.Usecases.CreateUser(user)
	if err.Code != usecases.NoError {
		t.Errorf("Unexpected creation error - " + err.Error.Error())
	}

	// Should read ok
	user, err = registry.Usecases.ReadUser(userName)
	if err.Code != usecases.NoError {
		t.Errorf("Unexpected - User not exists after create !!")
	}
	// Check state is false
	if user.Enabled != false {
		t.Errorf("Unexpected - State not false!!")
	}

	// Update
	user.Enabled = true
	user1, err := registry.Usecases.UpdateUser(user)
	if err.Code != usecases.NoError {
		t.Errorf("Unexpected - Update Error !!")
	}
	// Check state is true
	if user1.Enabled != true {
		t.Errorf("Unexpected - State not true!!")
	}
}

func TestDeleteUser(t *testing.T) {
	registry := createTestRegistry()
	userName := "deleteTestUser"

	// Get Current User
	_, err := registry.Usecases.ReadUser(userName)
	if err.Code == usecases.NoError {
		t.Errorf("Unexpected - User exists before we test!!!")
	}

	// Create - Then Read
	user := entities.User{}
	user.Username = userName
	user.Enabled = false

	_, err = registry.Usecases.CreateUser(user)
	if err.Code != usecases.NoError {
		t.Errorf("Unexpected creation error - " + err.Error.Error())
	}

	// Should read ok
	user, err = registry.Usecases.ReadUser(userName)
	if err.Code != usecases.NoError {
		t.Errorf("Unexpected - User not exists after create")
	}

	// Check Names Match && state match
	if user.Username != userName || user.Enabled != false {
		t.Errorf("Unexpected - Read user does not match written ")
	}

	// Delete
	registry.Usecases.DeleteUser(userName)

	// Should not read ok
	_, err = registry.Usecases.ReadUser(userName)
	if err.Code == usecases.NoError {
		t.Errorf("Unexpected - User not Deleted")
	}

}

func TestListUser(t *testing.T) {
	registry := createTestRegistry()
	userName := "listTestUser"

	// Get Current User List
	numberUser := len(registry.Usecases.ListUsers("", -1, -1))

	// Create - Then Read
	user := entities.User{}
	user.Username = userName
	user.Enabled = false

	_, err := registry.Usecases.CreateUser(user)
	if err.Code != usecases.NoError {
		t.Errorf("Unexpected creation error - " + err.Error.Error())
	}

	numberUserAfter := len(registry.Usecases.ListUsers("", -1, -1))
	if numberUserAfter != (numberUser + 1) {
		t.Errorf("Expected User List to go - it didnt")
	}

}
