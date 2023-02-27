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

// hashInfo(filepath) - calculate hashes of the apk file
func hashInfo(path string) (string, error) {
	h256 := sha256.New()
	digestsha256, err := getFileHash(h256, path)
	if err != nil {
		return "", err
	}
	h1 := sha1.New()
	digestsha1, err := getFileHash(h1, path)
	if err != nil {
		return "", err
	}
	hmd5 := md5.New()
	digestmd5, err := getFileHash(hmd5, path)
	if err != nil {
		return "", err
	}

	fmt.Printf("md5:\t\t")
	cyan.Printf("%x\n", digestmd5)
	fmt.Printf("sha1:\t\t")
	cyan.Printf("%x\n", digestsha1)
	fmt.Printf("sha256:\t\t")
	cyan.Printf("%x\n", digestsha256)

	return fmt.Sprintf("%x", digestsha256), nil
}
