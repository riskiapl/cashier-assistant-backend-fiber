package services

import (
	"log"
	"time"

	"github.com/riskiapl/fiber-app/database"
	"github.com/riskiapl/fiber-app/repository"
)

type CronService struct {
	authRepo *repository.AuthRepository
}

func NewCronService() *CronService {
	return &CronService{
		authRepo: repository.NewAuthRepository(database.DB),
	}
}

// CleanupExpiredData removes pending members and OTPs older than the specified duration
func (s *CronService) CleanupExpiredData() {
	// Delete expired pending members
	deletedPendingMembers, err := s.authRepo.DeleteExpiredPendingMembers(24 * time.Hour)
	if err != nil {
		log.Printf("Error deleting expired pending members: %v", err)
	} else {
		log.Printf("Deleted %d expired pending members", deletedPendingMembers)
	}

	// Delete expired OTPs
	deletedOTPs, err := s.authRepo.DeleteExpiredOTPs(24 * time.Hour)
	if err != nil {
		log.Printf("Error deleting expired OTPs: %v", err)
	} else {
		log.Printf("Deleted %d expired OTPs", deletedOTPs)
	}

	// Delete expired reset password tokens
	deletedResetTokens, err := s.authRepo.DeleteExpiredResetPasswordTokens(720 * time.Hour) // 30 days (720 hours)
	if err != nil {
		log.Printf("Error deleting expired reset password tokens: %v", err)
	} else {
		log.Printf("Deleted %d expired reset password tokens", deletedResetTokens)
	}
}
