package constants

import "testing"

func TestPermissionClass(t *testing.T) {
	tests := []struct {
		name       string
		permission string
		want       string
	}{
		{"dangerous", "android.permission.CAMERA", PermissionDangerous},
		{"system", "android.permission.INTERNET", PermissionSystem},
		{"custom", "com.example.permission.INTERNAL", PermissionCustom},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PermissionClass(tt.permission); got != tt.want {
				t.Fatalf("PermissionClass(%q) = %q, want %q", tt.permission, got, tt.want)
			}
		})
	}
}
