package analyzer

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"os"
)

// SetHashes calculates MD5, SHA1, and SHA256 hashes of the APK file.
// It reads the file once into memory and computes all hashes from that data.
// This is more efficient than reading the file multiple times for each hash.
func (app *AndroidApp) SetHashes(path string) error {
	// Read the entire file into memory once
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file %q: %w", path, err)
	}

	// Compute SHA256
	sha256Hash := sha256.Sum256(data)
	app.Hashes.Sha256 = fmt.Sprintf("%x", sha256Hash)

	// Compute SHA1
	sha1Hash := sha1.Sum(data)
	app.Hashes.Sha1 = fmt.Sprintf("%x", sha1Hash)

	// Compute MD5
	md5Hash := md5.Sum(data)
	app.Hashes.Md5 = fmt.Sprintf("%x", md5Hash)

	return nil
}
