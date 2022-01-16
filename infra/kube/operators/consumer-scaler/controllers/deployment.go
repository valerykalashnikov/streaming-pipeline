package controllers

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (r *ScalerReconciler) scaleDeployment(dep *appsv1.Deployment, maxConsumers int32) error {
	patch := client.MergeFrom(dep.DeepCopy())
	newSize := int32(1)
	dep.Spec.Replicas = &newSize
	err := r.Patch(context.TODO(), dep, patch)
	return err
}
