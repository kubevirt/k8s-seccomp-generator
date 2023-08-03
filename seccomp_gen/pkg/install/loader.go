/*
 * This install package is aware of how to deploy the system in the kubernetes cluster.
 * As far as this package is concerned, it does not know about `falco`
 */
package install

import (
	"context"
	"fmt"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func DeployLoaderJob(clientset *kubernetes.Clientset, distro Distro) error {
	// apply the loader job to the cluster
	loaderJob := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: "loader-job",
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Name: "loader-pod"},
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{
						{Name: "var", VolumeSource: corev1.VolumeSource{
							HostPath: &corev1.HostPathVolumeSource{Path: "/var"}}},
						{Name: "usr", VolumeSource: corev1.VolumeSource{
							HostPath: &corev1.HostPathVolumeSource{Path: "/usr"}}},
						{Name: "lib", VolumeSource: corev1.VolumeSource{
							HostPath: &corev1.HostPathVolumeSource{Path: "/lib"}}},
						{Name: "etc", VolumeSource: corev1.VolumeSource{
							HostPath: &corev1.HostPathVolumeSource{Path: "/etc"}}},
					},
					Containers: []corev1.Container{
						{Name: "falco-loader", Image: "nithishdev/falco-loader:" + distro.String(), ImagePullPolicy: corev1.PullAlways,
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "var",
									MountPath: "/var",
								},
								{
									Name:      "usr",
									MountPath: "/usr",
								},
								{
									Name:      "lib",
									MountPath: "/lib",
								},
								{
									Name:      "etc",
									MountPath: "/etc",
								},
							},
							SecurityContext: &corev1.SecurityContext{Privileged: BoolPointer(true)},
						},
					},
					RestartPolicy: corev1.RestartPolicyNever,
				},
			},
		},
	}

	fmt.Println("Deploying the `loader job...")
	_, err := clientset.BatchV1().Jobs("default").Create(context.TODO(), loaderJob, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("Cannot deploy the loader job: %s", err.Error())
	}
	fmt.Println("Waiting for the `loader` job to complete. This might take up to 6 mins...")
	// wait till the falco-loader pod reaches the `completed` or `crash-loop backoff` state
  // TODO: Make sure there are no edge cases which leads to waiting forever
	waitForJobToComplete("loader-job", clientset)

	return nil
}
