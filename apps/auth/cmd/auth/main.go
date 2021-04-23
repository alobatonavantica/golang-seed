package main

import (
	"golang-seed/apps/auth/pkg/config"
	"golang-seed/apps/auth/pkg/handler/authhand"
	"golang-seed/apps/auth/pkg/handler/clientshand"
	"golang-seed/apps/auth/pkg/handler/usershand"
	"golang-seed/apps/auth/pkg/models"
	"golang-seed/apps/auth/pkg/service/clientsserv"
	"golang-seed/apps/auth/pkg/service/usersserv"
	"golang-seed/apps/auth/pkg/store"
	"golang-seed/pkg/messages"
	"golang-seed/pkg/server"
	"golang-seed/pkg/server/middleware"
	"net/http"

	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/manage"
	oauth2server "github.com/go-oauth2/oauth2/v4/server"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func main() {
	if err := config.ParseSettings(); err != nil {
		log.Fatal(err)
	}

	if err := models.ConnectRepo(); err != nil {
		log.Fatal(err)
	}

	if err := messages.Init("apps/auth/config", "es"); err != nil {
		log.Fatal(err)
	}

	server := server.Init(config.Settings.Name, config.Settings.Port)
	server.ConfigureRouting()

	registerRoutes(server.RoutingRouter())

	server.Run()
}

func registerRoutes(r *mux.Router) {
	// Set up your services first.
	clientsService := clientsserv.NewClientsService()
	usersService := usersserv.NewUsersService()

	// auth handler
	manager := manage.NewDefaultManager()
	// To setup the token duration user manager.SetAuthorizeCodeTokenCfg, manager.SetImplicitTokenCfg,
	// manager.SetPasswordTokenCfg, manager.SetClientTokenCfg, manager.SetRefreshTokenCfg
	// token store
	manager.MustTokenStorage(store.NewTokenStore())

	// client store
	manager.MapClientStorage(store.NewClientStore())

	// auth server
	srv := oauth2server.NewDefaultServer(manager)
	srv.SetAllowGetAccessRequest(true)
	srv.SetAllowedGrantType(oauth2.PasswordCredentials, oauth2.ClientCredentials, oauth2.Refreshing)
	srv.SetClientInfoHandler(oauth2server.ClientBasicHandler)
	srv.SetAllowedResponseType(oauth2.Token)

	authHandler := authhand.NewAuthHandler(srv, usersService)
	r.HandleFunc("/oauth/authorize", authHandler.Authorize)
	r.HandleFunc("/oauth/token", authHandler.Token)

	srv.SetInternalErrorHandler(authHandler.InternalErrorHandler)
	srv.SetPasswordAuthorizationHandler(authHandler.PasswordAuthorizationHandler)

	// clients handler
	clientsHandler := clientshand.NewClientsHandler(clientsService)
	s := r.PathPrefix("/clients").Subrouter()
	s.Use(middleware.AuthenticationHandler(authHandler.ValidateToken))
	s.Handle("/{id}", middleware.Middleware(
		middleware.ErrorHandler(clientsHandler.Get),
		middleware.AuthorizeHandler("read:client", authHandler.ValidatePermission)),
	).Methods(http.MethodGet)
	s.Handle("/search/list", middleware.Middleware(
		middleware.ErrorHandler(clientsHandler.GetAll),
		middleware.AuthorizeHandler("read:clients", authHandler.ValidatePermission)),
	).Methods(http.MethodGet)
	s.Handle("/search/paged", middleware.Middleware(
		middleware.ErrorHandler(clientsHandler.GetAllPaged),
		middleware.AuthorizeHandler("read:clients", authHandler.ValidatePermission)),
	).Methods(http.MethodGet)
	s.Handle("", middleware.Middleware(
		middleware.ErrorHandler(clientsHandler.Create),
		middleware.AuthorizeHandler("create:client", authHandler.ValidatePermission)),
	).Methods(http.MethodPost)
	s.Handle("/{id}", middleware.Middleware(
		middleware.ErrorHandler(clientsHandler.Update),
		middleware.AuthorizeHandler("update:client", authHandler.ValidatePermission)),
	).Methods(http.MethodPut)
	s.Handle("/{id}", middleware.Middleware(
		middleware.ErrorHandler(clientsHandler.Delete),
		middleware.AuthorizeHandler("delete:client", authHandler.ValidatePermission)),
	).Methods(http.MethodDelete)

	// users handler
	usersHandler := usershand.NewUsersHandler(usersService)
	s = r.PathPrefix("/users").Subrouter()
	s.Use(middleware.AuthenticationHandler(authHandler.ValidateToken))
	s.Handle("/{id}", middleware.Middleware(
		middleware.ErrorHandler(usersHandler.Get),
		middleware.AuthorizeHandler("get:user", authHandler.ValidatePermission)),
	).Methods(http.MethodGet)
	s.Handle("/search/list", middleware.Middleware(
		middleware.ErrorHandler(usersHandler.GetAll),
		middleware.AuthorizeHandler("get:users", authHandler.ValidatePermission)),
	).Methods(http.MethodGet)
	s.Handle("/search/paged", middleware.Middleware(
		middleware.ErrorHandler(usersHandler.GetAllPaged),
		middleware.AuthorizeHandler("get:users", authHandler.ValidatePermission)),
	).Methods(http.MethodGet)
	s.Handle("", middleware.Middleware(
		middleware.ErrorHandler(usersHandler.Create),
		middleware.AuthorizeHandler("create:user", authHandler.ValidatePermission)),
	).Methods(http.MethodPost)
	s.Handle("/{id}", middleware.Middleware(
		middleware.ErrorHandler(usersHandler.Update),
		middleware.AuthorizeHandler("update:user", authHandler.ValidatePermission)),
	).Methods(http.MethodPut)
	s.Handle("/{id}", middleware.Middleware(
		middleware.ErrorHandler(usersHandler.Delete),
		middleware.AuthorizeHandler("delete:user", authHandler.ValidatePermission)),
	).Methods(http.MethodDelete)
}
