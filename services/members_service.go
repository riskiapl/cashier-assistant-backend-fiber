package services

import (
	"golang.org/x/crypto/bcrypt"

	"github.com/riskiapl/fiber-app/database"
	"github.com/riskiapl/fiber-app/repository"
	"github.com/riskiapl/fiber-app/types"
)

type MemberService struct {
	repo *repository.MemberRepository
}

func NewMemberService() *MemberService {
	return &MemberService{
		repo: repository.NewMemberRepository(database.DB),
	}
}

func (s *MemberService) GetMembers(limit, offset int) (*types.MembersResponse, error) {
	members, count, err := s.repo.GetMembers(limit, offset)
	if err != nil {
		return nil, err
	}

	var memberResponses []types.MemberResponse
	for _, member := range members {
		memberResponses = append(memberResponses, types.MemberResponse{
			ID:        member.ID,
			Username:  member.Username,
			Email:     member.Email,
			Status:    member.Status,
			Avatar:    member.Avatar,
			CreatedAt: member.CreatedAt,
			UpdatedAt: member.UpdatedAt,
		})
	}

	return &types.MembersResponse{
		Members: memberResponses,
		Count:   count,
	}, nil
}

func (s *MemberService) GetMember(id uint) (*types.MemberResponse, error) {
	member, err := s.repo.GetMemberByID(id)
	if err != nil {
		return nil, err
	}

	return &types.MemberResponse{
		ID:        member.ID,
		Username:  member.Username,
		Email:     member.Email,
		Status:    member.Status,
		Avatar:    member.Avatar,
		CreatedAt: member.CreatedAt,
		UpdatedAt: member.UpdatedAt,
	}, nil
}

func (s *MemberService) UpdateMember(id uint, req *types.UpdateMemberRequest) (*types.MemberResponse, error) {
	member, err := s.repo.GetMemberByID(id)
	if err != nil {
		return nil, err
	}

	if req.Username != "" {
		member.Username = req.Username
	}
	if req.Email != "" {
		member.Email = req.Email
	}
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		member.Password = string(hashedPassword)
		member.PlainPassword = req.Password
	}
	if req.Status != "" {
		member.Status = req.Status
	}
	if req.Avatar != "" {
		member.Avatar = req.Avatar
	}
	// Add handling for Name, PhoneNumber and Address fields
	if req.Name != nil {
		member.Name = req.Name
	}
	if req.PhoneNumber != nil {
		member.PhoneNumber = req.PhoneNumber
	}
	if req.Address != nil {
		member.Address = req.Address
	}

	member.ActionType = "updated"

	if err := s.repo.UpdateMember(member); err != nil {
		return nil, err
	}

	return &types.MemberResponse{
		ID:          member.ID,
		Username:    member.Username,
		Email:       member.Email,
		Status:      member.Status,
		Avatar:      member.Avatar,
		CreatedAt:   member.CreatedAt,
		UpdatedAt:   member.UpdatedAt,
		Name:        member.Name,
		PhoneNumber: member.PhoneNumber,
		Address:     member.Address,
	}, nil
}

func (s *MemberService) DeleteMember(id uint) error {
	return s.repo.DeleteMember(id)
}
