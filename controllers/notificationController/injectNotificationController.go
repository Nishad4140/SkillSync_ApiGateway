package notificationController

import (
	"github.com/Nishad4140/SkillSync_ProtoFiles/pb"
	"github.com/go-chi/chi"
	"google.golang.org/grpc"
)

type NotificationController struct {
	Conn   pb.NotificationServiceClient
	Secret string
}

func NewNotificationServiceClient(conn *grpc.ClientConn, secret string) *NotificationController {
	return &NotificationController{
		Conn:   pb.NewNotificationServiceClient(conn),
		Secret: secret,
	}
}

func (n *NotificationController) InitializeNotificationController(r *chi.Mux){
	
}