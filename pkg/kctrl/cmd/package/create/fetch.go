package create

import "fmt"

const (
	Imgpkg int = iota
	Helm
	HelmRelease
	Inline
)

var fetchTypeNames = []string{"Imgpkg(recommended)", "Helm", "Helm Release", "Inline"}

type FetchStep struct {
}

func (fetch FetchStep) PreInteract(o *CreateOptions) error {
	str := `Now, we have to add the configuration which makes up the package for distribution. 
Configuration can be fetched from different types of sources.`
	o.Ui.PrintBlock([]byte(str))
	return nil
}

func (fetch FetchStep) Interact(o *CreateOptions) error {
	var fetchOptionSelected int
	fetchOptionSelected, err := o.Ui.AskForChoice("Fetch Configuration Types", fetchTypeNames)
	if err != nil {

	}
	switch fetchOptionSelected {
	case Imgpkg:
		fmt.Println("Imgpkg called")
		imgpkgStep := ImgpkgStep{}
		imgpkgStep.Run(o)
	}
	return nil
}

func (fetch FetchStep) PostInteract(o *CreateOptions) error {

	return nil
}

func (fetch FetchStep) Run(o *CreateOptions) error {
	fetch.PreInteract(o)
	fetch.Interact(o)
	fetch.PostInteract(o)
	return nil
}
