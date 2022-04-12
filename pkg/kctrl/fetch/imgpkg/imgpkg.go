package imgpkg

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/rohitagg2020/kctrl/pkg/kctrl/common"
	"github.com/rohitagg2020/kctrl/pkg/kctrl/fetch/imgpkg/upstream"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"strings"
)

type ImgpkgStep struct {
	Ui           ui.UI  `yaml:"_"`
	PkgName      string `yaml:"_"`
	ImgpkgBundle v1alpha1.AppFetchImgpkgBundle
}

func NewImgPkgStep(ui ui.UI, pkgName string) *ImgpkgStep {
	imgpkg := ImgpkgStep{
		Ui:      ui,
		PkgName: pkgName,
	}
	return &imgpkg
}

func (imgpkg ImgpkgStep) PreInteract() error {
	return nil
}

func (imgpkg *ImgpkgStep) Interact() error {
	var isImgpkgCreated bool
	input, err := imgpkg.Ui.AskForText("Is the imgpkg bundle already created(y/n)")
	if err != nil {
		//TODO
	}
	for {
		var isValidInput bool
		isImgpkgCreated, isValidInput = common.ValidateInputYesOrNo(input)
		if isValidInput {
			break
		} else {
			input, err = imgpkg.Ui.AskForText("Invalid input.")
			if err != nil {
				return err
			}
		}
	}

	if isImgpkgCreated {

	} else {
		createImgPkgStep := CreateImgPkgStep{}
		createImgPkgStep.Run()
		imgpkg.ImgpkgBundle.Image = createImgPkgStep.Image
	}
	return nil
}

func (imgpkg ImgpkgStep) PostInteract() error {

	return nil
}

func (imgpkg *ImgpkgStep) Run() error {
	imgpkg.Ui = *common.GetUI()
	imgpkg.PreInteract()
	imgpkg.Interact()
	imgpkg.PostInteract()
	return nil
}

type CreateImgPkgStep struct {
	Ui      ui.UI
	Image   string
	PkgName string
}

func (createImgPkgStep *CreateImgPkgStep) Run() {
	createImgPkgStep.Ui = *common.GetUI()
	createImgPkgStep.PreInteract()
	createImgPkgStep.Interact()
	createImgPkgStep.PostInteract()
}

func (createImgpkg CreateImgPkgStep) PreInteract() error {
	str := fmt.Sprintf(`# Cool, lets create the imgpkg bundle first
# Creating directory %s/bundle/config
# mkdir -p %s/bundle/config`, common.GetPkgLocation(createImgpkg.PkgName), common.GetPkgLocation(createImgpkg.PkgName))
	createImgpkg.Ui.PrintBlock([]byte(str))
	output, err := common.Execute("mkdir", []string{"-p", common.GetPkgLocation(createImgpkg.PkgName) + "/bundle/config"})
	if err != nil {
		return err
	}
	createImgpkg.Ui.PrintBlock([]byte(output))
	str = fmt.Sprintf(`
# Creating directory %s/bundle/.imgpkg
# mkdir -p %s/bundle/.imgpkg`, common.GetPkgLocation(createImgpkg.PkgName), common.GetPkgLocation(createImgpkg.PkgName))
	createImgpkg.Ui.PrintBlock([]byte(str))
	output, err = common.Execute("mkdir", []string{"-p", common.GetPkgLocation(createImgpkg.PkgName) + "/bundle/.imgpkg"})
	if err != nil {
		return err
	}
	createImgpkg.Ui.PrintBlock([]byte(output))
	return nil
}

func (createImgpkg CreateImgPkgStep) Interact() error {
	upstreamStep := upstream.UpstreamStep{}
	upstreamStep.Run()
	str := `
# If you wish to use default values, then skip next step. Otherwise, we can use the ytt(a templating and overlay tool) to provide custom values.`
	createImgpkg.Ui.PrintBlock([]byte(str))
	var useYttAsTemplate bool
	for {
		input, err := createImgpkg.Ui.AskForText("Do you want to use ytt as a templating and overlay tool(y/n)")
		if err != nil {
			return err
		}
		var isInputValid bool
		useYttAsTemplate, isInputValid = common.ValidateInputYesOrNo(input)
		if isInputValid {
			break
		}
	}
	if useYttAsTemplate {
		yttPath, err := createImgpkg.Ui.AskForText("Enter the path where ytt files are located:")
		yttPath = yttPath
		if err != nil {
			return err
		}
		str = fmt.Sprintf(`# Copying the ytt files inside the package.
# cp -r %s %s/bundle/config`, yttPath, common.GetPkgLocation(createImgpkg.PkgName))
		createImgpkg.Ui.PrintBlock([]byte(str))
		common.Execute("cp", []string{"-r", yttPath, common.GetPkgLocation(createImgpkg.PkgName) + "/bundle/config"})
	}

	return nil
}

