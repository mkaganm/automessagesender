package scheduler

import (
	"github.com/robfig/cron/v3"
)

// CronClient represents the cron job service
type CronClient struct {
	cron *cron.Cron
}

// NewCronJobService creates a new CronClient
func NewCronJobService() *CronClient {
	return &CronClient{
		cron: cron.New(cron.WithSeconds()),
	}
}

// AddJob adds a new job to the cron scheduler
func (s *CronClient) AddJob(schedule string, job func()) error {
	_, err := s.cron.AddFunc(schedule, job)
	return err
}

// Start starts the cron scheduler
func (s *CronClient) Start() {
	s.cron.Start()
}
