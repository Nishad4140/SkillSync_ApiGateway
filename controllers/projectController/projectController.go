package projectcontroller

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/Nishad4140/SkillSync_ApiGateway/helper"
	helperstruct "github.com/Nishad4140/SkillSync_ApiGateway/helperStruct"
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

	var viewGig helperstruct.FilterQuery

	viewGig.Page, _ = strconv.Atoi(queryParams.Get("page"))
	viewGig.Limit, _ = strconv.Atoi(queryParams.Get("limit"))
	viewGig.Query = queryParams.Get("query")
	viewGig.Filter = queryParams.Get("filter")
	viewGig.SortBy = queryParams.Get("sort_by")
	viewGig.SortDesc, _ = strconv.ParseBool(queryParams.Get("sort_desc"))

	gigs, err := project.Conn.GetAllGigs(context.Background(), &pb.GigFilterQuery{
		Page:     int32(viewGig.Page),
		Limit:    int32(viewGig.Limit),
		Query:    viewGig.Query,
		Filter:   viewGig.Filter,
		SortBy:   viewGig.SortBy,
		SortDesc: viewGig.SortDesc,
	})
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

func (project *Projectcontroller) clientAddRequest(w http.ResponseWriter, r *http.Request) {
	var req *pb.AddClientGigRequest
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
	if !helper.CheckDate(req.DeliveryDate) {
		http.Error(w, "enter a valid date", http.StatusBadRequest)
		return
	}

	date, err := helper.ConvertStringToDate(req.DeliveryDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !date.After(time.Now()) {
		http.Error(w, "please enter a valid date", http.StatusBadRequest)
		return
	}
	if req.CategoryId == 0 {
		http.Error(w, "enter a valid category id", http.StatusBadRequest)
		return
	}
	if req.SkillId == 0 {
		http.Error(w, "enter a valid skill id", http.StatusBadRequest)
		return
	}
	if req.Price == 0 {
		http.Error(w, "please enter a valid price", http.StatusBadRequest)
		return
	}
	clientID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "error while retrieving the freelancer id", http.StatusBadRequest)
		return
	}
	req.ClientId = clientID

	if _, err := project.Conn.ClientAddRequest(context.Background(), req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"request added successfully"}`))
}

func (project *Projectcontroller) clientUpdateRequest(w http.ResponseWriter, r *http.Request) {
	var req *pb.ClientRequestResponse
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
	if !helper.CheckDate(req.DeliveryDate) {
		http.Error(w, "enter a valid date", http.StatusBadRequest)
		return
	}
	date, err := helper.ConvertStringToDate(req.DeliveryDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !date.After(time.Now()) {
		http.Error(w, "please enter a valid date", http.StatusBadRequest)
		return
	}
	if req.CategoryId == 0 {
		http.Error(w, "enter a valid category id", http.StatusBadRequest)
		return
	}
	if req.SkillId == 0 {
		http.Error(w, "enter a valid skill id", http.StatusBadRequest)
		return
	}
	if req.Price == 0 {
		http.Error(w, "please enter a valid price", http.StatusBadRequest)
		return
	}
	clientID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "error while retrieving the freelancer id", http.StatusBadRequest)
		return
	}
	req.ClientId = clientID
	queryParams := r.URL.Query()
	request_id := queryParams.Get("request_id")
	req.Id = request_id
	if _, err := project.Conn.ClientUpdateRequest(context.Background(), req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"request updated successfully"}`))
}

