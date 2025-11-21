package analyzer

import (
	"crypto/x509/pkix"
	"encoding/hex"
	"errors"
	"time"

	"github.com/avast/apkverifier"
)

// CertificateInfo represents certificate information
type CertificateInfo struct {
	Serial     string   `json:"serial"`
	Thumbprint string   `json:"thumbprint"`
	ValidFrom  string   `json:"valid-from"`
	ValidTo    string   `json:"valid-to"`
	Subject    CertName `json:"subject"`
	Issuer     CertName `json:"issuer"`
}

// CertName represents issuer and subject details
type CertName struct {
	Country      string `json:"country"`
	Organization string `json:"organization"`
	OrgUnit      string `json:"organizational-unit"`
	Locality     string `json:"locality"`
	Province     string `json:"province"`
	CommonName   string `json:"common-name"`
	Raw          string `json:"raw"`
}

// setCertName sets certificate name details
func (cn *CertName) setCertName(name pkix.Name) {
	cn.Country = firstElement(name.Country)
	cn.Organization = firstElement(name.Organization)
	cn.OrgUnit = firstElement(name.OrganizationalUnit)
	cn.Locality = firstElement(name.Locality)
	cn.Province = firstElement(name.Province)
	cn.CommonName = name.CommonName
}

// firstElement returns the first element of a string slice, or an empty string if the slice is empty
func firstElement(slice []string) string {
	if len(slice) > 0 {
		return slice[0]
	}
	return ""
}

// SetCertInfo retrieves and sets certificate information
func (app *AndroidApp) SetCertInfo(filepath string) error {
	res, err := apkverifier.ExtractCerts(filepath, nil)
	if err != nil {
		return err
	}

	cert, certx := apkverifier.PickBestApkCert(res)
	if cert == nil {
		return errors.New("no certificate found")
	}

	app.Certificate.Serial = hex.EncodeToString(cert.SerialNumber.Bytes())
	app.Certificate.Thumbprint = cert.Sha1
	app.Certificate.ValidFrom = cert.ValidFrom.Format(time.DateTime)
	app.Certificate.ValidTo = cert.ValidTo.Format(time.DateTime)

	if certx != nil {
		app.Certificate.Subject.setCertName(certx.Subject)
		app.Certificate.Issuer.setCertName(certx.Issuer)
	}
	app.Certificate.Subject.Raw = cert.Subject
	app.Certificate.Issuer.Raw = cert.Issuer

	return nil
}
