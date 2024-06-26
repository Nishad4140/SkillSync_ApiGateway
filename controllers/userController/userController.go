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
	helperstruct "github.com/Nishad4140/SkillSync_ApiGateway/helperStruct"
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

	category, err := user.Conn.GetCategoryById(context.Background(), &pb.GetCategoryByIdRequest{
		Id: req.CategoryId,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
	if address.Country == "" {
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

func (user *UserController) getClientProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "error while retirieving the client id", http.StatusBadRequest)
		return
	}
	client, err := user.Conn.GetClientById(context.Background(), &pb.GetUserById{
		Id: userID,
	})
	if err != nil {
		http.Error(w, "error in retrieving the user info", http.StatusBadRequest)
		return
	}
	address, err := user.Conn.ClientGetAddress(context.Background(), &pb.GetUserById{
		Id: userID,
	})
	if err != nil {
		http.Error(w, "error in retrieving the user address", http.StatusBadRequest)
		return
	}
	imageData, err := user.Conn.ClientGetProfileImage(context.Background(), &pb.GetUserById{
		Id: userID,
	})
	if err != nil {
		http.Error(w, "error in retrieving the user profile image", http.StatusBadRequest)
		return
	}
	res := helperstruct.ClientProfile{
		Id:      client.Id,
		Name:    client.Name,
		Email:   client.Email,
		Phone:   client.Phone,
		Image:   imageData.Url,
		Address: address,
	}
	jsonData, err := json.Marshal(res)
	if err != nil {
		http.Error(w, "error while marshalling the data", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func (user *UserController) uploadClientProfileImage(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "unable to parse form", http.StatusBadRequest)
		return
	}
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "unable to fetch the file from request", http.StatusBadRequest)
		return
	}
	defer file.Close()
	filebyte, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "unable to read the file", http.StatusInternalServerError)
		return
	}
	clientID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "error while retrieving the client id", http.StatusBadRequest)
		return
	}
	req := &pb.ImageRequest{
		ObjectName: fmt.Sprintf("%s-profile", clientID),
		ImageData:  filebyte,
		UserId:     clientID,
	}
	res, err := user.Conn.ClientUploadProfileImage(context.Background(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	jsonData, err := json.Marshal(res)
	if err != nil {
		http.Error(w, "error while marshalling the data", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func (user *UserController) clientEditName(w http.ResponseWriter, r *http.Request) {
	var req *pb.EditNameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	clientID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "error while retrieving the client id", http.StatusBadRequest)
		return
	}
	req.UserId = clientID
	if !helper.CheckString(req.Name) {
		http.Error(w, "please enter a valid name", http.StatusBadRequest)
		return
	}
	if _, err := user.Conn.ClientEditName(context.Background(), req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"Name Updated Successfully"}`))
}

func (user *UserController) clientEditPhone(w http.ResponseWriter, r *http.Request) {
	var req *pb.EditPhoneRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	clientID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "error while retrieving the client id", http.StatusBadRequest)
		return
	}
	req.UserId = clientID
	if !helper.CheckStringNumber(req.Phone) {
		http.Error(w, "please enter a valid phone number", http.StatusBadRequest)
		return
	}
	if _, err := user.Conn.ClientEditPhone(context.Background(), req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"Phone number Updated Successfully"}`))
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
	if address.Country == "" {
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

func (user *UserController) uploadFreelancerProfileImage(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "unable to parse form", http.StatusBadRequest)
		return
	}
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "unable to fetch the file from request", http.StatusBadRequest)
		return
	}
	defer file.Close()
	filebyte, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "unable to read the file", http.StatusInternalServerError)
		return
	}
	freelancerID, ok := r.Context().Value("freelancerID").(string)
	if !ok {
		http.Error(w, "error while retrieving the freelancer id", http.StatusBadRequest)
		return
	}
	req := &pb.ImageRequest{
		ObjectName: fmt.Sprintf("%s-profile", freelancerID),
		ImageData:  filebyte,
		UserId:     freelancerID,
	}
	res, err := user.Conn.FreelancerUploadProfileImage(context.Background(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	jsonData, err := json.Marshal(res)
	if err != nil {
		http.Error(w, "error while marshalling the data", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func (user *UserController) freelancerEditName(w http.ResponseWriter, r *http.Request) {
	var req *pb.EditNameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	freelancerID, ok := r.Context().Value("freelancerID").(string)
	if !ok {
		http.Error(w, "error while retrieving the freelancer id", http.StatusBadRequest)
		return
	}
	req.UserId = freelancerID
	if !helper.CheckString(req.Name) {
		http.Error(w, "please enter a valid name", http.StatusBadRequest)
		return
	}
	if _, err := user.Conn.FreelancerEditName(context.Background(), req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"Name Updated Successfully"}`))
}

func (user *UserController) freelancerEditPhone(w http.ResponseWriter, r *http.Request) {
	var req *pb.EditPhoneRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	freelancerID, ok := r.Context().Value("freelancerID").(string)
	if !ok {
		http.Error(w, "error while retrieving the freelancer id", http.StatusBadRequest)
		return
	}
	req.UserId = freelancerID
	if !helper.CheckStringNumber(req.Phone) {
		http.Error(w, "please enter a valid phone number", http.StatusBadRequest)
		return
	}
	if _, err := user.Conn.FreelancerEditPhone(context.Background(), req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"Phone number Updated Successfully"}`))
}

func (user *UserController) freelancerAddSkill(w http.ResponseWriter, r *http.Request) {
	var req *pb.SkillRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	freelancerID, ok := r.Context().Value("freelancerID").(string)
	if !ok {
		http.Error(w, "error while retrieving the freelancer id", http.StatusBadRequest)
		return
	}
	req.UserId = freelancerID

	if _, err := user.Conn.FreelancerAddSkill(context.Background(), req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"Skill Added Successfully"}`))
}

func (user *UserController) freelancerDeleteSkill(w http.ResponseWriter, r *http.Request) {
	freelancerID, ok := r.Context().Value("freelancerID").(string)
	if !ok {
		http.Error(w, "error while retrieving the freelancer id", http.StatusBadRequest)
		return
	}
	req := &pb.SkillRequest{
		UserId: freelancerID,
	}
	queryParams := r.URL.Query()
	skillId, err := strconv.Atoi(queryParams.Get("skill_id"))
	if err != nil {
		http.Error(w, "error while parsing the client id to in", http.StatusBadRequest)
		return
	}
	req.SkillId = int32(skillId)
	if _, err := user.Conn.FreelancerDeleteSkill(context.Background(), req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"Skill Deleted Successfully"}`))
}

