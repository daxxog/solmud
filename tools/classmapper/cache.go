package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	CACHE_VERSION    = "1.0"
	CACHE_INDEX_FILE = "cache_index.json"
	CHECKSUMS_FILE   = "checksums.json"
)

type CacheManager struct {
	cacheDir  string
	checksums map[string]FileChecksum
	verbose   bool
}

type FileChecksum struct {
	Size    int64     `json:"size"`
	ModTime time.Time `json:"mod_time"`
	SHA256  string    `json:"sha256"`
}

type CacheIndex struct {
	Version   string    `json:"version"`
	CreatedAt time.Time `json:"created_at"`
	LastUsed  time.Time `json:"last_used"`
	CacheDir  string    `json:"cache_dir"`
}

func NewCacheManager(cacheDir string, verbose bool) (*CacheManager, error) {
	cm := &CacheManager{
		cacheDir:  cacheDir,
		checksums: make(map[string]FileChecksum),
		verbose:   verbose,
	}

	// Create cache directory if it doesn't exist
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create cache directory: %w", err)
	}

	// Initialize or load cache index
	if err := cm.initializeCacheIndex(); err != nil {
		return nil, fmt.Errorf("failed to initialize cache: %w", err)
	}

	// Load existing checksums
	if err := cm.LoadChecksums(); err != nil {
		if cm.verbose {
			fmt.Fprintf(os.Stderr, "  Warning: failed to load existing checksums: %v\n", err)
		}
		// Don't fail - start with empty cache
	}

	return cm, nil
}

func (cm *CacheManager) initializeCacheIndex() error {
	indexPath := filepath.Join(cm.cacheDir, CACHE_INDEX_FILE)

	// Check if index exists
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		// Create new index
		index := CacheIndex{
			Version:   CACHE_VERSION,
			CreatedAt: time.Now(),
			LastUsed:  time.Now(),
			CacheDir:  cm.cacheDir,
		}

		data, err := json.MarshalIndent(index, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal cache index: %w", err)
		}

		if err := os.WriteFile(indexPath, data, 0644); err != nil {
			return fmt.Errorf("failed to write cache index: %w", err)
		}

		if cm.verbose {
			fmt.Fprintf(os.Stderr, "  Created new cache index in %s\n", cm.cacheDir)
		}
	} else {
		// Update last used time
		data, err := os.ReadFile(indexPath)
		if err != nil {
			return fmt.Errorf("failed to read cache index: %w", err)
		}

		var index CacheIndex
		if err := json.Unmarshal(data, &index); err != nil {
			return fmt.Errorf("failed to unmarshal cache index: %w", err)
		}

		index.LastUsed = time.Now()

		updatedData, err := json.MarshalIndent(index, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal updated cache index: %w", err)
		}

		if err := os.WriteFile(indexPath, updatedData, 0644); err != nil {
			return fmt.Errorf("failed to update cache index: %w", err)
		}
	}

	return nil
}

func (cm *CacheManager) LoadChecksums() error {
	checksumsPath := filepath.Join(cm.cacheDir, CHECKSUMS_FILE)

	if _, err := os.Stat(checksumsPath); os.IsNotExist(err) {
		// No existing checksums file
		if cm.verbose {
			fmt.Fprintf(os.Stderr, "  No existing checksums file at %s\n", checksumsPath)
		}
		return nil
	}

	data, err := os.ReadFile(checksumsPath)
	if err != nil {
		return fmt.Errorf("failed to read checksums file: %w", err)
	}

	if err := json.Unmarshal(data, &cm.checksums); err != nil {
		return fmt.Errorf("failed to unmarshal checksums: %w", err)
	}

	if cm.verbose {
		fmt.Fprintf(os.Stderr, "  Loaded %d checksums from %s\n", len(cm.checksums), checksumsPath)
	}

	return nil
}

func (cm *CacheManager) SaveChecksums() error {
	checksumsPath := filepath.Join(cm.cacheDir, CHECKSUMS_FILE)

	if cm.verbose {
		fmt.Fprintf(os.Stderr, "  Saving %d checksums to %s\n", len(cm.checksums), checksumsPath)
		for file, checksum := range cm.checksums {
			fmt.Fprintf(os.Stderr, "    %s: size=%d, sha256=%s\n", file, checksum.Size, checksum.SHA256[:8]+"...")
			break // Just show first one
		}
	}

	data, err := json.MarshalIndent(cm.checksums, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal checksums: %w", err)
	}

	if err := os.WriteFile(checksumsPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write checksums file: %w", err)
	}

	return nil
}

