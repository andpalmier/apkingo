package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
)

// setHashes calculates hashes of the APK file
func (androidapp *AndroidApp) setHashes(path string) error {
	file, err := os.Open(path)
	if err != nil {
		var OpenFileErr = errors.New("cannot open file path")
		return OpenFileErr
	}
	defer file.Close()

	h256 := sha256.New()
	if _, err = io.Copy(h256, file); err != nil {
		var Sha256Err = errors.New("cannot compute sha256 hash")
		return Sha256Err
	}
	androidapp.Hashes.Sha256 = fmt.Sprintf("%x", h256.Sum(nil))
	if _, err = file.Seek(0, io.SeekStart); err != nil {
		return err
	}

	h1 := sha1.New()
	if _, err = io.Copy(h1, file); err != nil {
		var Sha1Err = errors.New("cannot compute sha1 hash")
		return Sha1Err
	}
	androidapp.Hashes.Sha1 = fmt.Sprintf("%x", h1.Sum(nil))
	if _, err = file.Seek(0, io.SeekStart); err != nil {
		return err
	}

	hmd5 := md5.New()
	if _, err = io.Copy(hmd5, file); err != nil {
		var Md5Err = errors.New("cannot compute md5 hash")
		return Md5Err
	}
	androidapp.Hashes.Md5 = fmt.Sprintf("%x", hmd5.Sum(nil))

	return nil
}
