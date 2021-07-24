package main

import (
	"crypto/sha512"
	"encoding/hex"
	"io/ioutil"
	"os"
	"regexp"
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
	for _, re := range regexes {
		m := re.FindAllStringSubmatch(name, -1)
		if len(m) != 1 || len(m[0]) != 4 {
			continue
		}
		return &file{
			Type: m[0][2] + "_" + m[0][3],
			Target: m[0][1],
		}
	}
	return nil
}

func calculateSha512(path string) string {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return ""
	}
	hash := sha512.Sum512(f)
	return hex.EncodeToString(hash[:])
}