func (cm *CacheManager) CalculateChecksum(filePath string) (FileChecksum, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return FileChecksum{}, fmt.Errorf("failed to stat file %s: %w", filePath, err)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return FileChecksum{}, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	hash := sha256.Sum256(data)

	return FileChecksum{
		Size:    info.Size(),
		ModTime: info.ModTime(),
		SHA256:  hex.EncodeToString(hash[:]),
	}, nil
}

func (cm *CacheManager) IsCacheValid(classFile string) bool {
	stored, exists := cm.checksums[classFile]
	if !exists {
		return false
	}

	current, err := cm.CalculateChecksum(classFile)
	if err != nil {
		if cm.verbose {
			fmt.Fprintf(os.Stderr, "  Warning: failed to calculate checksum for %s: %v\n", classFile, err)
		}
		return false
	}

	// Check if file hasn't changed
	return current.Size == stored.Size &&
		current.ModTime.Equal(stored.ModTime) &&
		current.SHA256 == stored.SHA256
}

func (cm *CacheManager) GetCachedOutput(classFile string) ([]byte, error) {
	if !cm.IsCacheValid(classFile) {
		return nil, fmt.Errorf("cache invalid")
	}

	cacheFileName := cm.getCacheFileName(classFile)
	cacheFilePath := filepath.Join(cm.cacheDir, cacheFileName)

	data, err := os.ReadFile(cacheFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read cached output: %w", err)
	}

	return data, nil
}

func (cm *CacheManager) StoreCachedOutput(classFile string, output []byte) error {
	// Calculate and store checksum
	checksum, err := cm.CalculateChecksum(classFile)
	if err != nil {
		return fmt.Errorf("failed to calculate checksum for storage: %w", err)
	}

	cm.checksums[classFile] = checksum

	// Save cache file
	cacheFileName := cm.getCacheFileName(classFile)
	cacheFilePath := filepath.Join(cm.cacheDir, cacheFileName)

	if err := os.WriteFile(cacheFilePath, output, 0644); err != nil {
		return fmt.Errorf("failed to write cache file: %w", err)
	}

	return nil
}

func (cm *CacheManager) getCacheFileName(classFile string) string {
	// Extract class name from path
	className := filepath.Base(classFile)
	// Remove .class extension and add .javap.cache extension
	className = className[:len(className)-6] + ".javap.cache"
	return className
}

func (cm *CacheManager) Cleanup() error {
	if cm.verbose {
		fmt.Fprintf(os.Stderr, "  Cleaning up cache...\n")
	}

	// Get list of all cache files
	entries, err := os.ReadDir(cm.cacheDir)
	if err != nil {
		return fmt.Errorf("failed to read cache directory: %w", err)
	}

	removed := 0
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		fileName := entry.Name()
		if fileName == CACHE_INDEX_FILE || fileName == CHECKSUMS_FILE {
			continue
		}

		// Check if this is a javap cache file and if its source still exists
		if filepath.Ext(fileName) == ".cache" {
			// Extract class name from cache file name
			className := fileName[:len(fileName)-12] + ".class" // Remove .javap.cache, add .class

			// Try to find the source file in known locations
			found := false
			searchPaths := []string{
				filepath.Dir(cm.cacheDir), // Parent directory
				"srcAllDummysRemoved/bin",
				"bytecode/client",
			}

			for _, searchPath := range searchPaths {
				fullPath := filepath.Join(searchPath, className)
				if _, err := os.Stat(fullPath); err == nil {
					found = true
					break
				}
			}

			if !found {
				// Source file no longer exists, remove cache
				cachePath := filepath.Join(cm.cacheDir, fileName)
				if err := os.Remove(cachePath); err != nil {
					if cm.verbose {
						fmt.Fprintf(os.Stderr, "  Warning: failed to remove stale cache file %s: %v\n", fileName, err)
					}
				} else {
					removed++
					if cm.verbose {
						fmt.Fprintf(os.Stderr, "  Removed stale cache file: %s\n", fileName)
					}
				}

				// Also remove from checksums
				delete(cm.checksums, filepath.Join(filepath.Dir(cm.cacheDir), className))
			}
		}
	}

	if cm.verbose && removed > 0 {
		fmt.Fprintf(os.Stderr, "  Cleaned up %d stale cache files\n", removed)
	}

	// Save updated checksums
	return cm.SaveChecksums()
}

func (cm *CacheManager) GetStats() (hits int, misses int, total int) {
	total = len(cm.checksums)
	// We don't track hits directly, so return basic stats
	return 0, 0, total
}
