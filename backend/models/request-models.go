package models

type LoginModel struct {
	// Username for authentication
	User string `json:"user" validate:"required"`
	// Password for authentication
	Pass string `json:"pass" validate:"required"`
}

type ListPodsV1RequestModel struct {
	// Namespace to filter pods
	Namespace string `query:"namespace" example:"default"`
}

type ListPodsV2RequestModel struct {
	// Namespace to filter pods
	Namespace string `query:"namespace" example:"default"`
}

type ListContainersRequestModel struct {
	// Namespace of the pod containing the containers
	Namespace string `query:"namespace" example:"default"`
	// Name of the pod containing the containers
	PodName string `query:"pod_name" example:"mypod"`
}

type ListDeploymentsRequestModel struct {
	// Namespace to filter deployments
	Namespace string `query:"namespace" example:"default"`
}

type CreateDeploymentRequestModel struct {
	// Namespace for the deployment
	Namespace string `json:"namespace" validate:"required" example:"default"`
	// Name for the deployment
	Name string `json:"name" validate:"required" example:"mydeployment"`
	// Docker image to use in the deployment
	Image string `json:"image" validate:"required" example:"nginx"`
	// Number of replicas for the deployment
	Replicas int32 `json:"replicas" validate:"gte=1,lte=32,required" example:"2"`
	// CPU request for each pod in the deployment (default: 100m)
	CPURequest string `json:"cpu_request" example:"100m"`
	// Memory request for each pod in the deployment (default: 256Mi)
	MemoryRequest string `json:"memory_request" example:"256Mi"`
	// CPU limit for each pod in the deployment (default: 200m)
	CPULimit string `json:"cpu_limit" example:"200m"`
	// Memory limit for each pod in the deployment (default: 512Mi)
	MemoryLimit string `json:"memory_limit" example:"512Mi"`
}

// TODO: code dup but don't know how to avoid it here

type UpdateDeploymentRequestModel struct {
	// Namespace for the deployment
	Namespace string `json:"namespace" validate:"required" example:"default"`
	// Name for the deployment
	Name string `json:"name" validate:"required" example:"mydeployment"`
	// Number of replicas for the deployment
	Replicas int32 `json:"replicas" validate:"gte=0,lte=32" example:"2"`
	// CPU request for each pod in the deployment
	CPURequest string `json:"cpu_request" example:"100m"`
	// Memory request for each pod in the deployment
	MemoryRequest string `json:"memory_request" example:"256Mi"`
	// CPU limit for each pod in the deployment
	CPULimit string `json:"cpu_limit" example:"200m"`
	// Memory limit for each pod in the deployment
	MemoryLimit string `json:"memory_limit" example:"512Mi"`
}

type DeleteDeploymentRequestModel struct {
	// Namespace of the deployment to delete
	Namespace string `json:"namespace" validate:"required" example:"default"`
	// Name of the deployment to delete
	Name string `json:"name" validate:"required" example:"mydeployment"`
}

type GetPodMetricsV1RequestModel struct {
	// Namespace to filter pod metrics
	Namespace string `query:"namespace" example:"default"`
}

type GetPodMetricsV2RequestModel struct {
	// Name of the pod for which to retrieve metrics
	PodName string `query:"pod_name" example:"mypod"`
	// Start time for metric collection in RFC3339 format
	StartTime string `query:"start_time" example:"2024-08-24T20:00:00.000Z"`
	// End time for metric collection in RFC3339 format
	EndTime string `query:"end_time" example:"2024-08-24T20:30:00.000Z"`
}

type DeletePodMetricsRequestModel struct {
	// Start time of metrics to delete in RFC3339 format
	StartTime string `json:"start_time" example:"2024-08-24T20:00:00.000Z"`
	// End time of metrics to delete in RFC3339 format
	EndTime string `json:"end_time" example:"2024-08-24T20:30:00.000Z"`
}

type CreateServiceRequestModel struct {
	// Namespace for the service
	Namespace string `json:"namespace" validate:"required" example:"default"`
	// Name for the service
	Name string `json:"name" validate:"required" example:"nginx-service"`
	// Type of service (ClusterIP or NodePort) - load balancers not yet supported
	Type string `json:"type" validate:"required,oneof=ClusterIP NodePort"`
	// Port of the service
	Port int32 `json:"port" validate:"required,gte=1,lte=65535" example:"80"`
	// External node port for NodePort service type
	NodePort int32 `json:"node_port" validate:"gte=30000,lte=32767" example:"30080"`
	// Target port on the container
	TargetPort int32 `json:"target_port" validate:"gte=1,lte=65535" example:"80"`
	// List of external IPs to expose the service on (not used yet)
	ExternalIPs []string `json:"external_ips" example:"10.1.2.30,10.2.2.30"`
	// Selector to match pods for the service. To create service for specific deployment use {"app": deploymentname}.
	Selector map[string]string `json:"selector" validate:"required" example:"{\"app\": \"nginx\"}"`
}

type ListServicesRequestModel struct {
	// Namespace to filter services
	Namespace string `query:"namespace" example:"default"`
}

type DeleteServiceRequestModel struct {
	// Namespace of the service to delete
	Namespace string `json:"namespace" validate:"required" example:"default"`
	// Name of the service to delete
	Name string `json:"name" validate:"required" example:"myservice"`
}
