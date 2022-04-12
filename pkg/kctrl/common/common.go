package common

import (
	"bytes"
	"errors"
	"fmt"
	goCLIUI "github.com/cppforlife/go-cli-ui/ui"
	"os"
	"os/exec"
)

var ui goCLIUI.UI

type Step interface {
	PreInteract() error
	PostInteract() error
	Interact() error
	Run() error
}

func GetUI() *goCLIUI.UI {
	if ui == nil {
		confUI := goCLIUI.NewConfUI(goCLIUI.NewNoopLogger())
		ui = confUI
	}
	return &ui
}

func Execute(cmd string, args []string) (string, error) {
	cmdFound := isExecutableInstalled(cmd)
	if !cmdFound {
		fmt.Printf("Executable \"%s\" not installed", cmd)
	}
	command := exec.Command(cmd, args...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	command.Stdout = &out
	command.Stderr = &stderr

	err := command.Run()
	if err != nil {
		return "", errors.New(fmt.Sprintf("%s: %s", fmt.Sprint(err), stderr.String()))
	}
	return out.String(), nil
}

func isExecutableInstalled(cmd string) bool {
	_, err := exec.LookPath(cmd)
	if err != nil {
		//TODO should log an error here instead of sending it up
		return false
	}
	return true

}

func ValidateInputYesOrNo(input string) (bool, bool) {
	if input == "y" || input == "Y" {
		return true, true
	} else if input == "n" || input == "N" {
		return false, true
	}
	return false, false
}

func GetPkgLocation(pkgName string) string {
	var pkgLocation string
	pkgLocation, _ = os.UserHomeDir()
	return pkgLocation + "/.kctrl" + pkgName
}
