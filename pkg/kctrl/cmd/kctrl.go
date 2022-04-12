// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"io"

	"github.com/cppforlife/cobrautil"
	"github.com/cppforlife/go-cli-ui/ui"
	cmdcore "github.com/rohitagg2020/kctrl/pkg/kctrl/cmd/core"
	cmdpkg "github.com/rohitagg2020/kctrl/pkg/kctrl/cmd/package"
	pkgcreate "github.com/rohitagg2020/kctrl/pkg/kctrl/cmd/package/create"
	"github.com/rohitagg2020/kctrl/pkg/kctrl/logger"
	"github.com/rohitagg2020/kctrl/pkg/kctrl/version"
	"github.com/spf13/cobra"
)

type KctrlOptions struct {
	ui            *ui.ConfUI
	logger        *logger.UILogger
	configFactory cmdcore.ConfigFactory
	depsFactory   cmdcore.DepsFactory

	UIFlags         UIFlags
	LoggerFlags     LoggerFlags
	KubeAPIFlags    cmdcore.KubeAPIFlags
	KubeconfigFlags cmdcore.KubeconfigFlags
}

func NewKctrlOptions(ui *ui.ConfUI, configFactory cmdcore.ConfigFactory,
	depsFactory cmdcore.DepsFactory) *KctrlOptions {

	return &KctrlOptions{ui: ui, logger: logger.NewUILogger(ui),
		configFactory: configFactory, depsFactory: depsFactory}
}

func NewDefaultKctrlCmd(ui *ui.ConfUI) *cobra.Command {
	configFactory := cmdcore.NewConfigFactoryImpl()
	depsFactory := cmdcore.NewDepsFactoryImpl(configFactory, ui)
	options := NewKctrlOptions(ui, configFactory, depsFactory)
	flagsFactory := cmdcore.NewFlagsFactory(configFactory, depsFactory)
	return NewKctrlCmd(options, flagsFactory)
}

func NewKctrlCmd(o *KctrlOptions, flagsFactory cmdcore.FlagsFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kctrl",
		Short: "kctrl helps to manage packages and repositories on your Kubernetes cluster",

		RunE: cobrautil.ShowHelp,

		// Affects children as well
		SilenceErrors: true,
		SilenceUsage:  true,

		// Disable docs header
		DisableAutoGenTag: true,
		Version:           version.Version,

		// TODO bash completion
	}

	cmd.SetOutput(uiBlockWriter{o.ui}) // setting output for cmd.Help()

	cmd.SetUsageTemplate(cobrautil.HelpSectionsUsageTemplate([]cobrautil.HelpSection{
		cmdcore.PackageHelpGroup,
		cmdcore.RestOfCommandsHelpGroup,
	}))

	pkgOpts := cmdcore.PackageCommandTreeOpts{BinaryName: "kctrl", PositionalArgs: false, Color: true, JSON: true}

	SetGlobalFlags(o, cmd, flagsFactory, pkgOpts)

	ConfigurePathResolvers(o, cmd, flagsFactory)

	cmd.AddCommand(NewVersionCmd(NewVersionOptions(o.ui, o.depsFactory), flagsFactory))

	pkgCmd := cmdpkg.NewCmd()
	AddPackageCommands(o, pkgCmd, flagsFactory, pkgOpts)

	cmd.AddCommand(pkgCmd)

	ConfigureGlobalFlags(o, cmd, flagsFactory, pkgOpts.PositionalArgs)

	return cmd
}

func SetGlobalFlags(o *KctrlOptions, cmd *cobra.Command, flagsFactory cmdcore.FlagsFactory, opts cmdcore.PackageCommandTreeOpts) {
	o.UIFlags.Set(cmd, flagsFactory, opts)
	o.LoggerFlags.Set(cmd, flagsFactory)
	o.KubeAPIFlags.Set(cmd, flagsFactory)
	o.KubeconfigFlags.Set(cmd, flagsFactory)
}

func ConfigurePathResolvers(o *KctrlOptions, cmd *cobra.Command, flagsFactory cmdcore.FlagsFactory) {
	o.configFactory.ConfigurePathResolver(o.KubeconfigFlags.Path.Value)
	o.configFactory.ConfigureContextResolver(o.KubeconfigFlags.Context.Value)
	o.configFactory.ConfigureYAMLResolver(o.KubeconfigFlags.YAML.Value)
}

func ConfigureGlobalFlags(o *KctrlOptions, cmd *cobra.Command, flagsFactory cmdcore.FlagsFactory, positionalNameArg bool) {
	finishDebugLog := func(cmd *cobra.Command) {
		origRunE := cmd.RunE
		if origRunE != nil {
			cmd.RunE = func(cmd2 *cobra.Command, args []string) error {
				defer o.logger.DebugFunc("CommandRun").Finish()
				return origRunE(cmd2, args)
			}
		}
	}

	configureGlobal := cobrautil.WrapRunEForCmd(func(*cobra.Command, []string) error {
		o.UIFlags.ConfigureUI(o.ui)
		o.LoggerFlags.Configure(o.logger)
		o.KubeAPIFlags.Configure(o.configFactory)
		return nil
	})

	// Last one runs first
	// TODO: Add validation for number of arguments when positionalNameArg is true
	if positionalNameArg {
		cobrautil.VisitCommands(cmd, finishDebugLog, cobrautil.ReconfigureCmdWithSubcmd,
			configureGlobal, cobrautil.WrapRunEForCmd(cobrautil.ResolveFlagsForCmd))
	} else {
		cobrautil.VisitCommands(cmd, finishDebugLog, cobrautil.ReconfigureCmdWithSubcmd, configureGlobal, cobrautil.WrapRunEForCmd(cobrautil.ResolveFlagsForCmd))
	}
}

func AddPackageCommands(o *KctrlOptions, cmd *cobra.Command, flagsFactory cmdcore.FlagsFactory, opts cmdcore.PackageCommandTreeOpts) {
	cmd.AddCommand(pkgcreate.NewCreateCmd(pkgcreate.NewCreateOptions(o.logger, "kapp-a", opts), flagsFactory))
}

func AttachGlobalFlags(o *KctrlOptions, cmd *cobra.Command, flagsFactory cmdcore.FlagsFactory, opts cmdcore.PackageCommandTreeOpts) {
	SetGlobalFlags(o, cmd, flagsFactory, opts)
	ConfigurePathResolvers(o, cmd, flagsFactory)
	ConfigureGlobalFlags(o, cmd, flagsFactory, opts.PositionalArgs)
}

func AttachKctrlPackageCommandTree(cmd *cobra.Command, confUI *ui.ConfUI, opts cmdcore.PackageCommandTreeOpts) {
	configFactory := cmdcore.NewConfigFactoryImpl()
	depsFactory := cmdcore.NewDepsFactoryImpl(configFactory, confUI)
	options := NewKctrlOptions(confUI, configFactory, depsFactory)
	flagsFactory := cmdcore.NewFlagsFactory(configFactory, depsFactory)

	AddPackageCommands(options, cmd, flagsFactory, opts)
	AttachGlobalFlags(options, cmd, flagsFactory, opts)
}

type uiBlockWriter struct {
	ui ui.UI
}

var _ io.Writer = uiBlockWriter{}

func (w uiBlockWriter) Write(p []byte) (n int, err error) {
	w.ui.PrintBlock(p)
	return len(p), nil
}
