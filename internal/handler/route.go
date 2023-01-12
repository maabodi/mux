package handler

import (
	ac "muxapp/internal/handler/controller/activitycontroller"
	avc "muxapp/internal/handler/controller/avatarcontroller"
	nc "muxapp/internal/handler/controller/newscontroller"
	rc "muxapp/internal/handler/controller/rolecontroller"
	uc "muxapp/internal/handler/controller/usercontroller"
	upc "muxapp/internal/handler/controller/userprofilecontroller"
	wc "muxapp/internal/handler/controller/wilayahcontroller"
	"muxapp/internal/helper"

	"github.com/gorilla/mux"
)

// ConfigureRouter setup the router
func ConfigureRouter() *mux.Router {
	router := mux.NewRouter()

	// user route
	router.HandleFunc("/api/register", uc.CreateUser).Methods("POST")
	router.HandleFunc("/api/login", helper.Authenticate).Methods("POST")

	router.HandleFunc("/api/users", uc.IndexUser).Methods("GET")
	router.Handle("/api/user/{id}", helper.AuthMiddleware(uc.UpdateUser)).Methods("PUT")
	router.Handle("/api/user/{id}", helper.AuthMiddleware(uc.DeleteUser)).Methods("DELETE")
	router.Handle("/api/user/{id}", helper.AuthMiddleware(uc.ShowUser)).Methods("GET")
	router.HandleFunc("/api/user/get/{id}", uc.GetOneForDashboard).Methods("GET")
	router.HandleFunc("/api/user/get/role/{role}", uc.IndexUserByRole).Methods("GET")

	// user profile route
	router.HandleFunc("/api/user_profiles", upc.IndexUserProfile).Methods("GET")
	router.HandleFunc("/api/user_profile/{id}", upc.CreateUserProfile).Methods("POST")
	router.HandleFunc("/api/user_profile/{id}", upc.ShowUserProfile).Methods("GET")
	router.HandleFunc("/api/user_profile/{id}", upc.UpdateUserProfile).Methods("PUT")
	router.HandleFunc("/api/user_profile/{id}", upc.DeleteUserProfile).Methods("DELETE")

	// role route
	router.HandleFunc("/api/roles", rc.IndexRole).Methods("GET")
	router.Handle("/api/role", helper.AuthMiddleware(rc.CreateRole)).Methods("POST")
	router.Handle("/api/role/{id}", helper.AuthMiddleware(rc.UpdateRole)).Methods("PUT")
	router.Handle("/api/role/{id}", helper.AuthMiddleware(rc.DeleteRole)).Methods("DELETE")
	router.Handle("/api/role/{id}", helper.AuthMiddleware(rc.ShowRole)).Methods("GET")

	// news route
	router.HandleFunc("/api/news/all/{limit}/{offset}", nc.IndexNews).Methods("GET")
	router.Handle("/api/news/{id}", helper.AuthMiddleware(nc.ShowNews)).Methods("GET")
	router.Handle("/api/news/{id}", helper.AuthMiddleware(nc.UpdateNews)).Methods("PUT")
	router.Handle("/api/news/{id}", helper.AuthMiddleware(nc.DeleteNews)).Methods("DELETE")
	router.Handle("/api/news", helper.AuthMiddleware(nc.CreateNews)).Methods("POST")
	router.HandleFunc("/api/news/s/{keyword}", nc.SearchNews).Methods("GET")
	router.HandleFunc("/api/news/check/limit", nc.GetNewsLimit).Methods("GET")
	// router.Handle("/api/news", helper.AuthMiddleware(nc.CreateNews)).Methods("POST")

	// avatar route
	router.HandleFunc("/api/avatars", avc.IndexAvatar).Methods("GET")
	router.HandleFunc("/api/avatar/{id}", avc.ChangeAvatar).Methods("PUT")

	// wilayah route
	router.HandleFunc("/api/provinsi", wc.Index).Methods("GET")

	// activity route
	// router.HandleFunc("/api/activity", ac.Index).Methods("GET")
	router.HandleFunc("/api/activity/{limit}/{offset}", ac.IndexPage).Methods("GET")
	router.HandleFunc("/api/activity", ac.Create).Methods("POST")
	router.HandleFunc("/api/activity/{id}", ac.Update).Methods("PUT")
	router.HandleFunc("/api/activity/{id}", ac.Delete).Methods("DELETE")

	return router
}
