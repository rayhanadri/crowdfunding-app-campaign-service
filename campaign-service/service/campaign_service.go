package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/rayhanadri/crowdfunding-app-campaign-service/gen/go/campaign/v1"
	"github.com/rayhanadri/crowdfunding-app-campaign-service/helper"
	"github.com/rayhanadri/crowdfunding-app-campaign-service/models"
	"github.com/rayhanadri/crowdfunding-app-campaign-service/repository"
)

// CampaignService interface defines the Service methods for campaign operations
type CampaignService interface {
	CreateCampaign(ctx context.Context, req *campaign.CreateCampaignRequest) (*campaign.CreateCampaignResponse, error)
	GetCampaignByID(ctx context.Context, req *campaign.GetCampaignByIDRequest) (*campaign.GetCampaignByIDResponse, error)
	DeleteCampaignByID(ctx context.Context, req *campaign.DeleteCampaignByIDRequest) (*campaign.DeleteCampaignByIDResponse, error)
	UpdateCampaignByID(ctx context.Context, req *campaign.UpdateCampaignByIDRequest) (*campaign.UpdateCampaignByIDResponse, error)
	GetCampaignsByUserID(ctx context.Context, req *campaign.GetCampaignsByUserIDRequest) (*campaign.GetCampaignsByUserIDResponse, error)
}

// campaignService is the struct implementation of CampaignService
type campaignService struct {
	campaign.UnimplementedCampaignServiceServer
	campaignRepo repository.CampaignRepository
}

// NewCampaignService initializes and returns a new campaignService instance with a given Campaign repository
func NewCampaignService(campaignRepo repository.CampaignRepository) *campaignService {
	return &campaignService{campaignRepo: campaignRepo}
}

func (s *campaignService) CreateCampaign(ctx context.Context, req *campaign.CreateCampaignRequest) (*campaign.CreateCampaignResponse, error) {
	// Create a new uuid
	uuid := uuid.New()

	// Prepare a struct for campaign
	campaignPayload := models.CampaignDB{
		ID:           uuid.String(),
		UserID:       req.UserId,
		Title:        req.Title,
		Description:  req.Description,
		TargetAmount: req.TargetAmount,
		Deadline:     req.Deadline.AsTime(),
		Category:     helper.MapCategoryDB(int32(req.Category)),
		MinDonation:  req.MinDonation,
	}

	// Insert campaign to database
	campaignInterface, err := s.campaignRepo.CreateCampaign(campaignPayload)
	if err != nil {
		return nil, err
	}

	// Cast the campaignInterface type to models.CampaignDB
	createdCampaign, ok := campaignInterface.(models.CampaignDB)
	if !ok {
		return nil, fmt.Errorf("failed to cast created campaign")
	}

	res := &campaign.CreateCampaignResponse{
		CreatedCampaign: []*campaign.Campaign{
			{
				Id:              createdCampaign.ID,
				UserId:          createdCampaign.UserID,
				Title:           createdCampaign.Title,
				Description:     createdCampaign.Description,
				TargetAmount:    createdCampaign.TargetAmount,
				CollectedAmount: createdCampaign.CollectedAmount,
				Deadline:        timestamppb.New(createdCampaign.Deadline),
				Status:          campaign.CampaignStatus(helper.MapStatusProto(createdCampaign.Status)),
				Category:        campaign.CampaignCategory(helper.MapCateogryProto(createdCampaign.Category)),
				MinDonation:     createdCampaign.MinDonation,
				CreatedAt:       timestamppb.New(createdCampaign.CreatedAt),
				UpdatedAt:       timestamppb.New(createdCampaign.UpdatedAt),
			},
		},
	}
	return res, nil
}

func (s *campaignService) GetCampaignByID(ctx context.Context, req *campaign.GetCampaignByIDRequest) (*campaign.GetCampaignByIDResponse, error) {
	// Get campaign by id
	campaignInterface, err := s.campaignRepo.GetCampaignByID(req.Id)
	if err != nil {
		return nil, err
	}

	// Cast the campaignInterface type to models.CampaignDB
	getCampaign, ok := campaignInterface.(models.CampaignDB)
	if !ok {
		return nil, fmt.Errorf("failed to cast campaign")
	}

	res := &campaign.GetCampaignByIDResponse{
		Campaign: []*campaign.Campaign{
			{
				Id:              getCampaign.ID,
				UserId:          getCampaign.UserID,
				Title:           getCampaign.Title,
				Description:     getCampaign.Description,
				TargetAmount:    getCampaign.TargetAmount,
				CollectedAmount: getCampaign.CollectedAmount,
				Deadline:        timestamppb.New(getCampaign.Deadline),
				Status:          campaign.CampaignStatus(helper.MapStatusProto(getCampaign.Status)),
				Category:        campaign.CampaignCategory(helper.MapCateogryProto(getCampaign.Category)),
				MinDonation:     getCampaign.MinDonation,
				CreatedAt:       timestamppb.New(getCampaign.CreatedAt),
				UpdatedAt:       timestamppb.New(getCampaign.UpdatedAt),
			},
		},
	}
	return res, nil
}

