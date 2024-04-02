package projectcontroller

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/Nishad4140/SkillSync_ApiGateway/helper"
	"github.com/Nishad4140/SkillSync_ProtoFiles/pb"
)

func (project *Projectcontroller) freelancerCreateGig(w http.ResponseWriter, r *http.Request) {
	var req *pb.CreateGigRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

}

func (project *Projectcontroller) adminAddProjectType(w http.ResponseWriter, r *http.Request) {
	var req *pb.AddPackageTypeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !helper.CheckString(req.PackageType) {
		http.Error(w, "please enter a valid type", http.StatusBadRequest)
		return
	}
	if _, err := project.Conn.AddPackageType(context.Background(), req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"package type added successfully"}`))
}

func (project *Projectcontroller) adminEditProjectType(w http.ResponseWriter, r *http.Request) {
	var req *pb.PackageTypeResponse
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.Id == 0 {
		http.Error(w, "please enter a valid id", http.StatusBadRequest)
		return
	}
	if !helper.CheckString(req.PackageType) {
		http.Error(w, "please enter a valid type", http.StatusBadRequest)
		return
	}
	if _, err := project.Conn.EditPackageType(context.Background(), req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"package type edited successfully"}`))
}

func (project *Projectcontroller) getAllPackageTypes(w http.ResponseWriter, r *http.Request) {
	packageTypes, err := project.Conn.GetPackageType(context.Background(), nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	packageData := []*pb.PackageTypeResponse{}

	for {
		types, err := packageTypes.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		packageData = append(packageData, types)
	}
	jsonData, err := json.Marshal(packageData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
