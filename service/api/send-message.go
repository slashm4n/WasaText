package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type SendMessageRequest struct {
	Conversation_id int
	To_user         string `json:"to_user"`
	Message         string `json:"message"`
}

func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
	var req SendMessageRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid JSON data", http.StatusBadRequest)
		rt.baseLogger.Error("error decoding JSON:", err.Error())
		return
	}
	rt.baseLogger.Info("received request to send message to the user ", req.To_user)

	// Verify the db is ok
	if err = rt.db.Ping(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Get the receiver id
	var to_user_id int
	to_user_id, _ = rt.db.GetIdFromName(req.To_user)
	if to_user_id == 0 {
		rt.baseLogger.Error("the receipent of the message `", req.To_user, "` does not exist ")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Retrieve the conversation_id based on the user id's
	req.Conversation_id, _ = rt.db.GetConversationId(user.Id, to_user_id)

	// If the conversation does not exist, create it
	if req.Conversation_id == 0 {
		req.Conversation_id, err = rt.db.GetNextConversationId()
		if err != nil {
			rt.baseLogger.Error("error while retrieving the next conversation id ", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = rt.db.CreateConversation(req.Conversation_id, user.Id, to_user_id)
		if err != nil {
			rt.baseLogger.Error("error while creating the new conversation ", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		rt.baseLogger.Info("the conversation does not exist so a new one is created with conversation id ", req.Conversation_id)
	} else {
		rt.baseLogger.Info("retrieved conversation id ", req.Conversation_id)
	}

	// Send the message
	err = rt.db.SendMessage(req.Conversation_id, user.Id, req.Message, 0)
	if err != nil {
		rt.baseLogger.Error("error while sending a message " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)

	// Exit
	rt.baseLogger.Info("message sent")
}
