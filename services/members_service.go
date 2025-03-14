package services

import (
	"golang.org/x/crypto/bcrypt"

	"github.com/riskiapl/fiber-app/repository"
	"github.com/riskiapl/fiber-app/types"
)

type MemberService struct {
	repo *repository.MemberRepository
}

func NewMemberService(repo *repository.MemberRepository) *MemberService {
	return &MemberService{repo: repo}
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
	member.ActionType = "updated"

	if err := s.repo.UpdateMember(member); err != nil {
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

func (s *MemberService) DeleteMember(id uint) error {
	return s.repo.DeleteMember(id)
}
