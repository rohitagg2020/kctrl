package create

type UpstreamStep struct {
}

const (
	Github int = iota
	HelmChart
	Image
)

func (u UpstreamStep) PreInteract(o *CreateOptions) error {
	str := `In Carvel, An upstream source is the location from where we want to sync the software configuration.
Different types of upstream are:`

	o.Ui.PrintBlock([]byte(str))
	return nil
}

func (u UpstreamStep) PostInteract(o *CreateOptions) error {
	panic("implement me")
}

func (u UpstreamStep) Interact(o *CreateOptions) error {
	upstreamTypeSelected, err := o.Ui.AskForChoice("Different types of upstream are", []string{"Github Release", "HelmChart", "Image"})
	if err != nil {

	}
	switch upstreamTypeSelected {
	case Github:
		githubStep := GithubStep{}
		githubStep.Run(o)
	}
	panic("implement me")
}

func (u UpstreamStep) Run(o *CreateOptions) error {
	u.PreInteract(o)
	u.Interact(o)
	u.PostInteract(o)
	return nil
}
