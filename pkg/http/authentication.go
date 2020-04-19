package http

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strings"
)

type customClaim struct {
	jwt.StandardClaims
	UserName string
}

type authenticationInterceptor struct {
	jwtKey       []byte
	authRequired map[string]bool
}

func (a *authenticationInterceptor) getJwtKey() []byte {
	return a.jwtKey
}

func (a *authenticationInterceptor) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization,Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		requestPath := r.URL.Path
		if _, ok := a.authRequired[requestPath]; !ok {
			next.ServeHTTP(w, r)
			return
		}

		tokenHeader := r.Header.Get("Authorization")
		success, cClaim, errorMsg := a.validateToken(tokenHeader)

		if !success {
			w.Header().Add("WWW-Authenticate", "Bearer")
			http.Error(w, errorMsg, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user", cClaim.UserName)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func (a *authenticationInterceptor) validateToken(tokenHeader string) (success bool, claim *customClaim, errorMsg string) {
	if tokenHeader == "" {
		return false, nil, "Authorization header is empty or does not exist."
	}

	splitted := strings.Split(tokenHeader, " ")
	if len(splitted) != 2 {
		return false, nil, "Invalid or malformed token, expected two space-delimited words."
	}
	if splitted[0] != "Bearer" {
		return false, nil, "Wrong authorization type, want \"Bearer\" got something else."
	}

	tokenString := splitted[1]
	cClaim := &customClaim{}
	token, err := jwt.ParseWithClaims(tokenString, cClaim, func(token *jwt.Token) (interface{}, error) {
		return a.jwtKey, nil
	})

	if validationError, ok := err.(*jwt.ValidationError); ok {
		if validationError.Errors&jwt.ValidationErrorSignatureInvalid == jwt.ValidationErrorSignatureInvalid {
			return false, nil, "The signature of the token is invalid."
		}
		if validationError.Errors&jwt.ValidationErrorExpired == jwt.ValidationErrorExpired {
			return false, nil, "The token is expired."
		}
		return false, nil, "The token could not be parsed."
	}

	if token == nil || !token.Valid {
		//no sure how this line can be reached
		return false, nil, "The token could not be parsed."
	}
	return true, cClaim, ""
}

type jwtAuthenticator struct {
	jwtKey          []byte
	accountProvider func() ([]Account, error)
}

func (j *jwtAuthenticator) Handle(w http.ResponseWriter, r *http.Request) {
	tokenHeader := r.Header.Get("Authorization")
	const httpWwwAuthHeader = "WWW-Authenticate"
	if tokenHeader == "" {
		w.Header().Add(httpWwwAuthHeader, "Basic")
		http.Error(w, "No authorization header found", http.StatusUnauthorized)
		return
	}

	splitted := strings.Split(tokenHeader, " ")
	if len(splitted) != 2 || splitted[0] != "Basic" {
		w.Header().Add(httpWwwAuthHeader, "Basic")
		http.Error(w, "Wrong header found: expected Authorization: Basic", http.StatusUnauthorized)
		return
	}

	credentials, err := decodeAndSplitAuthString(splitted[1])
	if err != nil {
		w.Header().Add(httpWwwAuthHeader, "Basic")
		http.Error(w, "Extracting and decoding the credentials failed.", http.StatusUnauthorized)
		return
	}

	correctAccount := j.findCorrectAccount(credentials[0], credentials[1])
	if correctAccount == nil {
		w.Header().Add(httpWwwAuthHeader, "Basic")
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	claims := &customClaim{UserName: correctAccount.UserName}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claims)
	tokenString, _ := token.SignedString(j.jwtKey)
	response := map[string]string{"token": tokenString}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
	}
}

func (j *jwtAuthenticator) findCorrectAccount(username string, password string) *Account {
	accounts, err := j.accountProvider()
	if err != nil {
		log.Printf("%v", accounts)
		return nil
	}
	for _, account := range accounts {
		hashedPassword := []byte(account.EncryptedPassword)
		if account.UserName == username {
			passwordError := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
			if passwordError == nil {
				return &account
			}
		}
	}
	return nil
}

func decodeAndSplitAuthString(authString string) ([]string, error) {
	decoded, err := base64.StdEncoding.DecodeString(authString)
	if err != nil {
		return nil, fmt.Errorf("could not decode auth header: %v", err)
	}
	str := string(decoded)
	if !strings.Contains(str, ":") {
		return nil, fmt.Errorf("the decoded string must at least contain one \":\", but did not")
	}
	return strings.Split(str, ":"), nil
}
