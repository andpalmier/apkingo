// Package report provides formatted output for APK analysis results.
package report

import (
	"fmt"
	"strings"

	"github.com/andpalmier/apkingo/internal/analyzer"
	"github.com/andpalmier/apkingo/internal/ui"
	"github.com/andpalmier/apkingo/internal/utils"
	"github.com/fatih/color"
)

type Reporter struct {
	printer *ui.Printer
}

func NewReporter(printer *ui.Printer) *Reporter {
	return &Reporter{printer: printer}
}

func (r *Reporter) PrintBanner() {
	r.printer.PrintBanner()
}

func (r *Reporter) PrintGeneralInfo(app *analyzer.AndroidApp) {
	r.printer.PrintSectionHeader("General Info")
	if app.Errors.General != nil {
		r.printer.PrintItalic("General information not available")
		return
	}

	minSDK := "not found"
	if app.MinimumSDK != 0 {
		minSDK = fmt.Sprintf("%d (%s)", app.MinimumSDK, utils.AndroidName[int(app.MinimumSDK)])
	}

	targetSDK := "not found"
	if app.TargetSDK != 0 {
		targetSDK = fmt.Sprintf("%d (%s)", app.TargetSDK, utils.AndroidName[int(app.TargetSDK)])
	}

	r.printer.PrintKV("Name", app.Name)
	r.printer.PrintKV("Package Name", app.PackageName)
	r.printer.PrintKV("Version", app.Version)
	r.printer.PrintKV("Main Activity", app.MainActivity)
	r.printer.PrintKV("Minimum SDK", minSDK)
	r.printer.PrintKV("Target SDK", targetSDK)
}

func (r *Reporter) PrintHash(hashes analyzer.Hashes) {
	r.printer.PrintSectionHeader("Hash Values")
	r.printer.PrintKV("MD5", hashes.Md5)
	r.printer.PrintKV("SHA1", hashes.Sha1)
	r.printer.PrintKV("SHA256", hashes.Sha256)
}

func (r *Reporter) PrintPlayStoreInfo(app *analyzer.AndroidApp) {
	r.printer.PrintSectionHeader("Play Store")
	if app.Errors.PlayStore != nil || app.PlayStore == nil {
		r.printer.PrintItalic("App not found in Play Store")
		return
	}
	if app.PlayStore != nil {
		r.printer.PrintKV("URL", app.PlayStore.Url)
		r.printer.PrintKV("Version", app.PlayStore.Version)
		r.printer.PrintKV("Released", app.PlayStore.Release)
		r.printer.PrintKV("Updated", app.PlayStore.Updated.Format("Jan 2, 2006"))
		r.printer.PrintKV("Genre", app.PlayStore.Genre)
		r.printer.PrintKV("Summary", app.PlayStore.Summary)
		r.printer.PrintKV("Installs", app.PlayStore.Installs)
		r.printer.PrintKV("Score", fmt.Sprintf("%v", app.PlayStore.Score))
		r.printer.PrintKV("Developer ID", app.PlayStore.Developer.Id)
		r.printer.PrintKV("Developer Name", app.PlayStore.Developer.Name)
		r.printer.PrintKV("Developer Email", app.PlayStore.Developer.Mail)
		r.printer.PrintKV("Developer URL", app.PlayStore.Developer.URL)
	}
}

func (r *Reporter) PrintCertInfo(app *analyzer.AndroidApp) {
	r.printer.PrintSectionHeader("Certificate")
	if app.Errors.Cert != nil {
		r.printer.PrintItalic("Certificate information not available")
		return
	}
	certinfo := app.Certificate

	// Format Issuer and Subject manually
	issuer := fmt.Sprintf("C=%s, O=%s, OU=%s, CN=%s",
		certinfo.Issuer.Country, certinfo.Issuer.Organization, certinfo.Issuer.OrgUnit, certinfo.Issuer.CommonName)
	subject := fmt.Sprintf("C=%s, O=%s, OU=%s, CN=%s",
		certinfo.Subject.Country, certinfo.Subject.Organization, certinfo.Subject.OrgUnit, certinfo.Subject.CommonName)

	r.printer.PrintKV("Serial", certinfo.Serial)
	r.printer.PrintKV("Thumbprint", certinfo.Thumbprint)
	r.printer.PrintKV("Valid From", certinfo.ValidFrom)
	r.printer.PrintKV("Valid To", certinfo.ValidTo)
	r.printer.PrintKV("Issuer", issuer)
	r.printer.PrintKV("Subject", subject)
}

