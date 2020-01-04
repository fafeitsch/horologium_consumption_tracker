package http

import (
	"crypto/rand"
	"encoding/json"
	stretch "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewServerValid(t *testing.T) {
	assert := stretch.New(t)
	jwtKey := make([]byte, 255)
	if _, err := rand.Read(jwtKey); err != nil {
		t.Fatalf("%v", err)
	}
	endpoint := func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Was here"))
	}
	encryptedPassword, _ := bcrypt.GenerateFromPassword([]byte("p4ssw0rd"), 12)
	accountProvider := func() ([]Account, error) {
		return []Account{{UserName: "Vincent", EncryptedPassword: string(encryptedPassword)}}, nil
	}
	config := ServerConfig{
		BindAddress:     "127.0.0.1",
		Port:            25565,
		AccountProvider: accountProvider,
		JwtKey:          jwtKey,
	}
	server, err := NewServer(&config, endpoint)
	require.NoError(t, err)
	impl, ok := server.(*serverImpl)
	require.True(t, ok, "the returned server implementation must be of type serverImpl")
	loginRecorder := httptest.NewRecorder()
	loginRequest := httptest.NewRequest("GET", "/login", strings.NewReader(""))
	loginRequest.Header.Add("Authorization", "Basic "+"VmluY2VudDpwNHNzdzByZA==")
	impl.router.ServeHTTP(loginRecorder, loginRequest)

	assert.Equal(200, loginRecorder.Result().StatusCode, "status code")
	result := make(map[string]string)
	assert.NoError(json.NewDecoder(loginRecorder.Body).Decode(&result), "the token must be decodeable")

	protectedRecorder := httptest.NewRecorder()
	protectedRequest := httptest.NewRequest("POST", "/api", strings.NewReader("{}"))
	protectedRequest.Header.Add("Authorization", "Bearer "+result["token"])
	impl.router.ServeHTTP(protectedRecorder, protectedRequest)

	assert.Equal("Was here", protectedRecorder.Body.String(), "protected recorder must be called")
}

func TestNewServer_Invalid(t *testing.T) {
	c := &ServerConfig{}
	_, err := NewServer(c, func(w http.ResponseWriter, r *http.Request) {})
	require.EqualError(t, err, "server configuration invalid: account provider must be given")
}

func TestServerConfig_validateServerConfig(t *testing.T) {
	assert := stretch.New(t)
	tests := []struct {
		name        string
		BindAddress string
		Port        int
		Accounts    []Account
		JwtKey      []byte
		wantMsg     string
	}{
		{name: "Port Too Small", BindAddress: "", Port: -22, Accounts: []Account{}, JwtKey: []byte{}, wantMsg: "the port has value -22 but should be greater than zero"},
		{name: "Port Too High", BindAddress: "", Port: 65545, Accounts: []Account{}, JwtKey: []byte{}, wantMsg: "the port has value 65545 but should be less than or equal to 65535"},
		{name: "No JWT Key", BindAddress: "", Port: 4321, Accounts: []Account{{UserName: "admin", EncryptedPassword: "sdasfad"}}, JwtKey: []byte{}, wantMsg: "the jwt key has length 0, which is not allowed – a random key of length 265 is recommended"},
		{name: "Nil JWT Key", BindAddress: "", Port: 4321, Accounts: []Account{{UserName: "admin", EncryptedPassword: "sdasfad"}}, JwtKey: nil, wantMsg: "the jwt key has length 0, which is not allowed – a random key of length 265 is recommended"},
		{name: "Valid", BindAddress: "", Port: 4321, Accounts: []Account{{UserName: "admin", EncryptedPassword: "sdasfad"}}, JwtKey: []byte{1}, wantMsg: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accountProvider := func() ([]Account, error) {
				return tt.Accounts, nil
			}
			c := &ServerConfig{
				BindAddress:     tt.BindAddress,
				Port:            tt.Port,
				AccountProvider: accountProvider,
				JwtKey:          tt.JwtKey,
			}
			err := c.validateServerConfig()
			if tt.wantMsg != "" {
				assert.EqualError(err, tt.wantMsg, "an error must be thrown")
			} else {
				assert.NoError(err, "no error must be thrown")
			}
		})
	}
}

func TestServerConfig_validateServerConfigNilAccount(t *testing.T) {
	c := &ServerConfig{
		BindAddress:     "127.0.0.1",
		Port:            5991,
		AccountProvider: nil,
		JwtKey:          []byte("password"),
	}
	err := c.validateServerConfig()
	require.EqualError(t, err, "account provider must be given")
}
