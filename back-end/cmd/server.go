package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) server() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Get("/media/*", func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/media/", http.FileServer(http.Dir("../media"))).ServeHTTP(w, r)
	})

	mux.Get("/", app.Home)

	// login & register
	mux.Post("/login", app.LoginHandler)
	mux.Post("/logout", app.LogoutHandler)
	mux.Post("/register", app.RegisterHandler)

	// messenger
	mux.Get("/messenger/users", app.GetUsers)
	mux.Get("/messenger/messages/{targetID}", app.LoadMessages)
	mux.HandleFunc("/messenger/ws", app.MessengerEndpoint)
	go app.ListenToMessengerWsChannel()

	// profile
	mux.Get("/profile/{userID}", app.GetProfile)
	mux.Put("/profile", app.PutProfile)
	mux.Get("/profile/image/{userID}", app.ProfileImageHandler)
	mux.Get("/profile/followers", app.ProfileFollowersHandler)
	mux.Post("/followUser", app.FollowUserHandler)
	mux.Get("/profile/followStatus/{userID}", app.GetFollowStatus)
	mux.Get("/followers", app.GetFollowers)

	// users
	mux.Get("/users", app.GetAllUsers)

	// notifications
	mux.Get("/notifications", app.GetNotificationsHandler)
	mux.Put("/notificationsSeen", app.NotificationsSeenHandler)
	mux.Post("/notificationAction", app.NotificationActionHandler)
	mux.Delete("/notifications", app.ClearAllNotifications)

	// posts
	mux.Get("/posts/user", app.GetUserPosts)
	mux.Get("/posts/group", app.GetGroupPosts)
	mux.Post("/posts/group", app.PostGroupPosts)
	mux.Post("/posts/user", app.PostUserPost)
	mux.Delete("/posts/user", app.DeleteUserPost)

	mux.Get("/user/post/{postID}", app.GetUserPostAndComments)
	mux.Get("/group/post/{postID}", app.GetGroupPostAndComments)

	mux.Post("/comments/{type}", app.PostComment)
	mux.Delete("/comments/{type}", app.DeleteComment)

	mux.Get("/posts/group", app.GetGroupPosts)
	mux.Post("/posts/group", app.PostGroupPosts)
	mux.Delete("/posts/group", app.DeleteGroupPosts)

	// groups
	mux.Get("/groups", app.GetGroups)
	mux.Post("/groups", app.PostGroups)
	mux.Delete("/groups", app.DeleteGroupHandler)

	mux.Get("/group/image/{groupID}", app.GroupImageHandler)
	mux.Get("/group", app.GetGroup)
	mux.Post("/group/join", app.GroupJoinHandler)
	mux.Post("/group/invite", app.GroupInviteHandler)
	mux.Delete("/group/join", app.RemoveGroupJoinHandler)
	mux.Post("/group/event", app.GroupCreateEventHandler)
	mux.Get("/group/events", app.GetGroupEvents)
	mux.Post("/group/eventStatus", app.PostGoingStatus)
	mux.Get("/group/event", app.GetEventHandler)
	mux.Delete("/group/leave", app.LeaveGroupHandler)
	mux.Delete("/group/rmuser", app.RemoveUserFromGroupHandler)

	//groupChat
	mux.Get("/groupChat/{groupID}", app.LoadGroupChatMessages)
	mux.HandleFunc("/groupChat/ws", app.GroupChatEndpoint)
	go app.ListenToGroupChatWsChannel()

	handler := app.enableCORS(mux)
	return handler
}

func (app *application) enableCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, X-CSRF-Token, Authorization")
			return
		} else {
			h.ServeHTTP(w, r)
		}
	})
}
