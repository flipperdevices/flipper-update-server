package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/caarlos0/env/v6"
	"github.com/coreos/go-semver/semver"
	"github.com/flipper-zero/flipper-update-server/github"
	"github.com/gin-gonic/gin"
)

var cfg config
var gh *github.Github
var latestDirectory directory

func main() {
	if err := env.Parse(&cfg); err != nil {
		log.Fatalln("Config", err)
	}

	if !isExistingDir(cfg.ArtifactsPath) {
		log.Fatalln(cfg.ArtifactsPath, "is not an existing directory")
	}

	var err error

	gh, err = github.New(cfg.Github)
	if err != nil {
		log.Fatalln("GitHub", err)
	}

	err = regenDirectory()
	if err != nil {
		log.Fatalln("Regen", err)
	}

	log.Println("Server started")

	r := gin.New()

	r.GET("/directory.json", serveDirectory)
	r.GET("/:channel/:target/:type", serveLatest)
	r.POST("/reindex", handleReindex)

	log.Fatal(r.Run(":8080"))
}

func regenDirectory() error {
	err := gh.Sync()
	if err != nil {
		return err
	}

	devChannel := channel{
		ID:          "development",
		Title:       "Development Channel",
		Description: "Latest builds, not yet tested by Flipper QA, be careful",
	}
	rcChannel := channel{
		ID:          "release-candidate",
		Title:       "Release Candidate Channel",
		Description: "This is going to be released soon, undergoing QA tests now",
	}
	releaseChannel := channel{
		ID:          "release",
		Title:       "Stable Release Channel",
		Description: "Stable releases, tested by Flipper QA",
	}

	dirs := make(map[string]struct{})
	err = filepath.Walk(cfg.ArtifactsPath, func(path string, c os.FileInfo, err error) error {
		if !c.IsDir() {
			return nil
		}
		delete(dirs, filepath.Dir(path))
		dirs[path] = struct{}{}
		return nil
	})
	if err != nil {
		return err
	}

	for path := range dirs {
		name := strings.TrimPrefix(path, strings.TrimPrefix(cfg.ArtifactsPath, "./"))
		name = strings.TrimLeft(name, "/")
		if arrayContains(cfg.Excluded, name) {
			continue
		}

		ver, isBranch := gh.Lookup(name)
		if isBranch {
			continue
		}
		if ver == nil {
			log.Println("Deleting", name)
			err = removeWithParents(filepath.Join(cfg.ArtifactsPath, name))
			if err != nil {
				log.Println("Can't delete", name, err)
			}
			continue
		}

		v := version{
			Version:   ver.Version,
			Changelog: ver.Changelog,
			Timestamp: Time(ver.Date),
			Files:     scanFiles(name),
		}

		if name == cfg.Github.DevelopmentBranch {
			devChannel.Versions = append(devChannel.Versions, v)
			continue
		}
		if ver.Rc {
			rcChannel.Versions = append(rcChannel.Versions, v)
		} else {
			releaseChannel.Versions = append(releaseChannel.Versions, v)
		}
	}

	latestDirectory = directory{
		Channels: []channel{devChannel, rcChannel, releaseChannel},
	}
	for k := range latestDirectory.Channels {
		c := &latestDirectory.Channels[k]
		sort.Slice(c.Versions, func(i, j int) bool {
			v1, err := semver.NewVersion(c.Versions[i].Version)
			v2, err := semver.NewVersion(c.Versions[j].Version)
			if err != nil {
				return c.Versions[i].Timestamp.Time().Before(c.Versions[j].Timestamp.Time())
			}
			return !v1.LessThan(*v2)
		})
		if len(c.Versions) > 1 {
			c.Versions = c.Versions[:1]
		}
	}

	return nil
}

func scanFiles(folder string) (files []file) {
	content, err := ioutil.ReadDir(filepath.Join(cfg.ArtifactsPath, folder))
	if err != nil {
		return
	}
	for _, c := range content {
		if c.IsDir() {
			continue
		}
		f := parseFilename(c.Name())
		if f == nil {
			continue
		}
		f.URL = cfg.BaseURL + filepath.Join(folder, c.Name())
		bin, err := ioutil.ReadFile(filepath.Join(cfg.ArtifactsPath, folder, c.Name()))
		if err == nil {
			f.Sha256 = calculateSha256(bin)
		}
		files = append(files, *f)
	}
	return
}
