package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	_ "github.com/mattn/go-sqlite3"
)

type SetMyUserNameRequest struct {
	New_name string `json:"new_name"`
}

func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Authenticate user through session token
	var user User
	var err error
	user, err = AuthenticateUser(r)
	if err != nil {
		rt.baseLogger.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	rt.baseLogger.Info("authenticated user `", user.Name, "`, id ", user.Id, "")

	// Read the request
	var req SetMyUserNameRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid JSON data", http.StatusBadRequest)
		rt.baseLogger.Error("error decoding JSON:", err.Error())
		return
	}
	rt.baseLogger.Info("received request to change user name to '" + req.New_name + "'")

	// Verify the db is ok
	if err := rt.db.Ping(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Check the new name pattern
	if len(req.New_name) < 3 || len(req.New_name) > 16 {
		rt.baseLogger.Error("user name is too long or too short. User name is: " + req.New_name)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Update the name
	err = rt.db.SetMyUserName(user.Id, req.New_name)
	if err != nil {
		rt.baseLogger.Error("error trying to update user name: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	rt.baseLogger.Info("name succesfully updated.")
}
