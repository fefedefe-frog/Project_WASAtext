package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {

	//Login operation
	rt.router.POST("/session", rt.wrap(rt.postSession))										//fatto

	//User operations
	rt.router.GET("/users", rt.wrap(rt.BearerAuth(rt.getUsers)))								//fatto
	rt.router.GET("/users/:usr_id", rt.wrap(rt.BearerAuth(rt.getUserInfo)))					//fatto
	rt.router.PATCH("/users/:usr_id", rt.wrap(rt.BearerAuth(rt.patchChangeUserName)))			//fatto
	rt.router.PATCH("/users/:usr_id/propic", rt.wrap(rt.BearerAuth(rt.patchChangeUserPhoto)))	//fatto

	//Chat operations
	rt.router.GET("/chats", rt.wrap(rt.BearerAuth(rt.getUserChats)))							//TODO
	rt.router.GET("/chats/:chat_id/messages", rt.wrap(rt.BearerAuth(rt.getChatMessages)))		//TODO
	rt.router.GET("/chats/:chat_id", rt.wrap(rt.BearerAuth(rt.getChatInfo)))					//TODO

	//Group operations
	//TODO
	rt.router.DELETE("/chats/:chat_id", rt.wrap(rt.BearerAuth(rt.leaveGroup)))					//TODO

	//Message operations
	//TODO

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
