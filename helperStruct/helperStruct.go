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
