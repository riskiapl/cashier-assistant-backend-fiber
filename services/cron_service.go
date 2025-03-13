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
func (s *CronService) CleanupExpiredData(expirationDuration time.Duration) {
	// Delete expired pending members
	deletedPendingMembers, err := s.authRepo.DeleteExpiredPendingMembers(expirationDuration)
	if err != nil {
		log.Printf("Error deleting expired pending members: %v", err)
	} else {
		log.Printf("Deleted %d expired pending members", deletedPendingMembers)
	}

	// Delete expired OTPs
	deletedOTPs, err := s.authRepo.DeleteExpiredOTPs(expirationDuration)
	if err != nil {
		log.Printf("Error deleting expired OTPs: %v", err)
	} else {
		log.Printf("Deleted %d expired OTPs", deletedOTPs)
	}
}
