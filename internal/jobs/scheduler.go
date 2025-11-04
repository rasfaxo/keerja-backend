package jobs

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

// Job represents a background job
type Job interface {
	// Name returns the job name
	Name() string

	// Run executes the job
	Run(ctx context.Context) error

	// Schedule returns the cron schedule (e.g., "0 * * * *" for every hour)
	Schedule() string
}

// Scheduler manages background jobs
type Scheduler struct {
	cron   *cron.Cron
	jobs   map[string]Job
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
	mu     sync.RWMutex
}

// NewScheduler creates a new job scheduler
func NewScheduler() *Scheduler {
	ctx, cancel := context.WithCancel(context.Background())

	return &Scheduler{
		cron:   cron.New(cron.WithSeconds()),
		jobs:   make(map[string]Job),
		ctx:    ctx,
		cancel: cancel,
	}
}

// Register registers a new job
func (s *Scheduler) Register(job Job) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	name := job.Name()

	// Check if job already registered
	if _, exists := s.jobs[name]; exists {
		return fmt.Errorf("job %s already registered", name)
	}

	// Add job to cron
	_, err := s.cron.AddFunc(job.Schedule(), func() {
		s.runJob(job)
	})

	if err != nil {
		return fmt.Errorf("failed to add job %s to cron: %w", name, err)
	}

	s.jobs[name] = job
	fmt.Printf("Registered background job: %s (schedule: %s)\n", name, job.Schedule())

	return nil
}

// Start starts the scheduler
func (s *Scheduler) Start() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.jobs) == 0 {
		fmt.Println("No background jobs registered")
		return
	}

	s.cron.Start()
	fmt.Printf("Background job scheduler started with %d jobs\n", len(s.jobs))
}

// Stop gracefully stops the scheduler
func (s *Scheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	fmt.Println("Stopping background job scheduler...")

	// Cancel context to stop running jobs
	s.cancel()

	// Stop accepting new jobs
	ctx := s.cron.Stop()

	// Wait for cron to finish
	<-ctx.Done()

	// Wait for all running jobs to complete
	s.wg.Wait()

	fmt.Println("Background job scheduler stopped")
}

// RunNow runs a specific job immediately
func (s *Scheduler) RunNow(jobName string) error {
	s.mu.RLock()
	job, exists := s.jobs[jobName]
	s.mu.RUnlock()

	if !exists {
		return fmt.Errorf("job %s not found", jobName)
	}

	fmt.Printf("Running job %s immediately...\n", jobName)
	s.runJob(job)

	return nil
}

// runJob executes a job with error handling
func (s *Scheduler) runJob(job Job) {
	s.wg.Add(1)
	defer s.wg.Done()

	name := job.Name()
	startTime := time.Now()

	fmt.Printf("Starting job: %s\n", name)

	// Run with timeout context
	ctx, cancel := context.WithTimeout(s.ctx, 30*time.Minute)
	defer cancel()

	// Execute job
	if err := job.Run(ctx); err != nil {
		duration := time.Since(startTime)
		fmt.Printf("Job %s failed after %v: %v\n", name, duration, err)
		return
	}

	duration := time.Since(startTime)
	fmt.Printf("Job %s completed successfully in %v\n", name, duration)
}

// GetRegisteredJobs returns list of registered job names
func (s *Scheduler) GetRegisteredJobs() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	names := make([]string, 0, len(s.jobs))
	for name := range s.jobs {
		names = append(names, name)
	}

	return names
}

// GetJobInfo returns information about a specific job
func (s *Scheduler) GetJobInfo(jobName string) (map[string]interface{}, error) {
	s.mu.RLock()
	job, exists := s.jobs[jobName]
	s.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("job %s not found", jobName)
	}

	return map[string]interface{}{
		"name":     job.Name(),
		"schedule": job.Schedule(),
	}, nil
}
