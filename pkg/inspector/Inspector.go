package inspector

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"slices"
	"sync"
)

// Config holds inspector configuration
type Config struct {
	File string
}

// Inspector tracks and records dependency rule usage patterns
type Inspector struct {
	cfg Config

	mu sync.Mutex

	// from -> to list
	usedRules map[string][]string
}

// New creates a new Inspector instance
func New(cfg Config) *Inspector {
	return &Inspector{
		cfg:       cfg,
		mu:        sync.Mutex{},
		usedRules: make(map[string][]string),
	}
}

// RecordUsage records a dependency rule usage from -> to
func (i *Inspector) RecordUsage(frm, to string) {
	i.mu.Lock()
	defer i.mu.Unlock()

	i.usedRules[frm] = append(i.usedRules[frm], to)
}

// IsRecorded checks if a dependency rule usage is recorded
func (i *Inspector) IsRecorded(frm, to string) bool {
	i.mu.Lock()
	defer i.mu.Unlock()

	return slices.Contains(i.usedRules[frm], to)
}

// Save writes recorded dependency rule usage to file
func (i *Inspector) Save() error {
	if i.cfg.File == "" {
		return nil
	}

	// Since golang.org/x/tools/go/analysis triggers Analyzer.Run multiple times,
	// we need to open the file in append mode.
	file, err := os.OpenFile(i.cfg.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	i.mu.Lock()
	defer i.mu.Unlock()

	for frm, toList := range i.usedRules {
		for _, to := range toList {
			if err := writer.Write([]string{frm, to}); err != nil {
				return fmt.Errorf("failed to write record: %w", err)
			}
		}
	}

	return nil
}

// Load creates an Inspector instance from saved records
func Load(cfg Config) (*Inspector, error) {
	if cfg.File == "" {
		return nil, errors.New("file is not specified")
	}

	file, err := os.Open(cfg.File)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	usedRules := make(map[string][]string)

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = 2
	for {
		record, err := reader.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, fmt.Errorf("error reading file: %w", err)
		}
		frm, to := record[0], record[1]
		usedRules[frm] = append(usedRules[frm], to)
	}

	inspector := New(cfg)
	inspector.usedRules = usedRules
	return inspector, nil
}
