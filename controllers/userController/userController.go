package usercontroller

import (
	"encoding/json"
	"net/http"

	"github.com/Nishad4140/SkillSync_ApiGateway/helper"
	"github.com/Nishad4140/SkillSync_ProtoFiles/pb"
)

func (user *UserController) clientSignup(w http.ResponseWriter, r *http.Request) {
	if cookie, _ := r.Cookie("ClientToken"); cookie != nil {
		http.Error(w, "you are already logged in", http.StatusConflict)
		return
	}
	var req pb.UserSignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.Email == "" {
		http.Error(w, "please enter a mail id", http.StatusBadRequest)
		return
	}
	if !helper.CheckString(req.Name) {
		http.Error(w, "please enter a name", http.StatusBadRequest)
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

	// if req.OTP == "" {

	// }
}
