// Package report provides formatted output for APK analysis results.
package report

import (
	"fmt"
	"sort"
	"strings"

	"github.com/andpalmier/apkingo/v3/internal/analyzer"
	"github.com/andpalmier/apkingo/v3/internal/constants"
	"github.com/andpalmier/apkingo/v3/internal/ui"
	"github.com/andpalmier/apkingo/v3/internal/utils"
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
		name := utils.AndroidName[int(app.MinimumSDK)]
		if name == "" {
			name = "unknown Android version"
		}
		minSDK = fmt.Sprintf("%d (%s)", app.MinimumSDK, name)
	}

	targetSDK := "not found"
	if app.TargetSDK != 0 {
		name := utils.AndroidName[int(app.TargetSDK)]
		if name == "" {
			name = "unknown Android version"
		}
		targetSDK = fmt.Sprintf("%d (%s)", app.TargetSDK, name)
	}

	r.printer.PrintKV("Name", app.Name)
	if app.NameLocale != "" {
		r.printer.PrintKV(fmt.Sprintf("Name (%s)", app.Locale), app.NameLocale)
	}
	r.printer.PrintKV("Package Name", app.PackageName)
	r.printer.PrintKV("Version Name", app.VersionName)
	if app.VersionCode != 0 {
		r.printer.PrintKV("Version Code", fmt.Sprintf("%d", app.VersionCode))
	}
	r.printer.PrintKV("Main Activity", app.MainActivity)
	r.printer.PrintKV("Minimum SDK", minSDK)
	r.printer.PrintKV("Target SDK", targetSDK)
	if len(app.Architectures) > 0 {
		r.printer.PrintKV("Architectures", strings.Join(app.Architectures, ", "))
	} else {
		r.printer.PrintKV("Architectures", "none (pure Java/Kotlin)")
	}
}

func (r *Reporter) PrintHash(hashes analyzer.Hashes) {
	r.printer.PrintSectionHeader("Hash Values")
	r.printer.PrintKV("MD5", hashes.Md5)
	r.printer.PrintKV("SHA1", hashes.Sha1)
	r.printer.PrintKV("SHA256", hashes.Sha256)
}

func (r *Reporter) PrintPlayStoreInfo(app *analyzer.AndroidApp) {
	r.printer.PrintSectionHeader("Play Store")
	if app.NoPlayStore {
		r.printer.PrintItalic("Analysis without Play Store")
		return
	}
	if app.Errors.PlayStore != nil || app.PlayStore == nil {
		r.printer.PrintItalic("App not found in Play Store")
		return
	}
	r.printer.PrintKV("URL", app.PlayStore.Url)
	r.printer.PrintKV("Version", app.PlayStore.Version)
	r.printer.PrintKV("Released", app.PlayStore.Release)
	updated := "not found"
	if !app.PlayStore.Updated.IsZero() {
		updated = app.PlayStore.Updated.Format("Jan 2, 2006")
	}
	r.printer.PrintKV("Updated", updated)
	r.printer.PrintKV("Genre", app.PlayStore.Genre)
	r.printer.PrintKV("Summary", app.PlayStore.Summary)
	r.printer.PrintKV("Installs", app.PlayStore.Installs)
	r.printer.PrintKV("Score", fmt.Sprintf("%v", app.PlayStore.Score))
	r.printer.PrintKV("Developer ID", app.PlayStore.Developer.Id)
	r.printer.PrintKV("Developer Name", app.PlayStore.Developer.Name)
	r.printer.PrintKV("Developer Email", app.PlayStore.Developer.Mail)
	r.printer.PrintKV("Developer URL", app.PlayStore.Developer.URL)
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
	r.printer.PrintKV("SHA256", certinfo.Sha256)
	r.printer.PrintKV("Valid From", certinfo.ValidFrom)
	r.printer.PrintKV("Valid To", certinfo.ValidTo)
	if certinfo.Verified {
		r.printer.PrintKV("Verified", "true")
	} else {
		r.printer.PrintKVRedBold("Verified", "false")
	}
	r.printer.PrintKV("Signing Scheme", certinfo.SigningScheme)
	r.printer.PrintKV("Signers", fmt.Sprintf("%d", certinfo.SignerCount))
	if certinfo.Expired {
		r.printer.PrintKVRedBold("Expired", "true")
	} else {
		r.printer.PrintKV("Expired", "false")
	}
	if certinfo.Debug {
		r.printer.PrintKVRedBold("Debug Certificate", "true")
	} else {
		r.printer.PrintKV("Debug Certificate", "false")
	}
	if certinfo.VerificationError != "" {
		r.printer.PrintKVRedBold("Verification Error", certinfo.VerificationError)
	}
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
		_, _ = r.printer.GetCyan().Fprintf(r.printer.GetTabWriter(), "Trusted:\t%t\n", kinfo.Trusted)
	} else {
		r.printer.PrintKV("Trusted", fmt.Sprintf("%t", kinfo.Trusted))
	}
}

func (r *Reporter) PrintPermissions(permissions []string, classes map[string]string) {
	r.printer.PrintSectionHeader("Permissions")
	if len(permissions) == 0 {
		r.printer.PrintItalic("No permissions found")
		return
	}
	for _, permission := range permissions {
		switch classes[permission] {
		case constants.PermissionDangerous:
			_, _ = r.printer.GetRedBold().Fprintf(r.printer.GetOut(), " - %s [dangerous]\n", utils.SanitizeText(permission))
		case constants.PermissionCustom:
			_, _ = r.printer.GetYellow().Fprintf(r.printer.GetOut(), " - %s [custom]\n", utils.SanitizeText(permission))
		default:
			_, _ = fmt.Fprintf(r.printer.GetOut(), " - %s\n", utils.SanitizeText(permission))
		}
	}
}

