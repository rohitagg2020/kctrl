// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package pkg

import (
	cmdcore "github.com/rohitagg2020/kctrl/pkg/kctrl/cmd/core"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "package",
		Aliases: []string{"pkg", "p"},
		Short:   "Package",
		Annotations: map[string]string{
			cmdcore.PackageHelpGroup.Key: cmdcore.PackageHelpGroup.Value,
		},
	}
	return cmd
}