func (r *Reporter) PrintKoodousInfo(app *analyzer.AndroidApp) {
	r.printer.PrintSectionHeader("KOODOUS")
	if app.Errors.Koodous != nil || app.Koodous == nil {
		r.printer.PrintItalic("Information not available from Koodous")
		return
	}
	kinfo := app.Koodous

	r.printer.PrintKV("URL", kinfo.Url)
	r.printer.PrintKV("ID", kinfo.Id)
	r.printer.PrintKV("Submission Date", kinfo.SubmissionDate)
	r.printer.PrintKV("Icon Link", kinfo.IconLink)
	r.printer.PrintKV("Size", fmt.Sprintf("%d bytes", kinfo.Size))
	r.printer.PrintKV("Tags", strings.Join(kinfo.Tags, ", "))

	detected := fmt.Sprintf("%t", kinfo.Detected)
	if kinfo.Detected {
		r.printer.PrintKVRedBold("Detected", detected)
	} else {
		r.printer.PrintKV("Detected", detected)
	}

	rating := fmt.Sprintf("%d", kinfo.Rating)
	if kinfo.Rating < 0 {
		r.printer.PrintKVRedBold("Rating", rating)
	} else {
		r.printer.PrintKV("Rating", rating)
	}

	if kinfo.Corrupted {
		r.printer.PrintKVRedBold("Corrupted", fmt.Sprintf("%t", kinfo.Corrupted))
	} else {
		r.printer.PrintKV("Corrupted", fmt.Sprintf("%t", kinfo.Corrupted))
	}

	if kinfo.Trusted {
		r.printer.GetCyan().Fprintf(r.printer.GetTabWriter(), "Trusted:\t%t\n", kinfo.Trusted)
	} else {
		r.printer.PrintKV("Trusted", fmt.Sprintf("%t", kinfo.Trusted))
	}
}

func (r *Reporter) PrintPermissions(permissions []string) {
	r.printer.PrintSectionHeader("Permissions")
	if len(permissions) == 0 {
		r.printer.PrintItalic("No permissions found")
	} else {
		r.printer.PrintList(permissions)
	}
}

func (r *Reporter) PrintMetadata(metadata map[string]string) {
	r.printer.PrintSectionHeader("Metadata")
	if len(metadata) == 0 {
		r.printer.PrintItalic("No metadata found")
		return
	}

	// Only show metadata entries that have actual values
	hasContent := false
	for key, value := range metadata {
		// Skip empty values
		if value != "" && value != "<nil>" && value != "[]" {
			r.printer.PrintKV(key, value)
			hasContent = true
		}
	}

	// If all values were empty, show a message
	if !hasContent {
		r.printer.PrintItalic("No metadata found")
	}
}

