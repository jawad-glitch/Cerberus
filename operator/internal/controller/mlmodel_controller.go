package controller

import (
	"context"
	mlv1alpha1 "github.com/jawad-glitch/Cerberus/operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type MLModelReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=ml.cerberus.io,resources=mlmodels,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=ml.cerberus.io,resources=mlmodels/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=ml.cerberus.io,resources=mlmodels/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete

func (r *MLModelReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Step 1: Fetch the MLModel
	mlmodel := &mlv1alpha1.MLModel{}
	if err := r.Get(ctx, req.NamespacedName, mlmodel); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	logger.Info("Reconciling MLModel", "name", mlmodel.Name)

	// Step 2: Check if Deployment exists
	found := &appsv1.Deployment{}
	err := r.Get(ctx, types.NamespacedName{
		Name:      mlmodel.Name,
		Namespace: mlmodel.Namespace,
	}, found)

	if errors.IsNotFound(err) {
		// Step 3: Create the Deployment
		dep := r.buildDeployment(mlmodel)
		logger.Info("Creating Deployment", "name", dep.Name)
		if err := r.Create(ctx, dep); err != nil {
			logger.Error(err, "Failed to create Deployment")
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		return ctrl.Result{}, err
	}

	currentImage := found.Spec.Template.Spec.Containers[0].Image
	desiredImage := mlmodel.Spec.Image

	if currentImage != desiredImage {
		logger.Info("Image changed, updating deployment",
			"current", currentImage,
			"desired", desiredImage)
		found.Spec.Template.Spec.Containers[0].Image = desiredImage
		if err := r.Update(ctx, found); err != nil {
			logger.Error(err, "Failed to update Deployment")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

func (r *MLModelReconciler) buildDeployment(m *mlv1alpha1.MLModel) *appsv1.Deployment {
	replicas := m.Spec.Replicas
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
			Labels:    map[string]string{"app": m.Name, "managed-by": "cerberus"},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": m.Name},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": m.Name},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  m.Name,
							Image: m.Spec.Image,
							Ports: []corev1.ContainerPort{
								{ContainerPort: m.Spec.Port},
							},
						},
					},
				},
			},
		},
	}
}

func (r *MLModelReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&mlv1alpha1.MLModel{}).
		Complete(r)
}
