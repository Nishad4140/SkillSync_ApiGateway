package projectcontroller

import (
	"github.com/Nishad4140/SkillSync_ApiGateway/middleware"
	"github.com/Nishad4140/SkillSync_ProtoFiles/pb"
	"github.com/go-chi/chi"
	"google.golang.org/grpc"
)

type Projectcontroller struct {
	Conn   pb.ProjectServiceClient
	Secret string
}

func NewProjectServiceClient(conn *grpc.ClientConn, secret string) *Projectcontroller {
	return &Projectcontroller{
		Conn:   pb.NewProjectServiceClient(conn),
		Secret: secret,
	}
}

func (project *Projectcontroller) InitializeProjectController(r *chi.Mux) {
	r.Post("/freelancers/gigs", middleware.FreelancerMiddleware(project.freelancerCreateGig))

	r.Post("/admin/pakcage-types", middleware.AdminMiddleware(project.adminAddProjectType))
	r.Patch("/admin/package-types", middleware.AdminMiddleware(project.adminEditProjectType))

	r.Get("/package-types", project.getAllPackageTypes)
}
