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

	r.Post("/freelancer/signup", user.freelancerSignup)
	r.Post("/freelancer/login", user.freelancerLogin)
	r.Post("/freelancer/logout", middleware.FreelancerMiddleware(user.freelancerLogout))

	r.Post("/admin/login", user.adminLogin)
	r.Post("/admin/logout", middleware.AdminMiddleware(user.adminLogout))
}
