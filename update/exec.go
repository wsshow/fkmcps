package update

import (
	"bufio"
	"fkmcps/version"
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"unicode"

	"github.com/Masterminds/semver/v3"
)

func SelfUpdate(owner string, repo string, proxyURL string) (err error) {
	up := NewUpdater().WithProxy(proxyURL)
	info := version.Get()
	// Check for updates by comparing the current version with the latest release on GitHub.
	latest, yes, err := up.CheckForUpdates(semver.MustParse(info.Version), owner, repo)
	if err != nil {
		return err
	}
	if !yes {
		fmt.Printf("Current version is the latest: %s\n", info.Version)
		return nil
	}

	fmt.Printf("New version found: %s, downloading update...\n", latest.TagName)

	// Apply update
	if err = up.Apply(latest, findAsset, findChecksum); err != nil {
		return err
	}
	fmt.Printf("Update successful, current version: %s\n", latest.TagName)
	return nil
}

func findAsset(items []Asset) (idx int) {
	ext := "zip"
	suffix := fmt.Sprintf("%s_%s.%s", CapitalizeOS(), GetNormalizedArch(), ext)
	for i := range items {
		if strings.HasSuffix(items[i].BrowserDownloadURL, suffix) {
			return i
		}
	}
	return -1
}

func findChecksum(items []Asset) (algo Algorithm, expectedChecksum string, err error) {
	ext := "zip"
	suffix := fmt.Sprintf("%s_%s.%s", CapitalizeOS(), GetNormalizedArch(), ext)
	var checksumFileURL string
	for i := range items {
		if items[i].Name == "checksums.txt" {
			checksumFileURL = items[i].BrowserDownloadURL
			break
		}
	}
	if checksumFileURL == "" {
		return SHA256, "", ErrChecksumFileNotFound
	}

	resp, err := http.Get(checksumFileURL)
	if err != nil {
		return SHA256, "", err
	}
	defer resp.Body.Close()

	if !IsHttpSuccess(resp.StatusCode) {
		return "", "", fmt.Errorf("URL %q is unreachable", checksumFileURL)
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasSuffix(line, suffix) {
			continue
		}
		return SHA256, strings.Fields(line)[0], nil
	}
	if err = scanner.Err(); err != nil {
		return SHA256, "", err
	}
	return SHA256, "", ErrChecksumFileNotFound
}

// CapitalizeOS returns the capitalized OS name (e.g., "Windows", "Linux", "Darwin").
func CapitalizeOS() string {
	osName := runtime.GOOS
	if len(osName) == 0 {
		return ""
	}

	runes := []rune(osName)
	runes[0] = unicode.ToUpper(runes[0])

	return string(runes)
}

// GetNormalizedArch returns the normalized architecture name (e.g., "x86_64", "i386").
func GetNormalizedArch() string {
	arch := runtime.GOARCH

	switch arch {
	case "amd64":
		return "x86_64"
	case "386":
		return "i386"
	default:
		return arch
	}
}
