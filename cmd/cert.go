package main

import (
	"github.com/avast/apkverifier"
)

// getCertInfo() - retrieve certificate information from APK
func (androidapp *AndroidApp) getCertInfo() error {
	res, err := apkverifier.ExtractCerts(androidapp.Path, nil)
	if err != nil {
		return err
	}

	// this may print an error, but the certificate info are still retrieved
	cert, _ := apkverifier.PickBestApkCert(res)
	if cert == nil {
		return err
	} else {
		androidapp.Cert = cert.String()
		return nil
	}
}
