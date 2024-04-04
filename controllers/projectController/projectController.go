package projectcontroller

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/Nishad4140/SkillSync_ApiGateway/helper"
	"github.com/Nishad4140/SkillSync_ProtoFiles/pb"
)

func (project *Projectcontroller) freelancerCreateGig(w http.ResponseWriter, r *http.Request) {
	var req *pb.CreateGigRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(req.Title) < 15 {
		http.Error(w, "please enter the title not less than 15 words", http.StatusBadRequest)
		return
	}
	if len(req.Description) < 50 {
		http.Error(w, "please enter the description not less than 50 words", http.StatusBadRequest)
		return
	}
	category, err := project.UserConn.GetCategoryById(context.Background(), &pb.GetCategoryByIdRequest{
		Id: req.CategoryId,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if category.Category == "" {
		http.Error(w, "please enter a valid category", http.StatusBadRequest)
		return
	}
	if req.SkillId == 0 {
		http.Error(w, "please enter a valid skill", http.StatusBadRequest)
		return
	}
	if req.PackageTypeId == 0 {
		http.Error(w, "please enter a valid package type", http.StatusBadRequest)
		return
	}
	if req.Price == 0 {
		http.Error(w, "please enter a valid price", http.StatusBadRequest)
		return
	}
	freelancerID, ok := r.Context().Value("freelancerID").(string)
	if !ok {
		http.Error(w, "error while retrieving the freelancer id", http.StatusBadRequest)
		return
	}
	req.FreelancerId = freelancerID

	if _, err := project.Conn.CreateGig(context.Background(), req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"gig added successfully"}`))
}

func (project *Projectcontroller) freelancerUpdateGig(w http.ResponseWriter, r *http.Request) {
	var req *pb.GigResponse
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(req.Title) < 15 {
		http.Error(w, "please enter the title not less than 15 words", http.StatusBadRequest)
		return
	}
	if len(req.Description) < 50 {
		http.Error(w, "please enter the description not less than 50 words", http.StatusBadRequest)
		return
	}
	category, err := project.UserConn.GetCategoryById(context.Background(), &pb.GetCategoryByIdRequest{
		Id: req.CategoryId,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if category.Category == "" {
		http.Error(w, "please enter a valid category", http.StatusBadRequest)
		return
	}
	if req.SkillId == 0 {
		http.Error(w, "please enter a valid skill", http.StatusBadRequest)
		return
	}
	if req.PackageTypeId == 0 {
		http.Error(w, "please enter a valid package type", http.StatusBadRequest)
		return
	}
	if req.Price == 0 {
		http.Error(w, "please enter a valid price", http.StatusBadRequest)
		return
	}
	freelancerID, ok := r.Context().Value("freelancerID").(string)
	if !ok {
		http.Error(w, "error while retrieving the freelancer id", http.StatusBadRequest)
		return
	}
	req.FreelancerId = freelancerID
	queryParams := r.URL.Query()
	gigId := queryParams.Get("gig_id")
	req.Id = gigId

	if _, err := project.Conn.UpdateGig(context.Background(), req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"gig updated successfully"}`))
}

func (project *Projectcontroller) getGig(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	gigId := queryParams.Get("gig_id")
	if gigId != "" {
		gig, err := project.Conn.GetGig(context.Background(), &pb.GetById{
			Id: gigId,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		jsonData, err := json.Marshal(gig)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)
		return
	}
	gigs, err := project.Conn.GetAllGigs(context.Background(), nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	gigData := []*pb.GigResponse{}
	for {
		gig, err := gigs.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		gigData = append(gigData, gig)
	}
	jsonData, err := json.Marshal(gigData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func (project *Projectcontroller) freelancerGetAllGigs(w http.ResponseWriter, r *http.Request) {
	freelancerID, ok := r.Context().Value("freelancerID").(string)
	if !ok {
		http.Error(w, "error while retrieving the freelancer id", http.StatusBadRequest)
		return
	}
	req := &pb.GetByUserId{
		Id: freelancerID,
	}
	gigs, err := project.Conn.GetAllFreelancerGigs(context.Background(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	gigData := []*pb.GigResponse{}

	for {
		gig, err := gigs.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		gigData = append(gigData, gig)
	}
	jsonData, err := json.Marshal(gigData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(gigData) == 0 {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message":"no gig added"}`))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
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

	queryParams := r.URL.Query()
	packageTypeId, err := strconv.Atoi(queryParams.Get("package_id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	req.Id = int32(packageTypeId)

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
