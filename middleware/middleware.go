package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/Nishad4140/SkillSync_ApiGateway/jwt"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load("../.env"); err != nil {
		fmt.Println("secret cannot be retrieved", err)
	}
	secret = os.Getenv("secret")
}

var (
	secret string
)

func ClientMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				http.Error(w, "you do not have the authority to perform this operation", http.StatusUnauthorized)
				return
			}
		}()
		cookie, err := r.Cookie("ClientToken")
		if err != nil {
			http.Error(w, "please login", http.StatusUnauthorized)
			return
		}
		cookieVal := cookie.Value
		claims, err := jwt.ValidateToken(cookieVal, []byte(secret))
		if err != nil {
			http.Error(w, "error in cokie validation", http.StatusUnauthorized)
			return
		}
		userID := claims["userID"]
		ctx := context.WithValue(r.Context(), "userID", userID)
		next(w, r.WithContext(ctx))
	}
}

func FreelancerMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				http.Error(w, "you do not have the authority to perform this operation", http.StatusUnauthorized)
				return
			}
		}()
		cookie, err := r.Cookie("FreelancerToken")
		if err != nil {
			http.Error(w, "please login", http.StatusUnauthorized)
			return
		}
		cookieVal := cookie.Value
		claims, err := jwt.ValidateToken(cookieVal, []byte(secret))
		if err != nil {
			http.Error(w, "error in cokie validation", http.StatusUnauthorized)
			return
		}
		userID := claims["userID"]
		ctx := context.WithValue(r.Context(), "freelancerID", userID)
		next(w, r.WithContext(ctx))
	}
}

func AdminMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				http.Error(w, "you do not have the authority to perform this operation", http.StatusUnauthorized)
				return
			}
		}()
		cookie, err := r.Cookie("AdminToken")
		if err != nil {
			http.Error(w, "please login", http.StatusUnauthorized)
			return
		}
		cookieVal := cookie.Value
		claims, err := jwt.ValidateToken(cookieVal, []byte(secret))
		if err != nil {
			http.Error(w, "error in cokie validation", http.StatusUnauthorized)
			return
		}
		userID := claims["userID"]
		ctx := context.WithValue(r.Context(), "userID", userID)
		next(w, r.WithContext(ctx))
	}
}
