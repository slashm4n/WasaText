package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	rt.router.POST("/session", rt.doLogin)
	rt.router.PUT("/user", rt.setMyUserName)
	rt.router.GET("/conversations", rt.getMyConversations)
	rt.router.GET("/conversations/:to_user", rt.getConversation)
	rt.router.POST("/messages", rt.sendMessage)
	rt.router.DELETE("/messages/:msg_id", rt.deleteMessage)
	rt.router.POST("/messages/:msg_id", rt.forwardMessage)
	rt.router.PUT("/user/photo", rt.setMyPhoto)

	return rt.router
}
