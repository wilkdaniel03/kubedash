package httpapi

import (
	"errors"
	"reflect"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"

	"k8s.io/client-go/kubernetes"

	"github.com/kube-dash/kube-dash-backend/controller"
	"github.com/kube-dash/kube-dash-backend/database"
	"github.com/kube-dash/kube-dash-backend/models"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"
)

// function that checks if a given model has required fields
func hasRequiredFields(model interface{}) bool {
	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		// get tag value
		tag := field.Tag.Get("validate")
		if strings.Contains(tag, "required") {
			return true
		}
	}

	return false
}

// parse body or query parameters (depending on isBody), validate it agains the req struct
// and put it into the req struct if validation was successful
// otherwise make an bad request error
func parseValidateParams(c *fiber.Ctx, req interface{}, isBody bool) error {

	if isBody && len((*c).Request().Body()) == 0 {
		if hasRequiredFields(req) {
			(*c).Status(fiber.StatusBadRequest).JSON(
				fiber.Map{"error": "empty body"},
			)
			return errors.New("empty body")
		} else {
			return nil
		}
	}

	var err error
	if isBody {
		err = (*c).Bind().Body(req)
	} else {
		err = (*c).Bind().Query(req)
	}
	if err != nil {
		(*c).Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"error": "invalid or missing param", "param": err.Error()},
		)
	}

	return err

}

// make internal server error
func makeISE(c *fiber.Ctx, err error) {
	(*c).Status(fiber.StatusInternalServerError).JSON(
		fiber.Map{"error": err.Error()},
	)
}

// make bad request error
func makeBR(c *fiber.Ctx, err error) {
	(*c).Status(fiber.StatusBadRequest).JSON(
		fiber.Map{"error": err.Error()},
	)
}

// @Summary        List Available Pods (deprecated)
// @Description    Get all available pods in the cluster
// @Deprecated     true
// @Tags           Pods
// @Security       ApiKeyAuth
// @Param          request   query   models.ListPodsV1RequestModel   false   "Query parameters"
// @Produce        json
// @Success        200
// @Failure        400
// @Failure        401
// @Failure        500
// @Router         /api/v1/listpods [get]
func ApiV1ListPods(clientset *kubernetes.Clientset) fiber.Handler {
	return func(c fiber.Ctx) error {

		req := new(models.ListPodsV1RequestModel)
		err := parseValidateParams(&c, req, false)
		if err != nil {
			// will return status bad request set in parseValidateBody
			return nil
		}

		pods, err := controller.ListPodsV1(clientset, req.Namespace)
		if err != nil {
			makeISE(&c, err)
			return nil
		}

		return c.JSON(pods)
	}
}

// @Summary        List Available Pods
// @Description    Get all available pods in the cluster
// @Tags           Pods
// @Security       ApiKeyAuth
// @Param          request   query   models.ListPodsV2RequestModel   false   "Query parameters"
// @Produce        json
// @Success        200                {object}    models.ListPodsV2ResponseModel
// @Failure        400
// @Failure        401
// @Failure        500
// @Router         /api/v2/listpods [get]
func ApiV2ListPods(clientset *kubernetes.Clientset) fiber.Handler {
	return func(c fiber.Ctx) error {

		req := new(models.ListPodsV2RequestModel)
		err := parseValidateParams(&c, req, false)
		if err != nil {
			// will return status bad request set in parseValidateBody
			return nil
		}

		resp, err := controller.ListPodsV2(clientset, req)
		if err != nil {
			makeISE(&c, err)
			return nil
		}

		return c.JSON(resp)
	}
}

// @Summary        List Available Containers
// @Description    Get all available containers in the cluster
// @Tags           Containers
// @Security       ApiKeyAuth
// @Param          request   query   models.ListContainersRequestModel   false   "Query parameters"
// @Produce        json
// @Success        200                {object}    models.ListContainersReponseModel
// @Failure        400
// @Failure        401
// @Failure        500
// @Router         /api/v1/listcontainers [get]
func ApiV1ListContainers(clientset *kubernetes.Clientset) fiber.Handler {
	return func(c fiber.Ctx) error {

		req := new(models.ListContainersRequestModel)
		err := parseValidateParams(&c, req, false)
		if err != nil {
			// will return status bad request set in parseValidateBody
			return nil
		}

		resp, err := controller.ListContainers(clientset, req)
		if err != nil {
			makeISE(&c, err)
			return nil
		}

		return c.JSON(resp)
	}
}

