package upstream

import (
	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/rohitagg2020/kctrl/pkg/kctrl/common"
)

type GithubStep struct {
	RepoSlug                      string `yaml:"slug"`
	ReleaseTag                    string `yaml:"tag"`
	Ui                            ui.UI  `yaml:"-"`
	DisableAutoChecksumValidation bool   `yaml:"disableAutoChecksumValidation"`
}

func (g *GithubStep) PreInteract() error {
	return nil
}

func (g *GithubStep) PostInteract() error {

	return nil
}

func (g *GithubStep) Interact() error {
	repoSlug, err := g.Ui.AskForText("Enter slug for repository(org/repo)")
	if err != nil {
		return err
	}
	//repoSlug = "jetstack/cert-manager"
	g.RepoSlug = repoSlug
	releaseTag, err := g.getVersion()
	if err != nil {
		return err
	}
	g.ReleaseTag = releaseTag
	g.DisableAutoChecksumValidation = true
	return nil
}

func (g GithubStep) getVersion() (string, error) {
	var useLatestVersion bool
	for {
		input, err := g.Ui.AskForText("Do you want to use the latest released version(y/n)")
		if err != nil {
			return "", err
		}
		var isValidInput bool
		useLatestVersion, isValidInput = common.ValidateInputYesOrNo(input)
		if isValidInput {
			break
		} else {
			g.Ui.PrintBlock([]byte("Invalid input. Try Again"))
		}
	}

	//if useLatestVersion {
	if useLatestVersion {

	} else {
		g.Ui.PrintBlock([]byte("# Ok. Then we have to mention the specific release tag now which we want to package."))
		releaseTag, err := g.Ui.AskForText("Enter the release tag")
		if err != nil {
			return "", err
		}
		//releaseTag = "v1.5.3"
		return releaseTag, nil
	}
	//o.Ui.PrintBlock([])
	return "", nil
}

func (g *GithubStep) Run() error {
	g.Ui = *common.GetUI()
	g.PreInteract()
	g.Interact()
	g.PostInteract()
	return nil
}
