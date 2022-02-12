package main

import (
	"errors"
	"github.com/avast/apkverifier"
)

// SetCertInfo(cert) - store certificate info in the androidapp struct
func (androidapp *AndroidApp) setCertInfo(cert apkverifier.CertInfo) {
	androidapp.Certificate.Serial = cert.SerialNumber.String()
	androidapp.Certificate.Issuer = cert.Issuer
	androidapp.Certificate.Subject = cert.Subject
	androidapp.Certificate.Sha1 = cert.Sha1
	androidapp.Certificate.ValidFrom = cert.ValidFrom.Format("Jan 2, 2006")
	androidapp.Certificate.ValidTo = cert.ValidTo.Format("Jan 2, 2006")
}

// getCertInfo(path) - retrieve certificate information from apk
func (androidapp *AndroidApp) getCertInfo(filepath string) (*apkverifier.CertInfo, error) {
	res, err := apkverifier.ExtractCerts(filepath, nil)
	if err != nil {
		return &apkverifier.CertInfo{}, err
	}

	// this may print an error, but certificate info are still retrieved!!
	cert, _ := apkverifier.PickBestApkCert(res)
	if cert == nil {
		return &apkverifier.CertInfo{}, errors.New("no certificate found")
	} else {
		return cert, nil
	}
}
