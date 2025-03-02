package deploy

import (
	"context"
	"fmt"
	"os"

	"github.com/ThomasRubini/cloud-paas/internal/paas_backend/models"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

/*
Inputs:
- Deployement (in the cloud-pass sense, not the kubernetes sense)
  - Name
  - Image
  - TCP Port exposed inside the image
*/
func DeployApp(deployment models.DBEnvironment, imageTag string, namespace string, exposedPort int) error {
	/* TL;DR
	I. Build Deployement
		1. Get image
		2. Get namespace
	II. Build Service
		1. Get Listening port
	II. Build Ingress
		1. Get Host
		2. Route by header
	*/

	kubeConfPath := os.Getenv("KUBECONFIG")
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfPath)
	if err != nil {
		return fmt.Errorf("fail to build the k8s config. Error - %s", err)
	}

	// build the client set
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("fail to create the k8s client set. Error - %s", err)
	}

	var replicas int32 = 1
	var appName = deployment.ParentProject.Name + "-" + deployment.Environement

	deploymentsClient := clientSet.AppsV1().Deployments(apiv1.NamespaceDefault)

	kubeDeployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: appName,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": appName,
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": appName,
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  appName,
							Image: imageTag,
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: int32(exposedPort),
								},
							},
						},
					},
				},
			},
		},
	}

	// Create Deployment
	fmt.Println("Creating deployment...")
	result, err := deploymentsClient.Create(context.TODO(), kubeDeployment, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())

	panic("RUBEN CODE CA PLZ")

}
