package constants

import "strings"

const (
	// PermissionDangerous is a platform permission with a dangerous protection level.
	PermissionDangerous = "dangerous"
	// PermissionSystem is a platform permission not in the generated dangerous set.
	PermissionSystem = "system"
	// PermissionCustom is an app-defined permission.
	PermissionCustom = "custom"
)

// Generated from AOSP android-17.0.0_r1:
// https://android.googlesource.com/platform/frameworks/base/+/refs/tags/android-17.0.0_r1/core/res/AndroidManifest.xml
// Regenerate this list when the pinned Android version changes.
const dangerousPermissionList = `
android.permission.ACCEPT_HANDOVER android.permission.ACCESS_BACKGROUND_LOCATION
android.permission.ACCESS_COARSE_LOCATION android.permission.ACCESS_FINE_LOCATION
android.permission.ACCESS_LOCAL_NETWORK android.permission.ACCESS_MEDIA_LOCATION
android.permission.ACTIVITY_RECOGNITION android.permission.ANSWER_PHONE_CALLS
android.permission.BLUETOOTH_ADVERTISE android.permission.BLUETOOTH_CONNECT
android.permission.BLUETOOTH_SCAN android.permission.BODY_SENSORS
android.permission.BODY_SENSORS_BACKGROUND android.permission.CALL_PHONE android.permission.CAMERA
android.permission.EYE_TRACKING_COARSE android.permission.EYE_TRACKING_FINE
android.permission.FACE_TRACKING android.permission.GET_ACCOUNTS android.permission.HAND_TRACKING
android.permission.HEAD_TRACKING android.permission.NEARBY_WIFI_DEVICES
android.permission.POST_NOTIFICATIONS android.permission.PROCESS_OUTGOING_CALLS
android.permission.RANGING android.permission.READ_CALENDAR android.permission.READ_CALL_LOG
android.permission.READ_CELL_BROADCASTS android.permission.READ_CONTACTS
android.permission.READ_EXTERNAL_STORAGE android.permission.READ_MEDIA_AUDIO
android.permission.READ_MEDIA_IMAGES android.permission.READ_MEDIA_VIDEO
android.permission.READ_MEDIA_VISUAL_USER_SELECTED android.permission.READ_PHONE_NUMBERS
android.permission.READ_PHONE_STATE android.permission.READ_SMS android.permission.RECEIVE_MMS
android.permission.RECEIVE_SMS android.permission.RECEIVE_WAP_PUSH android.permission.RECORD_AUDIO
android.permission.SCENE_UNDERSTANDING_COARSE android.permission.SCENE_UNDERSTANDING_FINE
android.permission.SEND_SMS android.permission.USE_SIP android.permission.UWB_RANGING
android.permission.WRITE_CALENDAR android.permission.WRITE_CALL_LOG
android.permission.WRITE_CONTACTS android.permission.WRITE_EXTERNAL_STORAGE
com.android.voicemail.permission.ADD_VOICEMAIL
`

var dangerousPermissions = func() map[string]struct{} {
	names := strings.Fields(dangerousPermissionList)
	set := make(map[string]struct{}, len(names))
	for _, name := range names {
		set[name] = struct{}{}
	}
	return set
}()

// PermissionClass classifies a permission using the generated platform list.
func PermissionClass(name string) string {
	if _, ok := dangerousPermissions[name]; ok {
		return PermissionDangerous
	}
	if strings.HasPrefix(name, "android.permission.") {
		return PermissionSystem
	}
	return PermissionCustom
}