func (s *campaignService) DeleteCampaignByID(ctx context.Context, req *campaign.DeleteCampaignByIDRequest) (*campaign.DeleteCampaignByIDResponse, error) {
	// Delete campaign by id
	err := s.campaignRepo.DeleteCampaignByID(req.Id)
	if err != nil {
		return nil, err
	}

	return &campaign.DeleteCampaignByIDResponse{
		DeleteResponse: &emptypb.Empty{},
	}, nil
}

func (s *campaignService) UpdateCampaignByID(ctx context.Context, req *campaign.UpdateCampaignByIDRequest) (*campaign.UpdateCampaignByIDResponse, error) {
	// Prepare a struct for campaign
	campaignPayload := models.CampaignDB{
		Title:        req.Title,
		Description:  req.Description,
		TargetAmount: req.TargetAmount,
		Deadline:     req.Deadline.AsTime(),
		Status:       helper.MapStatusDB(int32(req.Status)),
		Category:     helper.MapCategoryDB(int32(req.Category)),
		MinDonation:  req.MinDonation,
	}

	// Check if user is trying to update status manually to completed
	if campaignPayload.Status == "completed" {
		return nil, status.Error(codes.PermissionDenied, "You cannot manually set status to COMPLETED")
	}

	// Update campaign by id
	campaignInterface, err := s.campaignRepo.UpdateCampaignByID(req.Id, req.UserId, campaignPayload)
	if err != nil {
		return nil, err
	}

	// Cast the campaignInterface type to models.CampaignDB
	updatedCampaign, ok := campaignInterface.(models.CampaignDB)
	if !ok {
		return nil, fmt.Errorf("failed to cast updated campaign")
	}

	return &campaign.UpdateCampaignByIDResponse{
		UpdatedCampaign: []*campaign.Campaign{
			{
				Id:              updatedCampaign.ID,
				UserId:          updatedCampaign.UserID,
				Title:           updatedCampaign.Title,
				Description:     updatedCampaign.Description,
				TargetAmount:    updatedCampaign.TargetAmount,
				CollectedAmount: updatedCampaign.CollectedAmount,
				Deadline:        timestamppb.New(updatedCampaign.Deadline),
				Status:          campaign.CampaignStatus(helper.MapStatusProto(updatedCampaign.Status)),
				Category:        campaign.CampaignCategory(helper.MapCateogryProto(updatedCampaign.Category)),
				MinDonation:     updatedCampaign.MinDonation,
				CreatedAt:       timestamppb.New(updatedCampaign.CreatedAt),
				UpdatedAt:       timestamppb.New(updatedCampaign.UpdatedAt),
			},
		},
	}, nil
}

func (s *campaignService) GetCampaignsByUserID(ctx context.Context, req *campaign.GetCampaignsByUserIDRequest) (*campaign.GetCampaignsByUserIDResponse, error) {
	// Get campaign by user id
	campaignInterface, err := s.campaignRepo.GetCampaignsByUserID(req.UserId)
	if err != nil {
		return nil, err
	}

	// Cast the campaignInterface type to models.CampaignDB
	getCampaign, ok := campaignInterface.([]models.CampaignDB)
	if !ok {
		return nil, fmt.Errorf("failed to cast campaign")
	}

	var campaignList []*campaign.Campaign
	for _, val := range getCampaign {
		campaignList = append(campaignList, &campaign.Campaign{
			Id:              val.ID,
			UserId:          val.UserID,
			Title:           val.Title,
			Description:     val.Description,
			TargetAmount:    val.TargetAmount,
			CollectedAmount: val.CollectedAmount,
			Deadline:        timestamppb.New(val.Deadline),
			Status:          campaign.CampaignStatus(helper.MapStatusProto(val.Status)),
			Category:        campaign.CampaignCategory(helper.MapCateogryProto(val.Category)),
			MinDonation:     val.MinDonation,
			CreatedAt:       timestamppb.New(val.CreatedAt),
			UpdatedAt:       timestamppb.New(val.UpdatedAt),
		},
		)
	}

	return &campaign.GetCampaignsByUserIDResponse{
		Campaign: campaignList,
	}, nil
}
