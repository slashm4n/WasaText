package api

import (
	"WasaText/service/database"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var err error

	// Verify the db is ok
	if err = rt.db.Ping(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Get and verify the session token
	session_token_str := r.Header["Authorization"][0]
	if len(session_token_str) <= 7 {
		rt.baseLogger.Error("session token non settato, ricevuto solo 'Bearer '")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	session_token_str = strings.Split(session_token_str, " ")[1]
	session_token, _ := strconv.ParseInt(session_token_str, 10, 64)
	var user = sessions[session_token]
	if user.Id == 0 {
		rt.baseLogger.Error("unrecognized o elapsed session token: " + fmt.Sprint(session_token))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	rt.baseLogger.Info("used session token ", session_token_str, " for the user `", user.Name, "` (id ", user.Id, ")")

	// Retrieve the conversation_id and check if it belongs to the user
	conversation_id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/conversations/"))
	if err != nil {
		rt.baseLogger.Error("error reading conversation_id: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	rt.baseLogger.Info("received request to get conversation with id ", conversation_id)
	belongTo, err := rt.db.IsUserInConversation(user.Id, conversation_id)
	if !belongTo {
		rt.baseLogger.Error("the conversation_id do not belong to the user")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the conversation
	var messages []database.Message
	messages, err = rt.db.GetConversation(conversation_id)
	if err != nil {
		rt.baseLogger.Error("error trying to get conversations (", err.Error(), ")")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	rt.baseLogger.Info("retrieved ", len(messages), " messages")

	// Prepare and send the response
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(messages)
	if err != nil {
		rt.baseLogger.Error("error during the encoding of the json (", err.Error(), ")")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
