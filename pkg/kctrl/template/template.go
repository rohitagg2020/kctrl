package template

import (
	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/rohitagg2020/kctrl/pkg/kctrl/common"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
)

const (
	Ytt int = iota
	HelmTemplate
)

type TemplateStep struct {
	Ui          ui.UI
	AppTemplate []v1alpha1.AppTemplate
}

func NewTemplateStep(ui ui.UI) *TemplateStep {
	templateStep := TemplateStep{
		Ui: ui,
	}
	return &templateStep
}

func (t TemplateStep) PreInteract() error {
	str := `# Next step is to define which templating tool to be used to render the package template. 
# A package template can be rendered by various different tools.`
	t.Ui.PrintBlock([]byte(str))
	return nil
}

func (t TemplateStep) PostInteract() error {
	return nil
}

func (t *TemplateStep) Interact() error {
	templateType, err := t.Ui.AskForChoice("Enter the templating tool to be used", []string{"ytt(recommended)", "helmTemplate"})
	if err != nil {
		return err
	}
	var appTemplateList []v1alpha1.AppTemplate
	switch templateType {
	case Ytt:
		yttTemplateStep := YttTemplateStep{}
		yttTemplateStep.Run()
		yttAppTemplate := v1alpha1.AppTemplate{
			Ytt: &yttTemplateStep.appTemplateYtt,
		}
		appTemplateList = append(appTemplateList, yttAppTemplate)
		t.AppTemplate = appTemplateList
	case HelmTemplate:
	}
	/*str := `# Next step is to use kbld to resolve the image references to use digest. kbld will resolve the image and store the digest in a images.yml file. To skip, press enter`
	t.Ui.PrintBlock([]byte(str))
	kbldPath, err := t.Ui.AskForText("Enter the path where kbld should store the digest(default: .imgpkg/images.yaml)")
	if err != nil {
		return err
	}
	kbldAppTemplate := v1alpha1.AppTemplate{
		Kbld: &v1alpha1.AppTemplateKbld{Paths: strings.Split(kbldPath, ",")},
	}
	appTemplateList = append(appTemplateList, kbldAppTemplate)

	*/
	return nil
}

func (t *TemplateStep) Run() error {
	t.Ui = *common.GetUI()
	t.PreInteract()
	t.Interact()
	t.PostInteract()
	return nil
}
