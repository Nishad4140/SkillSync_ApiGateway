package projectcontroller

import (
	"github.com/Nishad4140/SkillSync_ApiGateway/helper"
	"github.com/Nishad4140/SkillSync_ApiGateway/middleware"
	"github.com/Nishad4140/SkillSync_ProtoFiles/pb"
	"github.com/go-chi/chi"
	"google.golang.org/grpc"
)

type Projectcontroller struct {
	Conn     pb.ProjectServiceClient
	UserConn pb.UserServiceClient
	Secret   string
}

func NewProjectServiceClient(conn *grpc.ClientConn, secret string) *Projectcontroller {
	userConn, _ := helper.DialGrpc("localhost:4001")
	return &Projectcontroller{
		Conn:     pb.NewProjectServiceClient(conn),
		UserConn: pb.NewUserServiceClient(userConn),
		Secret:   secret,
	}
}

func (project *Projectcontroller) InitializeProjectController(r *chi.Mux) {
	r.Post("/freelancers/gigs", middleware.FreelancerMiddleware(project.freelancerCreateGig))
	r.Patch("/freelancers/gigs", middleware.FreelancerMiddleware(project.freelancerUpdateGig))
	r.Get("/freelancers/gigs", middleware.FreelancerMiddleware(project.freelancerGetAllGigs))
	r.Get("/freelancer/requests", middleware.FreelancerMiddleware(project.freelancerGetClientRequests))
	r.Post("/freelancer/requests/intrests", middleware.FreelancerMiddleware(project.freelancerShowIntrests))

	r.Post("/admin/pakcage-types", middleware.AdminMiddleware(project.adminAddProjectType))
	r.Patch("/admin/package-types", middleware.AdminMiddleware(project.adminEditProjectType))

	r.Post("/client/requests", middleware.ClientMiddleware(project.clientAddRequest))
	r.Patch("/client/requests", middleware.ClientMiddleware(project.clientUpdateRequest))
	r.Get("/client/requests", middleware.ClientMiddleware(project.clientGetRequest))
	r.Get("/client/request/intrests", middleware.ClientMiddleware(project.getClientRequestIntrests))
	r.Post("/client/request/intrests", middleware.ClientMiddleware(project.clientAcknowledgeIntrest))

	r.Get("/package-types", project.getAllPackageTypes)
	r.Get("/gigs", project.getGig)
}
