/*
Copyright The CloudNativePG Contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package upgrade implements the "instance upgrade" subcommand of the operator
package upgrade

import (
	"github.com/spf13/cobra"

	"github.com/cloudnative-pg/cloudnative-pg/internal/cmd/manager/instance/upgrade/execute"
	"github.com/cloudnative-pg/cloudnative-pg/internal/cmd/manager/instance/upgrade/prepare"
)

// NewCmd creates the "instance upgrade" subcommand
func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "upgrade",
	}

	cmd.AddCommand(prepare.NewCmd())
	cmd.AddCommand(execute.NewCmd())

	return cmd
}