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

	// Iterate through all files in the archive
	for _, file := range r.File {
		// Check if this is an APK file
		if !isAPKFile(file.Name) {
			continue
		}

		// Extract the APK file
		apkPath, err := extractFile(file, dir)
		if err != nil {
			// Clean up on error
			_ = os.RemoveAll(dir)
			return nil, fmt.Errorf("failed to extract %s: %w", file.Name, err)
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
func extractFile(file *zip.File, destDir string) (string, error) {
	// Get the filename from the path
	filename := filepath.Base(file.Name)

	// Sanitize filename to prevent path traversal
	if filename == "" || filename == "." || filename == ".." {
		return "", fmt.Errorf("invalid filename in archive: %s", file.Name)
	}

	// Create the destination file
	destPath := filepath.Join(destDir, filename)
	destFile, err := os.Create(destPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file %s: %w", destPath, err)
	}
	defer func() { _ = destFile.Close() }()

	// Open the source file from the archive
	srcFile, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file in archive: %w", err)
	}
	defer func() { _ = srcFile.Close() }()

	// Copy the content
	if _, err := io.Copy(destFile, srcFile); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return destPath, nil
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
