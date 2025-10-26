package controller

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"

	// kubernetes api -> coreapiv1, metaapiv1, don't confuse with internal apiv1
	"gorm.io/gorm"
	appsapiv1 "k8s.io/api/apps/v1"
	coreapiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metaapiv1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"

	"github.com/kube-dash/kube-dash-backend/models"
)

// validate CPURequest/CPULimit
func validateCPU(request string) bool {

	// will match: 0.1, 100, 100m but also 0.1m
	cpuPattern := "^\\d+(\\.\\d+)?(m)?$"
	match, err := regexp.MatchString(cpuPattern, request)
	if err != nil {
		return false
	}
	return match
}

// validate MemoryRequest/MemoryLimit
func validateMemory(request string) bool {

	// will match: 4294967296, 4Gi, 4G, 2048M, 2048Mi etc
	memoryPattern := "^\\d+(?:[EPTGMK]i?)?$"
	match, err := regexp.MatchString(memoryPattern, request)
	if err != nil {
		return false
	}
	return match
}

func NewClientSet(kubeconfigPath string) (*kubernetes.Clientset, *metricsv.Clientset, error) {

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, nil, fmt.Errorf("error building kubeconfig: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating clientset: %v", err)
	}

	metricsset, err := metricsv.NewForConfig(config)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating metricsset: %v", err)
	}

	return clientset, metricsset, nil
}

func ListPodsV1(clientset *kubernetes.Clientset, namespace string) (*coreapiv1.PodList, error) {

	pods, err := clientset.CoreV1().Pods(namespace).List(
		context.TODO(), metaapiv1.ListOptions{},
	)

	return pods, err

}

func ListPodsV2(
	clientset *kubernetes.Clientset,
	req *models.ListPodsV2RequestModel,
) (models.ListPodsV2ResponseModel, error) {

	pods, err := ListPodsV1(clientset, req.Namespace)

	resp := models.ListPodsV2ResponseModel{}
	resp.Pods = []models.ListPodsV2ResponseModelPod{}
	for _, poddata := range pods.Items {
		currentPod := models.ListPodsV2ResponseModelPod{}
		currentPod.Name = poddata.ObjectMeta.Name
		currentPod.Namespace = poddata.ObjectMeta.Namespace
		currentPod.Status = string(poddata.Status.Phase)
		resp.Pods = append(resp.Pods, currentPod)
	}

	return resp, err

}

func ListContainers(
	clientset *kubernetes.Clientset,
	req *models.ListContainersRequestModel,
) (models.ListContainersReponseModel, error) {

	pods, err := ListPodsV1(clientset, req.Namespace)

	resp := models.ListContainersReponseModel{}
	resp.Containers = []models.ListContainersReponseModelContainer{}
	for _, poddata := range pods.Items {

		// skip container if it doesn't belong to the desired Pod
		if req.PodName != "" && req.PodName != poddata.ObjectMeta.Name {
			continue
		}

		for _, containerdata := range poddata.Spec.Containers {

			currentContainer := models.ListContainersReponseModelContainer{}
			currentContainer.Namespace = poddata.ObjectMeta.Namespace
			currentContainer.Name = containerdata.Name
			currentContainer.Image = containerdata.Image

			// add ports of the container
			currentContainer.Ports = []models.ListContainersReponseModelPort{}
			for _, portdata := range containerdata.Ports {
				currentContainer.Ports = append(
					currentContainer.Ports,
					models.ListContainersReponseModelPort{
						Name:          portdata.Name,
						Protocol:      string(portdata.Protocol),
						ContainerPort: portdata.ContainerPort,
					},
				)
			}

			// add resources of the container
			currentResources := models.ListContainersReponseModelResources{}
			currentResources.CPULimit = containerdata.Resources.Limits.Cpu().String()
			currentResources.CPURequest = containerdata.Resources.Requests.Cpu().String()
			currentResources.MemoryLimit = containerdata.Resources.Limits.Memory().String()
			currentResources.MemoryRequest = containerdata.Resources.Requests.Memory().String()

			currentContainer.Resources = currentResources

			resp.Containers = append(resp.Containers, currentContainer)
		}
	}

	return resp, err

}

func ListNamespaces(clientset *kubernetes.Clientset) ([]string, error) {
	// get a list of all namespaces using the Kubernetes API
	namespaceList, err := clientset.CoreV1().Namespaces().List(
		context.TODO(), metaapiv1.ListOptions{},
	)
	if err != nil {
		return nil, errors.New("error getting namespace list")
	}

	// extract the namespace names from the list and add them to a slice
	var namespaces []string
	for _, ns := range namespaceList.Items {
		namespaces = append(namespaces, ns.Name)
	}

	// return a JSON array containing the namespace names
	return namespaces, nil
}

