package main

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var regexes = []*regexp.Regexp{
	regexp.MustCompile(`(?m)^flipper-z-(f[0-9]*)-(bootloader|firmware|full)-.*\.(dfu|bin|elf|hex)$`),
	regexp.MustCompile(`(?m)^(f[0-9]*)_(bootloader|firmware|full)\.(dfu|bin|elf|hex)$`),
}

func isExistingDir(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fi.IsDir()
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
