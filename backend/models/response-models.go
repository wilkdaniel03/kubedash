package models

type ListPodsV2ResponseModelPod struct {
	// The name of the pod.
	Name string `json:"name" example:"nginx-deploy-59849dcb58-tdknv"`
	// The namespace of the pod.
	Namespace string `json:"namespace" example:"default"`
	// The status of the pod.
	Status string `json:"status" example:"Running"`
}

type ListPodsV2ResponseModel struct {
	// A slice of ListPodsV2ResponseModelPod objects representing the pods.
	Pods []ListPodsV2ResponseModelPod `json:"pods"`
}

type ListDeploymentsResponseModelDeployment struct {
	// The namespace of the deployment.
	Namespace string `json:"namespace" example:"default"`
	// The name of the deployment.
	Name string `json:"name" example:"nginx-deployment"`
	// The number of replicas in the deployment.
	Replicas int32 `json:"replicas" example:"3"`
	// The number of ready replicas in the deployment.
	ReadyReplicas int32 `json:"ready_replicas" example:"3"`
	// The number of updated replicas in the deployment.
	UpdatedReplicas int32 `json:"updated_replicas" example:"3"`
	// The number of unavailable replicas in the deployment.
	UnavailableReplicas int32 `json:"unavailable_replicas" example:"0"`
	// The creation time of the deployment.
	CreationTime string `json:"creation_time" example:"2024-08-24T20:00:00.000Z"`
}

type ListDeploymentsResponseModel struct {
	// A slice of ListDeploymentsResponseModelDeployment objects representing the deployments.
	Deployments []ListDeploymentsResponseModelDeployment `json:"deployments"`
}

type ListContainersReponseModelPort struct {
	// The name of the port.
	Name string `json:"name" example:"myport"`
	// The protocol used by the port.
	Protocol string `json:"protocol" example:"TCP"`
	// The container port number.
	ContainerPort int32 `json:"container_port" example:"9512"`
}
type ListContainersReponseModelResources struct {
	// The CPU request for the container (e.g., "100m").
	CPURequest string `json:"cpu_request" example:"100m"`
	// The memory request for the container (e.g., "512Mi").
	MemoryRequest string `json:"memory_request" example:"512Mi"`
	// The CPU limit for the container (e.g., "200m").
	CPULimit string `json:"cpu_limit" example:"200m"`
	// The memory limit for the container (e.g., "1Gi").
	MemoryLimit string `json:"memory_limit" example:"1Gi"`
}
type ListContainersReponseModelContainer struct {
	// The name of the container.
	Name string `json:"name" example:"nginx-deployment"`
	// The namespace of the container.
	Namespace string `json:"namespace" example:"default"`
	// The image used by the container.
	Image string `json:"image" example:"nginx"`
	// Resource limits and requests for the container.
	Resources ListContainersReponseModelResources `json:"resources"`
	// A list of ListContainersReponseModelPort objects representing port mappings for the container.
	Ports []ListContainersReponseModelPort `json:"ports"`
}
type ListContainersReponseModel struct {
	// A list of ListContainersReponseModelContainer objects representing the containers.
	Containers []ListContainersReponseModelContainer `json:"containers"`
}

// TODO: maybe merge with ListContainersReponseModelPort
type ListServicesResponseModelPort struct {
	// The name of the port.
	Name string `json:"name" example:"myport"`
	// The protocol used by the port. Optional.
	Protocol string `json:"protocol,omitempty" example:"TCP"`
	// The service port number. Optional.
	Port int32 `json:"port,omitempty" example:"8080"`
	// The target port number for the service. Optional.
	TargetPort int32 `json:"target_port,omitempty" example:"8080"`
	// The node port number for the service. Optional.
	NodePort int32 `json:"node_port,omitempty" example:"30030"`
}

// ListServicesResponseModelService represents a single service in Kubernetes.
type ListServicesResponseModelService struct {
	// The name of the service.
	Name string `json:"name" example:"myservice"`
	// The namespace of the service.
	Namespace string `json:"namespace" example:"default"`
	// The type of the service.
	Type string `json:"type" example:"NodePort"`
	// A map of key-value pairs used to select pods for the service.
	Selector map[string]string `json:"selector" example:"{\"app\": \"nginx\"}"`
	// The cluster IP addresses assigned to the service.
	ClusterIPs []string `json:"cluster_ips" example:"10.1.2.30,10.2.2.30"`
	// A slice of ListServicesResponseModelPort objects representing port mappings for the service.
	Ports []ListServicesResponseModelPort `json:"ports"`
}

type ListServicesResponseModel struct {
	// A list of ListServicesResponseModelService containing services data.
	Services []ListServicesResponseModelService `json:"services"`
}
