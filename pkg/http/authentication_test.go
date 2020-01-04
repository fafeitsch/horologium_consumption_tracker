package http

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	stretch "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestJwtAuthenticatorImpl_Handle_Invalid(t *testing.T) {
	assert := stretch.New(t)
	jwtKey := make([]byte, 255)
	_, err := rand.Read(jwtKey)
	require.NoError(t, err)
	tests := []struct {
		name          string
		headers       map[string]string
		expMsg        string
		errorAccounts bool
	}{
		{name: "Missing Authorization Header", headers: map[string]string{}, expMsg: "No authorization header found"},
		{name: "Authorization Wrong Format", headers: map[string]string{"Authorization": "just_one_string"}, expMsg: "Wrong header found: expected Authorization: Basic"},
		{name: "Authorization Wrong Method", headers: map[string]string{"Authorization": "Digest 1234567890"}, expMsg: "Wrong header found: expected Authorization: Basic"},
		{name: "Authorization Wrong Password", headers: map[string]string{"Authorization": "Basic YWRtaW46ZnVjaHM="}, expMsg: "Invalid credentials"},
		{name: "Authorization No Colon", headers: map[string]string{"Authorization": "Basic YWRt"}, expMsg: "Extracting and decoding the credentials failed."},
		{name: "Error Getting Accounts", headers: map[string]string{"Authorization": "Basic YWRtaW46YWRtaW4="}, expMsg: "Invalid credentials", errorAccounts: true},
	}
	accountProvider := func() ([]Account, error) {
		return []Account{{UserName: "admin", EncryptedPassword: "$2y$12$mEMR.O0ES.A.xDHhzjVxZOuQ0Aj5iEaPlhtwLmE0fNL2y/2Eh74hW"}}, nil
	}
	handler := jwtAuthenticator{jwtKey: jwtKey, accountProvider: accountProvider}
	for _, tt := range tests {
		if tt.errorAccounts {
			handler.accountProvider = func() ([]Account, error) {
				return nil, fmt.Errorf("just a testing error")
			}
		} else {
			handler.accountProvider = accountProvider
		}
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest("POST", "/login", strings.NewReader(""))
			for key, value := range tt.headers {
				request.Header.Add(key, value)
			}
			recorder := httptest.NewRecorder()
			handler.Handle(recorder, request)
			assert.Equal("Basic", recorder.Header().Get("WWW-Authenticate"), "WWW-Authenticate header differs")
			assert.Equal(401, recorder.Result().StatusCode, "status code differs")
			assert.Equal(tt.expMsg+"\n", recorder.Body.String(), "message body differs")
		})
	}
}

