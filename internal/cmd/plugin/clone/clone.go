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

// Package clone implements a command to clone a cluster
package clone

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/yaml"

	apierrs "k8s.io/apimachinery/pkg/api/errors"
	corev1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	apiv1 "github.com/cloudnative-pg/cloudnative-pg/api/v1"
	"github.com/cloudnative-pg/cloudnative-pg/internal/cmd/plugin"
	"github.com/cloudnative-pg/cloudnative-pg/pkg/reconciler/persistentvolumeclaim"
)

// Clone implements clone subcommand
func Clone(ctx context.Context, clusterName, newClusterName string) error {
	var cluster apiv1.Cluster

	// Get the Cluster object
	err := plugin.Client.Get(ctx, client.ObjectKey{Namespace: plugin.Namespace, Name: clusterName}, &cluster)
	if err != nil {
		return fmt.Errorf("retrieving cluster %s in namespace %s: %w", clusterName, plugin.Namespace, err)
	}

	if err := ensureClusterDoesNotExist(ctx, newClusterName); err != nil {
		return fmt.Errorf("checking for cluster %s in namespace %s: %w", newClusterName, plugin.Namespace, err)
	}

	if cluster.Status.CurrentPrimary == "" {
		return fmt.Errorf("cluster %s has no primary node", clusterName)
	}

	source, err := persistentvolumeclaim.GetInstanceStorageSource(ctx, plugin.Client,
		cluster.Status.CurrentPrimary, plugin.Namespace)
	if err != nil {
		return err
	}

	newCluster := apiv1.Cluster{
		ObjectMeta: corev1.ObjectMeta{
			Name:      newClusterName,
			Namespace: plugin.Namespace},
		Spec: cluster.Spec,
	}
	newCluster.Spec.Bootstrap = &apiv1.BootstrapConfiguration{
		Recovery: &apiv1.BootstrapRecovery{
			VolumeSnapshots: &apiv1.DataSource{
				Storage:           source.DataSource,
				WalStorage:        source.WALSource,
				TablespaceStorage: source.TablespaceSource,
			},
			Database: cluster.GetApplicationDatabaseName(),
			Owner:    cluster.GetApplicationDatabaseOwner(),
			Secret:   &apiv1.LocalObjectReference{Name: cluster.GetApplicationSecretName()},
		},
	}

	b, err := yaml.Marshal(newCluster)
	if err != nil {
		return err
	}
	fmt.Print(string(b))
	fmt.Println("---")

	if err := plugin.Client.Create(ctx, &newCluster); err != nil {
		return fmt.Errorf("creating cluster %s in namespace %s: %w", newClusterName, plugin.Namespace, err)
	}

	return nil
}

// ensureClusterDoesNotExist checks whether the cluster exists and returns an error if it does
func ensureClusterDoesNotExist(ctx context.Context, clusterName string) error {
	var cluster apiv1.Cluster
	err := plugin.Client.Get(
		ctx,
		types.NamespacedName{Name: clusterName, Namespace: plugin.Namespace},
		&cluster,
	)
	if err == nil {
		return fmt.Errorf("cluster already exists")
	}
	if !apierrs.IsNotFound(err) {
		return err
	}
	return nil
}
