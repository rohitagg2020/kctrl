package create

import (
	"strings"
)

type GithubStep struct {
	RepoSlug   string
	ReleaseTag string
	Paths      []string
}

func (g *GithubStep) PreInteract(o *CreateOptions) error {
	return nil
}

func (g *GithubStep) PostInteract(o *CreateOptions) error {
	str := `# We have all the information needed to sync the upstream.
# Lets create vendir.yaml file with above inputs.`
	o.Ui.PrintBlock([]byte(str))

	return nil
}

func (g *GithubStep) Interact(o *CreateOptions) error {
	repoSlug, err := o.Ui.AskForText("Enter slug for repository(org/repo)")
	if err != nil {
		return err
	}
	g.RepoSlug = repoSlug
	releaseTag, err := g.getVersion(o)
	if err != nil {
		return err
	}
	g.ReleaseTag = releaseTag
	includeEverything, err := o.Ui.AskForText("Does you package encompasses everything mentioned in the github Release(y/n)")
	validateInput(includeEverything)
	var paths []string
	if !false {
		paths, err = g.getPaths(o)
		if err != nil {
			return err
		}
	}
	g.Paths = paths

	return nil
}

func (g GithubStep) getPaths(o *CreateOptions) ([]string, error) {
	str := `
# Now, we need to enter the specific paths which we want to include as package content. More than one paths can be added with comma separator.`
	o.Ui.PrintBlock([]byte(str))
	path, err := o.Ui.AskForText("Enter the paths which needs to be included as part of this package")
	if err != nil {
		return nil, err
	}
	paths := strings.Split(path, ",")
	return paths, nil
}

func (g GithubStep) getVersion(o *CreateOptions) (string, error) {
	useLatestVersion, err := o.Ui.AskForText("Do you want to use the latest released version(y/n)")
	if err != nil {
		return "", err
	}
	validateInput(useLatestVersion)
	//if useLatestVersion {
	if false {

	} else {
		o.Ui.PrintBlock([]byte("Ok. Then we have to mention the specific release tag now which we want to package."))
		releaseTag, err := o.Ui.AskForText("Enter the release tag")
		if err != nil {
			return "", err
		}
		return releaseTag, nil
	}
	//o.Ui.PrintBlock([])
	return "", nil
}

func (g *GithubStep) Run(o *CreateOptions) error {
	g.PreInteract(o)
	g.Interact(o)
	g.PostInteract(o)
	return nil
}