// @Summary        List Available Namespaces
// @Description    Returns the list of all available namespaces in the cluster
// @Tags           Namespaces
// @Security       ApiKeyAuth
// @Produce        json
// @Success        200                {array}    string
// @Failure        400
// @Failure        401
// @Failure        500
// @Router         /api/v1/listnamespaces [get]
func ApiV1ListNamespaces(clientset *kubernetes.Clientset) fiber.Handler {
	return func(c fiber.Ctx) error {

		namespaces, err := controller.ListNamespaces(clientset)

		if err != nil {
			makeISE(&c, err)
			return nil
		}

		// warn: breaking change {"namespaces": namespaces} -> namespaces
		return c.JSON(namespaces)
	}
}

// @Summary        List All Deployments
// @Description    Get all deployments in the cluster
// @Tags           Deployment
// @Security       ApiKeyAuth
// @Param          request   query   models.ListDeploymentsRequestModel   false   "Query parameters"
// @Produce        json
// @Success        200                {object}    models.ListDeploymentsResponseModel
// @Failure        400
// @Failure        401
// @Failure        500
// @Router         /api/v1/listdeployments [get]
func ApiV1ListDeployments(clientset *kubernetes.Clientset) fiber.Handler {
	return func(c fiber.Ctx) error {

		req := new(models.ListDeploymentsRequestModel)
		err := parseValidateParams(&c, req, false)
		if err != nil {
			// will return status bad request set in parseValidateBody
			return nil
		}

		pods, err := controller.ListDeployments(clientset, req)
		if err != nil {
			makeISE(&c, err)
			return nil
		}

		return c.JSON(pods)
	}
}

// @Summary        Create New Deployment
// @Description    Create a new deployment in the cluster with the given name, namespace and parameters.
// @Tags           Deployment
// @Security       ApiKeyAuth
// @Accept         json
// @Param          request   body   models.CreateDeploymentRequestModel   true   "Request Model of Create Deployment"
// @Produce        json
// @Success        200
// @Failure        400
// @Failure        401
// @Failure        500
// @Router         /api/v1/createdeployment [post]
func ApiV1CreateDeployment(clientset *kubernetes.Clientset) fiber.Handler {
	return func(c fiber.Ctx) error {

		req := new(models.CreateDeploymentRequestModel)
		err := parseValidateParams(&c, req, true)
		if err != nil {
			// will return status bad request set in parseValidateBody
			return nil
		}

		_, err = controller.CreateDeployment(clientset, req.Namespace, req)

		if err != nil {
			makeISE(&c, err)
			return nil
		}

		// Return a success response
		return c.JSON(fiber.Map{"status": "deployment created"})
	}
}

// @Summary        Update Existing Deployment
// @Description    Update the parameters of already existing deployment.
// @Tags           Deployment
// @Security       ApiKeyAuth
// @Accept         json
// @Param          request   body   models.UpdateDeploymentRequestModel   true   "Request Model of Update Deployment"
// @Produce        json
// @Success        200
// @Failure        400
// @Failure        401
// @Failure        500
// @Router         /api/v1/updatedeployment [post]
func ApiV1UpdateDeployment(clientset *kubernetes.Clientset) fiber.Handler {
	return func(c fiber.Ctx) error {

		req := new(models.UpdateDeploymentRequestModel)
		err := parseValidateParams(&c, req, true)
		if err != nil {
			// will return status bad request set in parseValidateBody
			return nil
		}

		err = controller.UpdateDeployment(clientset, req)

		if err != nil {
			makeISE(&c, err)
			return nil
		}

		// Return a success response
		return c.JSON(fiber.Map{"status": "deployment updated"})
	}
}

