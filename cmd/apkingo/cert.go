package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/avast/apkverifier"
)

// certInfo(path) - retrieve certificate information from apk
func certInfo(filepath string) error {
	res, err := apkverifier.ExtractCerts(filepath, nil)
	if err != nil {
		return err
	}

	// this may print an error, but certificate info are still retrieved
	cert, _ := apkverifier.PickBestApkCert(res)
	if cert == nil {
		return errors.New("no certificate found")
	}

	fmt.Printf("Serial:\t\t")
	printer(cert.SerialNumber.String())
	fmt.Printf("sha1:\t\t")
	printer(cert.Sha1)
	fmt.Printf("Issuer:\t\t")
	printer(cert.Issuer)
	fmt.Printf("Subject:\t")
	printer(cert.Subject)
	fmt.Printf("Valid from:\t")
	printer(cert.ValidFrom.Format(time.RFC822))
	fmt.Printf("Valid to:\t")
	printer(cert.ValidTo.Format(time.RFC822))

	return nil
}
