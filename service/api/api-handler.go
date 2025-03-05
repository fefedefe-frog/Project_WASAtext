package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {

	// Login operation
	rt.router.POST("/session", rt.wrap(rt.postSession)) // fatto

	// User operations
	rt.router.GET("/users", rt.wrap(rt.BearerAuth(rt.getUsers)))						// fatto
	rt.router.GET("/users/:usr_id", rt.wrap(rt.BearerAuth(rt.getUserInfo)))            // fatto
	rt.router.PUT("/profile", rt.wrap(rt.BearerAuth(rt.setMyUserName)))         	// fatto
	rt.router.PUT("/profile/propic", rt.wrap(rt.BearerAuth(rt.setMyPhoto))) 	// fatto

	// Chat operations
	rt.router.POST("/chats", rt.wrap(rt.BearerAuth(rt.startNewChat)))								// fatto
	rt.router.GET("/chats", rt.wrap(rt.BearerAuth(rt.getMyConversations)))								// fatto
	rt.router.GET("/chats/:chat_id", rt.wrap(rt.BearerAuth(rt.getConversation)))						// fatto
	rt.router.PUT("/chats/:chat_id/messages", rt.wrap(rt.BearerAuth(rt.putRetriveChatMessages)))	// fatto
	rt.router.GET("/chats/:chat_id/users", rt.wrap(rt.BearerAuth(rt.getChatParticipants)))			// fatto

	// Group operations
	rt.router.DELETE("/chats/:chat_id/users", rt.wrap(rt.BearerAuth(rt.leaveGroup)))       // fatto
	rt.router.PUT("/chats/:chat_id", rt.wrap(rt.BearerAuth(rt.setGroupName)))         	// fatto
	rt.router.PUT("/chats/:chat_id/propic", rt.wrap(rt.BearerAuth(rt.setGroupPhoto)))	// fatto
	rt.router.POST("/chats/:chat_id/users", rt.wrap(rt.BearerAuth(rt.addToGroup)))		// fatto

	// Message operations
	rt.router.POST("/chats/:chat_id/messages", rt.wrap(rt.BearerAuth(rt.sendMessage)))											// fatto
	rt.router.DELETE("/chats/:chat_id/messages", rt.wrap(rt.BearerAuth(rt.deleteMessage)))										// fatto
	rt.router.POST("/chats/:chat_id/messages/:msg_id", rt.wrap(rt.BearerAuth(rt.forwardMessage)))								// fatto
	rt.router.PUT("/chats/:chat_id/messages/:msg_id", rt.wrap(rt.BearerAuth(rt.putUpdateMessageStatus)))						// fatto
	rt.router.GET("/chats/:chat_id/messages/:msg_id/comments", rt.wrap(rt.BearerAuth(rt.getMessageComments)))					// fatto
	rt.router.POST("/chats/:chat_id/messages/:msg_id/comments", rt.wrap(rt.BearerAuth(rt.commentMessage)))						// fatto
	rt.router.DELETE("/chats/:chat_id/messages/:msg_id/comments/:commenter_id", rt.wrap(rt.BearerAuth(rt.uncommentMessage)))	// fatto

	return rt.router
}
