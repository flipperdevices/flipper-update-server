package github

type Config struct {
	GithubToken       string `env:"GH_TOKEN,required"`
	RepoOwner         string `env:"REPO_OWNER,required"`
	RepoName          string `env:"REPO_NAME,required"`
	DevelopmentBranch string `env:"DEV_BRANCH" envDefault:"dev"`
}
