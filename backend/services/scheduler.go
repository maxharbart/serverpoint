package services

import (
	"log"
	"sync"
	"time"

	"serverpoint/config"
	"serverpoint/models"

	"gorm.io/gorm"
)

type Scheduler struct {
	db       *gorm.DB
	cfg      *config.Config
	interval time.Duration
	workers  int
	stop     chan struct{}
}

func NewScheduler(db *gorm.DB, cfg *config.Config) *Scheduler {
	return &Scheduler{
		db:       db,
		cfg:      cfg,
		interval: 30 * time.Second,
		workers:  5,
		stop:     make(chan struct{}),
	}
}

func (s *Scheduler) Start() {
	log.Println("Scheduler started, collecting every", s.interval)
	go func() {
		// Run immediately on start
		s.collectAll()

		ticker := time.NewTicker(s.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				s.collectAll()
			case <-s.stop:
				log.Println("Scheduler stopped")
				return
			}
		}
	}()
}

func (s *Scheduler) Stop() {
	close(s.stop)
}

func (s *Scheduler) collectAll() {
	var servers []models.Server
	s.db.Find(&servers)

	if len(servers) == 0 {
		return
	}

	log.Printf("Collecting metrics for %d servers", len(servers))

	jobs := make(chan models.Server, len(servers))
	var wg sync.WaitGroup

	// Start worker pool
	for i := 0; i < s.workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for server := range jobs {
				s.collectOne(server)
			}
		}()
	}

	// Send jobs
	for _, server := range servers {
		jobs <- server
	}
	close(jobs)

	wg.Wait()
}

func (s *Scheduler) collectOne(server models.Server) {
	result := CollectMetrics(server, s.cfg)

	// Update server status
	s.db.Model(&server).Updates(map[string]interface{}{
		"online":   result.Online,
		"os":       result.OS,
		"hostname": result.Hostname,
		"uptime":   result.Uptime,
	})

	if !result.Online {
		return
	}

	// Save metric
	s.db.Create(&result.Metric)

	// Replace processes (delete old, insert new)
	s.db.Where("server_id = ?", server.ID).Delete(&models.Process{})
	for _, proc := range result.Processes {
		s.db.Create(&proc)
	}

	// Prune old metrics (keep last 24h)
	cutoff := time.Now().Add(-24 * time.Hour)
	s.db.Where("server_id = ? AND created_at < ?", server.ID, cutoff).Delete(&models.Metric{})
}