func (project *Projectcontroller) clientGetRequest(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	reqId := queryParams.Get("request_id")
	if reqId != "" {
		req, err := project.Conn.GetClientRequest(context.Background(), &pb.GetById{
			Id: reqId,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		jsonData, err := json.Marshal(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)
		return
	}
	clientID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "error while retrieving the freelancer id", http.StatusBadRequest)
		return
	}

	var viewReq helperstruct.FilterQuery

	viewReq.Page, _ = strconv.Atoi(queryParams.Get("page"))
	viewReq.Limit, _ = strconv.Atoi(queryParams.Get("limit"))
	viewReq.Query = queryParams.Get("query")
	viewReq.Filter = queryParams.Get("filter")
	viewReq.SortBy = queryParams.Get("sort_by")
	viewReq.SortDesc, _ = strconv.ParseBool(queryParams.Get("sort_desc"))

	req := &pb.RequestFilterQuery{
		UserId:   clientID,
		Page:     int32(viewReq.Page),
		Limit:    int32(viewReq.Limit),
		Query:    viewReq.Query,
		Filter:   viewReq.Filter,
		SortBy:   viewReq.SortBy,
		SortDesc: viewReq.SortDesc,
	}

	clientReqs, err := project.Conn.GetAllClientRequest(context.Background(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	reqData := []*pb.ClientRequestResponse{}

	for {
		gig, err := clientReqs.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		reqData = append(reqData, gig)
	}
	jsonData, err := json.Marshal(reqData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(reqData) == 0 {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message":"no request added"}`))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func (project *Projectcontroller) freelancerGetClientRequests(w http.ResponseWriter, r *http.Request) {
	freelancerID, ok := r.Context().Value("freelancerID").(string)
	if !ok {
		http.Error(w, "error while retrieving the freelancer id", http.StatusBadRequest)
		return
	}

	var viewReq helperstruct.FilterQuery

	queryParams := r.URL.Query()
	viewReq.Page, _ = strconv.Atoi(queryParams.Get("page"))
	viewReq.Limit, _ = strconv.Atoi(queryParams.Get("limit"))
	viewReq.Query = queryParams.Get("query")
	viewReq.Filter = queryParams.Get("filter")
	viewReq.SortBy = queryParams.Get("sort_by")
	viewReq.SortDesc, _ = strconv.ParseBool(queryParams.Get("sort_desc"))

	req := &pb.RequestFilterQuery{
		UserId:   freelancerID,
		Page:     int32(viewReq.Page),
		Limit:    int32(viewReq.Limit),
		Query:    viewReq.Query,
		Filter:   viewReq.Filter,
		SortBy:   viewReq.SortBy,
		SortDesc: viewReq.SortDesc,
	}

	reqClients, err := project.Conn.GetAllClientRequestForFreelancers(context.Background(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	reqData := []*pb.ClientRequestResponse{}

	for {
		req, err := reqClients.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		reqData = append(reqData, req)
	}
	jsonData, err := json.Marshal(reqData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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

func (project *Projectcontroller) freelancerShowIntrests(w http.ResponseWriter, r *http.Request) {
	var req *pb.IntrestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	freelancerID, ok := r.Context().Value("freelancerID").(string)
	if !ok {
		http.Error(w, "error while retrieving the freelancer id", http.StatusBadRequest)
		return
	}
	req.UserId = freelancerID

	queryParams := r.URL.Query()
	reqId := queryParams.Get("req_id")

	req.RequestId = reqId
	if _, err := project.Conn.ShowIntrest(context.Background(), req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"intrest added successfully"}`))
}

func (project *Projectcontroller) getClientRequestIntrests(w http.ResponseWriter, r *http.Request) {
	req := &pb.GetAllIntrestRequest{}

	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "error while retrieving the freelancer id", http.StatusBadRequest)
		return
	}
	req.UserId = userID

	queryParams := r.URL.Query()
	reqId := queryParams.Get("req_id")

	req.RequestId = reqId

	reqs, err := project.Conn.GetAllIntrest(context.Background(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var reqData []*pb.IntrestResponse
	for {
		req, err := reqs.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		reqData = append(reqData, req)
	}
	jsonData, err := json.Marshal(reqData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func (project *Projectcontroller) clientAcknowledgeIntrest(w http.ResponseWriter, r *http.Request) {
	req := &pb.IntrestAcknowledgmentRequest{}

	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "error while retrieving the freelancer id", http.StatusBadRequest)
		return
	}
	req.UserId = userID

	queryParams := r.URL.Query()
	intrestId := queryParams.Get("intrest_id")

	req.IntrestId = intrestId

	if _, err := project.Conn.ClientIntrestAcknowledgment(context.Background(), req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"intrest acknowledged added successfully"}`))
}