func TestJwtAuthenticatorImpl_Handle_CorrectPassword(t *testing.T) {
	assert := stretch.New(t)
	jwtKey := []byte("s3cr3tPassw0rd4You")
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin"), 12)
	accountProvider := func() ([]Account, error) {
		return []Account{{UserName: "admin", EncryptedPassword: string(hashedPassword)}}, nil
	}
	handler := jwtAuthenticator{jwtKey: jwtKey, accountProvider: accountProvider}
	request := httptest.NewRequest("POST", "/login", strings.NewReader(""))
	request.Header.Add("Authorization", "Basic YWRtaW46YWRtaW4=") //admin:admin
	recorder := httptest.NewRecorder()
	handler.Handle(recorder, request)
	assert.Equal(200, recorder.Result().StatusCode, "status code")
	assert.Equal("application/json", recorder.Header().Get("Content-Type"), "Content Type header")
	result := make(map[string]string)
	err := json.NewDecoder(recorder.Body).Decode(&result)
	require.NoError(t, err)
	claims := &customClaim{}
	token, err := jwt.ParseWithClaims(result["token"], claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	require.NoError(t, err)
	assert.True(token.Valid, "the token should be valid")
	assert.Equal("admin", claims.UserName, "username of custom claims differs")
}

func TestAuthenticationInterceptorImpl_validateToken(t *testing.T) {
	assert := stretch.New(t)
	interceptor := authenticationInterceptor{
		jwtKey: []byte("s3cr3tPassw0rd4You"),
	}
	tests := []struct {
		name         string
		tokenString  string
		wantSuccess  bool
		wantClaim    *customClaim
		wantErrorMsg string
	}{
		{name: "Empty Token", tokenString: "", wantSuccess: false, wantClaim: nil, wantErrorMsg: "Authorization header is empty or does not exist."},
		{name: "Uncomplete Token", tokenString: "Bearer", wantSuccess: false, wantClaim: nil, wantErrorMsg: "Invalid or malformed token, expected two space-delimited words."},
		{name: "Length Exceeding Token", tokenString: "Bearer 3245435 sdfs", wantSuccess: false, wantClaim: nil, wantErrorMsg: "Invalid or malformed token, expected two space-delimited words."},
		{name: "Wrong Auth Method", tokenString: "Digest 3245435", wantSuccess: false, wantClaim: nil, wantErrorMsg: "Wrong authorization type, want \"Bearer\" got something else."},
		{name: "Unparsable Token", tokenString: "Bearer 3245435", wantSuccess: false, wantClaim: nil, wantErrorMsg: "The token could not be parsed."},
		{name: "Expired Token", tokenString: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjk0Njc0ODMxNCwiVXNlck5hbWUiOiJhZG1pbiJ9.YKKNZkA6Y6ZUw1A-DnXYzFMImsonaUhM3OaiSKwjcTk", wantSuccess: false, wantClaim: nil, wantErrorMsg: "The token is expired."},
		{name: "Invalid Signature", tokenString: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjk0Njc0ODMxNCwiVXNlck5hbWUiOiJhZG1pbiJ9.Il8oXXFhh53JojYSeQF1F-ciYq2qfXvCjWRCTH-IMzU", wantSuccess: false, wantClaim: nil, wantErrorMsg: "The signature of the token is invalid."},
		{name: "Valid Token", tokenString: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VyTmFtZSI6ImFkbWluIn0.knJrVjCBpSKATvLx09b8Fj3cExHpd_n5Cq4sr0eMIC8", wantSuccess: true, wantClaim: &customClaim{
			UserName: "admin"}, wantErrorMsg: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSuccess, gotClaim, gotErrorMsg := interceptor.validateToken(tt.tokenString)
			assert.Equal(tt.wantSuccess, gotSuccess, "success differs")
			if !reflect.DeepEqual(gotClaim, tt.wantClaim) {
				t.Errorf("validateToken() gotClaim = %v, want %v", gotClaim, tt.wantClaim)
			}
			assert.Equal(tt.wantErrorMsg, gotErrorMsg, "error message")
		})
	}
}

func TestAuthenticationInterceptorImpl_Handler(t *testing.T) {
	assert := stretch.New(t)
	middleware := authenticationInterceptor{
		jwtKey:         []byte("s3cr3tPassw0rd4You"),
		noAuthRequired: map[string]bool{"/api/unprotected": true},
	}
	tests := []struct {
		name        string
		res         string
		validToken  bool
		wantHeaders map[string]string
		wantContent string
		wantCode    int
	}{
		{name: "Unprotected", res: "/api/unprotected", validToken: false, wantHeaders: map[string]string{}, wantContent: "Test handler", wantCode: 200},
		{name: "ProtectedSuccess", res: "/api/protected", validToken: true, wantHeaders: map[string]string{}, wantContent: "Test handler", wantCode: 200},
		{name: "ProtectedFail", res: "/api/protected", validToken: false, wantHeaders: map[string]string{"WWW-Authenticate": "Bearer"}, wantContent: "Authorization header is empty or does not exist.\n", wantCode: 401},
	}
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Test handler"))
		require.NoError(t, err)
	})
	validToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VyTmFtZSI6ImFkbWluIn0.knJrVjCBpSKATvLx09b8Fj3cExHpd_n5Cq4sr0eMIC8"
	for _, tt := range tests {
		tt.wantHeaders["Access-Control-Allow-Origin"] = "http://localhost:4200"
		tt.wantHeaders["Access-Control-Allow-Methods"] = "POST, GET, OPTIONS"
		tt.wantHeaders["Access-Control-Allow-Headers"] = "Authorization,Content-Type"
		t.Run(tt.name, func(t *testing.T) {
			optionsReq := httptest.NewRequest("OPTIONS", tt.res, strings.NewReader(""))
			optionsRec := httptest.NewRecorder()
			middleware.Handler(testHandler).ServeHTTP(optionsRec, optionsReq)
			assert.Equal(200, optionsRec.Result().StatusCode, "OPTIONS status code")
			request := httptest.NewRequest("GET", tt.res, strings.NewReader(""))
			if tt.validToken {
				request.Header.Add("Authorization", "Bearer "+validToken)
			}
			recorder := httptest.NewRecorder()
			middleware.Handler(testHandler).ServeHTTP(recorder, request)
			for key, value := range tt.wantHeaders {
				assert.NotEmpty(recorder.Header().Get(key), "header %s is missing", key)
				assert.Equal(value, recorder.Header().Get(key), "header %s differs", key)
			}
			assert.Equal(tt.wantCode, recorder.Result().StatusCode, "status code")
			assert.Equal(tt.wantContent, recorder.Body.String(), "response body")
		})
	}
}
