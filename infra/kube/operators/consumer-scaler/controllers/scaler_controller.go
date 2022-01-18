/*
Copyright 2022.

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
	"fmt"
	"math/rand"
	"time"

	"github.com/adjust/rmq"
	pipev1alhpa1 "github.com/valerykalashnikov/streaming-pipeline/infra/kube/operators/consumer-scaler/api/v1alpha1"
	"github.com/valerykalashnikov/streaming-pipeline/infra/kube/operators/consumer-scaler/consumer"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// ScalerReconciler reconciles a Scaler object
type ScalerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=streaming-pipeline.my.domain,resources=scalers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=streaming-pipeline.my.domain,resources=scalers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=streaming-pipeline.my.domain,resources=scalers/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Scaler object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile
func (r *ScalerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx).WithValues("Scaler", req.NamespacedName)

	// Fetch the Traveller instance
	scaler := &pipev1alhpa1.Scaler{}
	err := r.Get(context.TODO(), req.NamespacedName, scaler)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// Check if this Deployment already exists
	found := &appsv1.Deployment{}

	err = r.Get(context.TODO(), types.NamespacedName{
		Name:      scaler.Spec.Deployment.Name,
		Namespace: scaler.Namespace,
	}, found)
	if err != nil {
		log.Error(err, "Deployment Not ready - requeue")
		return ctrl.Result{}, err
	}
	calc := consumer.NewCalculator(scaler.Spec.Queue, r.createRMQConnection(scaler))
	err = r.scaleDeployment(found, calc, scaler.Spec.Deployment.MaxConsumers)
	if err != nil {
		log.Error(err, "Deployment Not ready - trying again")
		return ctrl.Result{}, err
	}
	// Scale up and down every 15 minutes
	result := ctrl.Result{
		RequeueAfter: 15 * time.Minute,
	}
	return result, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ScalerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&pipev1alhpa1.Scaler{}).
		Complete(r)
}

func (r *ScalerReconciler) createRMQConnection(scaler *pipev1alhpa1.Scaler) rmq.Connection {
	rand.Seed(time.Now().UnixNano())
	// generate the postfix with a lenght of 5 bytes
	b := make([]byte, 5)
	rand.Read(b)

	connName := "scaler" + fmt.Sprintf("%x", b)[:5]

	return rmq.OpenConnection(connName, "tcp", scaler.Spec.BrokerAddress, 2)
}
