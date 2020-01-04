package http

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type Account struct {
	UserName          string
	EncryptedPassword string
}

type ServerConfig struct {
	BindAddress     string
	Port            int
	AccountProvider func() ([]Account, error)
	JwtKey          []byte
}

func (c *ServerConfig) validateServerConfig() error {
	if c.Port < 0 {
		return fmt.Errorf("the port has value %d but should be greater than zero", c.Port)
	}
	if c.Port > 65535 {
		return fmt.Errorf("the port has value %d but should be less than or equal to 65535", c.Port)
	}
	if c.AccountProvider == nil {
		return fmt.Errorf("account provider must be given")
	}
	if len(c.JwtKey) == 0 {
		return fmt.Errorf("the jwt key has length 0, which is not allowed – a random key of length 265 is recommended")
	}
	return nil
}

type Server interface {
	StartServer() error
}

func NewServer(config *ServerConfig, apiHandler func(http.ResponseWriter, *http.Request)) (Server, error) {
	if err := config.validateServerConfig(); err != nil {
		return nil, fmt.Errorf("server configuration invalid: %v", err)
	}

	router := mux.NewRouter()
	middleware := authenticationInterceptor{
		jwtKey:         config.JwtKey,
		noAuthRequired: map[string]bool{"/login": true},
	}
	router.Use(middleware.Handler)

	authenticator := jwtAuthenticator{jwtKey: config.JwtKey, accountProvider: config.AccountProvider}
	router.HandleFunc("/api", apiHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/login", authenticator.Handle).Methods("GET", "OPTIONS")

	return &serverImpl{
		config: config,
		router: router,
	}, nil
}

type serverImpl struct {
	config *ServerConfig
	router *mux.Router
}

func (s *serverImpl) StartServer() error {
	return http.ListenAndServe(s.config.BindAddress+":"+strconv.Itoa(s.config.Port), s.router)
}