func (r *Reporter) PrintComponents(app *analyzer.AndroidApp) {
	r.printer.PrintSectionHeader("Components")
	hasComponents := false
	printComponents := func(label string, components []string) {
		if len(components) == 0 {
			return
		}
		hasComponents = true
		r.printer.Flush()
		_, _ = fmt.Fprintln(r.printer.GetOut())
		_, _ = r.printer.GetCyan().Fprintf(r.printer.GetOut(), "%s (%d):\n", label, len(components))
		r.printer.PrintList(components)
	}
	printComponents("Activities", app.Activities)
	printComponents("Services", app.Services)
	printComponents("Providers", app.Providers)
	printComponents("Receivers", app.Receivers)
	if !hasComponents {
		r.printer.PrintItalic("No components found")
	}
}

func (r *Reporter) PrintSecurityFlags(app *analyzer.AndroidApp) {
	r.printer.PrintSectionHeader("Security Flags")
	if app.Errors.Security != nil {
		r.printer.PrintItalic("Security flags not available")
		return
	}

	printBool := func(key string, value *bool) {
		if value == nil {
			r.printer.PrintKV(key, "not set")
			return
		}
		if *value {
			r.printer.PrintKVRedBold(key, "true")
		} else {
			r.printer.PrintKV(key, "false")
		}
	}

	printBool("Debuggable", app.Debuggable)
	printBool("Allow Backup", app.AllowBackup)
	printBool("Test Only", app.TestOnly)
	printBool("Uses Cleartext Traffic", app.UsesCleartextTraffic)
	if app.NetworkSecurityConfig == "" {
		r.printer.PrintKV("Network Security Config", "not set")
	} else {
		r.printer.PrintKV("Network Security Config", app.NetworkSecurityConfig)
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
	keys := make([]string, 0, len(metadata))
	for key := range metadata {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		value := metadata[key]
		// Skip empty values
		if value != "" && value != constants.NilValue && value != constants.EmptySlice {
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
	r.printer.PrintKV("URL", vtinfo.Url)
	submissionDate := "not found"
	if !vtinfo.SubmissDate.IsZero() {
		submissionDate = vtinfo.SubmissDate.String()
	}
	r.printer.PrintKV("Submission Date", submissionDate)
	r.printer.PrintKV("Submissions", fmt.Sprintf("%d", vtinfo.TimesSubmit))
	lastAnalysis := "not found"
	if !vtinfo.LastAnalysis.IsZero() {
		lastAnalysis = vtinfo.LastAnalysis.String()
	}
	r.printer.PrintKV("Last Analysis", lastAnalysis)

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

		if len(androguard.Libraries) > 0 {
			r.printer.Flush()
			_, _ = fmt.Fprintln(r.printer.GetOut())
			_, _ = r.printer.GetCyan().Fprintf(r.printer.GetOut(), "Libraries (%d):\n", len(androguard.Libraries))
			r.printer.PrintList(androguard.Libraries)
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
	r.PrintPermissions(app.Permissions, app.PermissionClasses)
	r.PrintComponents(app)
	r.PrintSecurityFlags(app)
	r.PrintMetadata(metadataToMap(app.Metadata))
	r.PrintCertInfo(app)
	r.PrintPlayStoreInfo(app)
	r.PrintKoodousInfo(app)
	r.PrintVTInfo(app)
	r.printer.Flush()
}

// PrintBatchSummary prints a summary of multiple APK analysis results.
func (r *Reporter) PrintBatchSummary(results map[string]*analyzer.AndroidApp, failed []string) {
	r.printer.PrintSectionHeader("Batch Analysis Summary")
	r.printer.Printf("Total APKs analyzed: %d, Failed: %d\n\n", len(results), len(failed))

	// Print header
	r.printer.Printf("%-40s %-10s %-15s %-10s\n", "Package", "Version", "VT Result", "Rating")
	r.printer.PrintText(strings.Repeat("-", 80))

	vtClean := 0
	vtMalicious := 0
	vtUnknown := 0

	keys := make([]string, 0, len(results))
	for key := range results {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		app := results[key]
		pkgName := app.PackageName
		if pkgName == "" {
			pkgName = "Unknown"
		}

		version := app.VersionName
		if version == "" {
			version = "N/A"
		}

		// Determine VT result
		vtResult := "N/A"
		if app.VirusTotal != nil && app.VirusTotal.AnalysStats != nil {
			malicious := app.VirusTotal.AnalysStats.Malicious
			total := malicious + app.VirusTotal.AnalysStats.Harmless +
				app.VirusTotal.AnalysStats.Suspicious + app.VirusTotal.AnalysStats.Undetected
			if total > 0 {
				vtResult = fmt.Sprintf("%d/%d", malicious, total)
				if malicious > 0 {
					vtMalicious++
				} else {
					vtClean++
				}
			}
		} else {
			vtUnknown++
		}

		// Determine Koodous result
		koodousResult := "N/A"
		if app.Koodous != nil {
			if app.Koodous.Detected {
				koodousResult = "Detected"
			} else {
				koodousResult = "Clean"
			}
		}

		r.printer.Printf("%-40s %-10s %-15s %-10s\n",
			truncateString(pkgName, 38),
			truncateString(version, 8),
			truncateString(vtResult, 13),
			truncateString(koodousResult, 8),
		)
	}

	r.printer.Flush()

	// Print summary statistics
	fmt.Println()
	r.printer.PrintSectionHeader("Summary Statistics")
	r.printer.Printf("VirusTotal: %d clean, %d malicious, %d unknown\n", vtClean, vtMalicious, vtUnknown)
	r.printer.Printf("Total APKs processed: %d, Failed: %d\n", len(results), len(failed))
}

// truncateString truncates a string to the specified length.
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
