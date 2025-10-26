package main

import (
	"errors"
	"flag"
	"fmt"
	"log"

	jwtware "github.com/gofiber/contrib/jwt"
	swagger "github.com/gofiber/swagger"
	_ "github.com/kube-dash/kube-dash-backend/docs"

	"github.com/go-playground/validator/v10"
	fiber "github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"

	"github.com/kube-dash/kube-dash-backend/common"
	"github.com/kube-dash/kube-dash-backend/controller"
	"github.com/kube-dash/kube-dash-backend/database"
	httpapi "github.com/kube-dash/kube-dash-backend/httpapi"
)

type structValidator struct {
	validate *validator.Validate
}

func (v *structValidator) Validate(out any) error {

	err := v.validate.Struct(out)
	if err == nil {
		return nil
	}

	// maybe there was an actual error with validation
	// in this case also return error
	if _, ok := err.(*validator.InvalidValidationError); ok {
		return err
	}

	// take the first parameter error
	validationErrors := err.(validator.ValidationErrors)
	firstErr := validationErrors[0]

	return errors.New(firstErr.StructField())

}

func JWTErrorHandler(c fiber.Ctx, err error) error {

	return c.Status(fiber.StatusUnauthorized).JSON(
		fiber.Map{"error": err.Error()},
	)
}

// @title kube-dash backend
// @version 1.0
// @description Backend API of kube-dash
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email olokelo@gmail.com
// @license.name MIT
// @license.url https://opensource.org/license/mit
// @host localhost:5000
// @BasePath /
func main() {

	app := fiber.New(fiber.Config{
		StructValidator: &structValidator{validate: validator.New()},
	})
	app.Use(logger.New())
	app.Use(cors.New())

	kubeconfigPath := flag.String(
		"kubeconfig", "kube-config.yml",
		"The path to kubeconfig file of the cluster",
	)

	secretkeyPath := flag.String(
		"secretkey", "secret.key",
		"The path to server secret key used for JWT generation",
	)

	devMode := flag.Bool("dev", false, "Run in development mode")
	swagMode := flag.Bool("swag", false, "Register /swagger endpoint")

	flag.Parse()

	clientset, metricsset, err := controller.NewClientSet(*kubeconfigPath)
	if err != nil {
		log.Fatal("could not initialize clientset, " +
			"point -kubeconfig parameter to your kubernetes config.yml",
		)
	}

	err = common.InitSSK(*secretkeyPath)
	if err != nil {
		log.Fatal(
			"could not read secret key, " +
				"generate it with: openssl rand -hex 32 > secret.key",
		)
	}
	ssk, err := common.GetSSK()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	database.StartDBPodMetricsCleaner(db)
	controller.StartPodMetricsMonitor(metricsset, db)

	// make sure to close the DB when the main goes out of scope
	defer func() {
		dbIns, _ := db.DB()
		_ = dbIns.Close()
	}()

	// Login route
	app.Post("/api/v1/login", httpapi.ApiV1Login)

	app.Get("/api/v1/accessible", httpapi.ApiV1Accessible)

	// JWT Middleware
	if !(*devMode) {
		app.Use(jwtware.New(jwtware.Config{
			SigningKey:   jwtware.SigningKey{Key: ssk},
			ErrorHandler: JWTErrorHandler,
		}))
	}

	// Restricted Routes
	app.Get("/api/v1/restricted", httpapi.ApiV1Restricted)
	app.Get("/api/v1/listpods", httpapi.ApiV1ListPods(clientset))
	app.Get("/api/v2/listpods", httpapi.ApiV2ListPods(clientset))
	app.Get("/api/v1/listcontainers", httpapi.ApiV1ListContainers(clientset))
	app.Get("/api/v1/listnamespaces", httpapi.ApiV1ListNamespaces(clientset))

	app.Get("/api/v1/listdeployments", httpapi.ApiV1ListDeployments(clientset))
	app.Post("/api/v1/createdeployment", httpapi.ApiV1CreateDeployment(clientset))
	app.Post("/api/v1/updatedeployment", httpapi.ApiV1UpdateDeployment(clientset))
	app.Post("/api/v1/deletedeployment", httpapi.ApiV1DeleteDeployment(clientset))

	app.Get("/api/v1/getpodmetrics", httpapi.ApiV1GetPodMetrics(metricsset))
	app.Get("/api/v2/getpodmetrics", httpapi.ApiV2GetPodMetrics(metricsset, db))
	app.Post("/api/v1/deletepodmetrics", httpapi.ApiV1DeletePodMetrics(db))

	app.Post("/api/v1/createservice", httpapi.ApiV1CreateService(clientset))
	app.Get("/api/v1/listservices", httpapi.ApiV1ListServices(clientset))
	app.Post("/api/v1/deleteservice", httpapi.ApiV1DeleteService(clientset))

	if *swagMode {
		app.Get("/swagger/*", swagger.New(swagger.Config{
			DocExpansion: "none",
		}))
	}

	// 404 handler
	app.Use(func(c fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(&fiber.Map{
			"error": "not found",
		})
	})

	fmt.Println("K8S Dashboard Backend listening on port 5000")
	log.Fatal(app.Listen(":5000"))

}
