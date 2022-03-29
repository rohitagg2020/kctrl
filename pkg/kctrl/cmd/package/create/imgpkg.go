package create

type ImgpkgStep struct {
}

func (imgpkg ImgpkgStep) PreInteract(o *CreateOptions) error {

	return nil
}

func validateInput(string) bool {
	return true
}
func (imgpkg ImgpkgStep) Interact(o *CreateOptions) error {
	isImgpkgCreated, err := o.Ui.AskForText("Is the imgpkg bundle already created(y/n)")
	if err != nil {

	}
	validateInput(isImgpkgCreated)

	if false {

	} else {
		createImgPkg(o)
	}
	return nil
}

func (imgpkg ImgpkgStep) PostInteract(o *CreateOptions) error {

	return nil
}

func (imgpkg ImgpkgStep) Run(o *CreateOptions) error {
	imgpkg.PreInteract(o)
	imgpkg.Interact(o)
	imgpkg.PostInteract(o)
	return nil
}

func createImgPkg(o *CreateOptions) {
	createImgPkgStep := createImgPkgStep{}
	createImgPkgStep.PreInteract(o)
	createImgPkgStep.Interact(o)
	createImgPkgStep.PostInteract(o)
}

type createImgPkgStep struct {
}

func (createImgpkg createImgPkgStep) PreInteract(o *CreateOptions) error {
	str := `Cool, lets create the imgpkg bundle first
Creating directory <DEFAULT_LOC>/pkg-a/bundle/config
mkdir -p <DEFAULT_LOC>/pkg-a/bundle/config`
	o.Ui.PrintBlock([]byte(str))
	Execute("mkdir", []string{"-p", "/Users/roaggarwal/.kctrl/pkg-a/bundle/config"})

	str = `Creating directory <DEFAULT_LOC>/pkg-a/bundle/config/.imgpkg
mkdir -p <DEFAULT_LOC>/pkg-a/bundle/config/.imgpkg`
	o.Ui.PrintBlock([]byte(str))
	Execute("mkdir", []string{"-p", "/Users/roaggarwal/.kctrl/pkg-a/bundle/config/.imgpkg"})
	return nil
}

func (createImgpkg createImgPkgStep) Interact(o *CreateOptions) error {
	upstreamStep := UpstreamStep{}
	upstreamStep.Run(o)
	return nil
}

func (createImgpkg createImgPkgStep) PostInteract(o *CreateOptions) error {
	return nil
}
