/*
Copyright 2023.

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

package controller

import (
	"context"

	"github.com/submariner-io/admiral/pkg/reporter"
	subctlclient "github.com/submariner-io/subctl/pkg/client"
	subctlservice "github.com/submariner-io/subctl/pkg/service"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/log"

	multiclusterv1beta1 "sha.ejaz/api/v1beta1"
	"sigs.k8s.io/mcs-api/pkg/apis/v1alpha1"
)

// DependencyListReconciler reconciles a DependencyList object
type DependencyListReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=multicluster.my.domain,resources=dependencylists,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=multicluster.my.domain,resources=dependencylists/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=multicluster.my.domain,resources=dependencylists/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the DependencyList object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.15.0/pkg/reconcile
func (r *DependencyListReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	log.Info("dependency list reconcile")

	dependencyList := multiclusterv1beta1.DependencyList{}

	if err := r.Client.Get(context.Background(), req.NamespacedName, &dependencyList); err != nil {
		if client.IgnoreNotFound(err) != nil {
			log.Error(err, "unable to get dependency list")
			return ctrl.Result{}, err
		}
	}

	log.Info("dependency list", "dependencyList", dependencyList)

	clusterName := dependencyList.Spec.ClusterName

	clusterMap := dependencyList.Spec.ClusterServiceMap
	serviceList := clusterMap[clusterName]

	dependencies := dependencyList.Spec.Dependencies
	servicesForExport := []string{}
	servicesLookup := map[string]bool{}

	for _, dependency := range dependencies {
		dependsOn := dependency.DependsOn
		for _, service := range serviceList {
			if service == dependsOn && !servicesLookup[service] {
				servicesForExport = append(servicesForExport, service)
				servicesLookup[service] = true
			}
		}
	}

	err := exportServices(servicesForExport, req.NamespacedName.Namespace)
	if err != nil {
		log.Error(err, "unable to export")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func exportServices(services []string, namespace string) error {
	config, err := config.GetConfig()
	if err != nil {
		return err
	}

	producer, err := subctlclient.NewProducerFromRestConfig(config)
	if err != nil {
		return err
	}

	scheme.Scheme.AddKnownTypes(v1alpha1.SchemeGroupVersion, &v1alpha1.ServiceExport{})

	for _, service := range services {
		err = subctlservice.Export(producer, namespace, service, reporter.Klog())
		if err != nil {
			return err
		}
	}

	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DependencyListReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&multiclusterv1beta1.DependencyList{}).
		Complete(r)
}
