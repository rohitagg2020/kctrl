// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	cmdcore "github.com/rohitagg2020/kctrl/pkg/kctrl/cmd/core"
	"github.com/rohitagg2020/kctrl/pkg/kctrl/logger"
	"github.com/spf13/cobra"
)

type LoggerFlags struct {
	Debug bool
}

func (f *LoggerFlags) Set(cmd *cobra.Command, flagsFactory cmdcore.FlagsFactory) {
	cmd.PersistentFlags().BoolVar(&f.Debug, "debug", false, "Include debug output")
}

func (f *LoggerFlags) Configure(logger *logger.UILogger) {
	logger.SetDebug(f.Debug)
}
