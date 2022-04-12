package fetch

import (
	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/rohitagg2020/kctrl/pkg/kctrl/common"
	"github.com/rohitagg2020/kctrl/pkg/kctrl/fetch/imgpkg"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
)

const (
	Imgpkg int = iota
	HelmChart
	Inline
)

var fetchTypeNames = []string{"Imgpkg(recommended)", "HelmChart", "Inline"}

type FetchStep struct {
	Ui       ui.UI `yaml:"_"`
	PkgName  string
	AppFetch []v1alpha1.AppFetch
}

func NewFetchStep(ui ui.UI, pkgName string) *FetchStep {
	fetchStep := FetchStep{
		Ui:      ui,
		PkgName: pkgName,
	}
	return &fetchStep
}

func (fetch FetchStep) PreInteract() error {
	str := `# Now, we have to add the configuration which makes up the package for distribution. 
# Configuration can be fetched from different types of sources.`
	fetch.Ui.PrintBlock([]byte(str))
	return nil
}

func (fetch *FetchStep) Interact() error {
	var appFetchList []v1alpha1.AppFetch
	var fetchOptionSelected int
	fetchOptionSelected, err := fetch.Ui.AskForChoice("Enter the fetch configuration type", fetchTypeNames)
	if err != nil {
		return err
	}
	switch fetchOptionSelected {
	case Imgpkg:
		imgpkgStep := imgpkg.NewImgPkgStep(*common.GetUI(), fetch.PkgName)
		imgpkgStep.Run()
		appFetchList = append(appFetchList, v1alpha1.AppFetch{
			ImgpkgBundle: &imgpkgStep.ImgpkgBundle})
	}

	fetch.AppFetch = appFetchList
	return nil
}

func (fetch FetchStep) PostInteract() error {
	return nil
}

func (fetch *FetchStep) Run() error {
	fetch.PreInteract()
	fetch.Interact()
	fetch.PostInteract()
	return nil
}