func (createImgpkg *CreateImgPkgStep) PostInteract() error {
	str := fmt.Sprintf(`# imgpkg bundle configuration is now complete. Let's use kbld to lock it down.
# kbld allows to build the imgpkg bundle with immutable image references.
# kbld scans a package configuration for any references to images and creates a mapping of image tags to a URL with a sha256 digest. 
# This mapping will then be placed into an images.yml lock file in your bundle/.imgpkg directory.
# Running kbld --file %s/bundle --imgpkg-lock-output %s/bundle/.imgpkg/images.yml`, common.GetPkgLocation(createImgpkg.PkgName), common.GetPkgLocation(createImgpkg.PkgName))
	createImgpkg.Ui.PrintBlock([]byte(str))

	output, err := common.Execute("kbld", []string{"--file", common.GetPkgLocation(createImgpkg.PkgName) + "/bundle", "--imgpkg-lock-output", common.GetPkgLocation(createImgpkg.PkgName) + "/bundle/.imgpkg/images.yml"})
	if err != nil {
		createImgpkg.Ui.PrintBlock([]byte(err.Error()))
		return err
	}
	//createImgpkg.Ui.PrintBlock([]byte(output))

	str = fmt.Sprintf(`
# Lets see how the images.yaml file looks like:
# Running cat %s/bundle/.imgpkg/images.yml`, common.GetPkgLocation(createImgpkg.PkgName))
	createImgpkg.Ui.PrintBlock([]byte(str))
	output, err = common.Execute("cat", []string{common.GetPkgLocation(createImgpkg.PkgName) + "/bundle/.imgpkg/images.yml"})
	if err != nil {
		return err
	}
	createImgpkg.Ui.PrintBlock([]byte(output))

	var pushBundle bool
	for {
		input, err := createImgpkg.Ui.AskForText("Do you want to push the bundle to the registry(y/n)")
		if err != nil {
			return err
		}
		var isValidInput bool
		pushBundle, isValidInput = common.ValidateInputYesOrNo(input)
		if isValidInput {
			break
		}
	}
	if pushBundle {
		registryAuthDetails, err := createImgpkg.PopulateRegistryAuthDetails()
		if err != nil {
			return err
		}
		//Can repoName be empty?
		repoName, err := createImgpkg.Ui.AskForText("Provide the repository name to which this bundle belong(e.g.: your org name)")
		if err != nil {
			return err
		}
		tagName, err := createImgpkg.Ui.AskForText("Do you want to provide the tag name(default: latest)")
		if err != nil {
			return err
		}
		bundleURL, err := createImgpkg.pushImgpkgBundleToRegistry(repoName, tagName, registryAuthDetails, common.GetPkgLocation(createImgpkg.PkgName)+"/bundle")
		if err != nil {
			return err
		}
		createImgpkg.Image = bundleURL

	}
	return nil
}

func (createImgPkgStep CreateImgPkgStep) pushImgpkgBundleToRegistry(repoName string, tagName string, authDetails RegistryAuthDetails, bundleLoc string) (string, error) {
	pushURL := authDetails.RegistryURL + "/" + repoName + ":" + tagName
	pushURL = "localhost:6000/rohitagg2020/certmanager:v1.5.3"
	str := fmt.Sprintf(`# Running imgpkg to push the bundle directory and indicate what project name and tag to give it.
# imgpkg push --bundle %s --file %s --json
`, pushURL, bundleLoc)
	createImgPkgStep.Ui.PrintBlock([]byte(str))

	output, err := common.Execute("imgpkg", []string{"push", "--bundle", pushURL, "--file", bundleLoc, "--registry-username", authDetails.Username, "--registry-password", authDetails.Password, "--json"})
	if err != nil {
		return "", err
	}

	bundleURL, err := getBundleURL(output, pushURL)
	//createImgPkgStep.Ui.Flush()
	createImgPkgStep.Ui.PrintBlock([]byte(output))
	strings.Trim(bundleURL, "'")
	return bundleURL, nil
}

type ImgpkgPushOutput struct {
	Lines  []string    `json:"Lines"`
	Tables interface{} `json:"Tables"`
	Blocks interface{} `json:"Blocks"`
}

func getBundleURL(output string, pushURL string) (string, error) {
	var imgPkgPushOutput ImgpkgPushOutput
	json.Unmarshal([]byte(output), &imgPkgPushOutput)

	for _, val := range imgPkgPushOutput.Lines {
		if strings.HasPrefix(val, "Pushed") {
			//TODO remove '' from the URL
			bundleURL := strings.Split(val, " ")[1]

			return strings.Trim(bundleURL, "'"), nil

		}
	}
	return "", errors.New("Unable to get the imgpkg URL")

}
