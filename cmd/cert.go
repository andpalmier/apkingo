package main

import (
	"errors"

	"github.com/avast/apkverifier"
)

// getCertInfo() - retrieve certificate information from APK
func (androidapp *AndroidApp) getCertInfo() error {
	res, err := apkverifier.ExtractCerts(androidapp.path, nil)
	if err != nil {
		return err
	}

	// this may print an error, but the certificate info are still retrieved
	cert, _ := apkverifier.PickBestApkCert(res)
	if cert == nil {
		return errors.New("no certificate found")
	}
	androidapp.cert = *cert
	return nil
}