func ListDeployments(
	clientset *kubernetes.Clientset,
	req *models.ListDeploymentsRequestModel,
) (models.ListDeploymentsResponseModel, error) {

	deploymentsClient := clientset.AppsV1().Deployments(req.Namespace)

	deployments, err := deploymentsClient.List(
		context.TODO(), metaapiv1.ListOptions{},
	)
	if err != nil {
		panic(err.Error())
	}

	resp := models.ListDeploymentsResponseModel{}
	resp.Deployments = []models.ListDeploymentsResponseModelDeployment{}
	for _, deploymentdata := range deployments.Items {
		currentDeployment := models.ListDeploymentsResponseModelDeployment{}
		currentDeployment.Namespace = deploymentdata.Namespace
		currentDeployment.Name = deploymentdata.Name
		currentDeployment.Replicas = deploymentdata.Status.Replicas
		currentDeployment.ReadyReplicas = deploymentdata.Status.ReadyReplicas
		currentDeployment.UpdatedReplicas = deploymentdata.Status.UpdatedReplicas
		currentDeployment.UnavailableReplicas = deploymentdata.Status.UnavailableReplicas
		currentDeployment.CreationTime =
			deploymentdata.CreationTimestamp.Time.UTC().Format(time.RFC3339)
		resp.Deployments = append(resp.Deployments, currentDeployment)
	}

	return resp, err

}

func CreateDeployment(
	clientset *kubernetes.Clientset, namespace string,
	req *models.CreateDeploymentRequestModel,
) (*appsapiv1.Deployment, error) {

	paramCPURequest := req.CPURequest
	if paramCPURequest == "" {
		paramCPURequest = "100m"
	}
	paramMemoryRequest := req.MemoryRequest
	if paramMemoryRequest == "" {
		paramMemoryRequest = "256Mi"
	}
	paramCPULimit := req.CPULimit
	if paramCPULimit == "" {
		paramCPULimit = "200m"
	}
	paramMemoryLimit := req.MemoryLimit
	if paramMemoryLimit == "" {
		paramMemoryLimit = "512Mi"
	}

	// make sure given quantites are correct
	// we don't want server panic
	if !(validateCPU(paramCPURequest) && validateCPU(paramCPULimit)) {
		return nil, fmt.Errorf("invalid CPU Resource format")
	}
	if !(validateMemory(paramMemoryRequest) && validateMemory(paramMemoryLimit)) {
		return nil, fmt.Errorf("invalid Memory Resource format")
	}

	CPURequest, err := resource.ParseQuantity(paramCPURequest)
	if err != nil {
		return nil, err
	}
	MemoryRequest, err := resource.ParseQuantity(paramMemoryRequest)
	if err != nil {
		return nil, err
	}
	CPULimit, err := resource.ParseQuantity(paramCPULimit)
	if err != nil {
		return nil, err
	}
	MemoryLimit, err := resource.ParseQuantity(paramMemoryLimit)
	if err != nil {
		return nil, err
	}

	deployment := &appsapiv1.Deployment{
		ObjectMeta: metaapiv1.ObjectMeta{
			Name: req.Name,
		},
		Spec: appsapiv1.DeploymentSpec{
			Replicas: &req.Replicas,
			Selector: &metaapiv1.LabelSelector{
				MatchLabels: map[string]string{
					"app": req.Name,
				},
			},
			Template: coreapiv1.PodTemplateSpec{
				ObjectMeta: metaapiv1.ObjectMeta{
					Labels: map[string]string{
						"app": req.Name,
					},
				},
				Spec: coreapiv1.PodSpec{
					Containers: []coreapiv1.Container{
						{
							Name:  req.Name,
							Image: req.Image,
							Resources: coreapiv1.ResourceRequirements{
								Requests: map[coreapiv1.ResourceName]resource.Quantity{
									coreapiv1.ResourceCPU:    CPURequest,
									coreapiv1.ResourceMemory: MemoryRequest,
								},
								Limits: map[coreapiv1.ResourceName]resource.Quantity{
									coreapiv1.ResourceCPU:    CPULimit,
									coreapiv1.ResourceMemory: MemoryLimit,
								},
							},
						},
					},
				},
			},
		},
	}

	// create the deployment in k8s using the client-go library
	dplmnt, err := clientset.AppsV1().Deployments(namespace).Create(
		context.TODO(), deployment, metaapiv1.CreateOptions{},
	)
	if err != nil {
		return nil, err
	}

	return dplmnt, nil

}

