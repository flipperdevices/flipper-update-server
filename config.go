package main

import "github.com/flipper-zero/flipper-update-server/github"

type config struct {
	Github        github.Config
	Excluded      []string `env:"EXCLUDED"`
	ArtifactsPath string   `env:"ARTIFACTS_PATH" envDefault:"/artifacts"`
	BaseURL       string   `env:"BASE_URL"`
}
