/*
Copyright 2023 The KubeStellar Authors.

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

package vcluster

import (
	"context"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	clog "sigs.k8s.io/controller-runtime/pkg/log"

	tenancyv1alpha1 "github.com/kubestellar/kubeflex/api/v1alpha1"
	"github.com/kubestellar/kubeflex/pkg/reconcilers/shared"
	"github.com/kubestellar/kubeflex/pkg/util"
)

const (
	ServiceName      = "vcluster"
	ServicePort      = 443
	kubeconfigSecret = "TODO"
)

// VClusterReconciler reconciles a OCM ControlPlane
type VClusterReconciler struct {
	*shared.BaseReconciler
}

func New(cl client.Client, scheme *runtime.Scheme, version string, clientSet *kubernetes.Clientset, dynamicClient *dynamic.DynamicClient, eventRecorder record.EventRecorder) *VClusterReconciler {
	return &VClusterReconciler{
		BaseReconciler: &shared.BaseReconciler{
			Client:        cl,
			Scheme:        scheme,
			ClientSet:     clientSet,
			DynamicClient: dynamicClient,
			EventRecorder: eventRecorder,
		},
	}
}

func (r *VClusterReconciler) Reconcile(ctx context.Context, hcp *tenancyv1alpha1.ControlPlane) (ctrl.Result, error) {
	var routeURL string
	_ = clog.FromContext(ctx)

	cfg, err := r.BaseReconciler.GetConfig(ctx)
	if err != nil {
		return r.UpdateStatusForSyncingError(hcp, err)
	}

	if err := r.BaseReconciler.ReconcileNamespace(ctx, hcp); err != nil {
		if util.IsTransientError(err) {
            return ctrl.Result{}, err
        }
		return r.UpdateStatusForSyncingError(hcp, err)
	}

	if cfg.IsOpenShift {
		if err = r.ReconcileAPIServerRoute(ctx, hcp, ServiceName, shared.SecurePort, cfg.Domain); err != nil {
			if util.IsTransientError(err) {
				return ctrl.Result{}, err
			}
			return r.UpdateStatusForSyncingError(hcp, err)
		}
		routeURL, err = r.GetAPIServerRouteURL(ctx, hcp)
		if err != nil {
			return r.UpdateStatusForSyncingError(hcp, err)
		}
		// re-queue until valid route URL is retrieved
		if routeURL == "" {
			return ctrl.Result{RequeueAfter: 3 * time.Second}, nil
		}
		cfg.ExternalURL = routeURL
	} else {
		if err := r.ReconcileAPIServerIngress(ctx, hcp, ServiceName, ServicePort, cfg.Domain); err != nil {
			if util.IsTransientError(err) {
				return ctrl.Result{}, err
			}
			return r.UpdateStatusForSyncingError(hcp, err)
		}
	}

	if err := r.ReconcileChart(ctx, hcp, cfg); err != nil {
		if util.IsTransientError(err) {
			return ctrl.Result{}, err
		}
		return r.UpdateStatusForSyncingError(hcp, err)
	}

	if err := r.ReconcileNodePortService(ctx, hcp); err != nil {
		if util.IsTransientError(err) {
			return ctrl.Result{}, err
		}
		return r.UpdateStatusForSyncingError(hcp, err)
	}

	if err := r.addOwnerReference(ctx, hcp); err != nil {
		if util.IsTransientError(err) {
			return ctrl.Result{}, err
		}
		return r.UpdateStatusForSyncingError(hcp, err)
	}

	if err := r.ReconcileUpdateClusterInfoJobRole(ctx, hcp); err != nil {
		if util.IsTransientError(err) {
			return ctrl.Result{}, err
		}
		return r.UpdateStatusForSyncingError(hcp, err)
	}

	if err := r.ReconcileUpdateClusterInfoJobRoleBinding(ctx, hcp); err != nil {
		if util.IsTransientError(err) {
			return ctrl.Result{}, err
		}
		return r.UpdateStatusForSyncingError(hcp, err)
	}

	if err := r.ReconcileUpdateClusterInfoJob(ctx, hcp, cfg, r.Version); err != nil {
		if util.IsTransientError(err) {
			return ctrl.Result{}, err
		}
		return r.UpdateStatusForSyncingError(hcp, err)
	}

	r.UpdateStatusWithSecretRef(hcp, util.VClusterKubeConfigSecret,
		util.KubeconfigSecretKeyVCluster, util.KubeconfigSecretKeyVClusterInCluster)

	if hcp.Spec.PostCreateHook != nil &&
		tenancyv1alpha1.HasConditionAvailable(hcp.Status.Conditions) {
		if err := r.ReconcileUpdatePostCreateHook(ctx, hcp); err != nil {
			if util.IsTransientError(err) {
				return ctrl.Result{}, err
			}
			return r.UpdateStatusForSyncingError(hcp, err)
		}
	}

	// update kubeconfig secret to add the incluster config
	if tenancyv1alpha1.HasConditionAvailable(hcp.Status.Conditions) {
		if err := r.ReconcileKubeconfigSecret(ctx, hcp); err != nil {
			if util.IsTransientError(err) {
				return ctrl.Result{}, err
			}
			return r.UpdateStatusForSyncingError(hcp, err)
		}
	}

	return r.UpdateStatusForSyncingSuccess(ctx, hcp)
}

// add owner ref to allow capturing lifecycle events for the OCM deployment
func (r *VClusterReconciler) addOwnerReference(ctx context.Context, hcp *tenancyv1alpha1.ControlPlane) error {
	namespace := util.GenerateNamespaceFromControlPlaneName(hcp.Name)
	statefulset := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      util.GetAPIServerDeploymentNameByControlPlaneType(string(hcp.Spec.Type)),
			Namespace: namespace,
		},
	}
	if err := r.Client.Get(context.TODO(), client.ObjectKeyFromObject(statefulset), statefulset, &client.GetOptions{}); err != nil {
		return err
	}

	if err := controllerutil.SetControllerReference(hcp, statefulset, r.Scheme); err != nil {
		return err
	}

	if err := r.Client.Update(context.TODO(), statefulset, &client.UpdateOptions{}); err != nil {
		return err
	}

	return nil
}
