package analyzer

import (
	"crypto/x509/pkix"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/andpalmier/apkingo/v3/internal/utils"
	"github.com/avast/apkverifier"
)

// CertificateInfo represents certificate information
type CertificateInfo struct {
	Serial            string   `json:"serial"`
	Thumbprint        string   `json:"thumbprint"`
	Sha256            string   `json:"sha256"`
	ValidFrom         string   `json:"valid-from"`
	ValidTo           string   `json:"valid-to"`
	Verified          bool     `json:"verified"`
	SigningScheme     string   `json:"signing-scheme"`
	SignerCount       int      `json:"signer-count"`
	VerificationError string   `json:"verification-error,omitempty"`
	Expired           bool     `json:"expired"`
	Debug             bool     `json:"debug"`
	Subject           CertName `json:"subject"`
	Issuer            CertName `json:"issuer"`
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
	res, verifyErr := apkverifier.Verify(filepath, nil)
	if len(res.SignerCerts) == 0 {
		certs, extractErr := apkverifier.ExtractCerts(filepath, nil)
		if extractErr != nil {
			if verifyErr != nil {
				return fmt.Errorf("signature verification failed: %w (certificate extraction failed: %v)", verifyErr, extractErr)
			}
			return extractErr
		}
		res.SignerCerts = certs
	}

	cert, certx := apkverifier.PickBestApkCert(res.SignerCerts)
	if cert == nil {
		if verifyErr != nil {
			return verifyErr
		}
		return errors.New("no certificate found")
	}

	app.Certificate.Serial = hex.EncodeToString(cert.SerialNumber.Bytes())
	app.Certificate.Thumbprint = cert.Sha1
	app.Certificate.Sha256 = cert.Sha256
	app.Certificate.ValidFrom = cert.ValidFrom.Format(time.DateTime)
	app.Certificate.ValidTo = cert.ValidTo.Format(time.DateTime)
	app.Certificate.Verified = verifyErr == nil
	app.Certificate.SigningScheme = signingSchemeName(res.SigningSchemeId)
	app.Certificate.SignerCount = len(res.SignerCerts)
	app.Certificate.Expired = time.Now().After(cert.ValidTo)
	app.Certificate.Debug = strings.Contains(cert.Subject, "CN=Android Debug")
	if verifyErr != nil {
		app.Certificate.VerificationError = utils.SanitizeText(verifyErr.Error())
	}

	if certx != nil {
		app.Certificate.Subject.setCertName(certx.Subject)
		app.Certificate.Issuer.setCertName(certx.Issuer)
	}
	app.Certificate.Subject.Raw = cert.Subject
	app.Certificate.Issuer.Raw = cert.Issuer

	return nil
}

func signingSchemeName(id int) string {
	switch id {
	case 1:
		return "v1"
	case 2:
		return "v2"
	case 3:
		return "v3"
	case 31:
		return "v3.1"
	default:
		return "unknown"
	}
}
