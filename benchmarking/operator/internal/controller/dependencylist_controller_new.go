package controller

import (
	"context"
	"time"

	"github.com/go-logr/logr"
	"github.com/submariner-io/admiral/pkg/reporter"
	subctlclient "github.com/submariner-io/subctl/pkg/client"
	subctlservice "github.com/submariner-io/subctl/pkg/service"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/log"

	multiclusterv1beta1 "sha.ejaz/api/v1beta1"
)

// DependencyListReconcilerNew reconciles a DependencyList object
type DependencyListReconcilerNew struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=endpoints,verbs=get;list;watch
// +kubebuilder:rbac:groups=multicluster.my.domain,resources=dependencylists,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=multicluster.my.domain,resources=dependencylists/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=multicluster.my.domain,resources=dependencylists/finalizers,verbs=update

func (r *DependencyListReconcilerNew) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("Starting reconcile for DependencyList")

	var dependencyList multiclusterv1beta1.DependencyList
	if err := r.Get(ctx, req.NamespacedName, &dependencyList); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	clusterName := dependencyList.Spec.ClusterName
	clusterMap := dependencyList.Spec.ClusterServiceMap
	serviceList := clusterMap[clusterName]
	dependencies := dependencyList.Spec.Dependencies

	servicesForExport := []string{}
	seen := map[string]bool{}

	for _, dep := range dependencies {
		for _, svc := range serviceList {
			if svc == dep.DependsOn && !seen[svc] {
				servicesForExport = append(servicesForExport, svc)
				seen[svc] = true
			}
		}
	}

	if len(servicesForExport) == 0 {
		log.Info("No services to export")
		return ctrl.Result{}, nil
	}

	requeue, err := r.exportServices(ctx, servicesForExport, req.Namespace, log)
	if err != nil {
		return ctrl.Result{}, err
	}

	if requeue {
		return ctrl.Result{RequeueAfter: 5 * time.Second}, nil
	}

	log.Info("Successfully reconciled DependencyList")
	return ctrl.Result{}, nil
}

func (r *DependencyListReconcilerNew) exportServices(
	ctx context.Context,
	services []string,
	namespace string,
	log logr.Logger,
) (bool, error) {

	cfg, err := config.GetConfig()
	if err != nil {
		return false, err
	}

	producer, err := subctlclient.NewProducerFromRestConfig(cfg)
	if err != nil {
		return false, err
	}

	var svcList corev1.ServiceList
	if err := r.List(ctx, &svcList, client.InNamespace(namespace)); err != nil {
		return false, err
	}

	for _, logicalName := range services {

		var matchedService *corev1.Service

		for i := range svcList.Items {
			svc := &svcList.Items[i]

			// ✅ Helm default naming: <namespace>-<logical-service>
			if svc.Name == namespace+"-"+logicalName {
				matchedService = svc
				break
			}

			// ✅ Backward compatibility (manual/default namespace case)
			if svc.Name == logicalName {
				matchedService = svc
				break
			}
		}

		if matchedService == nil {
			log.Info("Service not yet created, waiting",
				"logicalService", logicalName,
				"namespace", namespace)
			return true, nil
		}

		err = subctlservice.Export(
			producer,
			namespace,
			matchedService.Name,
			reporter.Klog(),
		)
		if err != nil {
			return false, err
		}

		log.Info("Service export successful",
			"logicalService", logicalName,
			"actualService", matchedService.Name,
			"namespace", namespace)
	}

	return false, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DependencyListReconcilerNew) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&multiclusterv1beta1.DependencyList{}).
		Complete(r)
}
