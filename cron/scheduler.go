package cron

import (
	"log"
	"time"

	"github.com/riskiapl/fiber-app/services"
	"github.com/robfig/cron/v3"
)

var cronScheduler *cron.Cron

// StartCronJobs initializes all the scheduled tasks
func StartCronJobs() {
	// Create a new cron scheduler with standard configuration (no seconds)
	cronScheduler = cron.New()

	// Create the cron service
	cronService := services.NewCronService()

	// Schedule job to run at minute 00 of every hour (standard cron expression: "0 * * * *")
	// Format: minute hour day-of-month month day-of-week
	_, err := cronScheduler.AddFunc("0 * * * *", func() {
		log.Println("Running scheduled cleanup of expired data...")
		cronService.CleanupExpiredData(24 * time.Hour)
	})

	if err != nil {
		log.Printf("Error scheduling cleanup job: %v", err)
		return
	}

	// Start the cron scheduler in a separate goroutine
	cronScheduler.Start()

	log.Println("Cron jobs initialized successfully")

	// For logging next run time
	entry := cronScheduler.Entries()[0]
	log.Printf("Next cleanup scheduled at: %s", entry.Next.Format("2006-01-02 15:04:05"))
}

// StopCronJobs stops all running cron jobs
func StopCronJobs() {
	if cronScheduler != nil {
		cronScheduler.Stop()
		log.Println("Cron jobs stopped")
	}
}
