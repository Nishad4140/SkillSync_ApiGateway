package usercontroller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Nishad4140/SkillSync_ApiGateway/helper"
	"github.com/Nishad4140/SkillSync_ApiGateway/jwt"
	"github.com/Nishad4140/SkillSync_ProtoFiles/pb"
)

func (user *UserController) clientSignup(w http.ResponseWriter, r *http.Request) {
	if cookie, _ := r.Cookie("ClientToken"); cookie != nil {
		http.Error(w, "you are already logged in", http.StatusConflict)
		return
	}
	var req pb.ClientSignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.Email == "" {
		http.Error(w, "please enter a mail id", http.StatusBadRequest)
		return
	}
	if !helper.CheckString(req.Name) {
		http.Error(w, "please enter a valid name", http.StatusBadRequest)
		return
	}
	if !helper.ValidMail(req.Email) {
		http.Error(w, "please enter a valid mail id", http.StatusBadRequest)
		return
	}
	if !helper.CheckStringNumber(req.Phone) {
		http.Error(w, "please enter a valid phone number", http.StatusBadRequest)
		return
	}
	if !helper.IsStrongPassword(req.Password) {
		http.Error(w, "please enter a strong password including lower case, upper case, number, special character", http.StatusBadRequest)
		return
	}

	if req.OTP == "" {
		_, err := user.NotificationConn.SendOTP(context.Background(), &pb.SendOTPRequest{
			Email: req.Email,
		})
		if err != nil {
			http.Error(w, "error in sending the otp", http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"message": "Please enter the OTP sent to your mail"})
		return
	} else {
		verifyOtp, err := user.NotificationConn.VerifyOTP(context.Background(), &pb.VerifyOTPRequest{
			Otp:   req.OTP,
			Email: req.Email,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if !verifyOtp.IsVerified {
			http.Error(w, "OTP verification failed, Please try again", http.StatusBadRequest)
			return
		}
	}
	res, err := user.Conn.ClientSignup(context.Background(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err = user.Conn.CreateProfile(context.Background(), &pb.GetUserById{
		Id: res.Id,
	})
	if err != nil {
		fmt.Println("error when creating profile")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	jsonData, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	cookieString, err := jwt.GenerateJWT(res.Id, false, []byte(user.Secret))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	cookie := &http.Cookie{
		Name:     "ClientToken",
		Value:    cookieString,
		Expires:  time.Now().Add(48 * time.Hour),
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func (user *UserController) freelancerSignup(w http.ResponseWriter, r *http.Request) {
	if cookie, _ := r.Cookie("FreelancerToken"); cookie != nil {
		http.Error(w, "you are already logged in", http.StatusConflict)
		return
	}
	var req pb.FreelancerSignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.Email == "" {
		http.Error(w, "please enter a mail id", http.StatusBadRequest)
		return
	}
	if !helper.CheckString(req.Name) {
		http.Error(w, "please enter a valid name", http.StatusBadRequest)
		return
	}
	if !helper.ValidMail(req.Email) {
		http.Error(w, "please enter a valid mail id", http.StatusBadRequest)
		return
	}
	if !helper.CheckStringNumber(req.Phone) {
		http.Error(w, "please enter a valid phone number", http.StatusBadRequest)
		return
	}
	if req.CategoryId == 0 {
		http.Error(w, "please enter the cateogy that you going to work", http.StatusBadRequest)
		return
	}
	if !helper.IsStrongPassword(req.Password) {
		http.Error(w, "please enter a strong password including lower case, upper case, number, special character", http.StatusBadRequest)
		return
	}

	if req.OTP == "" {
		_, err := user.NotificationConn.SendOTP(context.Background(), &pb.SendOTPRequest{
			Email: req.Email,
		})
		if err != nil {
			http.Error(w, "error in sending the otp", http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"message": "Please enter the OTP sent to your mail"})
		return
	} else {
		verifyOtp, err := user.NotificationConn.VerifyOTP(context.Background(), &pb.VerifyOTPRequest{
			Otp:   req.OTP,
			Email: req.Email,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if !verifyOtp.IsVerified {
			http.Error(w, "OTP verification failed, Please try again", http.StatusBadRequest)
			return
		}
	}
	// Check the Category here
	res, err := user.Conn.FreelancerSignup(context.Background(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Create Freelancer Profile here
	jsonData, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	cookieString, err := jwt.GenerateJWT(res.Id, false, []byte(user.Secret))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	cookie := &http.Cookie{
		Name:     "FreelancerToken",
		Value:    cookieString,
		Expires:  time.Now().Add(48 * time.Hour),
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func (user *UserController) clientLogin(w http.ResponseWriter, r *http.Request) {
	if cookie, _ := r.Cookie("ClientToken"); cookie != nil {
		http.Error(w, "you are already logged in", http.StatusConflict)
		return
	}
	var req pb.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	res, err := user.Conn.ClientLogin(context.Background(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	jsonData, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	cookieString, err := jwt.GenerateJWT(res.Id, false, []byte(user.Secret))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	cookie := &http.Cookie{
		Name:     "ClientToken",
		Value:    cookieString,
		Expires:  time.Now().Add(48 * time.Hour),
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func (user *UserController) freelancerLogin(w http.ResponseWriter, r *http.Request) {
	if cookie, _ := r.Cookie("FreelancerToken"); cookie != nil {
		http.Error(w, "you are already logged in", http.StatusConflict)
		return
	}
	var req pb.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	res, err := user.Conn.FreelancerLogin(context.Background(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	jsonData, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	cookieString, err := jwt.GenerateJWT(res.Id, false, []byte(user.Secret))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	cookie := &http.Cookie{
		Name:     "FreelancerToken",
		Value:    cookieString,
		Expires:  time.Now().Add(48 * time.Hour),
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func (user *UserController) adminLogin(w http.ResponseWriter, r *http.Request) {
	if cookie, _ := r.Cookie("AdminToken"); cookie != nil {
		http.Error(w, "you are already logged in", http.StatusConflict)
		return
	}
	var req pb.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	res, err := user.Conn.AdminLogin(context.Background(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	jsonData, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	cookieString, err := jwt.GenerateJWT(res.Id, true, []byte(user.Secret))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	cookie := &http.Cookie{
		Name:     "AdminToken",
		Value:    cookieString,
		Expires:  time.Now().Add(48 * time.Hour),
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func (user *UserController) clientLogout(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:     "ClientToken",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"Logged out successfully"}`))
}

func (user *UserController) freelancerLogout(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:     "FreelancerToken",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"Logged out successfully"}`))
}

func (user *UserController) adminLogout(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:     "AdminToken",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"Logged out successfully"}`))
}
