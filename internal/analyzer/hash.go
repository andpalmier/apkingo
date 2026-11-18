package analyzer

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

// SetHashes calculates MD5, SHA1, and SHA256 hashes of the APK file.
// It streams the file once through all three hash functions.
func (app *AndroidApp) SetHashes(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file %q: %w", path, err)
	}
	defer func() { _ = f.Close() }()

	md5Hash := md5.New()
	sha1Hash := sha1.New()
	sha256Hash := sha256.New()

	if _, err := io.Copy(io.MultiWriter(md5Hash, sha1Hash, sha256Hash), f); err != nil {
		return fmt.Errorf("failed to hash file %q: %w", path, err)
	}

	app.Hashes.Sha256 = fmt.Sprintf("%x", sha256Hash.Sum(nil))
	app.Hashes.Sha1 = fmt.Sprintf("%x", sha1Hash.Sum(nil))
	app.Hashes.Md5 = fmt.Sprintf("%x", md5Hash.Sum(nil))

	return nil
}