// @Summary        Delete Deployment
// @Description    Removes the deployment by given name and namespace.
// @Tags           Deployment
// @Security       ApiKeyAuth
// @Accept         json
// @Param          request   body   models.DeleteDeploymentRequestModel   true   "Request Model of Delete Deployment"
// @Produce        json
// @Success        200
// @Failure        400
// @Failure        401
// @Failure        500
// @Router         /api/v1/deletedeployment [post]
func ApiV1DeleteDeployment(clientset *kubernetes.Clientset) fiber.Handler {
	return func(c fiber.Ctx) error {

		req := new(models.DeleteDeploymentRequestModel)
		err := parseValidateParams(&c, req, true)
		if err != nil {
			// will return status bad request set in parseValidateBody
			return nil
		}

		err = controller.DeleteDeployment(clientset, req)

		if err != nil {
			makeISE(&c, err)
			return nil
		}

		// Return a success response
		return c.JSON(fiber.Map{"status": "deployment deleted"})
	}
}

// @Summary        Get Pod Metrics (deprecated)
// @Description    Get metrics for specific pod or all pods in the cluster
// @Deprecated     true
// @Tags           Metrics
// @Security       ApiKeyAuth
// @Accept         json
// @Param          request   query   models.GetPodMetricsV1RequestModel   false   "Query parameters"
// @Produce        json
// @Success        200
// @Failure        400
// @Failure        401
// @Failure        500
// @Router         /api/v1/getpodmetrics [get]
func ApiV1GetPodMetrics(metricsset *metricsv.Clientset) fiber.Handler {
	return func(c fiber.Ctx) error {

		// Parse the user-submitted data from the request body
		req := new(models.GetPodMetricsV1RequestModel)
		err := parseValidateParams(&c, req, false)
		if err != nil {
			// will return status bad request set in parseValidateBody
			return nil
		}

		metrics, err := controller.GetPodMetricsV1(metricsset, req.Namespace)
		if err != nil {
			makeISE(&c, err)
			return nil
		}

		return c.JSON(metrics)
	}
}

// @Summary        Get Pod Metrics
// @Description    Get metrics for specific pod or all pods in the cluster
// @Tags           Metrics
// @Security       ApiKeyAuth
// @Param          request   query   models.GetPodMetricsV2RequestModel   false   "Query parameters"
// @Produce        json
// @Success        200                {object}    models.DBClusterMetricsModel
// @Failure        400
// @Failure        401
// @Failure        500
// @Router         /api/v2/getpodmetrics [get]
func ApiV2GetPodMetrics(metricsset *metricsv.Clientset, db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {

		// Parse the user-submitted data from the request body
		req := new(models.GetPodMetricsV2RequestModel)
		err := parseValidateParams(&c, req, false)
		if err != nil {
			// will return status bad request set in parseValidateBody
			return nil
		}

		// set default value for podname, if not specified will select all pods
		podname := "%"
		if req.PodName != "" {
			podname = req.PodName
		}

		// if they're not set, GetPodMetricsV2 will return the most recent record
		var startTime *time.Time = nil
		var endTime *time.Time = nil

		// if either start_time or end_time is set, we should validate them
		if req.StartTime != "" || req.EndTime != "" {

			// parse time
			parsedStartTime, err := time.Parse(time.RFC3339, req.StartTime)
			if err != nil {
				makeBR(&c, errors.New("unable to parse start_time"))
				return nil
			}

			parsedEndTime, err := time.Parse(time.RFC3339, req.EndTime)
			if err != nil {
				makeBR(&c, errors.New("unable to parse end_time"))
				return nil
			}

			startTime = &parsedStartTime
			endTime = &parsedEndTime

		}

		metrics, err := controller.GetPodMetricsV2(
			metricsset, db, podname, startTime, endTime,
		)
		if err != nil {
			makeISE(&c, err)
			return nil
		}

		// WARN: breaking change {"metrics": metrics} -> metrics
		return c.JSON(metrics)
	}
}

