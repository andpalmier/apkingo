package analyzer

import (
	"github.com/andpalmier/apkingo/internal/utils"
	"github.com/shogo82148/androidbinary/apk"
)

// SetGeneralInfo sets general information about the APK
func (app *AndroidApp) SetGeneralInfo(apk *apk.Apk) error {
	var err error

	app.Name, err = apk.Label(nil)
	if err != nil {
		return err
	}

	app.PackageName = apk.PackageName()

	app.Version, err = apk.Manifest().VersionName.String()
	utils.LogError("error getting version information", err)

	app.MainActivity, err = apk.MainActivity()
	utils.LogError("error getting main activity information", err)

	app.MinimumSDK, err = apk.Manifest().SDK.Min.Int32()
	utils.LogError("error getting minimum SDK information", err)

	app.TargetSDK, err = apk.Manifest().SDK.Target.Int32()
	utils.LogError("error getting target SDK information", err)

	for _, n := range apk.Manifest().UsesPermissions {
		permission, _ := n.Name.String()
		if permission != "" {
			app.Permissions = append(app.Permissions, permission)
		}
	}

	for _, n := range apk.Manifest().App.MetaData {
		metadataName, _ := n.Name.String()
		metadataValue, _ := n.Value.String()
		if metadataName != "" {
			app.Metadata = append(app.Metadata, Metadata{
				Name:  metadataName,
				Value: metadataValue,
			})
		}
	}

	return nil
}
