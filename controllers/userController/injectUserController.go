package usercontroller

import (
	"github.com/Nishad4140/SkillSync_ApiGateway/helper"
	"github.com/Nishad4140/SkillSync_ApiGateway/middleware"
	"github.com/Nishad4140/SkillSync_ProtoFiles/pb"
	"github.com/go-chi/chi"
	"google.golang.org/grpc"
)

type UserController struct {
	Conn             pb.UserServiceClient
	NotificationConn pb.NotificationServiceClient
	Secret           string
}

func NewUserServiceClient(conn *grpc.ClientConn, secret string) *UserController {
	notificationConn, _ := helper.DialGrpc("ss-notification-service:4007")
	return &UserController{
		Conn:             pb.NewUserServiceClient(conn),
		NotificationConn: pb.NewNotificationServiceClient(notificationConn),
		Secret:           secret,
	}
}

func (user *UserController) InitializeUserControllers(r *chi.Mux) {
	r.Post("/client/signup", user.clientSignup)
	r.Post("/client/login", user.clientLogin)
	r.Post("/client/logout", middleware.ClientMiddleware(user.clientLogout))
	r.Get("/client/profile", middleware.ClientMiddleware(user.getClientProfile))
	r.Post("/client/profile/address", middleware.ClientMiddleware(user.clientAddAddress))
	r.Patch("/client/profile/address", middleware.ClientMiddleware(user.clientUpdateAddress))
	r.Get("/client/profile/address", middleware.ClientMiddleware(user.clientGetAddress))
	r.Post("/client/profile/uploadprofileimage", middleware.ClientMiddleware(user.uploadClientProfileImage))
	r.Patch("/client/profile/name", middleware.ClientMiddleware(user.clientEditName))
	r.Patch("/client/profile/phone", middleware.ClientMiddleware(user.clientEditPhone))
	r.Get("/client/notifications", middleware.ClientMiddleware(user.getAllClientNotifications))
	r.Get("/client/project/payment", middleware.CorsMiddleware(user.clientPaymentForProject))
	r.Get("/payment/verify", middleware.CorsMiddleware(user.verifyPayment))
	r.Get("/payment/verified", middleware.CorsMiddleware(user.paymentVerified))

	r.Post("/freelancer/signup", user.freelancerSignup)
	r.Post("/freelancer/login", user.freelancerLogin)
	r.Post("/freelancer/logout", middleware.FreelancerMiddleware(user.freelancerLogout))
	r.Get("/freelancer/profile", middleware.FreelancerMiddleware(user.getFreelancerProfile))
	r.Post("/freelancer/profile/address", middleware.FreelancerMiddleware(user.freelancerAddAddress))
	r.Patch("/freelancer/profile/address", middleware.FreelancerMiddleware(user.freelancerUpdateAddress))
	r.Get("/freelancer/profile/address", middleware.FreelancerMiddleware(user.freelancerGetAddress))
	r.Post("/freelancer/profile/uploadprofileimage", middleware.FreelancerMiddleware(user.uploadFreelancerProfileImage))
	r.Patch("/freelancer/profile/name", middleware.FreelancerMiddleware(user.freelancerEditName))
	r.Patch("/freelancer/profile/phone", middleware.FreelancerMiddleware(user.freelancerEditPhone))
	r.Post("/freelancer/profile/skill", middleware.FreelancerMiddleware(user.freelancerAddSkill))
	r.Delete("/freelancer/profile/skill", middleware.FreelancerMiddleware(user.freelancerDeleteSkill))
	r.Get("/freelancer/profile/skill", middleware.FreelancerMiddleware(user.freelancerGetAllSkill))
	r.Post("/freelancer/profile/experience", middleware.FreelancerMiddleware(user.freelancerAddExperience))
	r.Post("/freelancer/profile/title", middleware.FreelancerMiddleware(user.freelancerAddTitle))
	r.Post("/freelancer/education", middleware.FreelancerMiddleware(user.freelancerAddEducation))
	r.Patch("/freelancer/education", middleware.FreelancerMiddleware(user.freelancerEditEducation))
	r.Delete("/freelancer/education", middleware.FreelancerMiddleware(user.freelancerRemoveEducation))
	r.Get("/freelancer/notifications", middleware.FreelancerMiddleware(user.getAllFreelancerNotifications))

	r.Post("/admin/login", user.adminLogin)
	r.Post("/admin/logout", middleware.AdminMiddleware(user.adminLogout))
	r.Post("/admin/category", middleware.AdminMiddleware(user.addCategory))
	r.Patch("/admin/category", middleware.AdminMiddleware(user.updateCategory))
	r.Post("/admin/skill", middleware.AdminMiddleware(user.adminAddSkill))
	r.Patch("/admin/skill", middleware.AdminMiddleware(user.adminUpdateSkill))
	r.Post("/client/block", middleware.AdminMiddleware(user.blockClient))
	r.Post("/client/unblock", middleware.AdminMiddleware(user.unBlockClient))
	r.Post("/freelancer/block", middleware.AdminMiddleware(user.blockFreelancer))
	r.Post("/freelancer/unblock", middleware.AdminMiddleware(user.unBlockFreelancer))

	r.Get("/categories", user.getAllCategories)
	r.Get("/skills", user.getAllSkills)
}
