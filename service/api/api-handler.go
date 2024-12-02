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
	rt.router.POST("/chats", rt.wrap(rt.BearerAuth(rt.startNewChat)))							//TODO
	rt.router.GET("/chats", rt.wrap(rt.BearerAuth(rt.getUserChats)))							//fatto
	rt.router.GET("/chats/:chat_id", rt.wrap(rt.BearerAuth(rt.getChatInfo)))					//fatto
	rt.router.GET("/chats/:chat_id/messages", rt.wrap(rt.BearerAuth(rt.getChatMessages)))		//fatto

	//Group operations
	rt.router.DELETE("/chats/:chat_id/users", rt.wrap(rt.BearerAuth(rt.leaveGroup)))			//fatto
	rt.router.PATCH("/chats/:chat_id", rt.wrap(rt.BearerAuth(rt.changeGroupName)))				//fatto
	rt.router.PATCH("/chats/:chat_id/propic", rt.wrap(rt.BearerAuth(rt.changeGroupPhoto)))		//fatto
	rt.router.POST("/chats/:chat_id/users", rt.wrap(rt.BearerAuth(rt.postAddUserToGroup)))		//fatto

	//Message operations
	rt.router.POST("/chats/:chat_id/messages", rt.wrap(rt.BearerAuth(rt.sendMessage)))			//TODO
	rt.router.POST("/chats/:chat_id/messages", rt.wrap(rt.BearerAuth(rt.deleteMessage)))		//TODO
	rt.router.POST("/chats/:chat_id/messages", rt.wrap(rt.BearerAuth(rt.forwardMessage)))		//TODO
	rt.router.POST("/chats/:chat_id/messages", rt.wrap(rt.BearerAuth(rt.getMessageComments)))	//TODO
	rt.router.POST("/chats/:chat_id/messages", rt.wrap(rt.BearerAuth(rt.commentMessage)))		//TODO
	rt.router.POST("/chats/:chat_id/messages", rt.wrap(rt.BearerAuth(rt.uncommentMessage)))	//TODO

	// Special routes
	//rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
