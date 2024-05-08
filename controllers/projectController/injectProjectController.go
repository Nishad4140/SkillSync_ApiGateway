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
	userConn, _ := helper.DialGrpc("ss-user-service:4001")
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
	r.Post("/freelancer/project", middleware.FreelancerMiddleware(project.freelancerAcknowledgeProject))
	r.Patch("/freelancer/project", middleware.FreelancerMiddleware(project.freelancerUpdateProjectStatus))
	r.Get("/freelancer/project", middleware.FreelancerMiddleware(project.getProject))
	r.Get("/freelancer/project/all", middleware.FreelancerMiddleware(project.getAllProjects))
	r.Post("/freelancer/project/management", middleware.FreelancerMiddleware(project.freelancerProjectManagement))
	r.Patch("/freelancer/project/management", middleware.FreelancerMiddleware(project.freelancerUpdateModule))
	r.Get("/freelancer/project/management", middleware.FreelancerMiddleware(project.getProjectManagement))
	r.Post("/freelancer/project/file", middleware.FreelancerMiddleware(project.freelancerUploadFile))
	r.Get("/freelancer/project/file", middleware.FreelancerMiddleware(project.getFile))

	r.Post("/admin/pakcage-types", middleware.AdminMiddleware(project.adminAddProjectType))
	r.Patch("/admin/package-types", middleware.AdminMiddleware(project.adminEditProjectType))

	r.Post("/client/requests", middleware.ClientMiddleware(project.clientAddRequest))
	r.Patch("/client/requests", middleware.ClientMiddleware(project.clientUpdateRequest))
	r.Get("/client/requests", middleware.ClientMiddleware(project.clientGetRequest))
	r.Get("/client/request/intrests", middleware.ClientMiddleware(project.getClientRequestIntrests))
	r.Post("/client/request/intrests", middleware.ClientMiddleware(project.clientAcknowledgeIntrest))
	r.Post("/client/project", middleware.ClientMiddleware(project.clientCreateProject))
	r.Patch("/client/project", middleware.ClientMiddleware(project.clientUpdateProject))
	r.Get("/client/project", middleware.ClientMiddleware(project.getProject))
	r.Get("/client/project/all", middleware.FreelancerMiddleware(project.getAllProjects))
	r.Delete("/client/project", middleware.ClientMiddleware(project.clientRemoveProject))
	r.Get("/client/project/file", middleware.FreelancerMiddleware(project.getFile))

	r.Get("/package-types", project.getAllPackageTypes)
	r.Get("/gigs", project.getGig)
}