// @Summary        Delete Pod Metrics
// @Description    Delete metrics from the database. When no parameters are specified, all metrics will be deleted.
// @Tags           Metrics
// @Security       ApiKeyAuth
// @Accept         json
// @Param          request   body   models.DeletePodMetricsRequestModel   false   "Request Model of Deleting Pod Metrics"
// @Produce        json
// @Success        200   {object}  object  "Success"
// @Failure        400
// @Failure        401
// @Failure        500
// @Router         /api/v1/deletepodmetrics [post]
func ApiV1DeletePodMetrics(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {

		// Parse the user-submitted data from the request body
		req := new(models.DeletePodMetricsRequestModel)
		err := parseValidateParams(&c, req, true)
		if err != nil {
			// will return status bad request set in parseValidateBody
			return nil
		}

		var startTime *time.Time = nil
		var endTime *time.Time = nil

		if req.StartTime != "" {

			parsedStartTime, err := time.Parse(time.RFC3339, req.StartTime)
			if err != nil {
				makeBR(&c, errors.New("unable to parse start_time"))
				return nil
			}
			startTime = &parsedStartTime

		}

		if req.EndTime != "" {

			parsedEndTime, err := time.Parse(time.RFC3339, req.EndTime)
			if err != nil {
				makeBR(&c, errors.New("unable to parse end_time"))
				return nil
			}
			endTime = &parsedEndTime

		}

		err = database.DBDeletePodMetrics(db, startTime, endTime)
		if err != nil {
			makeISE(&c, err)
			return nil
		}

		return c.JSON(fiber.Map{
			"status": "OK",
		})
	}
}

// @Summary        Create Service
// @Description    Creates service in kubernetes cluster. Currently only NodePort and ClusterIP are supported, LoadBalancer coming later.
// @Tags           Services
// @Security       ApiKeyAuth
// @Accept         json
// @Param          request   body   models.CreateServiceRequestModel   false   "Request Model of Create Service"
// @Produce        json
// @Success        200   {object}  object  "Success"
// @Failure        400
// @Failure        401
// @Failure        500
// @Router         /api/v1/createservice [post]
func ApiV1CreateService(clientset *kubernetes.Clientset) fiber.Handler {
	return func(c fiber.Ctx) error {

		req := new(models.CreateServiceRequestModel)
		err := parseValidateParams(&c, req, true)
		if err != nil {
			return nil
		}

		err = controller.CreateService(clientset, req)

		if err != nil {
			makeISE(&c, err)
			return nil
		}

		// Return a success response
		return c.JSON(fiber.Map{"status": "service created"})
	}
}

// @Summary        List Available Services
// @Description    Get all available services in the cluster
// @Tags           Services
// @Security       ApiKeyAuth
// @Produce        json
// @Param          request   query   models.ListServicesRequestModel   false   "Query parameters"
// @Success        200                {object}    models.ListServicesResponseModel
// @Failure        400
// @Failure        401
// @Failure        500
// @Router         /api/v1/listservices [get]
func ApiV1ListServices(clientset *kubernetes.Clientset) fiber.Handler {
	return func(c fiber.Ctx) error {

		req := new(models.ListServicesRequestModel)
		err := parseValidateParams(&c, req, false)
		if err != nil {
			// will return status bad request set in parseValidateBody
			return nil
		}

		services, err := controller.ListServices(clientset, req)
		if err != nil {
			makeISE(&c, err)
			return nil
		}

		return c.JSON(services)
	}
}

// @Summary        Delete Service
// @Description    Removes the service by given name and namespace
// @Tags           Services
// @Security       ApiKeyAuth
// @Accept         json
// @Param          request   body   models.DeleteServiceRequestModel   false   "Request Model of Delete Service"
// @Produce        json
// @Success        200   {object}  object  "Success"
// @Failure        400
// @Failure        401
// @Failure        500
// @Router         /api/v1/deleteservice [post]
func ApiV1DeleteService(clientset *kubernetes.Clientset) fiber.Handler {
	return func(c fiber.Ctx) error {

		req := new(models.DeleteServiceRequestModel)
		err := parseValidateParams(&c, req, true)
		if err != nil {
			// will return status bad request set in parseValidateBody
			return nil
		}

		err = controller.DeleteService(clientset, req)
		if err != nil {
			makeISE(&c, err)
			return nil
		}

		return c.JSON(fiber.Map{"status": "service deleted"})
	}
}
