package create

import (
	"encoding/json"
	"fmt"
	"github.com/cppforlife/go-cli-ui/ui"
	cmdcore "github.com/rohitagg2020/kctrl/pkg/kctrl/cmd/core"
	"github.com/rohitagg2020/kctrl/pkg/kctrl/common"
	"github.com/rohitagg2020/kctrl/pkg/kctrl/fetch"
	"github.com/rohitagg2020/kctrl/pkg/kctrl/logger"
	"github.com/rohitagg2020/kctrl/pkg/kctrl/template"
	"github.com/spf13/cobra"
	v1alpha12 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging/v1alpha1"
	"sigs.k8s.io/yaml"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"os"
	"strings"
	"time"
)

type CreateOptions struct {
	logger         logger.Logger
	Name           string
	pkgCmdTreeOpts cmdcore.PackageCommandTreeOpts
}

func NewCreateOptions(logger logger.Logger, name string, pkgCmdTreeOpts cmdcore.PackageCommandTreeOpts) *CreateOptions {
	return &CreateOptions{logger: logger, Name: name, pkgCmdTreeOpts: pkgCmdTreeOpts}
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
	return cmd
}

func (create CreateStep) getStartBlock() []byte {
	str := fmt.Sprintf(`
# Lets start off on the package creation for %s
# First we need a directory to store all configurations
# Creating directory %s
# mkdir -p %s
`, common.GetPkgLocation(create.PkgName), common.GetPkgLocation(create.PkgName), common.GetPkgLocation(create.PkgName))
	return []byte(str)
}

type CreateStep struct {
	Ui           ui.UI
	PkgVersion   string
	PkgName      string
	FqName       string
	ValuesSchema ValuesSchema
	Fetch        fetch.FetchStep
	Template     template.TemplateStep
	Maintainers  []v1alpha1.Maintainer
}

type ValuesSchema struct {
	// +optional
	// +nullable
	// +kubebuilder:pruning:PreserveUnknownFields
	OpenAPIv3 runtime.RawExtension `json:"openAPIv3,omitempty" protobuf:"bytes,1,opt,name=openAPIv3"`
}

func (create CreateStep) PreInteract() error {
	create.Ui.PrintBlock(create.getStartBlock())
	output, err := common.Execute("mkdir", []string{"-p", common.GetPkgLocation(create.PkgName)})
	if err != nil {
		create.Ui.ErrorLinef("Error creating package directory.Error is: %s", err.Error())
		return err
	}
	create.Ui.PrintBlock([]byte(output))
	return nil
}

func (create *CreateStep) getPkgVersionBlock() []byte {
	str := fmt.Sprintf(`
# A package can have multiple versions. These versions are used by PackageInstall to install specific package into the Kubernetes cluster. 
`)
	return []byte(str)
}

func (create *CreateStep) getFQPkgNameBlock() []byte {
	str := `# A package name must be a fully qualified name. 
# It must consist of at least three segments separated by a '.'
# Cannot have a trailing '.' e.g. samplePackage.corp.com`
	return []byte(str)
}

func (create *CreateStep) Interact() error {
	//Get Package Version
	create.Ui.PrintBlock(create.getPkgVersionBlock())
	pkgVersion, err := create.Ui.AskForText("Enter the package version to be used")
	if err != nil {
		return err
	}
	create.PkgVersion = pkgVersion

	//Get Fully Qualified Name of the Package
	create.Ui.PrintBlock(create.getFQPkgNameBlock())
	fqName, err := create.Ui.AskForText("Enter the fully qualified package name(default: pkg-a.corp.com)")
	if err != nil {
		return err
	}
	create.FqName = fqName

	fetchConfiguration := fetch.NewFetchStep(*common.GetUI(), create.PkgName)
	fetchConfiguration.Run()
	create.Fetch = *fetchConfiguration

	templateConfiguration := template.NewTemplateStep(*common.GetUI())
	templateConfiguration.Run()
	create.Template = *templateConfiguration

	valuesSchema, err := create.getValueSchema()
	if err != nil {
		return err
	}
	create.ValuesSchema = valuesSchema

	maintainerNames, err := create.Ui.AskForText("Enter the Maintainer's Name. Multiple names can be provided by comma separated values")
	if err != nil {
		return err
	}
	var maintainers []v1alpha1.Maintainer
	for _, maintainerName := range strings.Split(maintainerNames, ",") {
		maintainers = append(maintainers, v1alpha1.Maintainer{Name: maintainerName})
	}
	create.Maintainers = maintainers
	return nil
}

