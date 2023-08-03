package install

import (
	"context"
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// TODO: Verify all the components are running properly after deployment
func DeployTracerComponents(clientset *kubernetes.Clientset) error {
  var clusterRole rbacv1.ClusterRole
  var clusterRoleBinding rbacv1.ClusterRoleBinding
  var serviceAccount corev1.ServiceAccount
  var tracerPod corev1.Pod
  var tracerService corev1.Service

  installDir := "../../install/"

  clusterRoleMan, err := os.ReadFile(installDir+"clusterrole.yaml")
  if err != nil {
    return fmt.Errorf("Cannot read file: %s", err.Error())
  }
  err = yaml.Unmarshal(clusterRoleMan, &clusterRole)
  if err != nil {
    return fmt.Errorf("Cannot Unmarshal yaml file: %s", err.Error())
  }
  clusterRoleBindingMan, err := os.ReadFile(installDir+"clusterrolebinding.yaml")
  if err != nil {
    return fmt.Errorf("Cannot read file: %s", err.Error())
  }
  err = yaml.Unmarshal(clusterRoleBindingMan, &clusterRoleBinding)
  if err != nil {
    return fmt.Errorf("Cannot Unmarshal yaml file: %s", err.Error())
  }
  serviceAccountMan, err := os.ReadFile(installDir+"serviceaccount.yaml")
  if err != nil {
    return fmt.Errorf("Cannot read file: %s", err.Error())
  }
  err = yaml.Unmarshal(serviceAccountMan, &serviceAccount)
  if err != nil {
    return fmt.Errorf("Cannot Unmarshal yaml file: %s", err.Error())
  }
  tracerPodMan, err := os.ReadFile(installDir+"syscalls-tracer-pod.yaml")
  if err != nil {
    return fmt.Errorf("Cannot read file: %s", err.Error())
  }
  err = yaml.Unmarshal(tracerPodMan, &tracerPod)
  if err != nil {
    return fmt.Errorf("Cannot Unmarshal yaml file: %s", err.Error())
  }
  tracerServiceMan, err := os.ReadFile(installDir+"syscalls-tracer-service.yaml")
  if err != nil {
    return fmt.Errorf("Cannot read file: %s", err.Error())
  }
  err = yaml.Unmarshal(tracerServiceMan, &tracerService)
  if err != nil {
    return fmt.Errorf("Cannot Unmarshal yaml file: %s", err.Error())
  }


  fmt.Println("Deploying tracer components...")
  namespace := "kubesecgen"
  clientset.CoreV1().Namespaces().Create(context.TODO(), &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: namespace}}, metav1.CreateOptions{})
  _, err = clientset.RbacV1().ClusterRoles().Create(context.TODO(), &clusterRole, metav1.CreateOptions{})
  if err != nil{
    return fmt.Errorf("Cannot apply clusterrole: %s", err.Error())
  }
  _, err = clientset.RbacV1().ClusterRoleBindings().Create(context.TODO(), &clusterRoleBinding, metav1.CreateOptions{})
  if err != nil{
    return fmt.Errorf("Cannot apply clusterrolebinding: %s", err.Error())
  }
  _, err = clientset.CoreV1().ServiceAccounts(namespace).Create(context.TODO(), &serviceAccount, metav1.CreateOptions{})
  if err != nil{
    return fmt.Errorf("Cannot apply serviceaccount: %s", err.Error())
  }
  _, err = clientset.CoreV1().Pods(namespace).Create(context.TODO(), &tracerPod, metav1.CreateOptions{})
  if err != nil{
    return fmt.Errorf("Cannot apply tracer pod: %s", err.Error())
  }
  _, err = clientset.CoreV1().Services(namespace).Create(context.TODO(), &tracerService, metav1.CreateOptions{})
  if err != nil{
    return fmt.Errorf("Cannot apply tracer service: %s", err.Error())
  }
  return nil
}
