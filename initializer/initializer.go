package initializer

import (
	"fmt"
	"os"

	"github.com/Nishad4140/SkillSync_ApiGateway/controllers/notificationController"
	usercontroller "github.com/Nishad4140/SkillSync_ApiGateway/controllers/userController"
	"github.com/Nishad4140/SkillSync_ApiGateway/helper"
	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
)

func Connect(r *chi.Mux) {
	if err := godotenv.Load("../.env"); err != nil {
		fmt.Println("error secret cannot be read")
	}
	secret := os.Getenv("SECRET")
	userConn, err := helper.DialGrpc(":4001")
	if err != nil {
		fmt.Println("cannot connect to user service", err)
	}
	notificationConn, err := helper.DialGrpc(":4007")
	if err != nil {
		fmt.Println("cannot connect to notification service", err)
	}
	userController := usercontroller.NewUserServiceClient(userConn, secret)
	notificationController := notificationController.NewNotificationServiceClient(notificationConn, secret)

	userController.InitializeUserControllers(r)
	notificationController.InitializeNotificationController(r)
}
