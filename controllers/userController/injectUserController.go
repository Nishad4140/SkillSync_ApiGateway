package usercontroller

import (
	"github.com/Nishad4140/SkillSync_ApiGateway/helper"
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
	r.Post("/freelancer/signup", user.freelancerSignup)
}
