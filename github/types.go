package github

import (
	"github.com/google/go-github/v37/github"
	"time"
)

type Github struct {
	c        *github.Client
	cfg      *Config
	releases map[string]*Version
	branches map[string]struct{}
	dev      *Version
}

type Version struct {
	Version   string
	Changelog string
	Date      time.Time
	Rc        bool
}
