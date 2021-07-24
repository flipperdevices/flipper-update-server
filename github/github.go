package github

import (
	"context"
	"errors"
	"github.com/google/go-github/v37/github"
	"golang.org/x/oauth2"
)

var ctx = context.Background()

func New(cfg Config) (*Github, error) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: cfg.GithubToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	c := github.NewClient(tc)
	gh := &Github{
		c:   c,
		cfg: &cfg,
	}

	err := gh.Sync()
	if err != nil {
		return nil, err
	}
	return gh, nil
}

func (gh *Github) Sync() error {
	releases, err := gh.fetchReleases()
	if err != nil {
		return errors.New("releases: " + err.Error())
	}
	gh.releases = releases

	branches, err := gh.fetchBranches()
	if err != nil {
		return errors.New("branches: " + err.Error())
	}
	gh.branches = branches

	dev, err := gh.fetchDev()
	if err != nil {
		return errors.New("dev: " + err.Error())
	}
	gh.dev = dev

	return nil
}

func (gh *Github) Lookup(ref string) (*Version, bool) {
	if ref == gh.cfg.DevelopmentBranch {
		return gh.dev, false
	}
	r, ok := gh.releases[ref]
	if ok {
		return r, false
	}
	_, ok = gh.branches[ref]
	return nil, ok
}

func (gh *Github) fetchReleases() (map[string]*Version, error) {
	releases, _, err := gh.c.Repositories.ListReleases(ctx, gh.cfg.RepoOwner, gh.cfg.RepoName, &github.ListOptions{
		Page:    1,
		PerPage: 100,
	})
	if err != nil {
		return nil, err
	}
	m := make(map[string]*Version)
	for _, r := range releases {
		m[r.GetTagName()] = &Version{
			Version:   r.GetTagName(),
			Changelog: r.GetBody(),
			Date:      r.GetCreatedAt().Time,
			Rc:        r.GetPrerelease(),
		}
	}
	return m, nil
}

func (gh *Github) fetchBranches() (map[string]struct{}, error) {
	branches, _, err := gh.c.Repositories.ListBranches(ctx, gh.cfg.RepoOwner, gh.cfg.RepoName, &github.BranchListOptions{
		ListOptions: github.ListOptions{
			Page:    1,
			PerPage: 100,
		},
	})
	if err != nil {
		return nil, err
	}
	m := make(map[string]struct{})
	for _, b := range branches {
		m[b.GetName()] = struct{}{}
	}
	return m, nil
}

func (gh *Github) fetchDev() (*Version, error) {
	commit, _, err := gh.c.Repositories.GetCommit(ctx, gh.cfg.RepoOwner, gh.cfg.RepoName, "HEAD")
	if err != nil {
		return nil, err
	}
	return &Version{
		Version:   commit.GetSHA()[0:7],
		Changelog: "Last commit: " + commit.Commit.GetMessage(),
		Date:      commit.Commit.Author.GetDate(),
	}, nil
}