func (r *Reporter) PrintVTInfo(app *analyzer.AndroidApp) {
	r.printer.PrintSectionHeader("VirusTotal")
	if app.Errors.VT != nil || app.VirusTotal == nil {
		r.printer.PrintItalic("App not found in VirusTotal")
		return
	}
	vtinfo := app.VirusTotal
	if vtinfo != nil {
		r.printer.PrintKV("URL", vtinfo.Url)
		r.printer.PrintKV("Submission Date", vtinfo.SubmissDate.String())
		r.printer.PrintKV("Submissions", fmt.Sprintf("%d", vtinfo.TimesSubmit))
		r.printer.PrintKV("Last Analysis", vtinfo.LastAnalysis.String())

		// Tags
		if len(vtinfo.Tags) > 0 {
			r.printer.PrintKV("Tags", strings.Join(vtinfo.Tags, ", "))
		}

		// Popular Threat Classification
		if vtinfo.PopularThreatCategory != "" {
			r.printer.PrintKVRedBold("Threat Category", vtinfo.PopularThreatCategory)
		}
		if vtinfo.PopularThreatName != "" {
			r.printer.PrintKVRedBold("Threat Name", vtinfo.PopularThreatName)
		}

		// Reputation
		if vtinfo.Reputation != 0 {
			repText := fmt.Sprintf("%d", vtinfo.Reputation)
			if vtinfo.Reputation < 0 {
				r.printer.PrintKVRed("Reputation", repText)
			} else {
				r.printer.PrintKV("Reputation", repText)
			}
		}

		// Crowdsourced Detections
		if vtinfo.TotalCrowdsourcedSigma > 0 {
			r.printer.PrintKVRedBold("Crowdsourced Sigma", fmt.Sprintf("%d detections", vtinfo.TotalCrowdsourcedSigma))
		}
		if vtinfo.TotalCrowdsourcedYara > 0 {
			r.printer.PrintKVRedBold("Crowdsourced YARA", fmt.Sprintf("%d rules matched", vtinfo.TotalCrowdsourcedYara))
		}
		if vtinfo.TotalCrowdsourcedIDSHits > 0 {
			r.printer.PrintKVRedBold("Crowdsourced IDS", fmt.Sprintf("%d hits", vtinfo.TotalCrowdsourcedIDSHits))
		}

		// Analysis Stats
		if vtinfo.AnalysStats != nil {
			r.printer.PrintSectionHeader("VT Analysis Stats")
			r.printer.PrintKV("Harmless", fmt.Sprintf("%d", vtinfo.AnalysStats.Harmless))
			if vtinfo.AnalysStats.Malicious > 0 {
				r.printer.PrintKVRedBold("Malicious", fmt.Sprintf("%d", vtinfo.AnalysStats.Malicious))
			} else {
				r.printer.PrintKV("Malicious", fmt.Sprintf("%d", vtinfo.AnalysStats.Malicious))
			}
			r.printer.PrintKV("Suspicious", fmt.Sprintf("%d", vtinfo.AnalysStats.Suspicious))
			r.printer.PrintKV("Undetected", fmt.Sprintf("%d", vtinfo.AnalysStats.Undetected))
			r.printer.PrintKV("Failure", fmt.Sprintf("%d", vtinfo.AnalysStats.Failure))
		}

		// Votes
		if vtinfo.Votes != nil {
			r.printer.PrintSectionHeader("VT Votes")
			votes := vtinfo.Votes
			r.printer.PrintKV("Harmless", fmt.Sprintf("%d", votes.Harmless))

			if votes.Malicious > 0 {
				r.printer.PrintKVRedBold("Malicious", fmt.Sprintf("%d", votes.Malicious))
			} else {
				r.printer.PrintKV("Malicious", fmt.Sprintf("%d", votes.Malicious))
			}
		}

		// Icon
		if vtinfo.Icon != nil {
			r.printer.PrintSectionHeader("VT Icon")
			icon := vtinfo.Icon
			r.printer.PrintKV("MD5", icon.Md5)
			r.printer.PrintKV("DHash", icon.Dhash)
		}

		// Androguard
		if vtinfo.Androguard != nil {
			androguard := vtinfo.Androguard
			r.printer.PrintSectionHeader("VT Androguard Analysis")

			// Package and Version Info
			if androguard.Package != "" {
				r.printer.PrintKV("Package", androguard.Package)
			}
			if androguard.AndroidVersionCode != "" {
				r.printer.PrintKV("Version Code", androguard.AndroidVersionCode)
			}
			if androguard.AndroidVersionName != "" {
				r.printer.PrintKV("Version Name", androguard.AndroidVersionName)
			}
			if androguard.MinSdkVersion != "" {
				r.printer.PrintKV("Min SDK", androguard.MinSdkVersion)
			}
			if androguard.TargetSdkVersion != "" {
				r.printer.PrintKV("Target SDK", androguard.TargetSdkVersion)
			}
			if androguard.MainActivity != "" {
				r.printer.PrintKV("Main Activity", androguard.MainActivity)
			}

			// Components
			if len(androguard.Activities) > 0 {
				r.printer.Flush()
				fmt.Fprintln(r.printer.GetOut())
				r.printer.GetCyan().Fprintf(r.printer.GetOut(), "Activities (%d):\n", len(androguard.Activities))
				r.printer.PrintList(androguard.Activities)
			}

			if len(androguard.Services) > 0 {
				r.printer.Flush()
				fmt.Fprintln(r.printer.GetOut())
				r.printer.GetCyan().Fprintf(r.printer.GetOut(), "Services (%d):\n", len(androguard.Services))
				r.printer.PrintList(androguard.Services)
			}

			if len(androguard.Providers) > 0 {
				r.printer.Flush()
				fmt.Fprintln(r.printer.GetOut())
				r.printer.GetCyan().Fprintf(r.printer.GetOut(), "Providers (%d):\n", len(androguard.Providers))
				r.printer.PrintList(androguard.Providers)
			}

			if len(androguard.Receivers) > 0 {
				r.printer.Flush()
				fmt.Fprintln(r.printer.GetOut())
				r.printer.GetCyan().Fprintf(r.printer.GetOut(), "Receivers (%d):\n", len(androguard.Receivers))
				r.printer.PrintList(androguard.Receivers)
			}

			if len(androguard.Libraries) > 0 {
				r.printer.Flush()
				fmt.Fprintln(r.printer.GetOut())
				r.printer.GetCyan().Fprintf(r.printer.GetOut(), "Libraries (%d):\n", len(androguard.Libraries))
				r.printer.PrintList(androguard.Libraries)
			}

			// Dangerous Permissions - highlighted
			if len(androguard.DangerPerm) > 0 {
				r.printer.Flush()
				fmt.Fprintln(r.printer.GetOut())
				r.printer.GetRed().Add(color.Bold).Fprintf(r.printer.GetOut(), "Dangerous Permissions:  %d found\n", len(androguard.DangerPerm))
				for _, perm := range androguard.DangerPerm {
					fmt.Fprintf(r.printer.GetOut(), "  • %s\n", r.printer.GetRed().Sprint(perm))
				}
			}
		}

	}
}

// metadataToMap converts a slice of Metadata to a map[string]string
func metadataToMap(metadata []analyzer.Metadata) map[string]string {
	result := make(map[string]string)
	for _, m := range metadata {
		result[m.Name] = m.Value
	}
	return result
}

func (r *Reporter) PrintAll(app *analyzer.AndroidApp) {
	r.PrintGeneralInfo(app)
	r.PrintHash(app.Hashes)
	r.PrintPermissions(app.Permissions)
	r.PrintMetadata(metadataToMap(app.Metadata))
	r.PrintCertInfo(app)
	r.PrintPlayStoreInfo(app)
	r.PrintKoodousInfo(app)
	r.PrintVTInfo(app)
	r.printer.Flush()
}
