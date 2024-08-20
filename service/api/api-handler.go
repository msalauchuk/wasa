package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {

	// Login enpoint
	rt.router.POST("/session", rt.wrap(rt.sessionHandler))

	// Search endpoint
	rt.router.GET("/users", rt.wrap(rt.searchUser))

	// User Endpoint
	rt.router.PUT("/users/:id", rt.wrap(rt.putNickname))
	rt.router.GET("/users/:id", rt.wrap(rt.getUserProfile))

	// Ban endpoint
	rt.router.PUT("/users/:id/banned_users/:banned_id", rt.wrap(rt.ban))
	rt.router.DELETE("/users/:id/banned_users/:banned_id", rt.wrap(rt.unban))

	// Followers endpoint
	rt.router.PUT("/users/:id/followers/:follower_id", rt.wrap(rt.follow))
	rt.router.DELETE("/users/:id/followers/:follower_id", rt.wrap(rt.unfollow))

	// Stream endpoint
	rt.router.GET("/users/:id/home", rt.wrap(rt.getFeed))

	// Photo Endpoint
	rt.router.POST("/users/:id/photos", rt.wrap(rt.postPhoto))
	rt.router.DELETE("/users/:id/photos/:photo_id", rt.wrap(rt.deletePhoto))
	rt.router.GET("/users/:id/photos/:photo_id", rt.wrap(rt.getPhoto))

	// Comments endpoint
	rt.router.POST("/users/:id/photos/:photo_id/comments", rt.wrap(rt.postComment))
	rt.router.DELETE("/users/:id/photos/:photo_id/comments/:comment_id", rt.wrap(rt.deleteComment))

	// Likes endpoint
	rt.router.PUT("/users/:id/photos/:photo_id/likes/:like_id", rt.wrap(rt.like))
	rt.router.DELETE("/users/:id/photos/:photo_id/likes/:like_id", rt.wrap(rt.unlike))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
