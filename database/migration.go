package database

import (
	"log"

	"github.com/riskiapl/fiber-app/models"
)

func Migrate() {
	log.Println("Running migrations...")

	err := DB.AutoMigrate(
		&models.Member{},
		&models.OTP{},
		&models.PendingMember{},
		&models.ResetPasswordToken{},
		// Tambahkan model lain di sini
	)

	if err != nil {
		log.Fatalf("Could not migrate: %v", err)
	}

	log.Println("Migration completed!")
}
