package create

import (
	"fmt"
	"github.com/cppforlife/go-cli-ui/ui"
	cmdcore "github.com/rohitagg2020/kctrl/pkg/kctrl/cmd/core"
	"github.com/rohitagg2020/kctrl/pkg/kctrl/logger"
	"github.com/spf13/cobra"
)

type CreateOptions struct {
	Ui             ui.UI
	logger         logger.Logger
	Name           string
	pkgCmdTreeOpts cmdcore.PackageCommandTreeOpts
}

func NewCreateOptions(ui ui.UI, logger logger.Logger, name string, pkgCmdTreeOpts cmdcore.PackageCommandTreeOpts) *CreateOptions {
	return &CreateOptions{Ui: ui, logger: logger, Name: name, pkgCmdTreeOpts: pkgCmdTreeOpts}
}

func NewCreateCmd(o *CreateOptions, flagsFactory cmdcore.FlagsFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Aliases: []string{"g"},
		Short:   "Create a package",
		Args:    cobra.ExactArgs(1),
		RunE:    func(_ *cobra.Command, args []string) error { return o.Run(args) },
		Example: cmdcore.Examples{
			cmdcore.Example{"Create a package",
				[]string{"package", "create", "pkg-a"},
			},
		}.Description("-p", o.pkgCmdTreeOpts),
		SilenceUsage: true,
		Annotations:  map[string]string{"table": ""},
	}

	/*if !o.pkgCmdTreeOpts.PositionalArgs {
		cmd.Flags().StringVarP(&o.Name, "package", "p", "", "Set package name (required)")
	} else {
		cmd.Use = "get PACKAGE_NAME or PACKAGE_NAME/VERSION"
	}
	*/
	return cmd
}

func getStartBlock() []byte {
	str := `Lets start off on the package creation for pkg-a
First we need a directory to store all configurations
Creating directory <DEFAULT_LOC>/pkg-a
mkdir -p <DEFAULT_LOC>/pkg-a`
	return []byte(str)
}

type CreateStep struct {
}

func (create CreateStep) PreInteract(o *CreateOptions) error {
	o.Ui.PrintBlock(getStartBlock())
	output, err := Execute("mkdir", []string{"-p", "/Users/roaggarwal/.kctrl/pkg-a"})
	if err != nil {
		fmt.Println("Error creating package directory.Error is: %s", err.Error())
		return err
	}
	o.Ui.PrintBlock([]byte(output))
	return nil
}

func (create CreateStep) Interact(o *CreateOptions) error {
	fetchStep := FetchStep{}
	fetchStep.Run(o)

	return nil
}

func (create CreateStep) PostInteract() error {
	return nil
}

func (o *CreateOptions) Run(args []string) error {

	createStep := CreateStep{}
	createStep.PreInteract(o)

	createStep.Interact(o)

	createStep.PostInteract()

	/*fmt.Println(o.Ui.IsInteractive())

	fmt.Println("Hello")
	fetchConfTypes := []string{
		"imgpkg",
		"helm Chart",
	}
	choice, err := o.Ui.AskForChoice("Enter the Fetch Configuration Types", fetchConfTypes)
	if err != nil {
		return err
	}
	fmt.Println(choice)
	pass, err := o.Ui.AskForPassword("Enter your registry password")
	if err != nil {

	}

	fmt.Println(pass)
	relVersion := "10.0"
	relVersion, err = o.Ui.AskForText("Enter the release version")
	if err != nil {
		return err
	}
	fmt.Println(relVersion)
	o.Ui.AskForConfirmation()*/
	return nil
}
