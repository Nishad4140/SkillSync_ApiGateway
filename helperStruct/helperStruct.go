package helperstruct

import "github.com/Nishad4140/SkillSync_ProtoFiles/pb"

type ClientProfile struct {
	Id      string
	Name    string
	Email   string
	Phone   string
	Image   string `json:"image,omitempty"`
	Address *pb.AddressResponse
}

type FreelancerProfile struct {
	Id                       string
	Name                     string
	Email                    string
	Phone                    string
	Image                    string `json:"image,omitempty"`
	ExperienceInCurrentField string `json:"experience_in_current_field,omitempty"`
	Category                 string
	Skills                   []*pb.SkillResponse
	Educations               []*pb.EducationResponse
	Address                  *pb.AddressResponse
}
