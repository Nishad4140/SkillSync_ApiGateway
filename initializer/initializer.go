package initializer

import (
	"fmt"
	"os"

	"github.com/Nishad4140/SkillSync_ApiGateway/controllers/notificationController"
	projectcontroller "github.com/Nishad4140/SkillSync_ApiGateway/controllers/projectController"
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

	userConn, err := helper.DialGrpc("ss-user-service:4001")
	if err != nil {
		fmt.Println("cannot connect to user service", err)
	}

	projectConn, err := helper.DialGrpc("ss-project-service:4002")
	if err != nil {
		fmt.Println("cannot connect to project service", err)
	}

	notificationConn, err := helper.DialGrpc("ss-notification-service:4007")
	if err != nil {
		fmt.Println("cannot connect to notification service", err)
	}

	userController := usercontroller.NewUserServiceClient(userConn, secret)
	projectController := projectcontroller.NewProjectServiceClient(projectConn, secret)
	notificationController := notificationController.NewNotificationServiceClient(notificationConn, secret)

	userController.InitializeUserControllers(r)
	projectController.InitializeProjectController(r)
	notificationController.InitializeNotificationController(r)
}
