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

package clone

import (
	"context"
	"github.com/spf13/cobra"

	"github.com/cloudnative-pg/cloudnative-pg/internal/cmd/plugin"
)

// NewCmd create the new "clone" subcommand
func NewCmd() *cobra.Command {
	cloneCmd := &cobra.Command{
		Use:   "clone [cluster] [new-cluster]",
		Short: "clone the cluster into a new one",
		Args:  plugin.RequiresArguments(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			clusterName := args[0]
			newClusterName := args[1]
			return Clone(ctx, clusterName, newClusterName)
		},
	}

	return cloneCmd
}
