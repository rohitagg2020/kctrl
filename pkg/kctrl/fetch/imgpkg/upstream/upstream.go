package upstream

import (
	"fmt"
	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/rohitagg2020/kctrl/pkg/kctrl/common"
	"gopkg.in/yaml.v3"

	"os"
	"strings"
)

const (
	GithubRelease int = iota
	HelmChart
	Image
)

type Content struct {
	Path          string
	GithubRelease GithubStep `yaml:"githubRelease,omitempty"`
	IncludePaths  []string   `yaml:"includePaths"`
}

type Directory struct {
	Path     string    `yaml:"path"`
	Contents []Content `yaml:"contents"`
}
type UpstreamStep struct {
	ApiVersion             string      `yaml:"apiVersion"`
	Kind                   string      `yaml:"kind"`
	MinimumRequiredVersion string      `yaml:"minimumRequiredVersion"`
	Directories            []Directory `yaml:"directories"`
	Ui                     ui.UI       `yaml:"-"`
	PkgName                string      `yaml:"-"`
}

func (u *UpstreamStep) populateUpstreamMetadata() {
	u.ApiVersion = "vendir.k14s.io/v1alpha1"
	u.Kind = "Config"
	u.MinimumRequiredVersion = "0.12.0"
}
func (u *UpstreamStep) PreInteract() error {
	str := `
# In Carvel, An upstream source is the location from where we want to sync the software configuration.
# Different types of upstream available are`

	u.Ui.PrintBlock([]byte(str))
	return nil
}

func (u *UpstreamStep) PostInteract() error {
	u.populateUpstreamMetadata()
	str := `# We have all the information needed to sync the upstream.
# Lets create vendir.yml file with above inputs.
`
	u.Ui.PrintBlock([]byte(str))
	data, err := yaml.Marshal(&u)
	if err != nil {
		fmt.Println("Unable to create vendir.yml")
		return err
	}
	f, err := os.Create(common.GetPkgLocation(u.PkgName) + "/bundle/vendir.yml")

	if err != nil {
		fmt.Println("File already exist")
		return err
	}

	defer f.Close()

	_, err = f.Write(data)
	if err != nil {
		return err
	}

	str = `cat vendir.yml
`
	u.Ui.PrintBlock([]byte(str))
	fmt.Println()
	resp, err := common.Execute("cat", []string{common.GetPkgLocation(u.PkgName) + "/bundle/vendir.yml"})
	if err != nil {
		fmt.Println("Unable to read vendir.yaml file")
		return err
	}

	u.Ui.PrintBlock([]byte(resp))
	str = fmt.Sprintf(`
# Next step is to Run vendir to sync the data from github to upstream folder.
# Running vendir sync --chdir %s/bundle
`, common.GetPkgLocation(u.PkgName))
	u.Ui.PrintBlock([]byte(str))
	resp, err = common.Execute("vendir", []string{"sync", "--chdir", common.GetPkgLocation(u.PkgName) + "/bundle"})
	if err != nil {
		fmt.Printf("Error while running vendir sync. Error is: %s", err.Error())
		return err
	}

	u.Ui.PrintBlock([]byte(resp))
	str = fmt.Sprintf(`# To ensure that data has been synced, lets do
# ls -l %s/bundle/config/`, common.GetPkgLocation(u.PkgName))
	u.Ui.PrintBlock([]byte(str))
	output, err := common.Execute("ls", []string{"-l", common.GetPkgLocation(u.PkgName) + "/bundle/config"})
	if err != nil {
		return err
	}
	u.Ui.PrintBlock([]byte(output))
	str = `
# After running vendir sync, there is one more file created i.e. bundle/vendir.lock.yml
# This lock file resolves the release tag to the specific GitHub release and declares that the config is the synchronization target path.
# Lets see its content
# cat bundle/vendir.lock.yml`
	u.Ui.PrintBlock([]byte(str))
	output, err = common.Execute("cat", []string{common.GetPkgLocation(u.PkgName) + "/bundle/vendir.lock.yml"})
	if err != nil {
		return err
	}
	u.Ui.PrintBlock([]byte(output))
	return nil
}

func (u *UpstreamStep) Interact() error {
	upstreamTypeSelected, err := u.Ui.AskForChoice("Enter the upstream type", []string{"Github Release", "HelmChart", "Image"})
	if err != nil {

	}
	var content Content
	switch upstreamTypeSelected {
	case GithubRelease:
		githubStep := GithubStep{}
		githubStep.Run()
		content.GithubRelease = githubStep

	}

	includedPaths, err := u.getIncludedPaths()
	if err != nil {
		return err
	}
	content.IncludePaths = includedPaths
	content.Path = "."

	u.Directories = []Directory{
		Directory{
			Path: "config",
			Contents: []Content{
				content,
			},
		},
	}

	return nil
}

func (u UpstreamStep) getIncludedPaths() ([]string, error) {
	var includeEverything bool
	for {
		input, err := u.Ui.AskForText("Does your package needs to include everything from the upstream(y/n)")
		if err != nil {
			return nil, err
		}
		var isValidInput bool
		includeEverything, isValidInput = common.ValidateInputYesOrNo(input)
		if isValidInput {
			break
		} else {
			u.Ui.PrintBlock([]byte("Invalid input. Try Again"))
		}
	}
	var paths []string
	var err error
	if includeEverything {

	} else {
		paths, err = u.getPaths()
		if err != nil {
			return nil, err
		}
	}
	return paths, nil
}

func (u UpstreamStep) getPaths() ([]string, error) {
	str := `# Now, we need to enter the specific paths which we want to include as package content. More than one paths can be added with comma separator.`
	u.Ui.PrintBlock([]byte(str))

	path, err := u.Ui.AskForText("Enter the paths which needs to be included as part of this package")
	if err != nil {
		return nil, err
	}
	//path = "cert-manager.yaml"
	paths := strings.Split(path, ",")
	return paths, nil
}

func (u *UpstreamStep) Run() error {
	u.Ui = *common.GetUI()
	u.PreInteract()
	u.Interact()
	u.PostInteract()
	return nil
}