func UpdateDeployment(
	clientset *kubernetes.Clientset,
	req *models.UpdateDeploymentRequestModel,
) error {

	deployment, err := clientset.AppsV1().Deployments(req.Namespace).Get(
		context.TODO(), req.Name, metaapiv1.GetOptions{},
	)
	if err != nil {
		return err
	}

	// Update the resource requirements of each container in the deployment
	for i := range deployment.Spec.Template.Spec.Containers {

		// Request resources
		if req.CPURequest != "" {
			if !validateCPU(req.CPURequest) {
				return fmt.Errorf("invalid CPU Resource format")
			}
			CPURequest, err := resource.ParseQuantity(req.CPURequest)
			if err != nil {
				return err
			}
			deployment.Spec.Template.Spec.Containers[i].Resources.Requests["cpu"] = CPURequest
		}
		if req.MemoryRequest != "" {
			if !validateMemory(req.MemoryRequest) {
				return fmt.Errorf("invalid Memory Resource format")
			}
			MemoryRequest, err := resource.ParseQuantity(req.MemoryRequest)
			if err != nil {
				return err
			}
			deployment.Spec.Template.Spec.Containers[i].Resources.Requests["memory"] = MemoryRequest
		}

		// Limit resources
		if req.CPULimit != "" {
			if !validateCPU(req.CPULimit) {
				return fmt.Errorf("invalid CPU Resource format")
			}
			CPULimit, err := resource.ParseQuantity(req.CPULimit)
			if err != nil {
				return err
			}
			deployment.Spec.Template.Spec.Containers[i].Resources.Limits["cpu"] = CPULimit
		}
		if req.MemoryLimit != "" {
			if !validateMemory(req.MemoryLimit) {
				return fmt.Errorf("invalid Memory Resource format")
			}
			MemoryLimit, err := resource.ParseQuantity(req.MemoryLimit)
			if err != nil {
				return err
			}
			deployment.Spec.Template.Spec.Containers[i].Resources.Limits["memory"] = MemoryLimit
		}
	}

	if req.Replicas != 0 {
		*deployment.Spec.Replicas = req.Replicas
	}

	// Update the deployment in k8s
	_, err = clientset.AppsV1().Deployments(req.Namespace).Update(
		context.TODO(), deployment, metaapiv1.UpdateOptions{},
	)
	if err != nil {
		return err
	}

	return nil

}

func DeleteDeployment(
	clientset *kubernetes.Clientset,
	req *models.DeleteDeploymentRequestModel,
) error {

	deploymentsClient := clientset.AppsV1().Deployments(req.Namespace)
	deletePolicy := metaapiv1.DeletePropagationForeground
	err := deploymentsClient.Delete(context.TODO(), req.Name, metaapiv1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
	if err != nil {
		return err
	}

	return nil

}

func GetPodMetricsV1(
	metricsset *metricsv.Clientset, namespace string,
) (*v1beta1.PodMetricsList, error) {

	metrics, err := metricsset.MetricsV1beta1().PodMetricses(namespace).List(
		context.TODO(), metaapiv1.ListOptions{},
	)
	if err != nil {
		return nil, err
	}

	return metrics, err

}

// this saves the pod metrics to a db
func savePodMetricsToDB(metricsset *metricsv.Clientset, db *gorm.DB) error {

	// get raw metrics for all namespaces
	metrics, err := GetPodMetricsV1(metricsset, "")
	if err != nil {
		return err
	}
	clusterMetricsRecord := models.DBClusterMetricsModel{}

	for _, podMetrics := range metrics.Items {

		podMetricsRecord := models.DBPodMetricsModel{
			Name: podMetrics.Name,
		}

		for _, cont := range podMetrics.Containers {
			containerMetricsRecord := models.DBContainerMetricsModel{
				Name:        cont.Name,
				CPUUsage:    cont.Usage.Cpu().MilliValue(),
				MemoryUsage: cont.Usage.Memory().Value(),
			}
			podMetricsRecord.Containers = append(
				podMetricsRecord.Containers, containerMetricsRecord,
			)
			//fmt.Println("Added container: ", cont.Name)
		}

		clusterMetricsRecord.Pods = append(
			clusterMetricsRecord.Pods, podMetricsRecord,
		)
	}

	db.Create(&clusterMetricsRecord)

	return nil
}

// start monitoring the pod	metrics periodically every 5 seconds
// and save the data to SQL database
func StartPodMetricsMonitor(
	metricsset *metricsv.Clientset, db *gorm.DB,
) error {

	// create ticker with ticks every 5 seconds
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for range ticker.C {
			err := savePodMetricsToDB(metricsset, db)
			if err != nil {
				// TODO: find a way to error handle the ticker
				return
			}
		}
	}()

	return nil
}

