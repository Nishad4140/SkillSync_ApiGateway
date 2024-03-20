package main

import (
	"fmt"
	"net/http"

	"github.com/Nishad4140/SkillSync_ApiGateway/initializer"
	"github.com/go-chi/chi"
)

func main() {
	r := chi.NewRouter()
	initializer.Connect(r)
	fmt.Println("API Gateway listening on the port 4000")
	http.ListenAndServe(":4000", r)
}
