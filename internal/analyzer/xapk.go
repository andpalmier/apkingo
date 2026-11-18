// Package analyzer provides APK analysis functionality.
package analyzer

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	maxAPKEntries   = 10000
	maxAPKEntrySize = int64(4 << 30) // 4 GiB per APK
	maxTotalAPKSize = int64(8 << 30) // 8 GiB per archive
)

// IsXAPK checks if the given file is an XAPK or APKS file.
// XAPK/APKS files are ZIP archives containing one or more APK files.
func IsXAPK(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".xapk" || ext == ".apks"
}

// ExtractAPKs extracts APK files from an XAPK/APKS archive.
// It returns a list of paths to the extracted APK files.
// The caller is responsible for cleaning up the extracted files.
func ExtractAPKs(xapkPath string) ([]string, error) {
	// Open the XAPK file
	r, err := zip.OpenReader(xapkPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open XAPK file: %w", err)
	}
	defer func() { _ = r.Close() }()

	// Create a temporary directory for extracted APKs
	dir, err := os.MkdirTemp("", "apkingo-xapk-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}

	var apkPaths []string
	var apkCount int
	var totalSize uint64

	// Iterate through all files in the archive
	for _, file := range r.File {
		// Check if this is an APK file
		if !isAPKFile(file.Name) {
			continue
		}

		apkCount++
		if apkCount > maxAPKEntries {
			_ = os.RemoveAll(dir)
			return nil, fmt.Errorf("XAPK contains more than %d APK files", maxAPKEntries)
		}
		if file.UncompressedSize64 > uint64(maxAPKEntrySize) {
			_ = os.RemoveAll(dir)
			return nil, fmt.Errorf("APK entry %q is too large: %d bytes", file.Name, file.UncompressedSize64)
		}

		// Extract the APK file
		apkPath, written, err := extractFile(file, dir)
		if err != nil {
			// Clean up on error
			_ = os.RemoveAll(dir)
			return nil, fmt.Errorf("failed to extract %s: %w", file.Name, err)
		}
		if err := addExtractedSize(&totalSize, written, uint64(maxTotalAPKSize)); err != nil {
			_ = os.RemoveAll(dir)
			return nil, err
		}

		apkPaths = append(apkPaths, apkPath)
	}

	if len(apkPaths) == 0 {
		_ = os.RemoveAll(dir)
		return nil, fmt.Errorf("no APK files found in XAPK archive")
	}

	return apkPaths, nil
}

// isAPKFile checks if the filename has an .apk extension.
func isAPKFile(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	return ext == ".apk"
}

// extractFile extracts a single file from the ZIP archive to the destination directory.
func extractFile(file *zip.File, destDir string) (string, int64, error) {
	// Get the filename from the path
	filename := filepath.Base(file.Name)

	// Sanitize filename to prevent path traversal
	if filename == "" || filename == "." || filename == ".." {
		return "", 0, fmt.Errorf("invalid filename in archive: %s", file.Name)
	}

	// Create the destination file
	destPath, err := uniqueDestPath(filepath.Join(destDir, filename))
	if err != nil {
		return "", 0, fmt.Errorf("failed to choose destination for %s: %w", file.Name, err)
	}
	destFile, err := os.OpenFile(destPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		return "", 0, fmt.Errorf("failed to create file %s: %w", destPath, err)
	}
	defer func() { _ = destFile.Close() }()

	// Open the source file from the archive
	srcFile, err := file.Open()
	if err != nil {
		return "", 0, fmt.Errorf("failed to open file in archive: %w", err)
	}
	defer func() { _ = srcFile.Close() }()

	// Copy the content
	written, err := io.Copy(destFile, io.LimitReader(srcFile, maxAPKEntrySize+1))
	if err != nil {
		return "", written, fmt.Errorf("failed to write file: %w", err)
	}
	if written > maxAPKEntrySize {
		return "", written, fmt.Errorf("APK entry %q exceeds size limit", file.Name)
	}

	return destPath, written, nil
}

func addExtractedSize(total *uint64, written int64, limit uint64) error {
	if written < 0 {
		return fmt.Errorf("invalid extracted size: %d", written)
	}
	size := uint64(written)
	if *total > limit || size > limit-*total {
		return fmt.Errorf("XAPK contains more than %d bytes of APK data", limit)
	}
	*total += size
	return nil
}

func uniqueDestPath(path string) (string, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return path, nil
	} else if err != nil {
		return "", err
	}

	ext := filepath.Ext(path)
	base := strings.TrimSuffix(path, ext)
	for i := 1; ; i++ {
		candidate := fmt.Sprintf("%s_%d%s", base, i, ext)
		if _, err := os.Stat(candidate); os.IsNotExist(err) {
			return candidate, nil
		} else if err != nil {
			return "", err
		}
	}
}

// CleanupExtractedFiles removes the temporary directory containing extracted APKs.
func CleanupExtractedFiles(paths []string) error {
	// Get unique directories
	dirs := make(map[string]bool)
	for _, path := range paths {
		dir := filepath.Dir(path)
		dirs[dir] = true
	}

	// Remove each directory
	for dir := range dirs {
		if err := os.RemoveAll(dir); err != nil {
			return fmt.Errorf("failed to remove directory %s: %w", dir, err)
		}
	}

	return nil
}
