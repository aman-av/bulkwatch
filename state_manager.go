package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

// State represents the application state for incremental processing
type State struct {
	LastProcessed map[string]string `json:"last_processed"` // exchange_dealtype -> date
	LastEnriched  map[string]string `json:"last_enriched"`  // exchange_dealtype -> date
	Statistics    StateStatistics   `json:"statistics"`
	mu            sync.RWMutex
}

// StateStatistics tracks processing statistics
type StateStatistics struct {
	TotalDealsProcessed int       `json:"total_deals_processed"`
	TotalDealsEnriched  int       `json:"total_deals_enriched"`
	LastRun             time.Time `json:"last_run"`
}

const stateFile = "state.json"

// Global state instance
var appState *State
var stateOnce sync.Once

// GetState returns the singleton state instance
func GetState() *State {
	stateOnce.Do(func() {
		appState = &State{
			LastProcessed: make(map[string]string),
			LastEnriched:  make(map[string]string),
		}
		appState.Load()
	})
	return appState
}

// Load reads state from file
func (s *State) Load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(stateFile)
	if err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist, use defaults
			return nil
		}
		return fmt.Errorf("failed to read state file: %w", err)
	}

	if err := json.Unmarshal(data, s); err != nil {
		return fmt.Errorf("failed to parse state file: %w", err)
	}

	return nil
}

// Save writes state to file
func (s *State) Save() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Statistics.LastRun = time.Now()

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	if err := os.WriteFile(stateFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write state file: %w", err)
	}

	return nil
}

// GetLastEnriched returns the last enriched date for a deal type
func (s *State) GetLastEnriched(exchange Exchange, dealType DealType) string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	key := fmt.Sprintf("%s_%s", exchange, dealType)
	return s.LastEnriched[key]
}

// SetLastEnriched updates the last enriched date
func (s *State) SetLastEnriched(exchange Exchange, dealType DealType, date string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := fmt.Sprintf("%s_%s", exchange, dealType)
	s.LastEnriched[key] = date
}

// GetLastProcessed returns the last processed date for a deal type
func (s *State) GetLastProcessed(exchange Exchange, dealType DealType) string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	key := fmt.Sprintf("%s_%s", exchange, dealType)
	return s.LastProcessed[key]
}

// SetLastProcessed updates the last processed date
func (s *State) SetLastProcessed(exchange Exchange, dealType DealType, date string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := fmt.Sprintf("%s_%s", exchange, dealType)
	s.LastProcessed[key] = date
}

// IncrementDealsProcessed increments the total deals processed counter
func (s *State) IncrementDealsProcessed(count int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Statistics.TotalDealsProcessed += count
}

// IncrementDealsEnriched increments the total deals enriched counter
func (s *State) IncrementDealsEnriched(count int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Statistics.TotalDealsEnriched += count
}

// NeedsEnrichment checks if a file needs to be enriched
func (s *State) NeedsEnrichment(exchange Exchange, dealType DealType, fileDate string) bool {
	lastEnriched := s.GetLastEnriched(exchange, dealType)

	// If never enriched, needs enrichment
	if lastEnriched == "" {
		return true
	}

	// If file date is after last enriched date, needs enrichment
	return fileDate > lastEnriched
}

// PrintStatus prints the current state status
func (s *State) PrintStatus() {
	s.mu.RLock()
	defer s.mu.RUnlock()

	fmt.Println("\n📊 Current State:")
	fmt.Println("================")

	fmt.Println("\nLast Enriched:")
	for key, date := range s.LastEnriched {
		fmt.Printf("  %s: %s\n", key, date)
	}

	fmt.Println("\nStatistics:")
	fmt.Printf("  Total Deals Processed: %d\n", s.Statistics.TotalDealsProcessed)
	fmt.Printf("  Total Deals Enriched: %d\n", s.Statistics.TotalDealsEnriched)
	if !s.Statistics.LastRun.IsZero() {
		fmt.Printf("  Last Run: %s\n", s.Statistics.LastRun.Format("2006-01-02 15:04:05"))
	}
}

// Made with Bob
