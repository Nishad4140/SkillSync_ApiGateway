package usercontroller

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
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
	_, err = user.Conn.ClientCreateProfile(context.Background(), &pb.GetUserById{
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
	category, err := user.Conn.GetCategoryById(context.Background(), &pb.GetCategoryByIdRequest{
		Id: req.CategoryId,
	})
	if err != nil {
		http.Error(w, "please enter a valid category", http.StatusBadRequest)
		return
	}
	if category.Category == "" {
		http.Error(w, "please enter a valid category", http.StatusBadRequest)
		return
	}
	res, err := user.Conn.FreelancerSignup(context.Background(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Create Freelancer Profile here
	_, err = user.Conn.FreelancerCreateProfile(context.Background(), &pb.GetUserById{
		Id: res.Id,
	})
	if err != nil {
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

func (user *UserController) addCategory(w http.ResponseWriter, r *http.Request) {
	var req *pb.AddCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !helper.CheckString(req.Category) {
		http.Error(w, "please enter a valid category", http.StatusBadRequest)
		return
	}
	_, err := user.Conn.AddCategory(context.Background(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"Category Added Successfully"}`))
}

func (user *UserController) updateCategory(w http.ResponseWriter, r *http.Request) {
	var req *pb.UpdateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !helper.CheckString(req.Category) {
		http.Error(w, "please enter a valid category", http.StatusBadRequest)
		return
	}
	queryParams := r.URL.Query()
	categoryId, err := strconv.Atoi(queryParams.Get("category_id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	req.Id = int32(categoryId)
	_, err = user.Conn.UpdateCategory(context.Background(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"Category Updated Successfully"}`))
}

func (user *UserController) getAllCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := user.Conn.GetAllCategory(context.Background(), nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	categoriesData := []*pb.UpdateCategoryRequest{}
	for {
		category, err := categories.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		categoriesData = append(categoriesData, category)
	}
	jsonData, err := json.Marshal(categoriesData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func (user *UserController) adminAddSkill(w http.ResponseWriter, r *http.Request) {
	var req *pb.AddSkillRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.CategoryId == 0 {
		http.Error(w, "Please enter a valid category id", http.StatusBadRequest)
		return
	}
	if !helper.CheckString(req.Skill) {
		http.Error(w, "please enter a valid skill name", http.StatusBadRequest)
		return
	}
	if _, err := user.Conn.AdminAddSkill(context.Background(), req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"Skill Added Successfully"}`))
}

func (user *UserController) adminUpdateSkill(w http.ResponseWriter, r *http.Request) {
	var req *pb.SkillResponse
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !helper.CheckString(req.Skill) {
		http.Error(w, "please enter a valid skill name", http.StatusBadRequest)
		return
	}
	queryParams := r.URL.Query()
	skillId, err := strconv.Atoi(queryParams.Get("skill_id"))
	if err != nil {
		http.Error(w, "error while converting the id", http.StatusBadRequest)
		return
	}
	req.Id = int32(skillId)
	if _, err := user.Conn.AdminUpdateSkill(context.Background(), req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"Skill Updated Successfully"}`))
}

func (user *UserController) getAllSkills(w http.ResponseWriter, r *http.Request) {
	skills, err := user.Conn.GetAllSkills(context.Background(), nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	skillsData := []*pb.SkillResponse{}
	for {
		skill, err := skills.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		skillsData = append(skillsData, skill)
	}
	jsonData, err := json.Marshal(skillsData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func (user *UserController) clientAddAddress(w http.ResponseWriter, r *http.Request) {
	var req *pb.AddAddressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !helper.CheckString(req.Country) {
		http.Error(w, "please enter a valid country name", http.StatusBadRequest)
		return
	}
	if !helper.CheckString(req.State) {
		http.Error(w, "please enter a valid state name", http.StatusBadRequest)
		return
	}
	if !helper.CheckString(req.District) {
		http.Error(w, "please enter a valid district name", http.StatusBadRequest)
		return
	}
	if !helper.CheckString(req.City) {
		http.Error(w, "please enter a valid city name", http.StatusBadRequest)
		return
	}
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "error while retirieving the client id", http.StatusBadRequest)
		return
	}
	req.UserId = userID
	if _, err := user.Conn.ClientAddAddress(context.Background(), req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"Address added Successfully"}`))
}

func (user *UserController) clientUpdateAddress(w http.ResponseWriter, r *http.Request) {
	var req *pb.AddressResponse
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !helper.CheckString(req.Country) {
		http.Error(w, "please enter a valid country name", http.StatusBadRequest)
		return
	}
	if !helper.CheckString(req.State) {
		http.Error(w, "please enter a valid state name", http.StatusBadRequest)
		return
	}
	if !helper.CheckString(req.District) {
		http.Error(w, "please enter a valid district name", http.StatusBadRequest)
		return
	}
	if !helper.CheckString(req.City) {
		http.Error(w, "please enter a valid city name", http.StatusBadRequest)
		return
	}
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "error while retirieving the client id", http.StatusBadRequest)
		return
	}
	req.UserId = userID
	queryParams := r.URL.Query()
	addressId := queryParams.Get("address_id")
	req.Id = addressId
	if _, err := user.Conn.ClientUpdateAddress(context.Background(), req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"Address Updated Successfully"}`))
}

func (user *UserController) clientGetAddress(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "error while retirieving the client id", http.StatusBadRequest)
		return
	}
	req := &pb.GetUserById{
		Id: userID,
	}
	address, err := user.Conn.ClientGetAddress(context.Background(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if address.Country != "" {
		w.Write([]byte(`{"message":"please add address"}`))
		return
	}
	jsonData, err := json.Marshal(address)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
}

func (user *UserController) freelancerAddAddress(w http.ResponseWriter, r *http.Request) {
	var req *pb.AddAddressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !helper.CheckString(req.Country) {
		http.Error(w, "please enter a valid country name", http.StatusBadRequest)
		return
	}
	if !helper.CheckString(req.State) {
		http.Error(w, "please enter a valid state name", http.StatusBadRequest)
		return
	}
	if !helper.CheckString(req.District) {
		http.Error(w, "please enter a valid district name", http.StatusBadRequest)
		return
	}
	if !helper.CheckString(req.City) {
		http.Error(w, "please enter a valid city name", http.StatusBadRequest)
		return
	}
	userID, ok := r.Context().Value("freelancerID").(string)
	if !ok {
		http.Error(w, "error while retirieving the freelancer id", http.StatusBadRequest)
		return
	}
	req.UserId = userID
	if _, err := user.Conn.FreelancerAddAddress(context.Background(), req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"Address added Successfully"}`))
}

func (user *UserController) freelancerUpdateAddress(w http.ResponseWriter, r *http.Request) {
	var req *pb.AddressResponse
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !helper.CheckString(req.Country) {
		http.Error(w, "please enter a valid country name", http.StatusBadRequest)
		return
	}
	if !helper.CheckString(req.State) {
		http.Error(w, "please enter a valid state name", http.StatusBadRequest)
		return
	}
	if !helper.CheckString(req.District) {
		http.Error(w, "please enter a valid district name", http.StatusBadRequest)
		return
	}
	if !helper.CheckString(req.City) {
		http.Error(w, "please enter a valid city name", http.StatusBadRequest)
		return
	}
	userID, ok := r.Context().Value("freelancerID").(string)
	if !ok {
		http.Error(w, "error while retirieving the freelancer id", http.StatusBadRequest)
		return
	}
	req.UserId = userID
	queryParams := r.URL.Query()
	addressId := queryParams.Get("address_id")
	req.Id = addressId
	if _, err := user.Conn.FreelancerUpdateAddress(context.Background(), req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"Address Updated Successfully"}`))
}

func (user *UserController) freelancerGetAddress(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("freelancerID").(string)
	if !ok {
		http.Error(w, "error while retirieving the freelancer id", http.StatusBadRequest)
		return
	}
	req := &pb.GetUserById{
		Id: userID,
	}
	address, err := user.Conn.FreelancerGetAddress(context.Background(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if address.Country != "" {
		w.Write([]byte(`{"message":"please add address"}`))
		return
	}
	jsonData, err := json.Marshal(address)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
}

func (user *UserController) blockClient(w http.ResponseWriter, r *http.Request) {
	queryParam := r.URL.Query()
	userId := queryParam.Get("client_id")
	if userId == "" {
		http.Error(w, "please provide the user id", http.StatusBadRequest)
		return
	}
	if _, err := user.Conn.BlockClient(context.Background(), &pb.GetUserById{
		Id: userId,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"User Blocked Successfully"}`))
}

func (user *UserController) unBlockClient(w http.ResponseWriter, r *http.Request) {
	queryParam := r.URL.Query()
	userId := queryParam.Get("client_id")
	if userId == "" {
		http.Error(w, "please provide the user id", http.StatusBadRequest)
		return
	}
	if _, err := user.Conn.UnBlockClient(context.Background(), &pb.GetUserById{
		Id: userId,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"User UnBlocked Successfully"}`))
}

func (user *UserController) blockFreelancer(w http.ResponseWriter, r *http.Request) {
	queryParam := r.URL.Query()
	userId := queryParam.Get("freelancer_id")
	if userId == "" {
		http.Error(w, "please provide the user id", http.StatusBadRequest)
		return
	}
	if _, err := user.Conn.BlockFreelancer(context.Background(), &pb.GetUserById{
		Id: userId,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"User Blocked Successfully"}`))
}

func (user *UserController) unBlockFreelancer(w http.ResponseWriter, r *http.Request) {
	queryParam := r.URL.Query()
	userId := queryParam.Get("freelancer_id")
	if userId == "" {
		http.Error(w, "please provide the user id", http.StatusBadRequest)
		return
	}
	if _, err := user.Conn.UnBlockFreelancer(context.Background(), &pb.GetUserById{
		Id: userId,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"User UnBlocked Successfully"}`))
}
