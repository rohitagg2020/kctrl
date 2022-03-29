package create

import (
	"fmt"
	"os/exec"
)

//var ui goCLIUI.UI

type Step interface {
	PreInteract(o *CreateOptions) error
	PostInteract(o *CreateOptions) error
	Interact(o *CreateOptions) error
}

/*
func GetUI() *goCLIUI.UI {
	if ui == nil {
		confUI := goCLIUI.NewConfUI(goCLIUI.NewNoopLogger())
		ui = confUI
	}
	return &ui
}

*/

func Execute(cmd string, args []string) (string, error) {
	cmdFound := isExecutableInstalled(cmd)
	if !cmdFound {
		fmt.Printf("Executable \"%s\" not installed", cmd)
	}
	command := exec.Command(cmd, args...)
	out, err := command.Output()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(out))
	return "", nil
}

func isExecutableInstalled(cmd string) bool {
	_, err := exec.LookPath(cmd)
	if err != nil {
		//TODO should log an error here instead of sending it up
		return false
	}
	return true

}
