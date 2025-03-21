package util

import (
	"fmt"
	"runtime"
	"testing"
)

func TestGetVersion(t *testing.T) {
	// Exécuter la fonction GetVersion
	versionInfo := GetVersion()

	// Vérifier que les champs sont correctement remplis
	t.Run("Check Version", func(t *testing.T) {
		if versionInfo.Version != version {
			t.Errorf("expected version %s, but got %s", version, versionInfo.Version)
		}
	})

	t.Run("Check GitCommit", func(t *testing.T) {
		if versionInfo.GitCommit != gitCommit {
			t.Errorf("expected git commit %s, but got %s", gitCommit, versionInfo.GitCommit)
		}
	})

	t.Run("Check BuildDate", func(t *testing.T) {
		if versionInfo.BuildDate != buildDate {
			t.Errorf("expected build date %s, but got %s", buildDate, versionInfo.BuildDate)
		}
	})

	t.Run("Check GoVersion", func(t *testing.T) {
		if versionInfo.GoVersion != runtime.Version() {
			t.Errorf("expected Go version %s, but got %s", runtime.Version(), versionInfo.GoVersion)
		}
	})

	t.Run("Check Compiler", func(t *testing.T) {
		if versionInfo.Compiler != runtime.Compiler {
			t.Errorf("expected compiler %s, but got %s", runtime.Compiler, versionInfo.Compiler)
		}
	})

	t.Run("Check Platform", func(t *testing.T) {
		expectedPlatform := fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
		if versionInfo.Platform != expectedPlatform {
			t.Errorf("expected platform %s, but got %s", expectedPlatform, versionInfo.Platform)
		}
	})
}
