package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"hash"
	"io/ioutil"
)

// getFileHash(h,f) - hash the file in the given path with the selected hash
func getFileHash(h hash.Hash, filepath string) ([]byte, error) {
	s, err := ioutil.ReadFile(androidapp.Path)
	if err != nil {
		return nil, err
	}
	h.Write(s)
	return h.Sum(nil), nil
}

// getHashValues() - calculate hash of the files with sha256, md5 and sha1
func (androidapp *AndroidApp) getHashValues() error {
	h256 := sha256.New()
	digestsha256, err := getFileHash(h256, androidapp.Path)
	if err != nil {
		return err
	}
	androidapp.HashSHA256 = digestsha256

	h1 := sha1.New()
	digestsha1, err := getFileHash(h1, androidapp.Path)
	if err != nil {
		return err
	}
	androidapp.HashSHA1 = digestsha1

	hmd5 := md5.New()
	digestmd5, err := getFileHash(hmd5, androidapp.Path)
	if err != nil {
		return err
	}
	androidapp.HashMD5 = digestmd5
	return nil
}
