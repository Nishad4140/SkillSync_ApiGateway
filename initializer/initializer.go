package initializer

import (
	"fmt"
	"os"

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
	userController := usercontroller.NewUserServiceClient(userConn, secret)

	userController.InitializeUserControllers(r)
}
