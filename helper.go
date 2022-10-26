package main

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var regexes = []*regexp.Regexp{
	regexp.MustCompile(`(?m)^flipper-z-(f[0-9]*|any)-(update|updater|bootloader|firmware|full|core2_firmware|scripts|resources|sdk)-.*\.(dfu|bin|elf|hex|tgz|json|zip)$`),
	regexp.MustCompile(`(?m)^(f[0-9]*)_(bootloader|firmware|full)\.(dfu|bin|elf|hex)$`),
}

func isExistingDir(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fi.IsDir()
}

func removeWithParents(path string) error {
	err := os.RemoveAll(path)
	if err != nil {
		return err
	}
	parent := filepath.Dir(path)
	empty, err := isDirEmpty(parent)
	if err != nil {
		return err
	}
	if empty {
		return removeWithParents(parent)
	}
	return nil
}

func isDirEmpty(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err
}

func arrayContains(stack []string, needle string) bool {
	for _, e := range stack {
		if e == needle {
			return true
		}
	}
	return false
}

func parseFilename(name string) *file {
	// TODO refactor this hardcoded crap
	if strings.HasPrefix(name, "qFlipper") {
		switch filepath.Ext(name) {
		case ".dmg":
			return &file{
				Target: "macos/amd64",
				Type:   "dmg",
			}
		case ".AppImage":
			return &file{
				Target: "linux/amd64",
				Type:   "AppImage",
			}
		case ".zip":
			return &file{
				Target: "windows/amd64",
				Type:   "portable",
			}
		case ".exe":
			return &file{
				Target: "windows/amd64",
				Type:   "installer",
			}
		}
	}

	for _, re := range regexes {
		m := re.FindAllStringSubmatch(name, -1)
		if len(m) != 1 || len(m[0]) != 4 {
			continue
		}
		return &file{
			Type:   m[0][2] + "_" + m[0][3],
			Target: m[0][1],
		}
	}
	return nil
}

func calculateSha256(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}
