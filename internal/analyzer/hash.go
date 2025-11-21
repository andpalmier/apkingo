package analyzer

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

// SetHashes calculates hashes of the APK file
func (app *AndroidApp) SetHashes(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("cannot open file path: %w", err)
	}
	defer file.Close()

	h256 := sha256.New()
	if _, err = io.Copy(h256, file); err != nil {
		return fmt.Errorf("cannot compute sha256 hash: %w", err)
	}
	app.Hashes.Sha256 = fmt.Sprintf("%x", h256.Sum(nil))
	if _, err = file.Seek(0, io.SeekStart); err != nil {
		return err
	}

	h1 := sha1.New()
	if _, err = io.Copy(h1, file); err != nil {
		return fmt.Errorf("cannot compute sha1 hash: %w", err)
	}
	app.Hashes.Sha1 = fmt.Sprintf("%x", h1.Sum(nil))
	if _, err = file.Seek(0, io.SeekStart); err != nil {
		return err
	}

	hmd5 := md5.New()
	if _, err = io.Copy(hmd5, file); err != nil {
		return fmt.Errorf("cannot compute md5 hash: %w", err)
	}
	app.Hashes.Md5 = fmt.Sprintf("%x", hmd5.Sum(nil))

	return nil
}
