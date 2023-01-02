package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"hash"
	"os"
)

// getFileHash(h, filepath) - hash the file in the given path with the selected hash
func getFileHash(h hash.Hash, filepath string) ([]byte, error) {
	file, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	h.Write(file)
	return h.Sum(nil), nil
}

// setHashValues(filepath) - calculate hashes of the file and store them as a string
// in the androidapp struct
func (androidapp *AndroidApp) setHashValues(path string) error {
	h256 := sha256.New()
	digestsha256, err := getFileHash(h256, path)
	if err != nil {
		return err
	}
	androidapp.Hashes.Sha256 = fmt.Sprintf("%x", digestsha256)

	h1 := sha1.New()
	digestsha1, err := getFileHash(h1, path)
	if err != nil {
		return err
	}
	androidapp.Hashes.Sha1 = fmt.Sprintf("%x", digestsha1)

	hmd5 := md5.New()
	digestmd5, err := getFileHash(hmd5, path)
	if err != nil {
		return err
	}
	androidapp.Hashes.Md5 = fmt.Sprintf("%x", digestmd5)

	return nil
}
