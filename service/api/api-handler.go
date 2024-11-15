package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {

	// Register routes
	//rt.router.GET("/", rt.getHelloWorld)
	//rt.router.GET("/context", rt.wrap(rt.getContextReply))

	//Login operation
	rt.router.POST("/session", rt.wrap(rt.postSession))

	//User operations
	rt.router.GET("/users", rt.wrap(rt.getUsers))
	rt.router.GET("/users/:usr_id", rt.wrap(rt.getUserInfo))
	rt.router.PATCH("/users/:usr_id", rt.wrap(rt.patchChangeUserName))
	rt.router.GET("/users/:usr_id/propic", rt.wrap(rt.patchChangeUserPhoto))

	//Chat operations
	rt.router.GET("/chats", rt.wrap(rt.getUserChat))
	rt.router.DELETE("/chats/:chat_id", rt.wrap(rt.leaveChat))
	rt.router.GET("/chats/:chat_id", rt.wrap(rt.getChatNameAndPhoto))
	rt.router.GET("/chats/:chat_id/messages", rt.wrap(rt.getChatMessages))

	//Group operations

	//Message operations


	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
