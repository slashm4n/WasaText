package api

import (
	"WasaText/service/database"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	_ "github.com/mattn/go-sqlite3"
)

func (rt *_router) getMyConversations(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var err error

	// Verify the db is ok
	if err = rt.db.Ping(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Authenticate user through session token
	var user User
	user, err = AuthenticateUser(r)
	if err != nil {
		rt.baseLogger.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//rt.baseLogger.Info("used session token ", session_token_str, " for the user `", user.Name, "` (id ", user.Id, ")")
	rt.baseLogger.Info("authenticated user `", user.Name, "`, id ", user.Id, "")

	// Get the conversations
	var conversations []database.Conversation
	conversations, err = rt.db.GetMyConversations(user.Id)
	if err != nil {
		rt.baseLogger.Error("error trying to get conversations: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	rt.baseLogger.Info("retrieved ", len(conversations), " conversations")

	// Prepare and send the response
	w.Header().Set("Content-Type", "application/json")
	err_encoder := json.NewEncoder(w).Encode(conversations)
	if err_encoder != nil {
		rt.baseLogger.Error("error during the encoding of the json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