func (create CreateStep) getValueSchema() (ValuesSchema, error) {
	valuesSchema := ValuesSchema{}
	var isValueSchemaSpecified bool
	var isValidInput bool
	input, err := create.Ui.AskForText("Do you want to specify the values Schema(y/n)")
	if err != nil {
		return valuesSchema, err
	}
	for {
		isValueSchemaSpecified, isValidInput = common.ValidateInputYesOrNo(input)
		if !isValidInput {
			input, err = create.Ui.AskForText("Invalid input.")
			if err != nil {
				return valuesSchema, err
			}
			continue
		}
		if isValueSchemaSpecified {
			valuesSchemaFileLocation, err := create.Ui.AskForText("Enter the values schema file location")
			if err != nil {
				return valuesSchema, err
			}
			valuesSchemaData, err := readDataFromFile(valuesSchemaFileLocation)
			if err != nil {
				return valuesSchema, err
			}
			valuesSchema = ValuesSchema{
				OpenAPIv3: runtime.RawExtension{
					Raw: valuesSchemaData,
				},
			}
		}
	}
	return valuesSchema, nil

}
func readDataFromFile(fileLocation string) ([]byte, error) {
	//TODO should we read it in a buffer
	data, err := os.ReadFile(fileLocation)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (create CreateStep) PostInteract() error {
	err := create.createAndPrintPackageCR()
	if err != nil {
		return err
	}
	err = create.createAndPrintPackageMetadataCR()
	if err != nil {
		return err
	}
	return nil
}

func (create CreateStep) createAndPrintPackageCR() error {
	str := `# Great, we have all the data needed to create the package.yml and packageMetadata.yml. 
# This is how the package.yml will look like
# cat package.yaml
`
	create.Ui.PrintBlock([]byte(str))
	pkg := create.populatePkgFromCreate()

	//TODO: remove this comment. Marshal will make yaml/json
	jsonPackageData, err := json.Marshal(&pkg)
	yaml.JSONToYAML(jsonPackageData)
	packageData, err := yaml.JSONToYAML(jsonPackageData)
	if err != nil {
		return err
	}
	writeToFile(common.GetPkgLocation(create.PkgName)+"/package.yaml", packageData)

	output, err := common.Execute("cat", []string{common.GetPkgLocation(create.PkgName) + "/package.yaml"})
	if err != nil {
		return err
	}
	create.Ui.PrintBlock([]byte(output))
	return nil
}

func (create CreateStep) createAndPrintPackageMetadataCR() error {
	str := `
# This is how the packageMetadata.yml will look like
# cat packageMetadata.yaml
`
	create.Ui.PrintBlock([]byte(str))
	pkgMetadata := create.populatePkgMetadataFromCreate()
	jsonPackageMetadataData, err := json.Marshal(&pkgMetadata)
	packageMetadataData, err := yaml.JSONToYAML(jsonPackageMetadataData)
	if err != nil {
		return err
	}
	writeToFile(common.GetPkgLocation(create.PkgName)+"/packageMetadata.yaml", packageMetadataData)

	output, err := common.Execute("cat", []string{common.GetPkgLocation(create.PkgName) + "/packageMetadata.yaml"})
	if err != nil {
		return err
	}

	create.Ui.PrintBlock([]byte(output))
	str = fmt.Sprintf(`# Both the files can be accessed from the following location: %s
`, common.GetPkgLocation(create.PkgName))
	create.Ui.PrintBlock([]byte(str))
	return nil
}

func writeToFile(path string, data []byte) error {
	err := os.WriteFile(path, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (create CreateStep) populatePkgMetadataFromCreate() v1alpha1.PackageMetadata {
	packageMetadataContent := v1alpha1.PackageMetadata{
		TypeMeta:   v1.TypeMeta{Kind: "PackageMetadata", APIVersion: "data.packaging.carvel.dev/v1alpha1"},
		ObjectMeta: v1.ObjectMeta{Name: create.FqName},
		Spec: v1alpha1.PackageMetadataSpec{
			DisplayName:      create.PkgName,
			LongDescription:  "A long description",
			ShortDescription: "A short description",
			ProviderName:     "VMWare",
			Maintainers:      create.Maintainers,
		},
	}
	return packageMetadataContent
}

func (create CreateStep) populatePkgFromCreate() v1alpha1.Package {
	packageContent := v1alpha1.Package{
		TypeMeta:   v1.TypeMeta{Kind: "Package", APIVersion: "data.packaging.carvel.dev/v1alpha1"},
		ObjectMeta: v1.ObjectMeta{Namespace: "default", Name: create.FqName + "." + create.PkgVersion},
		Spec: v1alpha1.PackageSpec{
			RefName:                         create.FqName,
			Version:                         create.PkgVersion,
			Licenses:                        []string{"Apache 2.0", "MIT"},
			ReleasedAt:                      v1.Time{time.Now()},
			CapactiyRequirementsDescription: "",
			ReleaseNotes:                    "",
			Template: v1alpha1.AppTemplateSpec{Spec: &v1alpha12.AppSpec{
				ServiceAccountName: "",
				Cluster:            nil,
				Fetch:              create.Fetch.AppFetch,
				Template:           create.Template.AppTemplate,
				Deploy: []v1alpha12.AppDeploy{
					v1alpha12.AppDeploy{Kapp: &v1alpha12.AppDeployKapp{}},
				},
				Paused:     false,
				Canceled:   false,
				SyncPeriod: nil,
				NoopDelete: false,
			}},
			ValuesSchema: v1alpha1.ValuesSchema{},
		},
	}
	return packageContent
}

func (o *CreateOptions) Run(args []string) error {
	//Configure GlobalUI
	createPkg := CreateStep{Ui: *common.GetUI()}
	createPkg.PreInteract()
	createPkg.Interact()
	createPkg.PostInteract()
	return nil
}