func GetPodMetricsV2(
	metricsset *metricsv.Clientset, db *gorm.DB,
	podname string,
	starttime *time.Time, endtime *time.Time,
) ([]models.DBClusterMetricsModel, error) {

	var clusterMetricsRecords []models.DBClusterMetricsModel
	//fmt.Println(podname)

	var dbtx *gorm.DB
	if starttime == nil && endtime == nil {

		// only most recent record
		dbtx = db.Order("created_at DESC").Limit(1)

	} else {

		// validate time range (max 4 hours for all and 72 hours for single pod)
		// this is mainly to not return enourmous sized json responses
		maxHours := 4 * time.Hour
		if podname != "%" {
			maxHours = 72 * time.Hour
		}
		if endtime.Sub(*starttime) > maxHours {
			return nil, errors.New("time range too wide")
		}

		// validate if starttime is not after than endtime
		if starttime.After(*endtime) {
			return nil, errors.New("start_time cannot be after end_time")
		}

		// all records within set timeframe
		dbtx = db.Where("created_at BETWEEN ? AND ?", starttime, endtime)
	}

	// podname doesn't have validation so it will return list of metrics records
	// but the list of pods will be empty
	err := dbtx.Preload("Pods", "name LIKE ?", podname).
		Preload("Pods.Containers").
		Find(&clusterMetricsRecords).
		Error
	if err != nil {
		return nil, err
	}

	return clusterMetricsRecords, nil
}

func CreateService(
	clientset *kubernetes.Clientset,
	req *models.CreateServiceRequestModel,
) error {

	service := &coreapiv1.Service{
		ObjectMeta: metaapiv1.ObjectMeta{
			Name: req.Name,
		},
		Spec: coreapiv1.ServiceSpec{
			Type: coreapiv1.ServiceType(req.Type),
			Ports: []coreapiv1.ServicePort{{
				Port: req.Port,
				TargetPort: intstr.IntOrString{
					Type: intstr.Int, IntVal: req.TargetPort},
				NodePort: req.NodePort,
			}},
			// loadbalancer may need this in the future
			ExternalIPs: req.ExternalIPs,
			// label selectors to match pods
			Selector: req.Selector,
		},
	}

	created_service, err := clientset.CoreV1().Services(req.Namespace).Create(
		context.TODO(), service, metaapiv1.CreateOptions{},
	)
	if err != nil {
		return err
	}

	//fmt.Println(created_service)
	_ = created_service

	return nil
}

func ListServices(
	clientset *kubernetes.Clientset,
	req *models.ListServicesRequestModel,
) (models.ListServicesResponseModel, error) {

	services, err := clientset.CoreV1().Services(req.Namespace).List(
		context.TODO(), metaapiv1.ListOptions{},
	)

	resp := models.ListServicesResponseModel{}
	resp.Services = []models.ListServicesResponseModelService{}
	for _, servicedata := range services.Items {
		currentService := models.ListServicesResponseModelService{}
		currentService.Name = servicedata.ObjectMeta.Name
		currentService.Namespace = servicedata.ObjectMeta.Namespace
		currentService.Type = string(servicedata.Spec.Type)
		currentService.Selector = servicedata.Spec.Selector
		currentService.ClusterIPs = servicedata.Spec.ClusterIPs
		//servicedata.Spec.ExternalIPs

		currentService.Ports = []models.ListServicesResponseModelPort{}
		for _, portdata := range servicedata.Spec.Ports {
			currentService.Ports = append(
				currentService.Ports,
				models.ListServicesResponseModelPort{
					Name:       portdata.Name,
					Protocol:   string(portdata.Protocol),
					Port:       portdata.Port,
					NodePort:   portdata.NodePort,
					TargetPort: portdata.TargetPort.IntVal,
				},
			)
		}

		resp.Services = append(resp.Services, currentService)
	}

	return resp, err

}

func DeleteService(
	clientset *kubernetes.Clientset,
	req *models.DeleteServiceRequestModel,
) error {

	err := clientset.CoreV1().Services(req.Namespace).Delete(
		context.TODO(), req.Name, metaapiv1.DeleteOptions{},
	)
	return err

}
