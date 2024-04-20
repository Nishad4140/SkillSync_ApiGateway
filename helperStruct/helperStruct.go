package helperstruct

import "github.com/Nishad4140/SkillSync_ProtoFiles/pb"

type ClientProfile struct {
	Id      string              `json:"id"`
	Name    string              `json:"name"`
	Email   string              `json:"email"`
	Phone   string              `json:"phone"`
	Image   string              `json:"image,omitempty"`
	Address *pb.AddressResponse `json:"address"`
}

type FreelancerProfile struct {
	Id                       string                  `json:"id"`
	Name                     string                  `json:"name"`
	Email                    string                  `json:"email"`
	Phone                    string                  `json:"phone"`
	Image                    string                  `json:"image,omitempty"`
	ExperienceInCurrentField string                  `json:"experience_in_current_field,omitempty"`
	Category                 string                  `json:"category"`
	Skills                   []*pb.SkillResponse     `json:"skills"`
	Educations               []*pb.EducationResponse `json:"educations"`
	Address                  *pb.AddressResponse     `json:"address"`
}

type FilterQuery struct {
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`
	Query    string `json:"query"`   //search key word
	Filter   string `json:"filter"`  //to specify the column name
	SortBy   string `json:"sort_by"` //to specify column to set the sorting
	SortDesc bool   `json:"sort_desc"`
}