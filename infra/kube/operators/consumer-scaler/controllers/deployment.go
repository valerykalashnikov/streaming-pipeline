package controllers

import (
	"context"

	pipev1alhpa1 "github.com/valerykalashnikov/streaming-pipeline/infra/kube/operators/consumer-scaler/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func labels() map[string]string {
	// Fetches and sets labels

	return map[string]string{
		"app": "consumer",
	}
}

func (r *ScalerReconciler) ensureDeployment(request reconcile.Request,
	scaler *pipev1alhpa1.Scaler,
	dep *appsv1.Deployment) (*reconcile.Result, error) {
	// See if deployment already exists and create if it doesn't
	found := &appsv1.Deployment{}
	err := r.Get(context.TODO(), types.NamespacedName{
		Name:      dep.Name,
		Namespace: scaler.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {
		// Create the deployment

		err = r.Create(context.TODO(), dep)
		if err != nil {
			// Deployment failed
			return &reconcile.Result{}, err
		} else {
			// Deployment was successful
			return nil, nil
		}
	} else if err != nil {
		// Error that isn't due to the deployment not existing
		return &reconcile.Result{}, err
	}

	return nil, nil
}

func (r *ScalerReconciler) backendDeployment(v *pipev1alhpa1.Scaler) *appsv1.Deployment {
	// scaler := consumer.NewScaler("")
	// size := scaler.ReplicaSize()
	v.Spec.Deployment.Spec.Template.Labels = labels()
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "consumer-deployment",
			Namespace: v.Namespace,
		},
		Spec: v.Spec.Deployment.Spec,
	}

	dep.Spec.Template.Spec.Overhead.Pods()
	controllerutil.SetControllerReference(v, dep, r.Scheme)
	return dep
}
