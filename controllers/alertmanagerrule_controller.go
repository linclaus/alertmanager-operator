/*


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

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	alertmanager "github.com/linclaus/alertmanager-operator/pkg/alertmanager"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	alertmanagerv1 "github.com/linclaus/alertmanager-operator/api/v1"
)

var LOG_FINALIZER = "alertmanagerRule"

// AlertmanagerRuleReconciler reconciles a AlertmanagerRule object
type AlertmanagerRuleReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=alertmanager.monitoring.io,resources=alertmanagerrules,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=alertmanager.monitoring.io,resources=alertmanagerrules/status,verbs=get;update;patch

func (r *AlertmanagerRuleReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	_ = r.Log.WithValues("AlertmanagerRule", req.NamespacedName)

	pr := &alertmanagerv1.AlertmanagerRule{}
	err := r.Get(ctx, req.NamespacedName, pr)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// AlertmanagerRule deleted
	if !pr.DeletionTimestamp.IsZero() {
		r.Log.V(1).Info("Deleting AlertmanagerRule")

		//delete rule
		err = alertmanager.DeleteAlertmanagerStrategy(*pr)
		if err != nil {
			r.updateStatus(pr, "Failed")
			return ctrl.Result{}, nil
		}

		//remove finalizer flag
		pr.Finalizers = removeString(pr.Finalizers, LOG_FINALIZER)
		if err = r.Update(ctx, pr); err != nil {
			return ctrl.Result{}, client.IgnoreNotFound(err)
		}
		r.Log.V(1).Info("AlertmanagerRule deleted")
		return ctrl.Result{}, nil
	}

	// Add finalizer
	if !containsString(pr.Finalizers, LOG_FINALIZER) {
		pr.Finalizers = append(pr.Finalizers, LOG_FINALIZER)
		if err = r.Update(ctx, pr); err != nil {
			return ctrl.Result{}, client.IgnoreNotFound(err)
		}
	}

	// AlertmanagerRule update
	r.Log.V(1).Info("Updating AlertmanagerRule")
	err = alertmanager.AddAlertmanagerStrategy(*pr)
	if err != nil {
		r.updateStatus(pr, "Failed")
		return ctrl.Result{}, nil
	}
	r.Log.V(1).Info("AlertmanagerRule updated")
	r.updateStatus(pr, "Successful")

	return ctrl.Result{}, nil
}

func (r *AlertmanagerRuleReconciler) updateStatus(pr *alertmanagerv1.AlertmanagerRule, status string) {
	pr.Status.Status = status
	if status == "Failed" {
		rty := pr.Status.RetryTimes
		if rty < 100 {
			pr.Status.RetryTimes = rty + 1
		}
	}
	r.Status().Update(context.Background(), pr)
}

func removeString(slice []string, s string) (result []string) {
	for _, item := range slice {
		if item == s {
			continue
		}
		result = append(result, item)
	}
	return
}

func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

func (r *AlertmanagerRuleReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&alertmanagerv1.AlertmanagerRule{}).
		Complete(r)
}
