package httpapi

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"

	"github.com/kube-dash/kube-dash-backend/common"
	"github.com/kube-dash/kube-dash-backend/models"
)

// @Summary Test unauthenticated endpoint
// @Description A check to see if user can reach public endpoints
// @Tags Test
// @Success 200   {object}  map[string]string   "Success"
// @Router /api/v1/accessible [get]
func ApiV1Accessible(c fiber.Ctx) error {
	return c.JSON(fiber.Map{"accessible": true})
}

// @Summary Test authenticated endpoint
// @Description A check to see if user can reach restricted endpoints
// @Tags Test
// @Security       ApiKeyAuth
// @Success 200   {object}  map[string]string   "Success"
// @Router /api/v1/restricted [get]
func ApiV1Restricted(c fiber.Ctx) error {

	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	if claims["usr"] == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{"error": "internal error"},
		)
	}

	username := claims["usr"].(string)
	return c.JSON(fiber.Map{"restricted": true, "username": username})
}

// @Summary Login endpoint
// @Description Returns a bearer token that has to be provided for authenticated endpoints
// @Tags Login
// @Param          request   query   models.LoginModel   false   "Query parameters"
// @Produce        json
// @Success 200
// @Failure 400
// @Failure 401
// @Failure 500
// @Router /api/v1/login [get]
func ApiV1Login(c fiber.Ctx) error {

	req := new(models.LoginModel)
	err := parseValidateParams(&c, req, true)
	if err != nil {
		// will return status bad request set in parseValidateBody
		return nil
	}

	// TODO: change hard coded user
	// Throws Unauthorized error
	if req.User != "john" || req.Pass != "doe" {
		return c.Status(fiber.StatusUnauthorized).JSON(
			fiber.Map{"error": "wrong credentials"},
		)
	}

	ssk, err := common.GetSSK()
	if err != nil {
		// something went very wrong as we don't even have secret key initialized
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{"error": "internal error"},
		)
	}

	// Create the Claims
	claims := jwt.MapClaims{
		"usr": req.User,
		"exp": time.Now().Add(time.Hour * 2).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString(ssk)
	if err != nil {
		// server is unable to sign token
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{"error": "internal error"},
		)
	}

	return c.JSON(fiber.Map{"token": t})
}
