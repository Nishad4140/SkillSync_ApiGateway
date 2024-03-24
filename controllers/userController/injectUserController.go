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
	notificationCOnn, _ := helper.DialGrpc("localhost:4007")
	return &UserController{
		Conn:             pb.NewUserServiceClient(conn),
		NotificationConn: pb.NewNotificationServiceClient(notificationCOnn),
		Secret:           secret,
	}
}

func (user *UserController) InitializeUserControllers(r *chi.Mux) {
	r.Post("/client/signup", user.clientSignup)
	r.Post("/client/login", user.clientLogin)
	r.Post("/client/logout", middleware.ClientMiddleware(user.clientLogout))
	r.Post("/client/profile/address", middleware.ClientMiddleware(user.clientAddAddress))
	r.Patch("/client/profile/address", middleware.ClientMiddleware(user.clientUpdateAddress))
	r.Get("/client/profile/address", middleware.ClientMiddleware(user.clientGetAddress))

	r.Post("/freelancer/signup", user.freelancerSignup)
	r.Post("/freelancer/login", user.freelancerLogin)
	r.Post("/freelancer/logout", middleware.FreelancerMiddleware(user.freelancerLogout))
	r.Post("/freelancer/profile/address", middleware.FreelancerMiddleware(user.freelancerAddAddress))
	r.Patch("/freelancer/profile/address", middleware.FreelancerMiddleware(user.freelancerUpdateAddress))
	r.Get("/freelancer/profile/address", middleware.FreelancerMiddleware(user.freelancerGetAddress))

	r.Post("/admin/login", user.adminLogin)
	r.Post("/admin/logout", middleware.AdminMiddleware(user.adminLogout))
	r.Post("/admin/category", middleware.AdminMiddleware(user.addCategory))
	r.Patch("/admin/category", middleware.AdminMiddleware(user.updateCategory))
	r.Get("/categories", user.getAllCategories)
	r.Post("/admin/skill", middleware.AdminMiddleware(user.adminAddSkill))
	r.Patch("/admin/skill", middleware.AdminMiddleware(user.adminUpdateSkill))
	r.Get("/skills", user.getAllSkills)
	r.Post("/client/block", middleware.AdminMiddleware(user.blockClient))
	r.Post("/client/unblock", middleware.AdminMiddleware(user.unBlockClient))
	r.Post("/freelancer/block", middleware.AdminMiddleware(user.blockFreelancer))
	r.Post("/freelancer/unblock", middleware.AdminMiddleware(user.unBlockFreelancer))
}