func (user *UserController) freelancerGetAllSkill(w http.ResponseWriter, r *http.Request) {
	freelancerID, ok := r.Context().Value("freelancerID").(string)
	if !ok {
		http.Error(w, "error while retrieving the freelancer id", http.StatusBadRequest)
		return
	}
	req := &pb.GetUserById{
		Id: freelancerID,
	}
	skills, err := user.Conn.FreelancerGetAllSkill(context.Background(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	skillData := []*pb.SkillResponse{}
	for {
		skill, err := skills.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		skillData = append(skillData, skill)
	}
	jsonData, err := json.Marshal(skillData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(skillData) == 0 {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message":"No Skill Added"}`))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func (user *UserController) freelancerAddExperience(w http.ResponseWriter, r *http.Request) {
	var req *pb.AddExperienceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if helper.CheckNegativeStringNumber(req.Experience) {
		http.Error(w, "please enter a valid experience", http.StatusBadRequest)
		return
	}
	if !helper.CheckNumberInString(req.Experience) {
		http.Error(w, "please enter a valid experience", http.StatusBadRequest)
		return
	}
	if !helper.CheckYear(req.Experience) {
		http.Error(w, "pleae enter a valid experience which contains the number of years", http.StatusBadRequest)
		return
	}
	userID, ok := r.Context().Value("freelancerID").(string)
	if !ok {
		http.Error(w, "error whil retrieving freelancerID", http.StatusBadRequest)
		return
	}
	req.UserId = userID
	if _, err := user.Conn.FreelancerAddExperience(context.Background(), req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"Experience Added Successfully"}`))
}

func (user *UserController) freelancerAddTitle(w http.ResponseWriter, r *http.Request) {
	var req *pb.AddTitleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !helper.CheckString(req.Title) {
		http.Error(w, "please enter a valid title", http.StatusBadRequest)
		return
	}
	userID, ok := r.Context().Value("freelancerID").(string)
	if !ok {
		http.Error(w, "error while retrieving the freelancer id", http.StatusBadRequest)
		return
	}
	req.UserId = userID
	if _, err := user.Conn.FreelancerAddTitle(context.Background(), req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"Title Added Successfully"}`))
}

func (user *UserController) getFreelancerProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("freelancerID").(string)
	if !ok {
		http.Error(w, "error while retrieving the freelancer id", http.StatusBadRequest)
		return
	}

	freelancer, err := user.Conn.GetFreelancerById(context.Background(), &pb.GetUserById{
		Id: userID,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	category, err := user.Conn.GetCategoryById(context.Background(), &pb.GetCategoryByIdRequest{
		Id: freelancer.CategoryId,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	address, err := user.Conn.FreelancerGetAddress(context.Background(), &pb.GetUserById{
		Id: userID,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	imageData, err := user.Conn.FreelancerGetProfileImage(context.Background(), &pb.GetUserById{
		Id: userID,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	skills, err := user.Conn.FreelancerGetAllSkill(context.Background(), &pb.GetUserById{
		Id: userID,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	skillData := []*pb.SkillResponse{}
	for {
		skill, err := skills.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		skillData = append(skillData, skill)
	}
	educations, err := user.Conn.FreelancerGetEducation(context.Background(), &pb.GetUserById{
		Id: userID,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	educationData := []*pb.EducationResponse{}
	for {
		education, err := educations.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		educationData = append(educationData, education)
	}
	profile, err := user.Conn.FreelancerGetProfile(context.Background(), &pb.GetUserById{
		Id: userID,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res := helperstruct.FreelancerProfile{
		Id:                       freelancer.Id,
		Name:                     freelancer.Name,
		Email:                    freelancer.Email,
		Phone:                    freelancer.Phone,
		Image:                    imageData.Url,
		ExperienceInCurrentField: profile.ExperienceInCurrentField,
		Category:                 category.Category,
		Skills:                   skillData,
		Educations:               educationData,
		Address:                  address,
	}
	jsonData, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func (user *UserController) freelancerAddEducation(w http.ResponseWriter, r *http.Request) {
	var req *pb.EducationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !helper.CheckString(req.Degree) {
		http.Error(w, "enter a valid degree", http.StatusBadRequest)
		return
	}
	if !helper.CheckString(req.Institution) {
		http.Error(w, "enter a valid Institution", http.StatusBadRequest)
		return
	}
	if !helper.CheckDate(req.StartDate) {
		http.Error(w, "enter a valid date", http.StatusBadRequest)
		return
	}
	if !helper.CheckDate(req.EndDate) {
		http.Error(w, "enter a valid date", http.StatusBadRequest)
		return
	}
	userID, ok := r.Context().Value("freelancerID").(string)
	if !ok {
		http.Error(w, "error while retrieving the freelancer id", http.StatusBadRequest)
		return
	}
	req.UserId = userID
	if _, err := user.Conn.FreelancerAddEducation(context.Background(), req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"Education Added Successfully"}`))
}

func (user *UserController) freelancerEditEducation(w http.ResponseWriter, r *http.Request) {
	var req *pb.EducationResponse
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !helper.CheckString(req.Degree) {
		http.Error(w, "enter a valid degree", http.StatusBadRequest)
		return
	}
	if !helper.CheckString(req.Institution) {
		http.Error(w, "enter a valid Institution", http.StatusBadRequest)
		return
	}
	if !helper.CheckDate(req.StartDate) {
		http.Error(w, "enter a valid date", http.StatusBadRequest)
		return
	}
	if !helper.CheckDate(req.EndDate) {
		http.Error(w, "enter a valid date", http.StatusBadRequest)
		return
	}
	userID, ok := r.Context().Value("freelancerID").(string)
	if !ok {
		http.Error(w, "error while retrieving the freelancer id", http.StatusBadRequest)
		return
	}
	req.UserId = userID
	queryParams := r.URL.Query()
	educationId := queryParams.Get("education_id")
	req.Id = educationId
	if _, err := user.Conn.FreelancerEditEducation(context.Background(), req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"Education Edited Successfully"}`))
}

func (user *UserController) freelancerRemoveEducation(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	educationId := queryParams.Get("education_id")
	if educationId == "" {
		http.Error(w, "please provide a valid education id", http.StatusBadRequest)
		return
	}
	if _, err := user.Conn.FreelancerRemoveEducation(context.Background(), &pb.EducationById{
		EducationId: educationId,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"Education Deleted Successfully"}`))
}

func (user *UserController) getAllClientNotifications(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "error while retrieving the freelancer id", http.StatusBadRequest)
		return
	}
	notifications, err := user.NotificationConn.GetAllNotification(context.Background(), &pb.GetNotificationsByUserId{
		UserId: userID,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	notificationData := []*pb.NotificationResponse{}
	for {
		notification, err := notifications.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		notificationData = append(notificationData, notification)
	}
	jsonData, err := json.Marshal(notificationData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if len(notificationData) == 0 {
		w.Write([]byte(`{"message":"you don't have any notifications yet"}`))
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

func (user *UserController) getAllFreelancerNotifications(w http.ResponseWriter, r *http.Request) {
	freelancerID, ok := r.Context().Value("freelancerID").(string)
	if !ok {
		http.Error(w, "error while retrieving the freelancer id", http.StatusBadRequest)
		return
	}
	notifications, err := user.NotificationConn.GetAllNotification(context.Background(), &pb.GetNotificationsByUserId{
		UserId: freelancerID,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	notificationData := []*pb.NotificationResponse{}
	for {
		notification, err := notifications.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		notificationData = append(notificationData, notification)
	}
	jsonData, err := json.Marshal(notificationData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if len(notificationData) == 0 {
		w.Write([]byte(`{"message":"you don't have any notifications yet"}`))
		return
	}
	w.Write(jsonData)
}

func (user *UserController) clientPaymentForProject(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	userId := queryParams.Get("user_id")
	projectId := queryParams.Get("project_id")
	url := fmt.Sprintf("http://payment-service:4009/user/project/payment?user_id=%s&project_id=%s", userId, projectId)
	req, err := http.NewRequest("GET", url, r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	client := &http.Client{}
	req.Header = r.Header
	res, err := client.Do(req)
	if err != nil || res == nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()
	for k, v := range res.Header {
		w.Header()[k] = v
	}
	w.WriteHeader(res.StatusCode)
	io.Copy(w, res.Body)
}
func (user *UserController) verifyPayment(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	userId := queryParams.Get("user_id")
	paymentRef := queryParams.Get("payment_ref")
	orderId := queryParams.Get("order_id")
	signature := queryParams.Get("signature")
	id := queryParams.Get("id")
	total := queryParams.Get("total")
	projectId := queryParams.Get("project_id")
	url := fmt.Sprintf("http://payment-service:4009/payment/verify?user_id=%s&payment_ref=%s&order_id=%s&signature=%s&id=%s&total=%s&project_id=%s", userId, paymentRef, orderId, signature, id, total, projectId)
	req, err := http.NewRequest("GET", url, r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	client := &http.Client{}
	req.Header = r.Header
	res, err := client.Do(req)
	if err != nil || res == nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()
	for k, v := range res.Header {
		w.Header()[k] = v
	}
	w.WriteHeader(res.StatusCode)
	io.Copy(w, res.Body)
}
func (user *UserController) paymentVerified(w http.ResponseWriter, r *http.Request) {
	req, err := http.NewRequest("GET", "http://payment-service:4009/payment/verified", r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	client := &http.Client{}
	req.Header = r.Header
	res, err := client.Do(req)
	if err != nil || res == nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()
	for k, v := range res.Header {
		w.Header()[k] = v
	}
	w.WriteHeader(res.StatusCode)
	io.Copy(w, res.Body)
}

func (user *UserController) addReviewForFreelancer(w http.ResponseWriter, r *http.Request) {
	var req *pb.UserReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "error whil retrieving companyId", http.StatusBadRequest)
		return
	}
	if !helper.CheckString(req.Description) {
		http.Error(w, "please enter a valid description", http.StatusBadRequest)
		return
	}
	if helper.CheckNegative(req.Rating) {
		http.Error(w, "please provide a valid rating within 5", http.StatusBadRequest)
		return
	}
	if req.Rating > 5 {
		http.Error(w, "please provide a valid rating within 5", http.StatusBadRequest)
		return
	}
	queryParams := r.URL.Query()
	freelancerId := queryParams.Get("freelancer_id")
	projectId := queryParams.Get("project_id")
	req.UserId = userID
	req.FreelancerId = freelancerId
	req.ProjectId = projectId
	if _, err := user.ReviewConn.UserAddReview(context.Background(), req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"user added review successfully"}`))
}
func (user *UserController) getReviewForFreelancer(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	freelancerId := queryParams.Get("freelancer_id")
	reviews, err := user.ReviewConn.GetReview(context.Background(), &pb.ReviewById{
		Id: freelancerId,
	})
	if freelancerId == "" {
		http.Error(w, "please select a freelancer to get reviews", http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	reviewData := []*pb.ReviewResponse{}
	for {
		review, err := reviews.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		reviewData = append(reviewData, review)
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if len(reviewData) == 0 {
		w.Write([]byte(`{"message":"no review yet"}`))
		return
	}
	jsonData, err := json.Marshal(reviewData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write(jsonData)
}
func (user *UserController) deleteReview(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	freelancerId := queryParams.Get("freelancer_id")
	projectId := queryParams.Get("project_id")
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "error whil retrieving companyId", http.StatusBadRequest)
		return
	}
	req := &pb.UserReviewRequest{
		UserId:       userID,
		FreelancerId: freelancerId,
		ProjectId:    projectId,
	}
	if _, err := user.ReviewConn.RemoveReview(context.Background(), req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"user review removed Successfully"}`))
}

func (user *UserController) reportClient(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	userId := queryParams.Get("client_id")
	if userId == "" {
		http.Error(w, "please select a user to report", http.StatusBadRequest)
		return
	}
	if _, err := user.Conn.ReportUser(context.Background(), &pb.GetUserById{
		Id: userId,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"user reported successfully"}`))
}
