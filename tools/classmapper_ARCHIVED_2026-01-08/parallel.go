package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

type ParseError struct {
	err       error
	file_path string
	worker_id int
}

type ProgressTracker struct {
	total        int
	completed    int
	start_time   time.Time
	verbose      bool
	mu           sync.Mutex
	worker_stats map[int]int
}

func NewProgressTracker(total int, verbose bool) *ProgressTracker {
	return &ProgressTracker{
		total:        total,
		completed:    0,
		start_time:   time.Now(),
		verbose:      verbose,
		worker_stats: make(map[int]int),
	}
}

func (self *ProgressTracker) Update(worker_id int, completed bool) {
	self.mu.Lock()
	defer self.mu.Unlock()

	if completed {
		self.completed++
		self.worker_stats[worker_id]++
	}

	if self.verbose && self.completed%5 == 0 {
		self.display_progress()
	}
}

func (self *ProgressTracker) display_progress() {
	percentage := float64(self.completed) / float64(self.total) * 100.0
	elapsed := time.Since(self.start_time).Seconds()

	bar_width := 50
	filled := int(percentage / 100.0 * float64(bar_width))
	bar := strings.Repeat("█", filled) + strings.Repeat(" ", bar_width-filled)

	fmt.Fprintf(os.Stderr, "\r  [%s] %d/%d files (%.1f%%) %.1fs",
		bar, self.completed, self.total, percentage, elapsed)
}

func (self *ProgressTracker) Finalize() {
	self.mu.Lock()
	defer self.mu.Unlock()

	if self.verbose && self.completed == self.total {
		elapsed := time.Since(self.start_time).Seconds()

		fmt.Fprintf(os.Stderr, "\n  Worker stats: ")
		var stats []string
		for id, count := range self.worker_stats {
			stats = append(stats, fmt.Sprintf("W%d:%d", id, count))
		}
		fmt.Fprintf(os.Stderr, "%s\n", strings.Join(stats, ", "))
		fmt.Fprintf(os.Stderr, "  Parallel parsing time: %.1fs\n", elapsed)
	}
}

func (self *BytecodeParser) ParseAllParallel(source_path string, workers int, progress *ProgressTracker) ([]ClassInfo, error) {
	files, err := os.ReadDir(source_path)
	if err != nil {
		return nil, fmt.Errorf(ERR_DIR_NOT_FOUND, source_path)
	}

	// Filter valid files
	var valid_files []string
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".bytecode.txt") {
			continue
		}
		valid_files = append(valid_files, filepath.Join(source_path, file.Name()))
	}

	if len(valid_files) == 0 {
		panic(fmt.Sprintf(ERR_NO_CLASSES_FOUND, source_path))
	}

	// Determine worker count
	num_workers := workers
	if num_workers <= 0 {
		num_workers = runtime.NumCPU() - 1
		if num_workers < 1 {
			num_workers = 1
		}
	}

	// Create semaphore for concurrency control
	sem := semaphore.NewWeighted(int64(num_workers))

	// Create channels
	jobs := make(chan string, len(valid_files))
	results := make(chan *ClassInfo, len(valid_files))
	errors := make(chan ParseError, len(valid_files))

	// Create worker pool
	var wg sync.WaitGroup
	ctx := context.Background()
	for w := 0; w < num_workers; w++ {
		wg.Add(1)
		go func(worker_id int) {
			defer wg.Done()
			for file_path := range jobs {
				// Acquire semaphore slot
				if err := sem.Acquire(ctx, 1); err != nil {
					errors <- ParseError{
						err:       err,
						file_path: file_path,
						worker_id: worker_id,
					}
					continue
				}

				class_info, err := self.parse_file(file_path)
				sem.Release(1) // Always release semaphore

				if err != nil {
					errors <- ParseError{
						err:       err,
						file_path: file_path,
						worker_id: worker_id,
					}
					if progress != nil {
						progress.Update(worker_id, false)
					}
					continue
				}

				if progress != nil {
					progress.Update(worker_id, true)
				}
				results <- class_info
			}
		}(w)
	}

	// Start job distribution in a goroutine
	go func() {
		for _, file_path := range valid_files {
			select {
			case jobs <- file_path:
			case <-errors:
				// Fail-fast: close jobs and return
				close(jobs)
				return
			}
		}
		close(jobs)
	}()

	// Start result/error collection
	go func() {
		wg.Wait()
		close(results)
		close(errors)
	}()

	// Collect results with fail-fast error handling
	var classes []ClassInfo
	for {
		select {
		case result, ok := <-results:
			if !ok {
				return classes, nil
			}
			classes = append(classes, *result)
		case err, ok := <-errors:
			if !ok {
				return classes, nil
			}
			// Fail-fast: return immediately with context
			return nil, fmt.Errorf("worker %d failed on %s: %w",
				err.worker_id, err.file_path, err.err)
		}
	}
}
