package controllers

import (
	"context"

	"github.com/valerykalashnikov/streaming-pipeline/infra/kube/operators/consumer-scaler/consumer"
	appsv1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (r *ScalerReconciler) scaleDeployment(dep *appsv1.Deployment, calc *consumer.Calculator, maxConsumers int32) error {
	patch := client.MergeFrom(dep.DeepCopy())
	newSize := calc.ReplicaSize()
	if newSize > maxConsumers {
		newSize = maxConsumers
	}
	size := int32(newSize)
	dep.Spec.Replicas = &size
	err := r.Patch(context.TODO(), dep, patch)
	return err
}
