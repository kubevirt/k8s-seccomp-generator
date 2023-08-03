package install

import (
	"context"
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func BoolPointer(b bool) *bool {
	return &b
}

func waitForJobToComplete(jobName string, client *kubernetes.Clientset) {
	for range time.Tick(time.Second * 10) {
		done, err := getJobStatus(jobName, client)
		if err != nil {
			fmt.Println(err.Error())
		}
		if done {
			return
		}
	}
}

// bool indicates whether or not we should stop waiting on the job. If the bool is true, then we should
// stop waiting on the job.
func getJobStatus(jobName string, k8sClient *kubernetes.Clientset) (bool, error) {
	job, err := k8sClient.BatchV1().Jobs("default").Get(context.Background(), jobName, metav1.GetOptions{})
	if err != nil {
		return false, err
	}
	if job.Status.Active == 0 && job.Status.Succeeded == 0 && job.Status.Failed == 0 {
		return false, fmt.Errorf("%s hasn't started yet", job.Name)
	}
	if job.Status.Active > 0 {
		return false, fmt.Errorf("%s is still running", job.Name)
	}
	if job.Status.Succeeded > 0 {
		return true, nil // Job ran successfully
	}
	return true, fmt.Errorf("%s has failed with error", job.Name)
}
